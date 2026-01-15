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
	"time"

	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/pkg/errors"
)

type PostDraft struct {
	ID              string    `bson:"_id"`
	CreatedAt       time.Time `bson:"created_at,omitempty"`
	UpdatedAt       time.Time `bson:"updated_at"`
	PostDraftFields `bson:",inline"`
}

type PostDraftUpdate struct {
	PostDraftFields `bson:",inline"`
}

type PostDraftFields struct {
	mongox.Model
	Author           string               `bson:"author"`
	Title            string               `bson:"title"`
	Summary          string               `bson:"summary"`
	Content          string               `bson:"content"`
	CoverImg         string               `bson:"cover_img"`
	Categories       []Category4PostDraft `bson:"categories"`
	Tags             []Tag4PostDraft      `bson:"tags"`
	IsDisplayed      bool                 `bson:"is_displayed"`
	StickyWeight     int                  `bson:"sticky_weight"`
	MetaDescription  string               `bson:"meta_description"`
	MetaKeywords     string               `bson:"meta_keywords"`
	WordCount        int                  `bson:"word_count"`
	IsCommentAllowed bool                 `bson:"is_comment_allowed"`
}

type Category4PostDraft struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
}

type Tag4PostDraft struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
}

type IPostDraftDao interface {
	Save(ctx context.Context, ID string, postDraftUpdate *PostDraftUpdate) (string, error)
	GetById(ctx context.Context, id string) (*PostDraft, error)
	DeleteById(ctx context.Context, id string) (int64, error)
	QueryPage(ctx context.Context, cond bson.D, findOptions *options.FindOptionsBuilder) ([]*PostDraft, int64, error)
}

var _ IPostDraftDao = (*PostDraftDao)(nil)

func NewPostDraftDao(db *mongox.Database) *PostDraftDao {
	return &PostDraftDao{coll: mongox.NewCollection[PostDraft](db, "post_draft")}
}

type PostDraftDao struct {
	coll *mongox.Collection[PostDraft]
}

func (d *PostDraftDao) QueryPage(ctx context.Context, cond bson.D, findOptions *options.FindOptionsBuilder) ([]*PostDraft, int64, error) {
	count, err := d.coll.Finder().Filter(cond).Count(ctx)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "failed to query the count of post draft: %v", cond)
	}
	postDrafts, err := d.coll.Finder().Filter(cond).Find(ctx, findOptions)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "failed to query post draft page: %v, %v", cond, findOptions)
	}
	return postDrafts, count, nil
}

func (d *PostDraftDao) DeleteById(ctx context.Context, id string) (int64, error) {
	deleteResult, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	if err != nil {
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}

func (d *PostDraftDao) GetById(ctx context.Context, id string) (*PostDraft, error) {
	postDraft, err := d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get post draft by id: %s", id)
	}
	return postDraft, nil
}

func (d *PostDraftDao) Save(ctx context.Context, ID string, postDraftUpdate *PostDraftUpdate) (string, error) {
	updateResult, err := d.coll.Updater().Filter(query.Id(ID)).Updates(update.SetFields(postDraftUpdate)).Upsert(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "failed to save post draft: %v", postDraftUpdate)
	}
	if updateResult.UpsertedCount == 0 && updateResult.ModifiedCount == 0 {
		return "", fmt.Errorf("UpsertedCount=0 || ModifiedCount=0, failed to save post draft: %v", postDraftUpdate)
	}
	if id, ok := updateResult.UpsertedID.(string); ok {
		return id, nil
	}
	return ID, nil
}
