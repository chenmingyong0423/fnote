package service

import (
	"context"
	"errors"
	"testing"

	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/domain"
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
)

type fakeAssetRepository struct {
	addCalls    int
	deleteCount int64
}

func (f *fakeAssetRepository) FindByIds(context.Context, []string) ([]*domain.Asset, error) {
	return nil, nil
}

func (f *fakeAssetRepository) Add(context.Context, *domain.Asset) (string, error) {
	f.addCalls++
	return "asset-id", nil
}

func (f *fakeAssetRepository) DeleteById(context.Context, string) (int64, error) {
	return f.deleteCount, nil
}

type fakeAssetFolderRepository struct {
	folder    *domain.AssetFolder
	pullCalls int
	putCalls  int
}

func (f *fakeAssetFolderRepository) FindByAssetTypeAndType(context.Context, string, string) ([]*domain.AssetFolder, error) {
	return nil, nil
}

func (f *fakeAssetFolderRepository) Add(context.Context, *domain.AssetFolder) (string, error) {
	return "", nil
}

func (f *fakeAssetFolderRepository) ModifyById(context.Context, *domain.AssetFolder) (int64, error) {
	return 1, nil
}

func (f *fakeAssetFolderRepository) FindById(context.Context, string) (*domain.AssetFolder, error) {
	return f.folder, nil
}

func (f *fakeAssetFolderRepository) DeleteById(context.Context, string) (int64, error) {
	return 1, nil
}

func (f *fakeAssetFolderRepository) AddSubFolder(context.Context, string, *domain.AssetFolder) (int64, string, error) {
	return 1, "subfolder-id", nil
}

func (f *fakeAssetFolderRepository) ModifySubFolderById(context.Context, string, *domain.AssetFolder) (int64, error) {
	return 1, nil
}

func (f *fakeAssetFolderRepository) DeleteSubFolderById(context.Context, string, string) (int64, error) {
	return 1, nil
}

func (f *fakeAssetFolderRepository) ModifyFolderNameById(context.Context, string, string) (int64, error) {
	return 1, nil
}

func (f *fakeAssetFolderRepository) PutAssetId(context.Context, string, string) (int64, error) {
	f.putCalls++
	return 1, nil
}

func (f *fakeAssetFolderRepository) PullAssetId(context.Context, string, string) (int64, error) {
	f.pullCalls++
	return 1, nil
}

func TestDeleteAssetRestoresFolderReferenceWhenAssetDeleteFails(t *testing.T) {
	folderRepo := &fakeAssetFolderRepository{}
	assetRepo := &fakeAssetRepository{deleteCount: 0}
	service := NewAssetService(folderRepo, assetRepo)

	err := service.DeleteAsset(context.Background(), "folder-id", "asset-id")
	if err == nil {
		t.Fatal("expected delete to fail")
	}
	if folderRepo.pullCalls != 1 {
		t.Fatalf("expected one pull, got %d", folderRepo.pullCalls)
	}
	if folderRepo.putCalls != 1 {
		t.Fatalf("expected compensation to restore the reference once, got %d", folderRepo.putCalls)
	}
}

func TestAddAssetRejectsMismatchedFolderType(t *testing.T) {
	folderRepo := &fakeAssetFolderRepository{folder: &domain.AssetFolder{AssetType: "image", Type: "post-editor"}}
	assetRepo := &fakeAssetRepository{}
	service := NewAssetService(folderRepo, assetRepo)

	_, err := service.AddAsset(context.Background(), "folder-id", &domain.Asset{AssetType: "video", Type: "post-editor"})
	if err == nil {
		t.Fatal("expected mismatched asset type to be rejected")
	}
	if assetRepo.addCalls != 0 {
		t.Fatalf("asset should not be persisted, got %d add calls", assetRepo.addCalls)
	}
}

func TestDeleteFolderHonorsCapabilitiesAndContent(t *testing.T) {
	tests := []struct {
		name   string
		folder *domain.AssetFolder
		status int
	}{
		{
			name:   "protected folder",
			folder: &domain.AssetFolder{SupportDelete: false},
			status: 403,
		},
		{
			name:   "non-empty folder",
			folder: &domain.AssetFolder{SupportDelete: true, Assets: []string{"asset-id"}},
			status: 409,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewAssetService(&fakeAssetFolderRepository{folder: tt.folder}, &fakeAssetRepository{})
			_, err := service.DeleteFolderById(context.Background(), "folder-id")
			if err == nil {
				t.Fatal("expected folder deletion to be rejected")
			}

			var responseErr apiwrap.ErrorResponseBody
			if !errors.As(err, &responseErr) {
				t.Fatalf("expected HTTP response error, got %T", err)
			}
			if responseErr.HttpCode != tt.status {
				t.Fatalf("expected status %d, got %d", tt.status, responseErr.HttpCode)
			}
		})
	}
}
