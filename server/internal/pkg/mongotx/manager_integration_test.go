package mongotx

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/chenmingyong0423/go-mongox/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestManagerCommitAndRollback(t *testing.T) {
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
	defer client.Disconnect(ctx)

	databaseName := "fnote_transaction_test"
	database := mongox.NewClient(client, &mongox.Config{}).NewDatabase(databaseName)
	t.Cleanup(func() { _ = database.Database().Drop(context.Background()) })
	manager := NewManager(database)
	collection := database.Database().Collection("records")

	rollbackErr := errors.New("rollback")
	err = manager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if _, insertErr := collection.InsertOne(txCtx, bson.M{"value": "rolled-back"}); insertErr != nil {
			return insertErr
		}
		return rollbackErr
	})
	if !errors.Is(err, rollbackErr) {
		t.Fatalf("expected rollback error, got %v", err)
	}
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("rollback left %d documents", count)
	}

	err = manager.WithinTransaction(ctx, func(txCtx context.Context) error {
		_, insertErr := collection.InsertOne(txCtx, bson.M{"value": "committed"})
		return insertErr
	})
	if err != nil {
		t.Fatalf("commit transaction: %v", err)
	}
	count, err = collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("commit stored %d documents, want 1", count)
	}
}
