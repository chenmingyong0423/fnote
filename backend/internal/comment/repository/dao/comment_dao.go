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

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Comment struct {
	Id string `bson:"_id"`
	// 文章信息
	PostInfo PostInfo4Comment `bson:"post_info"`
	// 评论的内容
	Content string `bson:"content"`
	// 用户信息
	UserInfo UserInfo4Comment `bson:"user_info"`

	// 该评论下的所有回复的内容
	Replies []CommentReply `bson:"replies"`
	Status  CommentStatus  `bson:"status"`
	// 评论时间
	CreateTime int64 `bson:"create_time"`
	// 修改时间
	UpdateTime int64 `bson:"update_time"`
}

type CommentStatus uint

const (
	// CommentStatusPending 审核中
	CommentStatusPending CommentStatus = iota
	// CommentStatusApproved 审核通过
	CommentStatusApproved
	// CommentStatusRejected 审核不通过
	CommentStatusRejected
)

type UserInfo4Reply UserInfo4Comment

type UserInfo4Comment struct {
	Name    string `bson:"name"`
	Email   string `bson:"email"`
	Ip      string `bson:"ip"`
	Website string `bson:"website"`
}

type PostInfo4Comment struct {
	// 文章 ID
	PostId string `bson:"post_id"`
	// 文章标题字段
	PostTitle string `bson:"post_title"`
}

type LatestComment struct {
	PostInfo4Comment `bson:"post_info"`
	Name             string `bson:"name"`
	Content          string `bson:"content"`
	CreateTime       int64  `bson:"create_time"`
}

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
	Status          CommentStatus  `bson:"status"`
	// 回复时间
	CreateTime int64 `bson:"create_time"`
	// 修改时间
	UpdateTime int64 `bson:"update_time"`
}

type ICommentDao interface {
	AddComment(ctx context.Context, comment Comment) (any, error)
	FindCommentById(ctx context.Context, cmtId string) (*Comment, error)
	AddCommentReply(ctx context.Context, cmtId string, commentReply CommentReply) error
	FineLatestCommentAndReply(ctx context.Context, cnt int) ([]LatestComment, error)
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

func (d *CommentDao) FineLatestCommentAndReply(ctx context.Context, cnt int) ([]LatestComment, error) {
	pipeline := mongo.Pipeline{
		{primitive.E{Key: "$match", Value: bson.M{"status": CommentStatusApproved}}},
		{primitive.E{Key: "$project", Value: bson.M{
			"combined": bson.M{
				"$concatArrays": []any{
					[]bson.M{{"post_info": "$post_info", "name": "$user_info.name", "content": "$content", "create_time": "$create_time"}},
					//"$replies",
					bson.M{
						"$map": bson.M{
							"input": bson.M{
								"$filter": bson.M{
									"input": "$replies",
									"as":    "replyItem",
									"cond":  bson.M{"$eq": []interface{}{"$$replyItem.status", CommentStatusApproved}},
								},
							},
							"as": "reply",
							"in": bson.M{
								"post_info":   "$post_info",
								"name":        "$$reply.user_info.name",
								"content":     "$$reply.content",
								"create_time": "$$reply.create_time",
							},
						},
					},
				},
			},
		}}},
		{primitive.E{Key: "$unwind", Value: "$combined"}},
		{primitive.E{Key: "$replaceRoot", Value: bson.M{"newRoot": "$combined"}}},
		{primitive.E{Key: "$sort", Value: bson.M{"create_time": -1}}},
		{primitive.E{Key: "$limit", Value: cnt}},
	}
	// 执行聚合查询
	cursor, err := d.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.Wrapf(err, "Fails to execute aggregation operation, pipeline=%v", pipeline)
	}
	defer cursor.Close(ctx)

	// 解析并输出结果
	var results []LatestComment
	if err = cursor.All(ctx, &results); err != nil {
		return nil, errors.Wrapf(err, "Fails to cursor.All, cursor=%v", cursor)
	}
	return results, nil
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
	err := d.coll.FindOne(ctx, bson.D{bson.E{Key: "_id", Value: cmtId}, bson.E{Key: "status", Value: CommentStatusApproved}}).Decode(comment)
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
