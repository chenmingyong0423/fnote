package service

import (
	"reflect"
	"testing"
)

func TestStaticFileIDs(t *testing.T) {
	got := staticFileIDs(
		"/static/00112233445566778899aabbccddeeff.png",
		"![image](/static/ffeeddccbbaa99887766554433221100.webp) /static/00112233445566778899aabbccddeeff.png",
	)
	want := []string{"00112233445566778899aabbccddeeff", "ffeeddccbbaa99887766554433221100"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("staticFileIDs() = %v, want %v", got, want)
	}
}

func TestDiffFileIDs(t *testing.T) {
	got := diffFileIDs([]string{"a", "b"}, []string{"b", "c"})
	if !reflect.DeepEqual(got, []string{"a"}) {
		t.Fatalf("diffFileIDs() = %v, want [a]", got)
	}
}
