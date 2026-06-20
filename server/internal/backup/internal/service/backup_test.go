package service

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestFileIDBackupRoundTrip(t *testing.T) {
	want := []byte{0x00, 0x11, 0x22, 0x33}
	document := bson.M{"file_id": bson.Binary{Subtype: 0, Data: want}}
	if err := encodeFileID(document); err != nil {
		t.Fatal(err)
	}
	if document["file_id"] != "00112233" {
		t.Fatalf("exported file id = %v", document["file_id"])
	}
	got, err := decodeFileID(document["file_id"])
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("decoded file id = %x, want %x", got, want)
	}
}

func TestDecodeFileIDRejectsLegacyDocument(t *testing.T) {
	_, err := decodeFileID(map[string]any{"Subtype": 0, "Data": "ABEiMw=="})
	if err == nil {
		t.Fatal("expected legacy document format to be rejected")
	}
}
