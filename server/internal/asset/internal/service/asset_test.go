package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/file"
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
)

type fakeAssetRepository struct {
	addCalls    int
	modifyCalls int
	deleteCount int64
	asset       *domain.Asset
	assets      []*domain.Asset
	foundIDs    []string
	fileAsset   *domain.Asset
	fileFindErr error
}

func (f *fakeAssetRepository) FindById(context.Context, string) (*domain.Asset, error) {
	if f.asset != nil {
		return f.asset, nil
	}
	return &domain.Asset{}, nil
}

func (f *fakeAssetRepository) FindByFileID(context.Context, string) (*domain.Asset, error) {
	return f.fileAsset, f.fileFindErr
}

func (f *fakeAssetRepository) FindByIds(_ context.Context, ids []string) ([]*domain.Asset, error) {
	f.foundIDs = ids
	return f.assets, nil
}

func (f *fakeAssetRepository) Add(context.Context, *domain.Asset) (string, error) {
	f.addCalls++
	return "asset-id", nil
}

func (f *fakeAssetRepository) ModifyById(context.Context, *domain.Asset) (int64, error) {
	f.modifyCalls++
	return 1, nil
}

func (f *fakeAssetRepository) DeleteById(context.Context, string) (int64, error) {
	return f.deleteCount, nil
}

type fakeAssetFolderRepository struct {
	folder    *domain.AssetFolder
	assetIDs  []string
	total     int64
	pullCalls int
	putCalls  int
}

type fakeFileUsageService struct {
	file.Service
	indexed int
	deleted int
	inUse   bool
}

type fakeTransactionRunner struct {
	calls int
}

func (f *fakeTransactionRunner) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	f.calls++
	return fn(ctx)
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

