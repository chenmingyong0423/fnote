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

	"github.com/chenmingyong0423/go-mongox/builder/query"

	"github.com/chenmingyong0423/go-mongox"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/chenmingyong0423/go-mongox/builder/update"

	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/aggregation"
	"github.com/chenmingyong0423/go-mongox/types"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type Comment struct {
	Id string `bson:"_id"`
	// 文章信息
	PostInfo PostInfo `bson:"post_info"`
	// 评论的内容
	Content string `bson:"content"`
	// 用户信息
	UserInfo UserInfo4Comment `bson:"user_info"`

	// 该评论下的所有回复的内容
	Replies []Reply       `bson:"replies"`
	Status  CommentStatus `bson:"status"`
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
	// CommentStatusHidden 隐藏
	CommentStatusHidden
	// CommentStatusRejected 审核不通过
	CommentStatusRejected
)

type UserInfo4Reply UserInfo

type UserInfo4Comment UserInfo

type UserInfo struct {
	Name    string `bson:"name"`
	Email   string `bson:"email"`
	Ip      string `bson:"ip"`
	Website string `bson:"website"`
}

type PostInfo struct {
	// 文章 ID
	PostId string `bson:"post_id"`
	// 文章标题字段
	PostTitle string `bson:"post_title"`
	// 文章链接
	PostUrl string `bson:"post_url"`
}

type LatestComment struct {
	PostInfo   `bson:"post_info"`
	Name       string `bson:"name"`
	Email      string `bson:"email"`
	Content    string `bson:"content"`
	CreateTime int64  `bson:"create_time"`
}

type Reply struct {
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

type ReplyWithPostInfo struct {
	Reply    `bson:"inline"`
	PostInfo `bson:"post_info"`
}

type AdminComment struct {
	Id string `bson:"_id"`
	// 评论的内容
	PostInfo   PostInfo      `bson:"post_info"`
	Content    string        `bson:"content"`
	UserInfo   UserInfo      `bson:"user_info"`
	Fid        string        `bson:"fid"`
	Type       int           `bson:"type"`
	Status     CommentStatus `bson:"status"`
	CreateTime int64         `bson:"create_time"`
}

type ICommentDao interface {
	AddComment(ctx context.Context, comment Comment) (string, error)
	FindApprovedCommentById(ctx context.Context, cmtId string) (*Comment, error)
	AddCommentReply(ctx context.Context, cmtId string, commentReply Reply) error
	FineLatestCommentAndReply(ctx context.Context, cnt int) ([]LatestComment, error)
	FindCommentsByPostIdAndCmtStatus(ctx context.Context, postId string, cmtStatus uint) ([]*Comment, error)
	AggregationQuerySkipAndSetLimit(ctx context.Context, cond bson.D, findOptions *options.FindOptions) ([]AdminComment, int64, error)
	FindCommentById(ctx context.Context, id string) (*Comment, error)
	UpdateCommentStatus(ctx context.Context, id string, commentStatus CommentStatus) error
	FindReplyByCIdAndRId(ctx context.Context, commentId string, replyId string) (*ReplyWithPostInfo, error)
	UpdateCommentReplyStatus(ctx context.Context, commentId string, replyId string, commentStatus CommentStatus) error
	FindCommentWithRepliesById(ctx context.Context, id string) (*Comment, error)
	DeleteById(ctx context.Context, id string) error
	DeleteReplyByCIdAndRId(ctx context.Context, commentId string, replyId string) error
}

func NewCommentDao(db *mongo.Database) *CommentDao {
	return &CommentDao{
		coll: mongox.NewCollection[Comment](db.Collection("comments")),
	}
}

var _ ICommentDao = (*CommentDao)(nil)

type CommentDao struct {
	coll *mongox.Collection[Comment]
}

func (d *CommentDao) DeleteReplyByCIdAndRId(ctx context.Context, commentId string, replyId string) error {
	result, err := d.coll.Updater().Filter(query.Id(commentId)).Updates(update.Pull("replies", bsonx.M("reply_id", replyId))).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to update one from comment, commentId=%s, replyId=%s", commentId, replyId)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("modifiedCount = 0, fails to update one from comment, commentId=%s, replyId=%s", commentId, replyId)
	}
	return nil
}

func (d *CommentDao) DeleteById(ctx context.Context, id string) error {
	result, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to delete one from comment, id=%s", id)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("deletedCount = 0, fails to delete one from comment, id=%s", id)
	}
	return nil
}

