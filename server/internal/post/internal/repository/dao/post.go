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
	"time"

	"github.com/chenmingyong0423/go-mongox"
	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"github.com/chenmingyong0423/go-mongox/builder/update"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	Id               string          `bson:"_id"`
	Author           string          `bson:"author"`
	Title            string          `bson:"title"`
	Summary          string          `bson:"summary"`
	Content          string          `bson:"content"`
	CoverImg         string          `bson:"cover_img"`
	Categories       []Category4Post `bson:"categories"`
	Tags             []Tag4Post      `bson:"tags"`
	IsDisplayed      bool            `bson:"is_displayed"`
	LikeCount        int             `bson:"like_count,omitempty"`
	CommentCount     int             `bson:"comment_count,omitempty"`
	VisitCount       int             `bson:"visit_count,omitempty"`
	StickyWeight     int             `bson:"sticky_weight"`
	MetaDescription  string          `bson:"meta_description"`
	MetaKeywords     string          `bson:"meta_keywords"`
	WordCount        int             `bson:"word_count"`
	IsCommentAllowed bool            `bson:"is_comment_allowed"`
	CreatedAt        time.Time       `bson:"created_at"`
	UpdatedAt        time.Time       `bson:"updated_at"`
}

type Category4Post struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
}

type Tag4Post struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
}

type IPostDao interface {
	GetFrontPosts(ctx context.Context, count int64) ([]*Post, error)
	QueryPostsPage(ctx context.Context, con bson.D, findOptions *options.FindOptions) ([]*Post, int64, error)
	GetPunishedPostById(ctx context.Context, sug string) (*Post, error)
	FindByIdAndIp(ctx context.Context, sug string, ip string) (*Post, error)
	AddLike(ctx context.Context, sug string, ip string) error
	DeleteLike(ctx context.Context, sug string, ip string) error
	IncreaseFieldById(ctx context.Context, id string, field string) error
	AddPost(ctx context.Context, post *Post) error
	DeleteById(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (*Post, error)
	DecreaseByField(ctx context.Context, id string, filedName string, cnt int) error
	SavePost(ctx context.Context, post *Post) error
	UpdateIsDisplayedById(ctx context.Context, id string, isDisplayed bool) error
	UpdateIsCommentAllowedById(ctx context.Context, id string, isCommentAllowed bool) error
	IncreasePostLikeCount(ctx context.Context, postId string) error
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

func (d *PostDao) IncreasePostLikeCount(ctx context.Context, postId string) error {
	updateResult, err := d.coll.Updater().Filter(query.Id(postId)).Updates(update.Inc("like_count", 1)).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to increase the like_count of post, id=%s", postId)
	}
	if updateResult.ModifiedCount == 0 {
		return fmt.Errorf("updateResult.ModifiedCount = 0, fails to increase the like_count of post, id=%s", postId)
	}
	return nil
}

func (d *PostDao) UpdateIsCommentAllowedById(ctx context.Context, id string, isCommentAllowed bool) error {
	result, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.BsonBuilder().Set("is_comment_allowed", isCommentAllowed).Set("updated_at", time.Now().Local()).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to update the is_comment_allowed of post, id=%s, isCommentAllowed=%v", id, isCommentAllowed)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("fails to update the is_comment_allowed of post, id=%s, isCommentAllowed=%v", id, isCommentAllowed)
	}
	return nil
}

func (d *PostDao) UpdateIsDisplayedById(ctx context.Context, id string, isDisplayed bool) error {
	result, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.BsonBuilder().Set("is_displayed", isDisplayed).Set("updated_at", time.Now().Local()).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to update the is_displayed of post, id=%s, isDisplayed=%v", id, isDisplayed)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("fails to update the is_displayed of post, id=%s, isDisplayed=%v", id, isDisplayed)
	}
	return nil
}

func (d *PostDao) SavePost(ctx context.Context, post *Post) error {
	result, err := d.coll.Updater().Filter(query.Id(post.Id)).UpdatesWithOperator("$set", post).UpdateOne(ctx, options.Update().SetUpsert(true))
	if err != nil {
		return errors.Wrapf(err, "fails to update a post, post=%v", post)
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		return fmt.Errorf("fails to update a post, post=%v", post)
	}
	return nil
}

func (d *PostDao) DecreaseByField(ctx context.Context, id string, filedName string, cnt int) error {
	filter := query.Id(id)
	u := update.BsonBuilder().Inc(filedName, -cnt).Set("updated_at", time.Now().Local()).Build()
	result, err := d.coll.Updater().Filter(filter).Updates(u).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to decrease the %s of post, id=%s, cnt=%d", filedName, id, cnt)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("fails to decrease the %s of post, id=%s, cnt=%d", filedName, id, cnt)
	}
	return nil
}

func (d *PostDao) FindById(ctx context.Context, id string) (*Post, error) {
	post, err := d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to find the document from post, id=%s", id)
	}
	return post, nil
}

func (d *PostDao) DeleteById(ctx context.Context, id string) error {
	result, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to delete a post, id=%s", id)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("fails to delete a post, id=%s", id)
	}
	return nil
}

func (d *PostDao) AddPost(ctx context.Context, post *Post) error {
	_, err := d.coll.Creator().InsertOne(ctx, post)
	if err != nil {
		return errors.Wrapf(err, "fails to insert a post, post=%v", post)
	}
	return nil
}

func (d *PostDao) IncreaseFieldById(ctx context.Context, id string, field string) error {
	result, err := d.coll.Updater().Filter(bsonx.Id(id)).Updates(update.Inc(field, 1)).UpdateOne(ctx)
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
		Filter(query.BsonBuilder().Id(id).Add("is_displayed", true).Build()).
		Updates(update.BsonBuilder().Pull("likes", ip).Inc("like_count", -1).Build()).
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
		Filter(query.BsonBuilder().Id(id).Add("is_displayed", true).Build()).
		Updates(update.BsonBuilder().Push("likes", ip).Inc("like_count", 1).Build()).
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
	post, err := d.coll.Finder().Filter(query.BsonBuilder().Id(id).Add("likes", ip).Build()).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to find the documents from post, id=%s, ip=%s", id, ip)
	}
	return post, nil
}

func (d *PostDao) GetPunishedPostById(ctx context.Context, id string) (*Post, error) {
	post, err := d.coll.Finder().Filter(query.BsonBuilder().Id(id).Add("is_displayed", true).Build()).FindOne(ctx)
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
	posts, err := d.coll.Finder().Filter(con).Find(ctx, findOptions)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "fails to find the documents from post, con=%v, findOptions=%v", con, findOptions)
	}
	return posts, cnt, nil
}

func (d *PostDao) GetFrontPosts(ctx context.Context, count int64) ([]*Post, error) {
	findOptions := options.Find().SetSort(
		bsonx.NewD().Add("sticky_weight", -1).Add("created_at", -1).Build()).SetLimit(count)
	posts, err := d.coll.Finder().Filter(bsonx.M("is_displayed", true)).Find(ctx, findOptions)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to find the documents from post, findOptions=%v", findOptions)
	}
	return posts, nil
}