func (f *fakeAssetFolderRepository) FindAssetIDPageById(context.Context, string, int64, int64) ([]string, int64, error) {
	return f.assetIDs, f.total, nil
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

func TestDeleteAssetReturnsErrorWithoutManualCompensation(t *testing.T) {
	folderRepo := &fakeAssetFolderRepository{}
	assetRepo := &fakeAssetRepository{deleteCount: 0}
	txRunner := &fakeTransactionRunner{}
	service := NewAssetService(folderRepo, assetRepo, &fakeFileUsageService{}, txRunner)

	err := service.DeleteAsset(context.Background(), "folder-id", "asset-id")
	if err == nil {
		t.Fatal("expected delete to fail")
	}
	if folderRepo.pullCalls != 1 {
		t.Fatalf("expected one pull, got %d", folderRepo.pullCalls)
	}
	if folderRepo.putCalls != 0 {
		t.Fatalf("manual compensation should not run, got %d put calls", folderRepo.putCalls)
	}
	if txRunner.calls != 1 {
		t.Fatalf("expected one transaction, got %d", txRunner.calls)
	}
}

func TestAddAssetRejectsMismatchedFolderType(t *testing.T) {
	folderRepo := &fakeAssetFolderRepository{folder: &domain.AssetFolder{AssetType: "image", Type: "legacy"}}
	assetRepo := &fakeAssetRepository{}
	service := NewAssetService(folderRepo, assetRepo, &fakeFileUsageService{}, &fakeTransactionRunner{})

	_, err := service.AddAsset(context.Background(), "folder-id", &domain.Asset{
		Content:   "/static/test.png",
		AssetType: domain.AssetTypeImage,
		Type:      domain.AssetUseTypePostEditor,
		Metadata:  map[string]any{"file_id": "00112233445566778899aabbccddeeff"},
	})
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
			service := NewAssetService(&fakeAssetFolderRepository{folder: tt.folder}, &fakeAssetRepository{}, &fakeFileUsageService{}, &fakeTransactionRunner{})
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
		Content:   "/static/test.png",
		AssetType: "image",
		Type:      "post-editor",
		Metadata:  map[string]any{"file_id": "00112233445566778899aabbccddeeff"},
	}
	assetRepo := &fakeAssetRepository{deleteCount: 1, asset: asset}
	service := NewAssetService(
		&fakeAssetFolderRepository{folder: &domain.AssetFolder{AssetType: "image", Type: "post-editor"}},
		assetRepo,
		fileUsage,
		&fakeTransactionRunner{},
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

func TestAddAssetValidatesClassificationAndImageURL(t *testing.T) {
	tests := []struct {
		name  string
		asset *domain.Asset
	}{
		{
			name:  "unsupported asset type",
			asset: &domain.Asset{AssetType: "video", Type: domain.AssetUseTypePostEditor},
		},
		{
			name:  "unsupported use type",
			asset: &domain.Asset{AssetType: domain.AssetTypeImage, Type: "cover"},
		},
		{
			name: "invalid image URL",
			asset: &domain.Asset{
				Content:   "javascript:alert(1)",
				AssetType: domain.AssetTypeImage,
				Type:      domain.AssetUseTypePostEditor,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assetRepo := &fakeAssetRepository{}
			service := NewAssetService(&fakeAssetFolderRepository{}, assetRepo, &fakeFileUsageService{}, &fakeTransactionRunner{})
			if _, err := service.AddAsset(context.Background(), "folder-id", tt.asset); err == nil {
				t.Fatal("expected invalid asset to be rejected")
			}
			if assetRepo.addCalls != 0 {
				t.Fatal("invalid asset should not be persisted")
			}
		})
	}
}

func TestAddAssetRejectsDuplicateFile(t *testing.T) {
	assetRepo := &fakeAssetRepository{fileAsset: &domain.Asset{Id: "existing-asset"}}
	service := NewAssetService(
		&fakeAssetFolderRepository{folder: &domain.AssetFolder{AssetType: domain.AssetTypeImage, Type: domain.AssetUseTypePostEditor}},
		assetRepo,
		&fakeFileUsageService{},
		&fakeTransactionRunner{},
	)
	asset := &domain.Asset{
		Content:   "/static/test.png",
		AssetType: domain.AssetTypeImage,
		Type:      domain.AssetUseTypePostEditor,
		Metadata:  map[string]any{"file_id": "00112233445566778899AABBCCDDEEFF"},
	}

	_, err := service.AddAsset(context.Background(), "folder-id", asset)
	if err == nil {
		t.Fatal("expected duplicate file to be rejected")
	}
	var responseErr apiwrap.ErrorResponseBody
	if !errors.As(err, &responseErr) || responseErr.HttpCode != http.StatusConflict {
		t.Fatalf("expected HTTP 409, got %v", err)
	}
	if asset.Metadata["file_id"] != "00112233445566778899aabbccddeeff" {
		t.Fatalf("file ID was not normalized: %v", asset.Metadata["file_id"])
	}
	if assetRepo.addCalls != 0 {
		t.Fatal("duplicate asset should not be persisted")
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
		&fakeTransactionRunner{},
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

func TestGetAssetsByFolderIdUsesPagedAssetIDs(t *testing.T) {
	folderRepo := &fakeAssetFolderRepository{assetIDs: []string{"asset-2", "asset-3"}, total: 5}
	assetRepo := &fakeAssetRepository{assets: []*domain.Asset{{Id: "asset-2"}, {Id: "asset-3"}}}
	service := NewAssetService(folderRepo, assetRepo, &fakeFileUsageService{}, &fakeTransactionRunner{})

	assets, total, err := service.GetAssetsByFolderId(context.Background(), "folder-id", 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if total != 5 || len(assets) != 2 {
		t.Fatalf("unexpected page: total=%d assets=%#v", total, assets)
	}
	if len(assetRepo.foundIDs) != 2 || assetRepo.foundIDs[0] != "asset-2" || assetRepo.foundIDs[1] != "asset-3" {
		t.Fatalf("unexpected queried IDs: %#v", assetRepo.foundIDs)
	}
}

func TestAddTextAsset(t *testing.T) {
	assetRepo := &fakeAssetRepository{}
	service := NewAssetService(
		&fakeAssetFolderRepository{folder: &domain.AssetFolder{AssetType: domain.AssetTypeText, Type: domain.AssetUseTypePostEditor}},
		assetRepo,
		&fakeFileUsageService{},
		&fakeTransactionRunner{},
	)

	_, err := service.AddAsset(context.Background(), "folder-id", &domain.Asset{
		Title:     "常用介绍",
		Content:   "这是一段可复用的文字。",
		AssetType: domain.AssetTypeText,
		Type:      domain.AssetUseTypePostEditor,
	})
	if err != nil {
		t.Fatalf("AddAsset() error = %v", err)
	}
	if assetRepo.addCalls != 1 {
		t.Fatalf("expected text asset to be persisted, got %d add calls", assetRepo.addCalls)
	}
}

func TestModifyTextAsset(t *testing.T) {
	assetRepo := &fakeAssetRepository{asset: &domain.Asset{
		Id:        "asset-id",
		AssetType: domain.AssetTypeText,
		Type:      domain.AssetUseTypePostEditor,
	}}
	service := NewAssetService(
		&fakeAssetFolderRepository{folder: &domain.AssetFolder{
			AssetType: domain.AssetTypeText,
			Type:      domain.AssetUseTypePostEditor,
			Assets:    []string{"asset-id"},
		}},
		assetRepo,
		&fakeFileUsageService{},
		&fakeTransactionRunner{},
	)

	_, err := service.ModifyAsset(context.Background(), "folder-id", &domain.Asset{
		Id:        "asset-id",
		Title:     "更新后的标题",
		Content:   "更新后的正文",
		AssetType: domain.AssetTypeText,
		Type:      domain.AssetUseTypePostEditor,
	})
	if err != nil {
		t.Fatalf("ModifyAsset() error = %v", err)
	}
	if assetRepo.modifyCalls != 1 {
		t.Fatalf("expected text asset to be modified, got %d calls", assetRepo.modifyCalls)
	}
}

func TestAddTextAssetRequiresTitle(t *testing.T) {
	assetRepo := &fakeAssetRepository{}
	service := NewAssetService(
		&fakeAssetFolderRepository{folder: &domain.AssetFolder{AssetType: domain.AssetTypeText, Type: domain.AssetUseTypePostEditor}},
		assetRepo,
		&fakeFileUsageService{},
		&fakeTransactionRunner{},
	)

	_, err := service.AddAsset(context.Background(), "folder-id", &domain.Asset{
		Content:   "正文",
		AssetType: domain.AssetTypeText,
		Type:      domain.AssetUseTypePostEditor,
	})
	if err == nil {
		t.Fatal("expected text asset without title to be rejected")
	}
	if assetRepo.addCalls != 0 {
		t.Fatalf("invalid text asset should not be persisted, got %d add calls", assetRepo.addCalls)
	}
}
