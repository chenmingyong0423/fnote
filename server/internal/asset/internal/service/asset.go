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
	"log/slog"
	"net/http"

	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/file"
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	"github.com/pkg/errors"
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
	GetAssetsByIDs(ctx context.Context, id string) ([]*domain.Asset, error)
	AddAsset(ctx context.Context, folderId string, asset *domain.Asset) (string, error)
	DeleteAssetById(ctx context.Context, id string) (int64, error)
	DeleteAsset(ctx context.Context, folderId string, assetId string) error
}

var _ IAssetService = (*AssetService)(nil)

func NewAssetService(repo repository.IAssetFolderRepository, assetRepo repository.IAssetRepository, fileService file.Service) *AssetService {
	return &AssetService{
		assetFolderRepo: repo,
		assetRepo:       assetRepo,
		fileService:     fileService,
	}
}

type AssetService struct {
	assetFolderRepo repository.IAssetFolderRepository
	assetRepo       repository.IAssetRepository
	fileService     file.Service
}

func (s *AssetService) DeleteAsset(ctx context.Context, folderId string, assetId string) error {
	asset, err := s.assetRepo.FindById(ctx, assetId)
	if err != nil {
		return err
	}
	fileID, hasFile, err := assetFileID(asset)
	if err != nil {
		return apiwrap.NewErrorResponseBody(http.StatusBadRequest, err.Error())
	}
	if hasFile {
		inUse, usageErr := s.fileService.HasOtherUsages(ctx, fileID, assetId, "asset")
		if usageErr != nil {
			return usageErr
		}
		if inUse {
			return apiwrap.NewErrorResponseBody(http.StatusConflict, "asset is still referenced by a post or draft")
		}
		if err = s.fileService.DeleteIndexFileMeta(ctx, fileID, assetId, "asset"); err != nil {
			return err
		}
	}
	// todo 后面考虑事务
	cnt, err := s.assetFolderRepo.PullAssetId(ctx, folderId, assetId)
	if err != nil {
		s.recoverFileUsage(ctx, fileID, hasFile, assetId)
		return err
	}
	if cnt == 0 {
		s.recoverFileUsage(ctx, fileID, hasFile, assetId)
		return errors.New("failed to pull assetId, ModifiedCount = 0")
	}
	cnt, err = s.assetRepo.DeleteById(ctx, assetId)
	if err != nil {
		s.recovery4PullAssetId(ctx, folderId, assetId)
		s.recoverFileUsage(ctx, fileID, hasFile, assetId)
		return err
	}
	if cnt == 0 {
		s.recovery4PullAssetId(ctx, folderId, assetId)
		s.recoverFileUsage(ctx, fileID, hasFile, assetId)
		return errors.New("failed to delete asset, DeletedCount = 0")
	}
	return nil
}

func (s *AssetService) AddAsset(ctx context.Context, folderId string, asset *domain.Asset) (string, error) {
	folder, err := s.assetFolderRepo.FindById(ctx, folderId)
	if err != nil {
		return "", err
	}
	if folder.AssetType != asset.AssetType || folder.Type != asset.Type {
		return "", apiwrap.NewErrorResponseBody(http.StatusBadRequest, "asset type does not match folder")
	}
	fileID, hasFile, err := assetFileID(asset)
	if err != nil {
		return "", apiwrap.NewErrorResponseBody(http.StatusBadRequest, err.Error())
	}

	// todo 后面考虑事务
	assetId, err := s.assetRepo.Add(ctx, asset)
	if err != nil {
		return "", err
	}
	cnt, err := s.assetFolderRepo.PutAssetId(ctx, folderId, assetId)
	if err != nil {
		s.recovery4AddAsset(ctx, assetId)
		return "", err
	}
	if cnt == 0 {
		s.recovery4AddAsset(ctx, assetId)
		return "", errors.New("failed to put assetId, DeletedCount = 0")
	}
	if hasFile {
		if err = s.fileService.IndexFileMeta(ctx, fileID, assetId, "asset"); err != nil {
			if _, pullErr := s.assetFolderRepo.PullAssetId(ctx, folderId, assetId); pullErr != nil {
				assetLogger(ctx).ErrorContext(ctx, "failed to rollback asset folder reference", "error", pullErr)
			}
			s.recovery4AddAsset(ctx, assetId)
			return "", err
		}
	}
	return assetId, nil
}

func assetFileID(asset *domain.Asset) ([]byte, bool, error) {
	if asset.AssetType != "image" {
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
	if err != nil {
		return nil, false, errors.New("image asset metadata.file_id is invalid")
	}
	return decoded, true, nil
}

func (s *AssetService) recoverFileUsage(ctx context.Context, fileID []byte, hasFile bool, assetID string) {
	if !hasFile {
		return
	}
	if err := s.fileService.IndexFileMeta(ctx, fileID, assetID, "asset"); err != nil {
		assetLogger(ctx).ErrorContext(ctx, "failed to recover file usage", "error", err)
	}
}

func (s *AssetService) recovery4AddAsset(ctx context.Context, assetId string) {
	deletedCnt, recoverErr := s.DeleteAssetById(ctx, assetId)
	if recoverErr != nil {
		l := assetLogger(ctx)
		l.ErrorContext(ctx, "failed to delete asset", "error", recoverErr)
	}
	if deletedCnt == 0 {
		l := assetLogger(ctx)
		l.ErrorContext(ctx, "failed to delete asset", "deleted_count", 0)
	}
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
	return s.assetFolderRepo.Add(ctx, assetFolder)
}

func (s *AssetService) GetFoldersByAssetTypeAndType(ctx context.Context, assertType string, typ string) ([]*domain.AssetFolder, error) {
	return s.assetFolderRepo.FindByAssetTypeAndType(ctx, assertType, typ)
}

func (s *AssetService) GetAssetsByIDs(ctx context.Context, id string) ([]*domain.Asset, error) {
	assetFolder, err := s.assetFolderRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if assetFolder.Assets != nil {
		return s.assetRepo.FindByIds(ctx, assetFolder.Assets)
	}
	return nil, nil
}

func (s *AssetService) DeleteAssetById(ctx context.Context, id string) (int64, error) {
	return s.assetRepo.DeleteById(ctx, id)
}

func (s *AssetService) recovery4PullAssetId(ctx context.Context, folderId string, assetId string) {
	cnt, err := s.assetFolderRepo.PutAssetId(ctx, folderId, assetId)
	if err != nil {
		l := assetLogger(ctx)
		l.ErrorContext(ctx, "failed to recover asset folder reference", "error", err)
	}
	if cnt == 0 {
		l := assetLogger(ctx)
		l.ErrorContext(ctx, "failed to recover asset folder reference", "modified_count", 0)
	}
}

func findChildFolder(folders []*domain.AssetFolder, id string) *domain.AssetFolder {
	for _, folder := range folders {
		if folder.Id == id {
			return folder
		}
	}
	return nil
}

func assetLogger(ctx context.Context) *slog.Logger {
	requestID := ""
	if getter, ok := ctx.(interface{ GetString(string) string }); ok {
		requestID = getter.GetString("X-Request-ID")
	}
	return slog.Default().With("X-Request-ID", requestID)
}
