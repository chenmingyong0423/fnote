package dao

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"
	"github.com/chenmingyong0423/go-mongox/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestAIConfigUpsertUsesExistingConfigsCollection(t *testing.T) {
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
	defer client.Disconnect(context.Background())
	db := mongox.NewClient(client, &mongox.Config{}).NewDatabase("fnote_ai_config_test_" + bson.NewObjectID().Hex())
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_ = db.Database().Drop(cleanupCtx)
	})
	dao := NewWebsiteConfigDao(db)
	cfg := domain.AIAssistantConfig{
		AssistantName: "博客助手", EnabledTools: []string{},
		Providers: []domain.AIProvider{{Id: "p1", ApiKey: "never-persist-plaintext", EncryptedApiKey: "encrypted-value"}},
	}
	if err = dao.UpsertAIAssistantConfig(ctx, cfg); err != nil {
		t.Fatal(err)
	}
	first, err := dao.FindByTyp(ctx, domain.AIAssistantConfigType)
	if err != nil || first.CreatedAt.IsZero() {
		t.Fatal("upsert must initialize config timestamps")
	}
	if err = dao.UpsertAIAssistantConfig(ctx, cfg); err != nil {
		t.Fatal("saving an unchanged config must succeed", err)
	}
	second, err := dao.FindByTyp(ctx, domain.AIAssistantConfigType)
	if err != nil || !second.CreatedAt.Equal(first.CreatedAt) {
		t.Fatal("updating config must preserve creation time")
	}
	collection := db.Database().Collection("configs")
	count, err := collection.CountDocuments(ctx, bson.M{"typ": domain.AIAssistantConfigType})
	if err != nil || count != 1 {
		t.Fatal("repeated saves must update a single document in configs")
	}
	raw, err := collection.FindOne(ctx, bson.M{"typ": domain.AIAssistantConfigType}).Raw()
	if err != nil || bytes.Contains(raw, []byte("never-persist-plaintext")) || !bytes.Contains(raw, []byte("encrypted-value")) {
		t.Fatal("stored BSON must include ciphertext and exclude the transient API key")
	}
}
