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
	"github.com/chenmingyong0423/fnote/server/internal/post_asset/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post_asset/internal/repository"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type IAssetFolderService interface {
	GetFoldersByAssetTypeAndType(ctx context.Context, assertType string, typ string) (*domain.AssetFolder, error)
}

var _ IAssetFolderService = (*AssetFolderService)(nil)

func NewAssetFolderService(repo repository.IAssetFolderRepository) *AssetFolderService {
	return &AssetFolderService{
		repo: repo,
	}
}

type AssetFolderService struct {
	repo repository.IAssetFolderRepository
}

func (s *AssetFolderService) GetFoldersByAssetTypeAndType(ctx context.Context, assertType string, typ string) (*domain.AssetFolder, error) {
	assetFolder, err := s.repo.FindByAssetTypeAndType(ctx, assertType, typ)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	return assetFolder, nil
}
