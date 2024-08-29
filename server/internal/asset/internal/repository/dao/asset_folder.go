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

package dao

import (
	"context"
	"github.com/chenmingyong0423/go-mongox"
	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"github.com/chenmingyong0423/go-mongox/builder/update"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type AssetFolder struct {
	mongox.Model `bson:",inline"`
	// 文件夹名称
	Name string `bson:"name"`
	// 文件夹归属的素材类型，image ···
	AssetType string `bson:"asset_type"`
	// 文件夹类型，post-editor ···
	Type          string               `bson:"type"`
	Assets        []primitive.ObjectID `bson:"assets,omitempty"`
	ChildFolders  []AssetFolder        `bson:"child_folders,omitempty"`
	SupportDelete bool                 `bson:"support_delete"`
	SupportEdit   bool                 `bson:"support_edit"`
	SupportAdd    bool                 `bson:"support_add"`
}

type IAssetFolderDao interface {
	FindByAssetTypeAndType(ctx context.Context, assertType string, typ string) ([]*AssetFolder, error)
	Add(ctx context.Context, assetFolder *AssetFolder) (primitive.ObjectID, error)
	ModifyById(ctx context.Context, assetFolder *AssetFolder) (int64, error)
	FindById(ctx context.Context, objectID primitive.ObjectID) (*AssetFolder, error)
	DeleteById(ctx context.Context, objectID primitive.ObjectID) (int64, error)
	AddSubFolder(ctx context.Context, id primitive.ObjectID, assetFolder *AssetFolder) (int64, error)
	ModifySubFolderById(ctx context.Context, objectID primitive.ObjectID, assetFolder *AssetFolder) (int64, error)
	DeleteSubFolderById(ctx context.Context, objectID primitive.ObjectID, subObjID primitive.ObjectID) (int64, error)
	ModifyNameById(ctx context.Context, objectID primitive.ObjectID, name string) (int64, error)
	PutAssetId(ctx context.Context, folderObjID primitive.ObjectID, assetObjID primitive.ObjectID) (int64, error)
	PullAssetId(ctx context.Context, folderObjID primitive.ObjectID, assetObjID primitive.ObjectID) (int64, error)
}

var _ IAssetFolderDao = (*AssetFolderDao)(nil)

func NewAssetFolderDao(db *mongo.Database) *AssetFolderDao {
	return &AssetFolderDao{coll: mongox.NewCollection[AssetFolder](db.Collection("asset_folders"))}
}

type AssetFolderDao struct {
	coll *mongox.Collection[AssetFolder]
}

func (d *AssetFolderDao) PullAssetId(ctx context.Context, folderObjID primitive.ObjectID, assetObjID primitive.ObjectID) (int64, error) {
	updateResult, err := d.coll.Updater().
		Filter(query.Id(folderObjID)).
		Updates(update.NewBuilder().Pull("assets", assetObjID).Set("updated_at", time.Now()).Build()).
		UpdateOne(ctx)
	if err != nil {
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}

func (d *AssetFolderDao) PutAssetId(ctx context.Context, folderObjID primitive.ObjectID, assetObjID primitive.ObjectID) (int64, error) {
	updateResult, err := d.coll.Updater().
		Filter(query.Id(folderObjID)).
		Updates(update.NewBuilder().Push("assets", assetObjID).Set("updated_at", time.Now()).Build()).
		UpdateOne(ctx)
	if err != nil {
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}

func (d *AssetFolderDao) ModifyNameById(ctx context.Context, objectID primitive.ObjectID, name string) (int64, error) {
	updateResult, err := d.coll.Updater().Filter(query.Id(objectID)).Updates(update.NewBuilder().Set("name", name).Set("updated_at", time.Now()).Build()).UpdateOne(ctx)
	if err != nil {
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}

func (d *AssetFolderDao) DeleteSubFolderById(ctx context.Context, objectID primitive.ObjectID, subObjID primitive.ObjectID) (int64, error) {
	updateResult, err := d.coll.Updater().Filter(query.Id(objectID)).Updates(update.NewBuilder().Pull("child_folders", bsonx.Id(subObjID)).Set("updated_at", time.Now()).Build()).UpdateOne(ctx)
	if err != nil {
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}

func (d *AssetFolderDao) ModifySubFolderById(ctx context.Context, objectID primitive.ObjectID, assetFolder *AssetFolder) (int64, error) {
	updateResult, err := d.coll.Updater().Filter(query.NewBuilder().Id(objectID).Eq("child_folders._id", assetFolder.ID).Build()).Updates(update.NewBuilder().Set("child_folders.$", assetFolder).Set("updated_at", assetFolder.UpdatedAt).Build()).UpdateOne(ctx)
	if err != nil {
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}

func (d *AssetFolderDao) AddSubFolder(ctx context.Context, id primitive.ObjectID, assetFolder *AssetFolder) (int64, error) {
	updateResult, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.NewBuilder().Push("child_folders", assetFolder).Set("updated_at", assetFolder.CreatedAt).Build()).UpdateOne(ctx)
	if err != nil {
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}

func (d *AssetFolderDao) DeleteById(ctx context.Context, objectID primitive.ObjectID) (int64, error) {
	deleteResult, err := d.coll.Deleter().Filter(query.Id(objectID)).DeleteOne(ctx)
	if err != nil {
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}

func (d *AssetFolderDao) FindById(ctx context.Context, objectID primitive.ObjectID) (*AssetFolder, error) {
	return d.coll.Finder().Filter(query.Id(objectID)).FindOne(ctx)
}

func (d *AssetFolderDao) ModifyById(ctx context.Context, assetFolder *AssetFolder) (int64, error) {
	updateResult, err := d.coll.Updater().Filter(query.Id(assetFolder.ID)).Updates(bsonx.M("$set", assetFolder)).UpdateOne(ctx)
	if err != nil {
		return 0, err
	}
	return updateResult.ModifiedCount, nil
}

func (d *AssetFolderDao) Add(ctx context.Context, assetFolder *AssetFolder) (objID primitive.ObjectID, err error) {
	var insertOneResult *mongo.InsertOneResult
	insertOneResult, err = d.coll.Creator().InsertOne(ctx, assetFolder)
	if err != nil {
		return
	}
	return insertOneResult.InsertedID.(primitive.ObjectID), nil
}

func (d *AssetFolderDao) FindByAssetTypeAndType(ctx context.Context, assertType string, typ string) ([]*AssetFolder, error) {
	assertFolders, err := d.coll.Finder().
		Filter(query.NewBuilder().Eq("asset_type", assertType).Eq("type", typ).Build()).
		Find(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find asset folder by asset type[%s] and type[%s]", assertType, typ)
	}
	return assertFolders, nil
}
