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
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/domain"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/types"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Comment struct {
	Id            string `bson:"_id"`
	types.Comment `bson:"inline"`
	// 该评论下的所有回复的内容
	Replies []CommentReply       `bson:"replies"`
	Status  domain.CommentStatus `bson:"status"`
	// 评论时间
	CreateTime int64 `bson:"created_at" bson:"create_time"`
	// 修改时间
	UpdateTime int64 `bson:"updated_at" bson:"update_time"`
}

type CommentReply struct {
	types.CommentReply `bson:"inline"`
	Status             domain.CommentStatus `bson:"status"`
	// 回复时间
	CreateTime int64 `bson:"created_at" bson:"create_time"`
	// 修改时间
	UpdateTime int64 `bson:"updated_at" bson:"update_time"`
}

type ICommentDao interface {
	AddComment(ctx context.Context, comment Comment) (any, error)
	FindCommentById(ctx context.Context, cmtId string) (*Comment, error)
	AddCommentReply(ctx context.Context, cmtId string, commentReply CommentReply) error
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

func (d *CommentDao) AddCommentReply(ctx context.Context, cmtId string, commentReply CommentReply) error {
	// 构建查询条件
	filter := bson.M{"_id": cmtId}

	// 构建更新操作
	update := bson.M{
		"$push": bson.M{"replies": commentReply},
	}
	result, err := d.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrapf(err, "fails to update one from %s, filter=%v, update=%v", d.coll.Name(), filter, update)
	}
	if result.ModifiedCount == 0 {
		return errors.Wrapf(err, "modifiedCount = 0, fails to update one from %s, filter=%v, update=%v", d.coll.Name(), filter, update)
	}
	return nil
}

func (d *CommentDao) FindCommentById(ctx context.Context, cmtId string) (*Comment, error) {
	comment := new(Comment)
	err := d.coll.FindOne(ctx, bson.D{bson.E{Key: "_id", Value: cmtId}, bson.E{Key: "status", Value: domain.CommentStatusApproved}}).Decode(comment)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to find the document from %s, cmtId=%s", d.coll.Name(), cmtId)
	}
	return comment, nil
}

func (d *CommentDao) AddComment(ctx context.Context, comment Comment) (any, error) {
	result, err := d.coll.InsertOne(ctx, comment)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to insert into %s, comment=%v", d.coll.Name(), comment)
	}
	return result.InsertedID, nil
}
