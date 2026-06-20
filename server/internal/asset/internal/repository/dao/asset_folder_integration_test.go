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

func TestFindAssetIDPageById(t *testing.T) {
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

	database := mongox.NewClient(client, &mongox.Config{}).NewDatabase("fnote_asset_pagination_test")
	t.Cleanup(func() { _ = database.Database().Drop(context.Background()) })
	folderID := bson.NewObjectID()
	assetIDs := []bson.ObjectID{bson.NewObjectID(), bson.NewObjectID(), bson.NewObjectID(), bson.NewObjectID(), bson.NewObjectID()}
	_, err = database.Database().Collection("asset_folders").InsertOne(ctx, bson.M{
		"_id":    folderID,
		"name":   "images",
		"assets": assetIDs,
	})
	if err != nil {
		t.Fatal(err)
	}

	page, err := NewAssetFolderDao(database).FindAssetIDPageById(ctx, folderID, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 5 {
		t.Fatalf("total = %d, want 5", page.Total)
	}
	if len(page.Assets) != 2 || page.Assets[0] != assetIDs[2] || page.Assets[1] != assetIDs[3] {
		t.Fatalf("assets = %#v, want %#v", page.Assets, assetIDs[2:4])
	}
}
