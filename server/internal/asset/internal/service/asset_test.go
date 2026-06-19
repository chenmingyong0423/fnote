package service

import (
	"context"
	"errors"
	"testing"

	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/file"
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
)

type fakeAssetRepository struct {
	addCalls    int
	deleteCount int64
	asset       *domain.Asset
}

func (f *fakeAssetRepository) FindById(context.Context, string) (*domain.Asset, error) {
	if f.asset != nil {
		return f.asset, nil
	}
	return &domain.Asset{}, nil
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

type fakeFileUsageService struct {
	file.Service
	indexed int
	deleted int
	inUse   bool
}

func (f *fakeFileUsageService) HasOtherUsages(context.Context, []byte, string, string) (bool, error) {
	return f.inUse, nil
}

func (f *fakeFileUsageService) IndexFileMeta(context.Context, []byte, string, string) error {
	f.indexed++
	return nil
}

func (f *fakeFileUsageService) DeleteIndexFileMeta(context.Context, []byte, string, string) error {
	f.deleted++
	return nil
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
	service := NewAssetService(folderRepo, assetRepo, &fakeFileUsageService{})

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
	service := NewAssetService(folderRepo, assetRepo, &fakeFileUsageService{})

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
			service := NewAssetService(&fakeAssetFolderRepository{folder: tt.folder}, &fakeAssetRepository{}, &fakeFileUsageService{})
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

func TestAddAndDeleteImageAssetMaintainFileUsage(t *testing.T) {
	fileUsage := &fakeFileUsageService{}
	asset := &domain.Asset{
		AssetType: "image",
		Type:      "post-editor",
		Metadata:  map[string]any{"file_id": "00112233445566778899aabbccddeeff"},
	}
	assetRepo := &fakeAssetRepository{deleteCount: 1, asset: asset}
	service := NewAssetService(
		&fakeAssetFolderRepository{folder: &domain.AssetFolder{AssetType: "image", Type: "post-editor"}},
		assetRepo,
		fileUsage,
	)

	if _, err := service.AddAsset(context.Background(), "folder-id", asset); err != nil {
		t.Fatalf("AddAsset() error = %v", err)
	}
	if fileUsage.indexed != 1 {
		t.Fatalf("expected file usage to be indexed once, got %d", fileUsage.indexed)
	}
	if err := service.DeleteAsset(context.Background(), "folder-id", "asset-id"); err != nil {
		t.Fatalf("DeleteAsset() error = %v", err)
	}
	if fileUsage.deleted != 1 {
		t.Fatalf("expected file usage to be deleted once, got %d", fileUsage.deleted)
	}
}

func TestDeleteImageAssetRejectsReferencedFile(t *testing.T) {
	fileUsage := &fakeFileUsageService{inUse: true}
	asset := &domain.Asset{
		AssetType: "image",
		Metadata:  map[string]any{"file_id": "00112233445566778899aabbccddeeff"},
	}
	service := NewAssetService(
		&fakeAssetFolderRepository{},
		&fakeAssetRepository{deleteCount: 1, asset: asset},
		fileUsage,
	)

	err := service.DeleteAsset(context.Background(), "folder-id", "asset-id")
	if err == nil {
		t.Fatal("expected referenced asset deletion to fail")
	}
	var responseErr apiwrap.ErrorResponseBody
	if !errors.As(err, &responseErr) || responseErr.HttpCode != 409 {
		t.Fatalf("expected HTTP 409, got %v", err)
	}
}
