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
	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/repository/dao"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAssetRepository interface {
	FindByIds(ctx context.Context, ids []string) ([]*domain.Asset, error)
	Add(ctx context.Context, asset *domain.Asset) (string, error)
	DeleteById(ctx context.Context, id string) (int64, error)
}

var _ IAssetRepository = (*AssetRepository)(nil)

func NewAssetRepository(dao dao.IAssetDao) *AssetRepository {
	return &AssetRepository{dao: dao}
}

type AssetRepository struct {
	dao dao.IAssetDao
}

func (r *AssetRepository) DeleteById(ctx context.Context, id string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	return r.dao.DeleteById(ctx, objectID)
}

func (r *AssetRepository) Add(ctx context.Context, asset *domain.Asset) (string, error) {
	objectID, err := r.dao.Add(ctx, r.toDao(asset))
	if err != nil {
		return "", err
	}
	return objectID.Hex(), nil
}

func (r *AssetRepository) FindByIds(ctx context.Context, ids []string) ([]*domain.Asset, error) {
	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objIDs = append(objIDs, objID)
	}
	assets, err := r.dao.FindByIds(ctx, objIDs)
	if err != nil {
		return nil, err
	}
	return r.toDomains(assets), nil
}

func (r *AssetRepository) toDomains(assets []*dao.Asset) []*domain.Asset {
	domains := make([]*domain.Asset, 0, len(assets))
	for _, asset := range assets {
		domains = append(domains, r.toDomain(asset))
	}
	return domains
}

func (r *AssetRepository) toDomain(asset *dao.Asset) *domain.Asset {
	return &domain.Asset{
		Id:          asset.ID.Hex(),
		Title:       asset.Title,
		Content:     asset.Content,
		Description: asset.Description,
		AssetType:   asset.AssetType,
		Type:        asset.Type,
		Metadata:    asset.Metadata,
	}
}

func (r *AssetRepository) toDao(asset *domain.Asset) *dao.Asset {
	return &dao.Asset{
		Title:       asset.Title,
		Content:     asset.Content,
		Description: asset.Description,
		AssetType:   asset.AssetType,
		Type:        asset.Type,
		Metadata:    asset.Metadata,
	}
}
