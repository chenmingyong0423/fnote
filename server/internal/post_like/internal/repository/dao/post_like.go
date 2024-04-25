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

type PostLike struct {
}

type IPostLikeDao interface {
}

var _ IPostLikeDao = (*PostLikeDao)(nil)

func NewPostLikeDao(db *mongo.Database) *PostLikeDao {
	return &PostLikeDao{coll: mongox.NewCollection[PostLike](db.Collection("post_like"))}
}

type PostLikeDao struct {
	coll *mongox.Collection[PostLike]
}
