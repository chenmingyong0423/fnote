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

	"github.com/chenmingyong0423/gkit/uuidx"
	"github.com/chenmingyong0423/go-mongox"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostDraft struct {
	ID               string               `bson:"_id"`
	CreatedAt        time.Time            `bson:"created_at,omitempty"`
	UpdatedAt        time.Time            `bson:"updated_at"`
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

func (m *PostDraft) DefaultId() {
	if m.ID == "" {
		m.ID = uuidx.RearrangeUUID4()
	}
}

func (m *PostDraft) DefaultCreatedAt() {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().Local()
	}
}

func (m *PostDraft) DefaultUpdatedAt() {
	m.UpdatedAt = time.Now().Local()
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
	Save(ctx context.Context, postDraft *PostDraft) (string, error)
	GetById(ctx context.Context, id string) (*PostDraft, error)
	DeleteById(ctx context.Context, id string) (int64, error)
}

var _ IPostDraftDao = (*PostDraftDao)(nil)

func NewPostDraftDao(db *mongo.Database) *PostDraftDao {
	return &PostDraftDao{coll: mongox.NewCollection[PostDraft](db.Collection("post_draft"))}
}

type PostDraftDao struct {
	coll *mongox.Collection[PostDraft]
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

func (d *PostDraftDao) Save(ctx context.Context, postDraft *PostDraft) (string, error) {
	updateResult, err := d.coll.Updater().Filter(query.Id(postDraft.ID)).Replacement(postDraft).Upsert(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "failed to save post draft: %v", postDraft)
	}
	if updateResult.UpsertedCount == 0 && updateResult.ModifiedCount == 0 {
		return "", fmt.Errorf("UpsertedCount=0 || ModifiedCount=0, failed to save post draft: %v", postDraft)
	}
	if id, ok := updateResult.UpsertedID.(string); ok {
		return id, nil
	}
	return postDraft.ID, nil
}
