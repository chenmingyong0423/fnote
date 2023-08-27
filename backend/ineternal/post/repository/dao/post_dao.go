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
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	Sug      string   `bson:"_id"`
	Author   string   `bson:"author"`
	Title    string   `bson:"title"`
	Summary  string   `bson:"summary"`
	Content  string   `bson:"content"`
	CoverImg string   `bson:"cover_img"`
	Category string   `bson:"category"`
	Tags     []string `bson:"tags"`
	// 0 - 草稿，1 - 私密，2 - 已发布
	Status          string   `bson:"status"`
	Likes           []string `bson:"likes"`
	LikeCount       int      `bson:"like_count"`
	Comments        int      `bson:"comments"`
	Visits          int      `bson:"visit"`
	Priority        int      `bson:"priority"`
	MetaDescription string   `bson:"meta_description"`
	MetaKeywords    string   `bson:"meta_keywords"`
	WordCount       int      `bson:"word_count"`
	AllowComment    bool     `bson:"allow_comment"`
	CreateTime      int64    `bson:"create_time"`
	UpdateTime      int64    `bson:"update_time"`
}
type IPostDao interface {
	GetLatest5Posts(ctx context.Context) ([]*Post, error)
	QueryPostsPage(ctx context.Context, con bson.D, findOptions *options.FindOptions) ([]*Post, int64, error)
	GetPostById(ctx context.Context, sug string) (*Post, error)
	IncreaseVisitsById(ctx context.Context, sug string) (int64, error)
	FindByIdAndIp(ctx context.Context, sug string, ip string) (*Post, error)
	AddLike(ctx context.Context, sug string, ip string) error
	DeleteLike(ctx context.Context, sug string, ip string) error
}

var _ IPostDao = (*PostDao)(nil)

func NewPostDao(coll *mongo.Collection) *PostDao {
	return &PostDao{
		coll: coll,
	}
}

type PostDao struct {
	db   *mongo.Database
	coll *mongo.Collection
}

func (d *PostDao) DeleteLike(ctx context.Context, sug string, ip string) error {
	result, err := d.coll.UpdateByID(ctx, sug, bson.D{
		{"$pull", bson.D{{"likes", ip}}},
		{"$inc", bson.D{{"like_count", -1}}},
	})
	if err != nil {
		return errors.Wrapf(err, "Fails to delete a like, sug=%s, ip=%s", sug, ip)
	}
	if result.ModifiedCount == 0 {
		return errors.Wrapf(err, "ModifiedCount = 0, fails to delete a like, sug=%s, ip=%s", sug, ip)
	}
	return nil
}

func (d *PostDao) AddLike(ctx context.Context, sug string, ip string) error {
	result, err := d.coll.UpdateByID(ctx, sug, bson.D{
		{"$push", bson.D{{"likes", ip}}},
		{"$inc", bson.D{{"like_count", 1}}},
	})
	if err != nil {
		return errors.Wrapf(err, "Fails to add a like, sug=%s, ip=%s", sug, ip)
	}
	if result.ModifiedCount == 0 {
		return errors.Wrapf(err, "ModifiedCount = 0, fails to add a like, sug=%s, ip=%s", sug, ip)
	}
	return nil
}

func (d *PostDao) FindByIdAndIp(ctx context.Context, sug string, ip string) (*Post, error) {
	post := new(Post)
	err := d.coll.FindOne(ctx, bson.D{{"_id", sug}, {"likes", ip}}).Decode(post)
	if err != nil {
		return nil, errors.Wrapf(err, "Fails to find the documents from %s, sug=%s, ip=%s", d.coll.Name(), sug, ip)
	}
	return post, nil
}

func (d *PostDao) IncreaseVisitsById(ctx context.Context, sug string) (int64, error) {
	result, err := d.coll.UpdateByID(ctx, sug, bson.D{{"$inc", bson.D{{"visits", 1}}}})
	if err != nil {
		return 0, errors.Wrapf(err, "the visits of post increases failed, id=%s", sug)
	}
	return result.UpsertedCount, nil
}

func (d *PostDao) GetPostById(ctx context.Context, sug string) (*Post, error) {
	post := new(Post)
	err := d.coll.FindOne(ctx, bson.M{"_id": sug}).Decode(post)
	if err != nil {
		return nil, errors.Wrapf(err, "Fails to find the documents from %s, sug=%s", d.coll.Name(), sug)
	}
	return post, nil
}

func (d *PostDao) QueryPostsPage(ctx context.Context, con bson.D, findOptions *options.FindOptions) ([]*Post, int64, error) {
	cnt, err := d.coll.CountDocuments(ctx, con)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "Fails to find the count of documents from %s, con=%v", d.coll.Name(), con)
	}
	cursor, err := d.coll.Find(ctx, con, findOptions)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "Fails to find the documents from %s, con=%v, findOptions=%v", d.coll.Name(), con, findOptions)
	}
	posts := make([]*Post, 5)
	err = cursor.All(ctx, &posts)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Fails to decode the result")
	}
	return posts, cnt, nil
}

func (d *PostDao) GetLatest5Posts(ctx context.Context) ([]*Post, error) {
	findOptions := options.Find().SetSort(bson.M{"create_time": -1}).SetLimit(5)
	cursor, err := d.coll.Find(ctx, bson.M{}, findOptions)
	defer cursor.Close(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Fails to find the documents from %s, findOptions=%v", d.coll.Name(), findOptions)
	}
	posts := make([]*Post, 5)
	err = cursor.All(ctx, &posts)
	if err != nil {
		return nil, errors.Wrap(err, "Fails to decode the result")
	}
	return posts, nil
}
