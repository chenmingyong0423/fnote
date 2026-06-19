package service

import (
	"reflect"
	"testing"
)

func TestStaticFileIDs(t *testing.T) {
	got := staticFileIDs("![image](/static/00112233445566778899aabbccddeeff.jpeg)")
	want := []string{"00112233445566778899aabbccddeeff"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("staticFileIDs() = %v, want %v", got, want)
	}
}
