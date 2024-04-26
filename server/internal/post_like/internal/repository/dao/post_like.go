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
	"fmt"

	"github.com/chenmingyong0423/go-mongox"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostLike struct {
	mongox.Model `bson:",inline"`
	PostId       string `bson:"post_id"`
	Ip           string `bson:"ip"`
	UserAgent    string `bson:"user_agent"`
}

type IPostLikeDao interface {
	Add(ctx context.Context, postLike *PostLike) (string, error)
	DeleteById(ctx context.Context, objectID primitive.ObjectID) error
}

var _ IPostLikeDao = (*PostLikeDao)(nil)

func NewPostLikeDao(db *mongo.Database) *PostLikeDao {
	return &PostLikeDao{coll: mongox.NewCollection[PostLike](db.Collection("post_likes"))}
}

type PostLikeDao struct {
	coll *mongox.Collection[PostLike]
}

func (d *PostLikeDao) DeleteById(ctx context.Context, objectID primitive.ObjectID) error {
	deleteResult, err := d.coll.Deleter().Filter(query.Id(objectID)).DeleteOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to delete post_likes by id: %s", objectID.Hex())
	}
	if deleteResult.DeletedCount == 0 {
		return fmt.Errorf("deleteResult.DeletedCount, failed to delete post_likes by id: %s", objectID.Hex())
	}
	return nil
}

func (d *PostLikeDao) Add(ctx context.Context, postLike *PostLike) (string, error) {
	insertOneResult, err := d.coll.Creator().InsertOne(ctx, postLike)
	if err != nil {
		return "", errors.Wrapf(err, "failed to insert one into post_likes: %v", postLike)
	}
	return insertOneResult.InsertedID.(primitive.ObjectID).Hex(), nil
}
