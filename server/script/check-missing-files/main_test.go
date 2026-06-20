package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveFilePathPrefersCurrentStaticPath(t *testing.T) {
	record := fileRecord{FileName: "image.png", FilePath: "/old/static/image.png"}
	got := resolveFilePath(record, "/fnote/static")
	want := filepath.Join("/fnote/static", "image.png")
	if got != want {
		t.Fatalf("resolveFilePath() = %q, want %q", got, want)
	}
}

func TestMissingReason(t *testing.T) {
	directory := t.TempDir()
	existing := filepath.Join(directory, "existing.png")
	if err := os.WriteFile(existing, []byte("image"), 0o644); err != nil {
		t.Fatal(err)
	}
	if reason := missingReason(existing); reason != "" {
		t.Fatalf("existing file reason = %q", reason)
	}
	if reason := missingReason(filepath.Join(directory, "missing.png")); reason != "文件不存在" {
		t.Fatalf("missing file reason = %q", reason)
	}
}