func (d *CommentDao) FindCommentWithRepliesById(ctx context.Context, id string) (*Comment, error) {
	comment, err := d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to find the document from comment, id=%s", id)
	}
	return comment, nil
}

func (d *CommentDao) UpdateCommentReplyStatus(ctx context.Context, commentId string, replyId string, commentStatus CommentStatus) error {
	filter := query.BsonBuilder().Id(commentId).Add("replies.reply_id", replyId).Build()
	updates := update.BsonBuilder().Set("replies.$.status", commentStatus).Set("replies.$.update_time", time.Now().Unix()).Build()
	result, err := d.coll.Updater().Filter(filter).Updates(updates).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to update one from comment, filter=%v, update=%v", filter, updates)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("modifiedCount = 0, fails to update one from comment, filter=%v, update=%v", filter, updates)
	}
	return nil
}

func (d *CommentDao) FindReplyByCIdAndRId(ctx context.Context, commentId string, replyId string) (*ReplyWithPostInfo, error) {
	pipeline := aggregation.StageBsonBuilder().
		Match(query.Id(commentId)).
		Unwind("$replies", nil).
		Match(bsonx.M("replies.reply_id", replyId)).
		AddFields(bsonx.M("replies.post_info", "$post_info")).
		ReplaceWith("$replies").Build()
	var result []ReplyWithPostInfo
	err := d.coll.Aggregator().Pipeline(pipeline).AggregateWithCallback(ctx, func(ctx context.Context, cursor *mongo.Cursor) error {
		return cursor.All(ctx, &result)
	})
	if err != nil {
		return nil, errors.Wrapf(err, "fails to execute aggregation operation, pipeline=%v", pipeline)
	}
	if len(result) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return &result[0], nil
}

func (d *CommentDao) UpdateCommentStatus(ctx context.Context, id string, commentStatus CommentStatus) error {
	result, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.BsonBuilder().Set("status", commentStatus).Set("update_time", time.Now().Unix()).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to update one from comment, id=%s, commentStatus=%d", id, commentStatus)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("modifiedCount = 0, fails to update one from comment, id=%s, commentStatus=%d", id, commentStatus)
	}
	return nil
}

