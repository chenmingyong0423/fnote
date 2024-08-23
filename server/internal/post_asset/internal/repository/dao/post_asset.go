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
	"github.com/chenmingyong0423/go-mongox"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostAsset struct {
	mongox.Model `bson:",inline"`
	// 素材标题
	Title string `bson:"title,omitempty"`
	// 素材内容
	Content string `bson:"content"`
	// 素材描述
	Description string `bson:"description,omitempty"`
	// 素材类型，image ···
	AssetType string `bson:"asset_type"`
	// 元数据
	Metadata map[string]any `bson:"metadata,omitempty"`
}

type IPostAssetDao interface {
}

var _ IPostAssetDao = (*PostAssetDao)(nil)

func NewPostAssetDao(db *mongo.Database) *PostAssetDao {
	return &PostAssetDao{coll: mongox.NewCollection[PostAsset](db.Collection("post_assets"))}
}

type PostAssetDao struct {
	coll *mongox.Collection[PostAsset]
}
