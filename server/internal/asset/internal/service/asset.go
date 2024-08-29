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
	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log/slog"
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

func NewAssetService(repo repository.IAssetFolderRepository, assetRepo repository.IAssetRepository) *AssetService {
	return &AssetService{
		assetFolderRepo: repo,
		assetRepo:       assetRepo,
	}
}

type AssetService struct {
	assetFolderRepo repository.IAssetFolderRepository
	assetRepo       repository.IAssetRepository
}

func (s *AssetService) DeleteAsset(ctx context.Context, folderId string, assetId string) error {
	// todo 后面考虑事务
	cnt, err := s.assetFolderRepo.PullAssetId(ctx, folderId, assetId)
	if err != nil {
		return err
	}
	if cnt == 0 {
		return errors.New("failed to pull assetId, ModifiedCount = 0")
	}
	cnt, err = s.assetRepo.DeleteById(ctx, assetId)
	if err != nil {
		s.recovery4PullAssetId(ctx, folderId, assetId)
		return err
	}
	if cnt == 0 {
		s.recovery4PullAssetId(ctx, folderId, assetId)
		return errors.New("failed to delete asset, DeletedCount = 0")
	}
	return nil
}

func (s *AssetService) AddAsset(ctx context.Context, folderId string, asset *domain.Asset) (string, error) {
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
	return assetId, nil
}

func (s *AssetService) recovery4AddAsset(ctx context.Context, assetId string) {
	deletedCnt, recoverErr := s.DeleteAssetById(ctx, assetId)
	if recoverErr != nil {
		l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
		l.ErrorContext(ctx, "failed to delete asset", recoverErr.Error())
	}
	if deletedCnt == 0 {
		l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
		l.ErrorContext(ctx, "failed to delete asset", "DeletedCount = 0")
	}
}

func (s *AssetService) GetAssetFolderById(ctx context.Context, id string) (*domain.AssetFolder, error) {
	return s.assetFolderRepo.FindById(ctx, id)
}

func (s *AssetService) ModifyFolderNameById(ctx context.Context, id string, name string) (int64, error) {
	return s.assetFolderRepo.ModifyFolderNameById(ctx, id, name)
}

func (s *AssetService) AddSubFolder(ctx context.Context, id string, assetFolder *domain.AssetFolder) (int64, string, error) {
	return s.assetFolderRepo.AddSubFolder(ctx, id, assetFolder)
}

func (s *AssetService) ModifySubFolderById(ctx context.Context, id string, assetFolder *domain.AssetFolder) (int64, error) {
	return s.assetFolderRepo.ModifySubFolderById(ctx, id, assetFolder)
}

func (s *AssetService) DeleteSubFolderById(ctx context.Context, id string, subId string) (int64, error) {
	return s.assetFolderRepo.DeleteSubFolderById(ctx, id, subId)
}

func (s *AssetService) DeleteFolderById(ctx context.Context, id string) (int64, error) {
	return s.assetFolderRepo.DeleteById(ctx, id)
}

func (s *AssetService) ModifyFolderById(ctx context.Context, assetFolder *domain.AssetFolder) (int64, error) {
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
	cnt, err := s.assetFolderRepo.PullAssetId(ctx, folderId, assetId)
	if err != nil {
		l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
		l.ErrorContext(ctx, "failed to recovery 4 PullAssetId", err.Error())
	}
	if cnt == 0 {
		l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
		l.ErrorContext(ctx, "failed to recovery 4 PullAssetId", "ModifiedCount = 0")
	}
}
