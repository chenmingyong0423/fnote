// Copyright 2024 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"
	"encoding/hex"
	"net/http"
	"net/url"
	"strings"

	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/file"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/mongotx"
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IAssetService interface {
	GetFoldersByAssetTypeAndType(ctx context.Context, assertType string, typ string) ([]*domain.AssetFolder, error)
	AddFolder(ctx context.Context, assetFolder *domain.AssetFolder) (string, error)
	ModifyFolderById(ctx context.Context, assetFolder *domain.AssetFolder) (int64, error)
	DeleteFolderById(ctx context.Context, id string) (int64, error)
	AddSubFolder(ctx context.Context, id string, assetFolder *domain.AssetFolder) (int64, string, error)
	ModifySubFolderById(ctx context.Context, id string, assetFolder *domain.AssetFolder) (int64, error)
	DeleteSubFolderById(ctx context.Context, id string, subId string) (int64, error)
	ModifyFolderNameById(ctx context.Context, id string, name string) (int64, error)
	GetAssetFolderById(ctx context.Context, id string) (*domain.AssetFolder, error)
	GetAssetsByFolderId(ctx context.Context, id string, pageNo, pageSize int64) ([]*domain.Asset, int64, error)
	AddAsset(ctx context.Context, folderId string, asset *domain.Asset) (string, error)
	DeleteAsset(ctx context.Context, folderId string, assetId string) error
}

var _ IAssetService = (*AssetService)(nil)

func NewAssetService(repo repository.IAssetFolderRepository, assetRepo repository.IAssetRepository, fileService file.Service, txManager mongotx.Runner) *AssetService {
	return &AssetService{
		assetFolderRepo: repo,
		assetRepo:       assetRepo,
		fileService:     fileService,
		txManager:       txManager,
	}
}

type AssetService struct {
	assetFolderRepo repository.IAssetFolderRepository
	assetRepo       repository.IAssetRepository
	fileService     file.Service
	txManager       mongotx.Runner
}

func (s *AssetService) DeleteAsset(ctx context.Context, folderId string, assetId string) error {
	return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		asset, err := s.assetRepo.FindById(txCtx, assetId)
		if err != nil {
			return err
		}
		fileID, hasFile, err := assetFileID(asset)
		if err != nil {
			return apiwrap.NewErrorResponseBody(http.StatusBadRequest, err.Error())
		}
		if hasFile {
			inUse, usageErr := s.fileService.HasOtherUsages(txCtx, fileID, assetId, "asset")
			if usageErr != nil {
				return usageErr
			}
			if inUse {
				return apiwrap.NewErrorResponseBody(http.StatusConflict, "asset is still referenced by a post or draft")
			}
			if err = s.fileService.DeleteIndexFileMeta(txCtx, fileID, assetId, "asset"); err != nil {
				return err
			}
		}

		cnt, err := s.assetFolderRepo.PullAssetId(txCtx, folderId, assetId)
		if err != nil {
			return err
		}
		if cnt == 0 {
			return errors.New("failed to pull assetId, ModifiedCount = 0")
		}
		cnt, err = s.assetRepo.DeleteById(txCtx, assetId)
		if err != nil {
			return err
		}
		if cnt == 0 {
			return errors.New("failed to delete asset, DeletedCount = 0")
		}
		return nil
	})
}

func (s *AssetService) AddAsset(ctx context.Context, folderId string, asset *domain.Asset) (string, error) {
	if err := validateAssetClassification(asset.AssetType, asset.Type); err != nil {
		return "", err
	}
	if asset.AssetType == domain.AssetTypeImage {
		if err := validateImageURL(asset.Content); err != nil {
			return "", err
		}
	}
	fileID, hasFile, err := assetFileID(asset)
	if err != nil {
		return "", apiwrap.NewErrorResponseBody(http.StatusBadRequest, err.Error())
	}

	var assetId string
	err = s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		folder, txErr := s.assetFolderRepo.FindById(txCtx, folderId)
		if txErr != nil {
			return txErr
		}
		if folder.AssetType != asset.AssetType || folder.Type != asset.Type {
			return apiwrap.NewErrorResponseBody(http.StatusBadRequest, "asset type does not match folder")
		}
		if hasFile {
			existing, findErr := s.assetRepo.FindByFileID(txCtx, asset.Metadata["file_id"].(string))
			if findErr != nil && !errors.Is(findErr, mongo.ErrNoDocuments) {
				return findErr
			}
			if existing != nil {
				return apiwrap.NewErrorResponseBody(http.StatusConflict, "file already exists as an asset")
			}
		}

		assetId, txErr = s.assetRepo.Add(txCtx, asset)
		if txErr != nil {
			return txErr
		}
		cnt, txErr := s.assetFolderRepo.PutAssetId(txCtx, folderId, assetId)
		if txErr != nil {
			return txErr
		}
		if cnt == 0 {
			return errors.New("failed to put assetId, ModifiedCount = 0")
		}
		if hasFile {
			if txErr = s.fileService.IndexFileMeta(txCtx, fileID, assetId, "asset"); txErr != nil {
				return txErr
			}
		}
		return nil
	})
	return assetId, err
}

