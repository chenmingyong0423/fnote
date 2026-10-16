package dao

import (
	"context"
	"errors"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/domain"
	"github.com/chenmingyong0423/go-mongox/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IAIAssistantDao interface {
	Identify(context.Context, string, string) error
	Create(context.Context, domain.Conversation) error
	Count(context.Context, string) (int64, error)
	List(context.Context, string) ([]domain.Conversation, error)
	Get(context.Context, string, string) (*domain.Conversation, error)
	Begin(context.Context, string, string, domain.Message, domain.Message) error
	SaveAnswer(context.Context, string, string, domain.Message, bool) error
}

type AIAssistantDao struct {
	visitors      *mongo.Collection
	conversations *mongo.Collection
}

func NewAIAssistantDao(db *mongox.Database) (*AIAssistantDao, error) {
	d := &AIAssistantDao{db.Database().Collection("ai_visitors"), db.Database().Collection("ai_conversations")}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := d.conversations.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{Key: "visitor_id", Value: 1}, {Key: "updated_at", Value: -1}}})
	if err != nil {
		return nil, err
	}
	// Single-instance deployment: in-flight answers left by the previous process
	// are recoverable history, never a permanent generation lock.
	_, err = d.conversations.UpdateMany(ctx, bson.M{"generation.status": "running"}, bson.M{"$set": bson.M{
		"generation.status": "interrupted", "messages.$[answer].status": "interrupted",
		"messages.$[answer].error": "服务重启，回答已中断，可继续提问", "updated_at": time.Now(),
	}}, options.UpdateMany().SetArrayFilters([]any{bson.M{"answer.status": "typing"}}))
	return d, err
}

func (d *AIAssistantDao) Identify(ctx context.Context, id, ip string) error {
	now := time.Now()
	_, err := d.visitors.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set":         bson.M{"last_ip": ip, "updated_at": now},
		"$setOnInsert": bson.M{"first_ip": ip, "created_at": now},
	}, options.UpdateOne().SetUpsert(true))
	return err
}

func (d *AIAssistantDao) Create(ctx context.Context, c domain.Conversation) error {
	_, err := d.conversations.InsertOne(ctx, c)
	return err
}

func (d *AIAssistantDao) Count(ctx context.Context, visitor string) (int64, error) {
	return d.conversations.CountDocuments(ctx, bson.M{"visitor_id": visitor})
}

func (d *AIAssistantDao) List(ctx context.Context, visitor string) ([]domain.Conversation, error) {
	cur, err := d.conversations.Find(ctx, bson.M{"visitor_id": visitor}, options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}}).SetLimit(100).SetProjection(bson.M{"messages": 0}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	result := []domain.Conversation{}
	err = cur.All(ctx, &result)
	return result, err
}

func (d *AIAssistantDao) Get(ctx context.Context, visitor, id string) (*domain.Conversation, error) {
	var c domain.Conversation
	err := d.conversations.FindOne(ctx, bson.M{"_id": id, "visitor_id": visitor}).Decode(&c)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrNotFound
	}
	return &c, err
}

func (d *AIAssistantDao) Begin(ctx context.Context, visitor, id string, user, answer domain.Message) error {
	filter := bson.M{"_id": id, "visitor_id": visitor, "generation.status": bson.M{"$ne": "running"},
		"messages.request_id": bson.M{"$ne": user.RequestID}, "$expr": bson.M{"$lte": bson.A{bson.M{"$size": "$messages"}, 98}}}
	changes := bson.M{"$push": bson.M{"messages": bson.M{"$each": bson.A{user, answer}}}, "$set": bson.M{
		"generation": domain.Generation{RequestID: user.RequestID, MessageID: answer.ID, Status: "running"}, "updated_at": time.Now(),
	}}
	result, err := d.conversations.UpdateOne(ctx, filter, changes)
	if err == nil && result.MatchedCount == 0 {
		return domain.ErrBusy
	}
	return err
}

func (d *AIAssistantDao) SaveAnswer(ctx context.Context, visitor, id string, answer domain.Message, final bool) error {
	fields := bson.M{"messages.$[answer].content": answer.Content, "messages.$[answer].references": answer.References,
		"messages.$[answer].status": answer.Status, "messages.$[answer].error": answer.Error, "updated_at": time.Now()}
	if final {
		fields["generation.status"] = answer.Status
	}
	filter := bson.M{"_id": id, "visitor_id": visitor, "generation.message_id": answer.ID, "generation.status": "running"}
	if final {
		filter["generation.status"] = bson.M{"$in": bson.A{"running", answer.Status}}
	}
	result, err := d.conversations.UpdateOne(ctx, filter, bson.M{"$set": fields},
		options.UpdateOne().SetArrayFilters([]any{bson.M{"answer.id": answer.ID}}))
	if err == nil && result.MatchedCount == 0 {
		return domain.ErrBusy
	}
	return err
}
