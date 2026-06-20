package dao

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/chenmingyong0423/go-mongox/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestDeleteUnusedByFileId(t *testing.T) {
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

	database := mongox.NewClient(client, &mongox.Config{}).NewDatabase("fnote_file_management_test")
	t.Cleanup(func() { _ = database.Database().Drop(context.Background()) })
	collection := database.Database().Collection("file_meta")
	unusedFileID := []byte{0x00, 0x11, 0x22, 0x33}
	usedFileID := []byte{0x44, 0x55, 0x66, 0x77}
	_, err = collection.InsertMany(ctx, []any{
		bson.M{"file_id": unusedFileID, "used_in": bson.A{}},
		bson.M{"file_id": usedFileID, "used_in": bson.A{bson.M{"entity_id": "post-1", "entity_type": "post"}}},
	})
	if err != nil {
		t.Fatal(err)
	}

	dao := NewFileDao(database)
	deleted, err := dao.DeleteUnusedByFileId(ctx, unusedFileID)
	if err != nil || deleted != 1 {
		t.Fatalf("delete unused file = (%d, %v), want (1, nil)", deleted, err)
	}
	deleted, err = dao.DeleteUnusedByFileId(ctx, usedFileID)
	if err != nil || deleted != 0 {
		t.Fatalf("delete used file = (%d, %v), want (0, nil)", deleted, err)
	}
}
