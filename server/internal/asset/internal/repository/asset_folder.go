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
	"github.com/chenmingyong0423/go-mongox"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type IAssetFolderRepository interface {
	FindByAssetTypeAndType(ctx context.Context, assertType string, typ string) ([]*domain.AssetFolder, error)
	Add(ctx context.Context, assetFolder *domain.AssetFolder) (string, error)
	ModifyById(ctx context.Context, assetFolder *domain.AssetFolder) (int64, error)
	FindById(ctx context.Context, id string) (*domain.AssetFolder, error)
	DeleteById(ctx context.Context, id string) (int64, error)
	AddSubFolder(ctx context.Context, id string, assetFolder *domain.AssetFolder) (int64, string, error)
	ModifySubFolderById(ctx context.Context, id string, assetFolder *domain.AssetFolder) (int64, error)
	DeleteSubFolderById(ctx context.Context, id string, subId string) (int64, error)
	ModifyFolderNameById(ctx context.Context, id string, name string) (int64, error)
	PutAssetId(ctx context.Context, folderId string, assetId string) (int64, error)
	PullAssetId(ctx context.Context, folderId string, assetId string) (int64, error)
}

var _ IAssetFolderRepository = (*AssetFolderRepository)(nil)

func NewAssetFolderRepository(dao dao.IAssetFolderDao) *AssetFolderRepository {
	return &AssetFolderRepository{dao: dao}
}

type AssetFolderRepository struct {
	dao dao.IAssetFolderDao
}

func (r *AssetFolderRepository) PullAssetId(ctx context.Context, folderId string, assetId string) (int64, error) {
	objID, err := primitive.ObjectIDFromHex(folderId)
	if err != nil {
		return 0, err
	}
	assetObjID, err := primitive.ObjectIDFromHex(assetId)
	if err != nil {
		return 0, err
	}
	return r.dao.PullAssetId(ctx, objID, assetObjID)
}

func (r *AssetFolderRepository) PutAssetId(ctx context.Context, folderId string, assetId string) (int64, error) {
	objID, err := primitive.ObjectIDFromHex(folderId)
	if err != nil {
		return 0, err
	}
	assetObjID, err := primitive.ObjectIDFromHex(assetId)
	if err != nil {
		return 0, err
	}
	return r.dao.PutAssetId(ctx, objID, assetObjID)
}

func (r *AssetFolderRepository) ModifyFolderNameById(ctx context.Context, id string, name string) (int64, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	return r.dao.ModifyNameById(ctx, objID, name)
}

func (r *AssetFolderRepository) DeleteSubFolderById(ctx context.Context, id string, subId string) (int64, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	subObjID, err := primitive.ObjectIDFromHex(subId)
	if err != nil {
		return 0, err
	}
	return r.dao.DeleteSubFolderById(ctx, objID, subObjID)
}

func (r *AssetFolderRepository) ModifySubFolderById(ctx context.Context, id string, assetFolder *domain.AssetFolder) (int64, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	subObjID, err := primitive.ObjectIDFromHex(assetFolder.Id)
	if err != nil {
		return 0, err
	}
	return r.dao.ModifySubFolderById(ctx, objID, &dao.AssetFolder{
		Model: mongox.Model{
			ID:        subObjID,
			UpdatedAt: time.Now(),
		},
		Name:          assetFolder.Name,
		AssetType:     assetFolder.AssetType,
		Type:          assetFolder.Type,
		SupportDelete: assetFolder.SupportDelete,
		SupportEdit:   assetFolder.SupportEdit,
		SupportAdd:    assetFolder.SupportAdd,
	})
}

func (r *AssetFolderRepository) AddSubFolder(ctx context.Context, id string, assetFolder *domain.AssetFolder) (int64, string, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, "", err
	}
	modifyCnt, err := r.dao.AddSubFolder(ctx, objID, &dao.AssetFolder{
		Model: mongox.Model{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
		},
		Name:          assetFolder.Name,
		AssetType:     assetFolder.AssetType,
		Type:          assetFolder.Type,
		SupportDelete: assetFolder.SupportDelete,
		SupportEdit:   assetFolder.SupportEdit,
		SupportAdd:    assetFolder.SupportAdd,
	})
	if err != nil {
		return 0, "", err
	}
	return modifyCnt, objID.Hex(), nil
}

func (r *AssetFolderRepository) DeleteById(ctx context.Context, id string) (int64, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	return r.dao.DeleteById(ctx, objID)
}

func (r *AssetFolderRepository) FindById(ctx context.Context, id string) (*domain.AssetFolder, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	assetFolder, err := r.dao.FindById(ctx, objID)
	if err != nil {
		return nil, err
	}
	return r.toDomain(assetFolder), nil
}

func (r *AssetFolderRepository) ModifyById(ctx context.Context, assetFolder *domain.AssetFolder) (int64, error) {
	objID, err := primitive.ObjectIDFromHex(assetFolder.Id)
	if err != nil {
		return 0, err
	}
	return r.dao.ModifyById(ctx, &dao.AssetFolder{
		Model: mongox.Model{
			ID:        objID,
			UpdatedAt: time.Now(),
		},
		Name:          assetFolder.Name,
		AssetType:     assetFolder.AssetType,
		Type:          assetFolder.Type,
		SupportDelete: assetFolder.SupportDelete,
		SupportEdit:   assetFolder.SupportEdit,
	})
}

func (r *AssetFolderRepository) Add(ctx context.Context, assetFolder *domain.AssetFolder) (string, error) {
	objectID, err := r.dao.Add(ctx, &dao.AssetFolder{
		Name:          assetFolder.Name,
		AssetType:     assetFolder.AssetType,
		Type:          assetFolder.Type,
		SupportDelete: assetFolder.SupportDelete,
		SupportEdit:   assetFolder.SupportEdit,
	})
	if err != nil {
		return "", err
	}
	return objectID.Hex(), nil
}

func (r *AssetFolderRepository) FindByAssetTypeAndType(ctx context.Context, assertType string, typ string) ([]*domain.AssetFolder, error) {
	assetFolders, err := r.dao.FindByAssetTypeAndType(ctx, assertType, typ)
	if err != nil {
		return nil, err
	}
	return r.toDomains(assetFolders), nil
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
		result.ChildFolders = make([]*domain.AssetFolder, len(assetFolder.ChildFolders))
		for i, v := range assetFolder.ChildFolders {
			result.ChildFolders[i] = r.toDomain(&v)
		}
	}

	return result
}

func (r *AssetFolderRepository) toDomains(assetFolders []*dao.AssetFolder) []*domain.AssetFolder {
	result := make([]*domain.AssetFolder, len(assetFolders))
	for i, v := range assetFolders {
		result[i] = r.toDomain(v)
	}
	return result
}
