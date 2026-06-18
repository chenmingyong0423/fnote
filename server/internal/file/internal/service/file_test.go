package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestRobotsTxtLifecycle(t *testing.T) {
	staticPath := t.TempDir()
	previousStaticPath := viper.GetString("system.static_path")
	viper.Set("system.static_path", staticPath)
	t.Cleanup(func() {
		viper.Set("system.static_path", previousStaticPath)
	})

	service := &FileService{}
	content, exists, err := service.GetRobotsTxt(context.Background())
	if err != nil {
		t.Fatalf("GetRobotsTxt() error = %v", err)
	}
	if exists || content != "" {
		t.Fatalf("GetRobotsTxt() = (%q, %t), want empty missing file", content, exists)
	}

	want := "User-agent: *\nAllow: /\n"
	if err = service.SaveRobotsTxt(context.Background(), want); err != nil {
		t.Fatalf("SaveRobotsTxt() error = %v", err)
	}

	content, exists, err = service.GetRobotsTxt(context.Background())
	if err != nil {
		t.Fatalf("GetRobotsTxt() after save error = %v", err)
	}
	if !exists || content != want {
		t.Fatalf("GetRobotsTxt() = (%q, %t), want (%q, true)", content, exists, want)
	}
}

func TestGetSitemap(t *testing.T) {
	staticPath := t.TempDir()
	previousStaticPath := viper.GetString("system.static_path")
	viper.Set("system.static_path", staticPath)
	t.Cleanup(func() {
		viper.Set("system.static_path", previousStaticPath)
	})

	service := &FileService{}
	content, exists, err := service.GetSitemap(context.Background())
	if err != nil || exists || content != "" {
		t.Fatalf("GetSitemap() = (%q, %t, %v), want empty missing file", content, exists, err)
	}

	want := "<?xml version=\"1.0\"?><urlset></urlset>"
	if err = os.WriteFile(filepath.Join(staticPath, "sitemap.xml"), []byte(want), 0o644); err != nil {
		t.Fatalf("write sitemap fixture: %v", err)
	}

	content, exists, err = service.GetSitemap(context.Background())
	if err != nil || !exists || content != want {
		t.Fatalf("GetSitemap() = (%q, %t, %v), want (%q, true, nil)", content, exists, err, want)
	}
}
