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

package repository

import (
	"context"
	"github.com/chenmingyong0423/fnote/server/internal/post_asset/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post_asset/internal/repository/dao"
)

type IAssetFolderRepository interface {
	FindByAssetTypeAndType(ctx context.Context, assertType string, typ string) (*domain.AssetFolder, error)
}

var _ IAssetFolderRepository = (*AssetFolderRepository)(nil)

func NewAssetFolderRepository(dao dao.IAssetFolderDao) *AssetFolderRepository {
	return &AssetFolderRepository{dao: dao}
}

type AssetFolderRepository struct {
	dao dao.IAssetFolderDao
}

func (r *AssetFolderRepository) FindByAssetTypeAndType(ctx context.Context, assertType string, typ string) (*domain.AssetFolder, error) {
	assetFolder, err := r.dao.FindByAssetTypeAndType(ctx, assertType, typ)
	if err != nil {
		return nil, err
	}
	return r.toDomain(assetFolder), nil
}

func (r *AssetFolderRepository) toDomain(assetFolder *dao.AssetFolder) *domain.AssetFolder {
	result := &domain.AssetFolder{
		Id:            assetFolder.ID.Hex(),
		Name:          assetFolder.Name,
		AssetType:     assetFolder.AssetType,
		Type:          assetFolder.Type,
		Assets:        nil,
		ChildFolders:  nil,
		SupportDelete: assetFolder.SupportDelete,
		SupportEdit:   assetFolder.SupportEdit,
	}
	// 转换 AssetFolder
	if assetFolder.Assets != nil {
		result.Assets = make([]string, len(assetFolder.Assets))
		for i, v := range assetFolder.Assets {
			result.Assets[i] = v.Hex()
		}
	}
	// 转换 ChildFolders
	if assetFolder.ChildFolders != nil {
		result.ChildFolders = make([]domain.AssetFolder, len(assetFolder.ChildFolders))
		for i, v := range assetFolder.ChildFolders {
			result.ChildFolders[i] = *r.toDomain(&v)
		}
	}

	return result
}