func assetFileID(asset *domain.Asset) ([]byte, bool, error) {
	if asset.AssetType != domain.AssetTypeImage {
		return nil, false, nil
	}
	value, ok := asset.Metadata["file_id"]
	if !ok {
		return nil, false, errors.New("image asset metadata.file_id is required")
	}
	fileID, ok := value.(string)
	if !ok || fileID == "" {
		return nil, false, errors.New("image asset metadata.file_id must be a non-empty string")
	}
	decoded, err := hex.DecodeString(fileID)
	if err != nil || len(decoded) != 16 {
		return nil, false, errors.New("image asset metadata.file_id is invalid")
	}
	asset.Metadata["file_id"] = hex.EncodeToString(decoded)
	return decoded, true, nil
}

func validateAssetClassification(assetType, useType string) error {
	if assetType != domain.AssetTypeImage {
		return apiwrap.NewErrorResponseBody(http.StatusBadRequest, "unsupported asset_type")
	}
	if useType != domain.AssetUseTypePostEditor {
		return apiwrap.NewErrorResponseBody(http.StatusBadRequest, "unsupported type")
	}
	return nil
}

func validateImageURL(content string) error {
	parsed, err := url.ParseRequestURI(content)
	if err != nil {
		return apiwrap.NewErrorResponseBody(http.StatusBadRequest, "invalid image URL")
	}
	isHTTPURL := parsed.Host != "" && (parsed.Scheme == "http" || parsed.Scheme == "https")
	isStaticPath := parsed.Host == "" && parsed.Scheme == "" && strings.HasPrefix(parsed.Path, "/static/")
	if !isHTTPURL && !isStaticPath {
		return apiwrap.NewErrorResponseBody(http.StatusBadRequest, "invalid image URL")
	}
	return nil
}

func (s *AssetService) GetAssetFolderById(ctx context.Context, id string) (*domain.AssetFolder, error) {
	return s.assetFolderRepo.FindById(ctx, id)
}

func (s *AssetService) ModifyFolderNameById(ctx context.Context, id string, name string) (int64, error) {
	folder, err := s.assetFolderRepo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	if !folder.SupportEdit {
		return 0, apiwrap.NewErrorResponseBody(http.StatusForbidden, "asset folder does not support editing")
	}
	return s.assetFolderRepo.ModifyFolderNameById(ctx, id, name)
}

func (s *AssetService) AddSubFolder(ctx context.Context, id string, assetFolder *domain.AssetFolder) (int64, string, error) {
	if err := validateAssetClassification(assetFolder.AssetType, assetFolder.Type); err != nil {
		return 0, "", err
	}
	folder, err := s.assetFolderRepo.FindById(ctx, id)
	if err != nil {
		return 0, "", err
	}
	if !folder.SupportAdd {
		return 0, "", apiwrap.NewErrorResponseBody(http.StatusForbidden, "asset folder does not support adding subfolders")
	}
	if folder.AssetType != assetFolder.AssetType || folder.Type != assetFolder.Type {
		return 0, "", apiwrap.NewErrorResponseBody(http.StatusBadRequest, "subfolder type does not match parent folder")
	}
	return s.assetFolderRepo.AddSubFolder(ctx, id, assetFolder)
}

