package service

import (
	"context"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chenmingyong0423/fnote/server/internal/file/internal/domain"
	"github.com/spf13/viper"
)

type failingFileRepository struct{}

func (failingFileRepository) Save(context.Context, *domain.File) error {
	return errors.New("save failed")
}
func (failingFileRepository) FindByFileId(context.Context, []byte) (*domain.File, error) {
	return nil, nil
}
func (failingFileRepository) PushIntoUsedIn(context.Context, []byte, string, string) error {
	return nil
}
func (failingFileRepository) PullUsedIn(context.Context, []byte, string, string) error { return nil }
func (failingFileRepository) FindByFileName(context.Context, string) (*domain.File, error) {
	return nil, nil
}
func (failingFileRepository) FindPageFilesByFileType(context.Context, domain.PageDTO) ([]*domain.File, int64, error) {
	return nil, 0, nil
}
func (failingFileRepository) DeleteUnusedByFileId(context.Context, []byte) (int64, error) {
	return 0, nil
}

type deletableFileRepository struct {
	failingFileRepository
	file        *domain.File
	deleteCount int64
	deleteCalls int
}

type batchDeleteFileRepository struct {
	failingFileRepository
	files       map[string]*domain.File
	deleteCalls int
}

func (r *batchDeleteFileRepository) FindByFileId(_ context.Context, fileID []byte) (*domain.File, error) {
	return r.files[hex.EncodeToString(fileID)], nil
}

func (r *batchDeleteFileRepository) DeleteUnusedByFileId(context.Context, []byte) (int64, error) {
	r.deleteCalls++
	return 1, nil
}

func (r *deletableFileRepository) FindByFileId(context.Context, []byte) (*domain.File, error) {
	return r.file, nil
}

func (r *deletableFileRepository) DeleteUnusedByFileId(context.Context, []byte) (int64, error) {
	r.deleteCalls++
	return r.deleteCount, nil
}

func TestUploadRemovesPhysicalFileWhenMetadataSaveFails(t *testing.T) {
	staticPath := t.TempDir()
	previousStaticPath := viper.GetString("system.static_path")
	viper.Set("system.static_path", staticPath)
	t.Cleanup(func() { viper.Set("system.static_path", previousStaticPath) })

	service := &FileService{repo: failingFileRepository{}}
	_, err := service.Upload(context.Background(), domain.FileDTO{
		FileName:       "image.png",
		FileExt:        ".png",
		CustomFileName: "rollback-test",
		Content:        []byte("image"),
	})
	if err == nil {
		t.Fatal("expected metadata save to fail")
	}
	if _, statErr := os.Stat(filepath.Join(staticPath, "rollback-test.png")); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("physical file was not rolled back: %v", statErr)
	}
}

func TestDeleteFileRemovesUnusedMetadataAndPhysicalFile(t *testing.T) {
	staticPath := t.TempDir()
	previousStaticPath := viper.GetString("system.static_path")
	viper.Set("system.static_path", staticPath)
	t.Cleanup(func() { viper.Set("system.static_path", previousStaticPath) })

	filePath := filepath.Join(staticPath, "unused.png")
	if err := os.WriteFile(filePath, []byte("image"), 0o644); err != nil {
		t.Fatal(err)
	}
	repo := &deletableFileRepository{
		file:        &domain.File{FileName: "unused.png", FilePath: "/old/static/unused.png"},
		deleteCount: 1,
	}
	service := &FileService{repo: repo}
	if err := service.DeleteFile(context.Background(), "00112233445566778899aabbccddeeff"); err != nil {
		t.Fatal(err)
	}
	if repo.deleteCalls != 1 {
		t.Fatalf("delete calls = %d, want 1", repo.deleteCalls)
	}
	if _, err := os.Stat(filePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("physical file still exists: %v", err)
	}
}

func TestDeleteFileRejectsReferencedFile(t *testing.T) {
	staticPath := t.TempDir()
	previousStaticPath := viper.GetString("system.static_path")
	viper.Set("system.static_path", staticPath)
	t.Cleanup(func() { viper.Set("system.static_path", previousStaticPath) })

	filePath := filepath.Join(staticPath, "used.png")
	if err := os.WriteFile(filePath, []byte("image"), 0o644); err != nil {
		t.Fatal(err)
	}
	repo := &deletableFileRepository{file: &domain.File{
		FileName: "used.png",
		FilePath: filePath,
		UsedIn:   []domain.FileUsage{{EntityId: "post-1", EntityType: "post"}},
	}}
	service := &FileService{repo: repo}
	if err := service.DeleteFile(context.Background(), "00112233445566778899aabbccddeeff"); err == nil {
		t.Fatal("expected referenced file deletion to fail")
	}
	if repo.deleteCalls != 0 {
		t.Fatalf("delete calls = %d, want 0", repo.deleteCalls)
	}
	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("referenced physical file was removed: %v", err)
	}
}

