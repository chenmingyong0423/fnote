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
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Asset struct {
	mongox.Model `bson:",inline"`
	// 素材标题
	Title string `bson:"title,omitempty"`
	// 素材内容
	Content string `bson:"content"`
	// 素材描述
	Description string `bson:"description,omitempty"`
	// 素材类型，image ···
	AssetType string `bson:"asset_type"`
	// 文件夹类型，post-editor ···
	Type string `bson:"type"`
	// 元数据
	Metadata map[string]any `bson:"metadata,omitempty"`
}

type IAssetDao interface {
	FindByIds(ctx context.Context, objIDs []primitive.ObjectID) ([]*Asset, error)
	Add(ctx context.Context, asset *Asset) (primitive.ObjectID, error)
	DeleteById(ctx context.Context, objectID primitive.ObjectID) (int64, error)
}

var _ IAssetDao = (*AssetDao)(nil)

func NewAssetDao(db *mongo.Database) *AssetDao {
	return &AssetDao{coll: mongox.NewCollection[Asset](db.Collection("assets"))}
}

type AssetDao struct {
	coll *mongox.Collection[Asset]
}

func (d *AssetDao) DeleteById(ctx context.Context, objectID primitive.ObjectID) (int64, error) {
	deleteResult, err := d.coll.Deleter().Filter(query.Id(objectID)).DeleteOne(ctx)
	if err != nil {
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}

func (d *AssetDao) Add(ctx context.Context, asset *Asset) (primitive.ObjectID, error) {
	insertOneResult, err := d.coll.Creator().InsertOne(ctx, asset)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return insertOneResult.InsertedID.(primitive.ObjectID), nil
}

func (d *AssetDao) FindByIds(ctx context.Context, objIDs []primitive.ObjectID) ([]*Asset, error) {
	return d.coll.Finder().Filter(query.In("_id", objIDs...)).Find(ctx)
}
