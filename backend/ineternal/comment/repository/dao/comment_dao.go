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
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/types"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type Comment struct {
	Id string `bson:"_id"`
	types.Comment
	// 该评论下的所有回复的内容
	Replies []CommentReply `bson:"replies"`
	// 评论状态：审核不通过 0 未审核 1 审核通过 2
	Status int `bson:"status"`
	// 评论时间
	CreateTime int64 `bson:"created_at" bson:"create_time"`
	// 修改时间
	UpdateTime int64 `bson:"updated_at" bson:"update_time"`
}

type UserInfo4Reply types.UserInfo4Comment

type CommentReply struct {
	ReplyId string `bson:"reply_id"`
	// 回复内容
	Content string `bson:"content"`
	// 被回复的回复 Id
	ReplyToId string `bson:"reply_to_id"`
	// 用户信息
	UserInfo UserInfo4Reply `bson:"user_info"`
	// 被回复用户的信息
	RepliedUserInfo UserInfo4Reply `bson:"replied_user_info"`
	// 回复时间
	CreateTime int64 `bson:"created_at" bson:"create_time"`
	// 修改时间
	UpdateTime int64 `bson:"updated_at" bson:"update_time"`
	// 评论状态：审核不通过 0 未审核 1 审核通过 2
	Status int `bson:"status"`
}

type ICommentDao interface {
	AddComment(ctx context.Context, comment Comment) (any, error)
}

func NewCommentDao(db *mongo.Database) *CommentDao {
	return &CommentDao{
		coll: db.Collection("comment"),
	}
}

var _ ICommentDao = (*CommentDao)(nil)

type CommentDao struct {
	coll *mongo.Collection
}

func (d *CommentDao) AddComment(ctx context.Context, comment Comment) (any, error) {
	result, err := d.coll.InsertOne(ctx, comment)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to insert into %s, comment=%v", d.coll.Name(), comment)
	}
	return result.InsertedID, nil
}