func (d *CommentDao) FindCommentById(ctx context.Context, id string) (*Comment, error) {
	comment, err := d.coll.Finder().Filter(query.Id(id)).FindOne(ctx, options.FindOne().SetProjection(bsonx.M("replies", 0)))
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (d *CommentDao) AggregationQuerySkipAndSetLimit(ctx context.Context, cond bson.D, findOptions *options.FindOptions) ([]AdminComment, int64, error) {
	pipeline := aggregation.StageBsonBuilder().
		Match(bson.D{}).
		Project(aggregation.ConcatArrays("combined", []any{
			bsonx.A(bsonx.NewD().Add("_id", "$_id").Add("post_info", "$post_info").Add("user_info", "$user_info").Add("content", "$content").Add("create_time", "$create_time").Add("type", 0).Add("status", "$status").Build()),
			aggregation.MapWithoutKey(
				"$replies",
				"reply",
				bsonx.NewD().Add("_id", "$$reply.reply_id").Add("fid", "$_id").Add("post_info", "$post_info").Add("user_info", "$$reply.user_info").Add("content", "$$reply.content").Add("create_time", "$$reply.create_time").Add("type", 1).Add("status", "$$reply.status").Build(),
			),
		}...)).
		Unwind("$combined", nil).
		ReplaceWith("$combined").
		Sort(bsonx.M("create_time", -1)).Build()

	var results []AdminComment

	// 执行聚合查询
	err := d.coll.Aggregator().Pipeline(pipeline).AggregateWithCallback(ctx, func(ctx context.Context, cursor *mongo.Cursor) error {
		// 解析并输出结果
		return cursor.All(ctx, &results)
	})
	if err != nil {
		return nil, 0, errors.Wrapf(err, "Fails to execute aggregation operation, pipeline=%v", pipeline)
	}

	return results, 0, nil
}

func (d *CommentDao) FindCommentsByPostIdAndCmtStatus(ctx context.Context, postId string, cmtStatus uint) ([]*Comment, error) {
	cond := bsonx.D(bsonx.E("post_info.post_id", postId), bsonx.E("status", cmtStatus))
	result, err := d.coll.Finder().Filter(cond).Find(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Fails to find the docs from comment, condition=%v", cond)
	}
	return result, nil
}

func (d *CommentDao) FineLatestCommentAndReply(ctx context.Context, cnt int) ([]LatestComment, error) {
	pipeline := aggregation.StageBsonBuilder().
		Match(bsonx.M("status", CommentStatusApproved)).
		Project(aggregation.ConcatArrays("combined", []any{
			bsonx.A(bsonx.NewD().Add("post_info", "$post_info").Add("name", "$user_info.name").Add("email", "$user_info.email").Add("content", "$content").Add("create_time", "$create_time").Build()),
			aggregation.MapWithoutKey(
				aggregation.FilterWithoutKey("$replies", aggregation.EqWithoutKey("$$replyItem.status", CommentStatusApproved), &types.FilterOptions{As: "replyItem"}),
				"reply",
				bsonx.NewD().Add("post_info", "$post_info").Add("name", "$$reply.user_info.name").Add("email", "$$reply.user_info.email").Add("content", "$$reply.content").Add("create_time", "$$reply.create_time").Build(),
			),
		}...)).
		Unwind("$combined", nil).
		ReplaceWith("$combined").
		Sort(bsonx.M("create_time", -1)).
		Limit(int64(cnt)).Build()
	//pipeline := mongo.Pipeline{
	//	{primitive.E{Key: "$match", Value: bson.M{"status": CommentStatusApproved}}},
	//	{primitive.E{Key: "$project", Value: bson.M{
	//		"combined": bson.M{
	//			"$concatArrays": []any{
	//				bson.A{bson.M{"post_info": "$post_info", "name": "$user_info.name", "content": "$content", "create_time": "$create_time", "email": "$user_info.email"}},
	//				//"$replies",
	//				bson.M{
	//					"$map": bson.M{
	//						"input": bson.M{
	//							"$filter": bson.M{
	//								"input": "$replies",
	//								"as":    "replyItem",
	//								"cond":  bson.M{"$eq": []interface{}{"$$replyItem.status", CommentStatusApproved}},
	//							},
	//						},
	//						"as": "reply",
	//						"in": bson.M{
	//							"post_info":   "$post_info",
	//							"name":        "$$reply.user_info.name",
	//							"content":     "$$reply.content",
	//							"create_time": "$$reply.create_time",
	// 							"email": "$$reply.user_info.email",
	//						},
	//					},
	//				},
	//			},
	//		},
	//	}}},
	//	{primitive.E{Key: "$unwind", Value: "$combined"}},
	//	{primitive.E{Key: "$replaceRoot", Value: bson.M{"newRoot": "$combined"}}},
	//	{primitive.E{Key: "$sort", Value: bson.M{"create_time": -1}}},
	//	{primitive.E{Key: "$limit", Value: cnt}},
	//}

	var results []LatestComment

	// 执行聚合查询
	err := d.coll.Aggregator().Pipeline(pipeline).AggregateWithCallback(ctx, func(ctx context.Context, cursor *mongo.Cursor) error {
		// 解析并输出结果
		return cursor.All(ctx, &results)
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Fails to execute aggregation operation, pipeline=%v", pipeline)
	}

	return results, nil
}

func (d *CommentDao) AddCommentReply(ctx context.Context, cmtId string, commentReply Reply) error {
	// 构建查询条件
	filter := query.Id(cmtId)
	// 构建更新操作
	updates := update.Push("replies", commentReply)

	result, err := d.coll.Updater().Filter(filter).Updates(updates).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to update one from comment, filter=%v, update=%v", filter, updates)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("modifiedCount = 0, fails to update one from comment, filter=%v, update=%v", filter, updates)
	}
	return nil
}

func (d *CommentDao) FindApprovedCommentById(ctx context.Context, cmtId string) (*Comment, error) {
	comment, err := d.coll.Finder().Filter(bsonx.D(bsonx.E("_id", cmtId), bsonx.E("status", CommentStatusApproved))).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to find the document from comment, cmtId=%s", cmtId)
	}
	return comment, nil
}

func (d *CommentDao) AddComment(ctx context.Context, comment Comment) (string, error) {
	result, err := d.coll.Creator().InsertOne(ctx, comment)
	if err != nil {
		return "", errors.Wrapf(err, "fails to insert into comment, comment=%v", comment)
	}
	return result.InsertedID.(string), nil
}
