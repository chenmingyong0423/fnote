// Copyright 2023 chenmingyong0423

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
	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"github.com/chenmingyong0423/go-mongox/builder/update"

	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	Sug              string            `bson:"_id"`
	Author           string            `bson:"author"`
	Title            string            `bson:"title"`
	Summary          string            `bson:"summary"`
	Content          string            `bson:"content"`
	CoverImg         string            `bson:"cover_img"`
	Category         string            `bson:"category"`
	Tags             []string          `bson:"tags"`
	Status           domain.PostStatus `bson:"status"`
	Likes            []string          `bson:"likes"`
	LikeCount        int               `bson:"like_count"`
	CommentCount     int               `bson:"comment_count"`
	VisitCount       int               `bson:"visit_count"`
	Priority         int               `bson:"priority"`
	MetaDescription  string            `bson:"meta_description"`
	MetaKeywords     string            `bson:"meta_keywords"`
	WordCount        int               `bson:"word_count"`
	IsCommentAllowed bool              `bson:"is_comment_allowed"`
	CreateTime       int64             `bson:"create_time"`
	UpdateTime       int64             `bson:"update_time"`
}
type IPostDao interface {
	GetLatest5Posts(ctx context.Context) ([]*Post, error)
	QueryPostsPage(ctx context.Context, con bson.D, findOptions *options.FindOptions) ([]*Post, int64, error)
	GetPunishedPostById(ctx context.Context, sug string) (*Post, error)
	FindByIdAndIp(ctx context.Context, sug string, ip string) (*Post, error)
	AddLike(ctx context.Context, sug string, ip string) error
	DeleteLike(ctx context.Context, sug string, ip string) error
	IncreaseFieldById(ctx context.Context, id string, field string) error
}

var _ IPostDao = (*PostDao)(nil)

func NewPostDao(db *mongo.Database) *PostDao {
	return &PostDao{
		coll: mongox.NewCollection[Post](db.Collection("posts")),
	}
}

type PostDao struct {
	coll *mongox.Collection[Post]
}

func (d *PostDao) IncreaseFieldById(ctx context.Context, id string, field string) error {
	// bson.D{bson.E{Key: "$inc", Value: bson.D{bson.E{Key: field, Value: 1}}}}
	result, err := d.coll.Updater().Filter(bsonx.Id(id)).Updates(update.Inc(bsonx.M(field, 1))).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to increase the %s of post, id=%s", field, id)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("fails to increase the %s of post, id=%s", field, id)
	}
	return nil
}

func (d *PostDao) DeleteLike(ctx context.Context, id string, ip string) error {
	result, err := d.coll.Updater().
		Filter(bsonx.Id(id)).
		Updates(update.BsonBuilder().Pull(bsonx.M("likes", ip)).Inc(bsonx.M("like_count", -1)).Build()).
		UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to delete a like, id=%s, ip=%s", id, ip)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("ModifiedCount = 0, fails to delete a like, id=%s, ip=%s", id, ip)
	}
	return nil
}

func (d *PostDao) AddLike(ctx context.Context, id string, ip string) error {
	result, err := d.coll.Updater().
		Filter(bsonx.Id(id)).
		Updates(update.BsonBuilder().Push(bsonx.M("likes", ip)).Inc(bsonx.M("like_count", 1)).Build()).
		UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to add a like, id=%s, ip=%s", id, ip)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("ModifiedCount = 0, fails to add a like, id=%s, ip=%s", id, ip)
	}
	return nil
}

func (d *PostDao) FindByIdAndIp(ctx context.Context, id string, ip string) (*Post, error) {
	// bson.D{bson.E{Key: "_id", Value: id}, bson.E{Key: "likes", Value: ip}}
	post, err := d.coll.Finder().Filter(query.BsonBuilder().Id(id).Add(bsonx.KV("likes", ip)).Build()).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to find the documents from post, id=%s, ip=%s", id, ip)
	}
	return post, nil
}

func (d *PostDao) GetPunishedPostById(ctx context.Context, id string) (*Post, error) {
	post, err := d.coll.Finder().Filter(query.BsonBuilder().Id(id).Add(bsonx.KV("status", domain.PostStatusPunished)).Build()).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to find the document from post, id=%s", id)
	}
	return post, nil
}

func (d *PostDao) QueryPostsPage(ctx context.Context, con bson.D, findOptions *options.FindOptions) ([]*Post, int64, error) {
	cnt, err := d.coll.Finder().Filter(con).Count(ctx)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "fails to find the count of documents from post, con=%v", con)
	}
	posts, err := d.coll.Finder().Filter(con).Options(findOptions).Find(ctx)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "fails to find the documents from post, con=%v, findOptions=%v", con, findOptions)
	}
	return posts, cnt, nil
}

func (d *PostDao) GetLatest5Posts(ctx context.Context) ([]*Post, error) {
	findOptions := options.Find().SetSort(bsonx.M("create_time", -1)).SetLimit(5)
	posts, err := d.coll.Finder().Filter(bsonx.M("status", domain.PostStatusPunished)).Options(findOptions).Find(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to find the documents from post, findOptions=%v", findOptions)
	}
	return posts, nil
}
