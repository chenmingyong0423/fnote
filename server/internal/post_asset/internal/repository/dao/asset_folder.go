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
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
}

type IAssetFolderDao interface {
	FindByAssetTypeAndType(ctx context.Context, assertType string, typ string) (*AssetFolder, error)
}

var _ IAssetFolderDao = (*AssetFolderDao)(nil)

func NewAssetFolderDao(db *mongo.Database) *AssetFolderDao {
	return &AssetFolderDao{coll: mongox.NewCollection[AssetFolder](db.Collection("asset_folders"))}
}

type AssetFolderDao struct {
	coll *mongox.Collection[AssetFolder]
}

func (d *AssetFolderDao) FindByAssetTypeAndType(ctx context.Context, assertType string, typ string) (*AssetFolder, error) {
	assertFOlder, err := d.coll.Finder().
		Filter(query.NewBuilder().Eq("asset_type", assertType).Eq("type", typ).Build()).
		FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find asset folder by asset type[%s] and type[%s]", assertType, typ)
	}
	return assertFOlder, nil
}