func TestDeleteFilesValidatesAllFilesBeforeDeletion(t *testing.T) {
	staticPath := t.TempDir()
	previousStaticPath := viper.GetString("system.static_path")
	viper.Set("system.static_path", staticPath)
	t.Cleanup(func() { viper.Set("system.static_path", previousStaticPath) })

	unusedID := "00112233445566778899aabbccddeeff"
	usedID := "ffeeddccbbaa99887766554433221100"
	repo := &batchDeleteFileRepository{files: map[string]*domain.File{
		unusedID: {FileName: "unused.png", FilePath: filepath.Join(staticPath, "unused.png")},
		usedID: {
			FileName: "used.png",
			FilePath: filepath.Join(staticPath, "used.png"),
			UsedIn:   []domain.FileUsage{{EntityId: "post-1", EntityType: "post"}},
		},
	}}
	service := &FileService{repo: repo}
	if err := service.DeleteFiles(context.Background(), []string{unusedID, usedID}); err == nil {
		t.Fatal("expected batch deletion to reject referenced file")
	}
	if repo.deleteCalls != 0 {
		t.Fatalf("delete calls = %d, want 0", repo.deleteCalls)
	}
}

func TestResolveStaticFilePathRejectsTraversal(t *testing.T) {
	if _, ok := resolveStaticFilePath("../outside.png"); ok {
		t.Fatal("expected traversal file name to be rejected")
	}
	if _, ok := resolveStaticFilePath(`folder\outside.png`); ok {
		t.Fatal("expected backslash file name to be rejected")
	}
}

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

func TestGenerateSitemapUsesCurrentPublicRoutes(t *testing.T) {
	staticPath := t.TempDir()
	previousStaticPath := viper.GetString("system.static_path")
	viper.Set("system.static_path", staticPath)
	t.Cleanup(func() { viper.Set("system.static_path", previousStaticPath) })
	t.Setenv("WEBSITE_BASE_HOST", "https://chenmingyong.cn/")
	t.Setenv("UPLOADER_HOST", "https://chenmingyong.cn/")

	service := &FileService{}
	posts := []byte(`[{"_id":"about-me","updated_at":1700000000},{"_id":"go-example","updated_at":1700000000,"tags":[{"name":"Go"}],"category":[{"name":"Backend"}]}]`)
	tags := []byte(`[{"Name":"Go","Route":"go","PostCount":0},{"Name":"Vue","Route":"vue","PostCount":10}]`)
	categories := []byte(`[{"Name":"Backend","Route":"backend"},{"Name":"Empty","Route":"empty"}]`)
	if err := service.GenerateSitemap(context.Background(), posts, categories, tags); err != nil {
		t.Fatalf("GenerateSitemap() error = %v", err)
	}

	content, exists, err := service.GetSitemap(context.Background())
	if err != nil || !exists {
		t.Fatalf("GetSitemap() = (_, %t, %v), want existing sitemap", exists, err)
	}
	if !strings.Contains(content, "<loc>https://chenmingyong.cn/about-me</loc>") {
		t.Fatal("sitemap does not contain the current /about-me route")
	}
	if strings.Contains(content, "<loc>https://chenmingyong.cn/about</loc>") {
		t.Fatal("sitemap contains the removed /about route")
	}
	if strings.Contains(content, "/search") {
		t.Fatal("sitemap contains a noindex search route")
	}
	if strings.Contains(content, "chenmingyong.cn//") {
		t.Fatal("sitemap contains a double slash after the host")
	}
	for _, route := range []string{"/navigation", "/posts/go-example", "/tags/go", "/categories/backend"} {
		if !strings.Contains(content, "<loc>https://chenmingyong.cn"+route+"</loc>") {
			t.Fatalf("sitemap is missing %s", route)
		}
	}
	for _, route := range []string{"/posts/about-me", "/tags/vue", "/categories/empty"} {
		if strings.Contains(content, "<loc>https://chenmingyong.cn"+route+"</loc>") {
			t.Fatalf("sitemap contains noncanonical or empty route %s", route)
		}
	}
}
