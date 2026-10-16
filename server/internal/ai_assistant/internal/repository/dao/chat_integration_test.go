package dao

import (
	"context"
	"errors"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/domain"
	"github.com/chenmingyong0423/go-mongox/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestMongoChatAtomicGenerationCheckpointsAndRestartRecovery(t *testing.T) {
	uri := os.Getenv("MONGODB_TEST_URI")
	if uri == "" {
		t.Skip("MONGODB_TEST_URI is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Disconnect(context.Background()) })
	db := mongox.NewClient(client, &mongox.Config{}).NewDatabase("fnote_ai_chat_test_" + bson.NewObjectID().Hex())
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = db.Database().Drop(cleanup)
	})
	d, err := NewAIAssistantDao(db)
	if err != nil {
		t.Fatal(err)
	}
	if err = d.Identify(ctx, "visitor", "192.0.2.1"); err != nil {
		t.Fatal(err)
	}
	if err = d.Identify(ctx, "visitor", "192.0.2.2"); err != nil {
		t.Fatal(err)
	}
	var visitor domain.Visitor
	if err = db.Database().Collection("ai_visitors").FindOne(ctx, bson.M{"_id": "visitor"}).Decode(&visitor); err != nil || visitor.FirstIP != "192.0.2.1" || visitor.LastIP != "192.0.2.2" {
		t.Fatal("visitor identity and first/latest IP must be retained", err)
	}
	c := domain.Conversation{ID: "conversation", VisitorID: "visitor", Title: "Test", Messages: []domain.Message{}, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err = d.Create(ctx, c); err != nil {
		t.Fatal(err)
	}
	var winners atomic.Int32
	var group sync.WaitGroup
	for _, request := range []string{"r1", "r2"} {
		group.Add(1)
		go func(id string) {
			defer group.Done()
			err := d.Begin(ctx, "visitor", c.ID, domain.Message{ID: id, RequestID: id, Role: "user", Content: "问题", Status: "complete"}, domain.Message{ID: "answer-" + id, Role: "assistant", Status: "typing"})
			if err == nil {
				winners.Add(1)
			} else if !errors.Is(err, domain.ErrBusy) {
				t.Error(err)
			}
		}(request)
	}
	group.Wait()
	if winners.Load() != 1 {
		t.Fatal("exactly one concurrent request may claim the document")
	}
	stored, err := d.Get(ctx, "visitor", c.ID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = d.Get(ctx, "other-visitor", c.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatal("ownership belongs in every Mongo query")
	}
	answer := stored.Messages[1]
	answer.Content = "部分回答"
	answer.References = []domain.Reference{{Title: "文章", Href: "/posts/public"}}
	if err = d.SaveAnswer(ctx, "visitor", c.ID, answer, false); err != nil {
		t.Fatal(err)
	}
	_, err = NewAIAssistantDao(db)
	if err != nil {
		t.Fatal(err)
	}
	stored, err = d.Get(ctx, "visitor", c.ID)
	if err != nil || stored.Generation.Status != "interrupted" || stored.Messages[1].Status != "interrupted" || stored.Messages[1].Content != "部分回答" || len(stored.Messages[1].References) != 1 {
		t.Fatal("restart must retain checkpoints and clear in-flight state", err)
	}
	if err = d.Begin(ctx, "visitor", c.ID, domain.Message{ID: "next", RequestID: "next", Role: "user"}, domain.Message{ID: "next-answer", Role: "assistant", Status: "typing"}); err != nil {
		t.Fatal("conversation must continue after restart", err)
	}
	answer = domain.Message{ID: "next-answer", Role: "assistant", Content: "完成", Status: "complete"}
	if err = d.SaveAnswer(ctx, "visitor", c.ID, answer, true); err != nil {
		t.Fatal(err)
	}
	if err = d.SaveAnswer(ctx, "visitor", c.ID, answer, true); err != nil {
		t.Fatal("retrying finalization after an ambiguous database response must be idempotent", err)
	}
	if err = d.Begin(ctx, "visitor", c.ID, domain.Message{RequestID: "next"}, domain.Message{ID: "duplicate"}); !errors.Is(err, domain.ErrBusy) {
		t.Fatal("duplicate requests must be rejected atomically")
	}
	list, err := d.List(ctx, "visitor")
	if err != nil || len(list) != 1 || len(list[0].Messages) != 0 {
		t.Fatal("history list should not load every message", err)
	}
}
