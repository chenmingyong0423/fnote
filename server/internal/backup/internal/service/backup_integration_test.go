package service

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/chenmingyong0423/go-mongox/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestFileIDBackupAndRecovery(t *testing.T) {
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

	databaseName := "fnote_backup_file_id_test"
	database := mongox.NewClient(client, &mongox.Config{}).NewDatabase(databaseName)
	t.Cleanup(func() { _ = database.Database().Drop(context.Background()) })
	want := []byte{0x00, 0x11, 0x22, 0x33}
	_, err = database.Database().Collection("file_meta").InsertOne(ctx, bson.M{
		"file_id":    want,
		"file_name":  "test.png",
		"created_at": time.Now().UTC(),
		"updated_at": time.Now().UTC(),
	})
	if err != nil {
		t.Fatal(err)
	}

	service := NewBackupService(database)
	dataDir := t.TempDir()
	if err = service.exportCollections(ctx, dataDir); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(dataDir, databaseName+"_file_meta.json"))
	if err != nil {
		t.Fatal(err)
	}
	var exported []map[string]any
	if err = json.Unmarshal(content, &exported); err != nil {
		t.Fatal(err)
	}
	if exported[0]["file_id"] != "00112233" {
		t.Fatalf("exported file id = %v", exported[0]["file_id"])
	}

	if err = service.DeleteAndInsertCollectionDoc(ctx, "file_meta", content); err != nil {
		t.Fatal(err)
	}
	var restored struct {
		FileID []byte `bson:"file_id"`
	}
	if err = database.Database().Collection("file_meta").FindOne(ctx, bson.M{}).Decode(&restored); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(restored.FileID, want) {
		t.Fatalf("restored file id = %x, want %x", restored.FileID, want)
	}
}