func (s *AssetService) ModifySubFolderById(ctx context.Context, id string, assetFolder *domain.AssetFolder) (int64, error) {
	if err := validateAssetClassification(assetFolder.AssetType, assetFolder.Type); err != nil {
		return 0, err
	}
	folder, err := s.assetFolderRepo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	existing := findChildFolder(folder.ChildFolders, assetFolder.Id)
	if existing == nil {
		return 0, apiwrap.NewErrorResponseBody(http.StatusNotFound, "asset subfolder not found")
	}
	if !existing.SupportEdit {
		return 0, apiwrap.NewErrorResponseBody(http.StatusForbidden, "asset subfolder does not support editing")
	}
	assetFolder.Assets = existing.Assets
	assetFolder.ChildFolders = existing.ChildFolders
	return s.assetFolderRepo.ModifySubFolderById(ctx, id, assetFolder)
}

func (s *AssetService) DeleteSubFolderById(ctx context.Context, id string, subId string) (int64, error) {
	folder, err := s.assetFolderRepo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	subfolder := findChildFolder(folder.ChildFolders, subId)
	if subfolder == nil {
		return 0, apiwrap.NewErrorResponseBody(http.StatusNotFound, "asset subfolder not found")
	}
	if !subfolder.SupportDelete {
		return 0, apiwrap.NewErrorResponseBody(http.StatusForbidden, "asset subfolder does not support deletion")
	}
	if len(subfolder.Assets) > 0 || len(subfolder.ChildFolders) > 0 {
		return 0, apiwrap.NewErrorResponseBody(http.StatusConflict, "asset subfolder is not empty")
	}
	return s.assetFolderRepo.DeleteSubFolderById(ctx, id, subId)
}

func (s *AssetService) DeleteFolderById(ctx context.Context, id string) (int64, error) {
	folder, err := s.assetFolderRepo.FindById(ctx, id)
	if err != nil {
		return 0, err
	}
	if !folder.SupportDelete {
		return 0, apiwrap.NewErrorResponseBody(http.StatusForbidden, "asset folder does not support deletion")
	}
	if len(folder.Assets) > 0 || len(folder.ChildFolders) > 0 {
		return 0, apiwrap.NewErrorResponseBody(http.StatusConflict, "asset folder is not empty")
	}
	return s.assetFolderRepo.DeleteById(ctx, id)
}

func (s *AssetService) ModifyFolderById(ctx context.Context, assetFolder *domain.AssetFolder) (int64, error) {
	if err := validateAssetClassification(assetFolder.AssetType, assetFolder.Type); err != nil {
		return 0, err
	}
	existing, err := s.assetFolderRepo.FindById(ctx, assetFolder.Id)
	if err != nil {
		return 0, err
	}
	if !existing.SupportEdit {
		return 0, apiwrap.NewErrorResponseBody(http.StatusForbidden, "asset folder does not support editing")
	}
	return s.assetFolderRepo.ModifyById(ctx, assetFolder)
}

func (s *AssetService) AddFolder(ctx context.Context, assetFolder *domain.AssetFolder) (string, error) {
	if err := validateAssetClassification(assetFolder.AssetType, assetFolder.Type); err != nil {
		return "", err
	}
	return s.assetFolderRepo.Add(ctx, assetFolder)
}

func (s *AssetService) GetFoldersByAssetTypeAndType(ctx context.Context, assertType string, typ string) ([]*domain.AssetFolder, error) {
	if err := validateAssetClassification(assertType, typ); err != nil {
		return nil, err
	}
	return s.assetFolderRepo.FindByAssetTypeAndType(ctx, assertType, typ)
}

func (s *AssetService) GetAssetsByFolderId(ctx context.Context, id string, pageNo, pageSize int64) ([]*domain.Asset, int64, error) {
	assetIDs, total, err := s.assetFolderRepo.FindAssetIDPageById(ctx, id, pageNo, pageSize)
	if err != nil {
		return nil, 0, err
	}
	if len(assetIDs) == 0 {
		return nil, total, nil
	}
	assets, err := s.assetRepo.FindByIds(ctx, assetIDs)
	return assets, total, err
}

func findChildFolder(folders []*domain.AssetFolder, id string) *domain.AssetFolder {
	for _, folder := range folders {
		if folder.Id == id {
			return folder
		}
	}
	return nil
}
