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

	"go.mongodb.org/mongo-driver/bson/primitive"

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
	mongox.Model `bson:",inline"`
	// 文章信息
	PostInfo PostInfo `bson:"post_info"`
	// 评论的内容
	Content string `bson:"content"`
	// 用户信息
	UserInfo UserInfo4Comment `bson:"user_info"`

	// 该评论下的所有回复的内容
	Replies        []Reply `bson:"replies"`
	ApprovalStatus bool    `bson:"approval_status"`
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
	ApprovalStatus  bool           `bson:"approval_status"`
	// 回复时间
	CreatedAt time.Time `bson:"created_at"`
	// 修改时间
	UpdatedAt time.Time `bson:"updated_at"`
}

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
	PostInfo  `bson:"post_info"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	Content   string    `bson:"content"`
	CreatedAt time.Time `bson:"created_at"`
}

type ReplyWithPostInfo struct {
	Reply    `bson:"inline"`
	PostInfo `bson:"post_info"`
}

type AdminComment struct {
	Id string `bson:"_id"`
	// 评论的内容
	PostInfo       PostInfo `bson:"post_info"`
	Content        string   `bson:"content"`
	UserInfo       UserInfo `bson:"user_info"`
	Fid            string   `bson:"fid"`
	Type           int      `bson:"type"`
	ApprovalStatus bool     `bson:"approval_status"`
	CreateTime     int64    `bson:"created_at"`
}

type ICommentDao interface {
	AddComment(ctx context.Context, comment *Comment) (string, error)
	FindApprovedCommentById(ctx context.Context, objectID primitive.ObjectID) (*Comment, error)
	AddCommentReply(ctx context.Context, objectID primitive.ObjectID, commentReply Reply) error
	FineLatestCommentAndReply(ctx context.Context, cnt int) ([]LatestComment, error)
	FindApprovedCommentsByPostId(ctx context.Context, postId string) ([]*Comment, error)
	AggregationQuerySkipAndSetLimit(ctx context.Context, cond bson.D, findOptions *options.FindOptions) ([]AdminComment, int64, error)
	FindCommentById(ctx context.Context, objectID primitive.ObjectID) (*Comment, error)
	UpdateCommentStatus2True(ctx context.Context, objectID primitive.ObjectID) error
	FindReplyByCIdAndRId(ctx context.Context, objectCommentID primitive.ObjectID, replyId string) (*ReplyWithPostInfo, error)
	UpdateCommentReplyStatus(ctx context.Context, objectID primitive.ObjectID, replyId string, commentStatus bool) error
	FindCommentWithRepliesById(ctx context.Context, objectID primitive.ObjectID) (*Comment, error)
	DeleteById(ctx context.Context, objectID primitive.ObjectID) error
	DeleteReplyByCIdAndRId(ctx context.Context, objectID primitive.ObjectID, replyId string) error
	CountOfToday(ctx context.Context) (int64, error)
	Find(ctx context.Context, findOptions *options.FindOptions) ([]*Comment, int64, error)
	UpdateCommentStatus2TrueByIds(ctx context.Context, ids []primitive.ObjectID) error
	FindByObjectIDs(ctx context.Context, ids []primitive.ObjectID) ([]*Comment, error)
	UpdateCReplyStatus2TrueByCidAndRIds(ctx context.Context, commentObjectID primitive.ObjectID, replyIds []string) error
	FindWithDisapprovedReplyByCidAndRIds(ctx context.Context, commentObjID primitive.ObjectID, replyIds []string) (*Comment, error)
	FindDisapprovedCommentByObjectIDs(ctx context.Context, commentObjectIDs []primitive.ObjectID) ([]*Comment, error)
	FindByAggregation(ctx context.Context, pipeline mongo.Pipeline) ([]*Comment, error)
	DeleteByIds(ctx context.Context, ids []primitive.ObjectID) error
	PullReplyByCIdAndRIds(ctx context.Context, commentId primitive.ObjectID, replyIds []string) error
	DeleteManyByPostId(ctx context.Context, postId string) error
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

func (d *CommentDao) DeleteManyByPostId(ctx context.Context, postId string) error {
	deleteResult, err := d.coll.Deleter().Filter(query.Eq("post_info.post_id", postId)).DeleteMany(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to delete comments, postId: %s", postId)
	}
	if deleteResult.DeletedCount == 0 {
		return fmt.Errorf("DeletedCount=0, failed to delete comments, postId: %s", postId)
	}
	return nil
}

func (d *CommentDao) PullReplyByCIdAndRIds(ctx context.Context, commentObjID primitive.ObjectID, replyIds []string) error {
	updateResult, err := d.coll.Updater().Filter(query.Id(commentObjID)).
		Updates(update.Pull("replies", query.In("reply_id", replyIds...))).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to pull reply by comment id: %v, reply ids: %v", commentObjID.Hex(), replyIds)
	}
	if updateResult.ModifiedCount == 0 {
		return fmt.Errorf("ModifiedCount=0, failed to pull reply by comment id: %v, reply ids: %v", commentObjID.Hex(), replyIds)
	}
	return nil
}

func (d *CommentDao) DeleteByIds(ctx context.Context, ids []primitive.ObjectID) error {
	deleteResult, err := d.coll.Deleter().Filter(query.In("_id", ids...)).DeleteMany(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to delete comment by ids: %v", ids)
	}
	if deleteResult.DeletedCount == 0 {
		return fmt.Errorf("DeletedCount=0, failed to delete comment by ids: %v", ids)
	}
	return nil
}

func (d *CommentDao) FindByAggregation(ctx context.Context, pipeline mongo.Pipeline) ([]*Comment, error) {
	return d.coll.Aggregator().Pipeline(pipeline).Aggregate(ctx)
}

func (d *CommentDao) FindDisapprovedCommentByObjectIDs(ctx context.Context, commentObjectIDs []primitive.ObjectID) ([]*Comment, error) {
	anyIds := make([]any, 0, len(commentObjectIDs))
	for _, objectID := range commentObjectIDs {
		anyIds = append(anyIds, objectID)
	}
	return d.coll.Finder().Filter(query.BsonBuilder().In("_id", anyIds...).Eq("approval_status", false).Build()).Find(ctx)
}

func (d *CommentDao) FindWithDisapprovedReplyByCidAndRIds(ctx context.Context, commentObjID primitive.ObjectID, replyIds []string) (*Comment, error) {
	pipeline := aggregation.StageBsonBuilder().
		Match(query.Id(commentObjID)).
		Project(
			bsonx.NewD().
				Add("_id", 1).
				Add("created_at", 1).
				Add("updated_at", 1).
				Add("post_info", 1).
				Add("content", 1).
				Add("user_info", 1).
				Add("approval_status", 1).
				Add("replies", aggregation.FilterWithoutKey("$replies",
					aggregation.AndWithoutKey(
						bsonx.D("$in", bsonx.A("$$reply.reply_id", replyIds)),
						aggregation.EqWithoutKey("$$reply.approval_status", false),
					), &types.FilterOptions{As: "reply"})).Build(),
		).Build()
	comments, err := d.coll.Aggregator().Pipeline(pipeline).Aggregate(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find comment by id=%s, pipeline=%v", commentObjID, pipeline)
	}
	if len(comments) == 0 {
		return nil, errors.Wrapf(mongo.ErrNoDocuments, "failed to find comment by id=%s", commentObjID)
	}
	return comments[0], nil
}

func (d *CommentDao) UpdateCReplyStatus2TrueByCidAndRIds(ctx context.Context, commentObjectID primitive.ObjectID, replyIds []string) error {
	now := time.Now().Local()
	updateResult, err := d.coll.Updater().Filter(query.Id(commentObjectID)).
		Updates(update.BsonBuilder().Set("updated_at", now).Set("replies.$[elem].approval_status", true).Set("replies.$[elem].updated_at", now).Build()).
		UpdateMany(ctx, options.Update().SetArrayFilters(options.ArrayFilters{Filters: []any{query.In("elem.reply_id", replyIds...)}}))
	if err != nil {
		return errors.Wrapf(err, "failed to update reply approval_status, commentId=%s, replyIds=%v", commentObjectID, replyIds)
	}
	if updateResult.ModifiedCount == 0 {
		return fmt.Errorf("UpsertedCount=0, failed to update reply approval_status, commentId=%s, replyIds=%v", commentObjectID, replyIds)
	}
	return nil
}

func (d *CommentDao) FindByObjectIDs(ctx context.Context, ids []primitive.ObjectID) ([]*Comment, error) {
	return d.coll.Finder().Filter(query.In("_id", ids...)).Find(ctx)
}

func (d *CommentDao) UpdateCommentStatus2TrueByIds(ctx context.Context, ids []primitive.ObjectID) error {
	anyIds := make([]any, 0, len(ids))
	for _, objectID := range ids {
		anyIds = append(anyIds, objectID)
	}
	updateResult, err := d.coll.Updater().Filter(query.BsonBuilder().In("_id", anyIds...).Eq("approval_status", false).Build()).Updates(update.BsonBuilder().Set("approval_status", true).Set("updated_at", time.Now().Local()).Build()).UpdateMany(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to update approval_status of comments, ids=%v", ids)
	}
	if updateResult.ModifiedCount == 0 {
		return fmt.Errorf("ModifiedCount=0, failed to update approval_status of comments, ids=%v", ids)
	}
	return nil
}

func (d *CommentDao) Find(ctx context.Context, findOptions *options.FindOptions) ([]*Comment, int64, error) {
	count, err := d.coll.Finder().Filter(bson.D{}).Count(ctx)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "failed to count from comment, cond={}, findOptions=%v", findOptions)
	}
	comments, err := d.coll.Finder().Filter(bson.D{}).Find(ctx, findOptions)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "failed to find from comment, cond={}, findOptions=%v", findOptions)
	}
	return comments, count, nil
}

func (d *CommentDao) CountOfToday(ctx context.Context) (int64, error) {
	startOfDayUnix, endOfDayUnix := d.getBeginSecondsAndEndSeconds()
	return d.coll.Finder().Filter(query.BsonBuilder().Gte("created_at", startOfDayUnix).Lte("created_at", endOfDayUnix).Build()).Count(ctx)
}

func (d *CommentDao) DeleteReplyByCIdAndRId(ctx context.Context, objectID primitive.ObjectID, replyId string) error {
	result, err := d.coll.Updater().Filter(query.Id(objectID)).Updates(update.Pull("replies", bsonx.M("reply_id", replyId))).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to update one from comment, commentId=%s, replyId=%s", objectID.Hex(), replyId)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("modifiedCount = 0, fails to update one from comment, commentId=%s, replyId=%s", objectID.Hex(), replyId)
	}
	return nil
}

func (d *CommentDao) DeleteById(ctx context.Context, objectID primitive.ObjectID) error {
	result, err := d.coll.Deleter().Filter(query.Id(objectID)).DeleteOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to delete one from comment, id=%s", objectID.Hex())
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("deletedCount = 0, fails to delete one from comment, id=%s", objectID.Hex())
	}
	return nil
}

func (d *CommentDao) FindCommentWithRepliesById(ctx context.Context, objectID primitive.ObjectID) (*Comment, error) {
	comment, err := d.coll.Finder().Filter(query.Id(objectID)).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to find the document from comment, id=%s", objectID.Hex())
	}
	return comment, nil
}

func (d *CommentDao) UpdateCommentReplyStatus(ctx context.Context, objectID primitive.ObjectID, replyId string, commentStatus bool) error {
	filter := query.BsonBuilder().Id(objectID).Add("replies.reply_id", replyId).Build()
	updates := update.BsonBuilder().Set("replies.$.approval_status", commentStatus).Set("replies.$.update_time", time.Now().Local().Unix()).Build()
	result, err := d.coll.Updater().Filter(filter).Updates(updates).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to update one from comment, filter=%v, update=%v", filter, updates)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("modifiedCount = 0, fails to update one from comment, filter=%v, update=%v", filter, updates)
	}
	return nil
}

func (d *CommentDao) FindReplyByCIdAndRId(ctx context.Context, objectCommentID primitive.ObjectID, replyId string) (*ReplyWithPostInfo, error) {
	pipeline := aggregation.StageBsonBuilder().
		Match(query.Id(objectCommentID)).
		Unwind("$replies", nil).
		Match(bsonx.M("replies.reply_id", replyId)).
		AddFields(bsonx.M("replies.post_info", "$post_info")).
		ReplaceWith("$replies").Build()
	var result []ReplyWithPostInfo
	err := d.coll.Aggregator().Pipeline(pipeline).AggregateWithParse(ctx, &result)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to execute aggregation operation, pipeline=%v", pipeline)
	}
	if len(result) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return &result[0], nil
}

func (d *CommentDao) UpdateCommentStatus2True(ctx context.Context, objectID primitive.ObjectID) error {
	result, err := d.coll.Updater().Filter(query.Id(objectID)).Updates(update.BsonBuilder().Set("approval_status", true).Set("updated_at", time.Now().Local()).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to update one from comment, id=%s, commentStatus=%v", objectID.Hex(), true)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("modifiedCount = 0, fails to update one from comment, id=%s, commentStatus=%v", objectID.Hex(), true)
	}
	return nil
}

func (d *CommentDao) FindCommentById(ctx context.Context, objectID primitive.ObjectID) (*Comment, error) {
	comment, err := d.coll.Finder().Filter(query.Id(objectID)).FindOne(ctx, options.FindOne().SetProjection(bsonx.M("replies", 0)))
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (d *CommentDao) AggregationQuerySkipAndSetLimit(ctx context.Context, cond bson.D, findOptions *options.FindOptions) ([]AdminComment, int64, error) {
	pipeline := aggregation.StageBsonBuilder().
		Match(bson.D{}).
		Project(aggregation.ConcatArrays("combined", []any{
			bsonx.A(bsonx.NewD().Add("_id", "$_id").Add("post_info", "$post_info").Add("user_info", "$user_info").Add("content", "$content").Add("created_at", "$created_at").Add("type", 0).Add("approval_status", "approval_status").Build()),
			aggregation.MapWithoutKey(
				"$replies",
				"reply",
				bsonx.NewD().Add("_id", "$$reply.reply_id").Add("fid", "$_id").Add("post_info", "$post_info").Add("user_info", "$$reply.user_info").Add("content", "$$reply.content").Add("created_at", "$$reply.created_at").Add("type", 1).Add("approval_status", "$$reply.approval_status").Build(),
			),
		}...)).
		Unwind("$combined", nil).
		ReplaceWith("$combined").
		Sort(bsonx.M("created_at", -1)).Build()

	var results []AdminComment

	// 执行聚合查询
	err := d.coll.Aggregator().Pipeline(pipeline).AggregateWithParse(ctx, &results)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "Fails to execute aggregation operation, pipeline=%v", pipeline)
	}

	return results, 0, nil
}

func (d *CommentDao) FindApprovedCommentsByPostId(ctx context.Context, postId string) ([]*Comment, error) {
	cond := bsonx.NewD().Add("post_info.post_id", postId).Add("approval_status", true).Build()
	result, err := d.coll.Finder().Filter(cond).Find(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Fails to find the docs from comment, condition=%v", cond)
	}
	return result, nil
}

func (d *CommentDao) FineLatestCommentAndReply(ctx context.Context, cnt int) ([]LatestComment, error) {
	pipeline := aggregation.StageBsonBuilder().
		Match(bsonx.M("approval_status", true)).
		Project(aggregation.ConcatArrays("combined", []any{
			bsonx.A(bsonx.NewD().Add("post_info", "$post_info").Add("name", "$user_info.name").Add("email", "$user_info.email").Add("content", "$content").Add("created_at", "$created_at").Build()),
			aggregation.MapWithoutKey(
				aggregation.FilterWithoutKey("$replies", aggregation.EqWithoutKey("$$replyItem.approval_status", true), &types.FilterOptions{As: "replyItem"}),
				"reply",
				bsonx.NewD().Add("post_info", "$post_info").Add("name", "$$reply.user_info.name").Add("email", "$$reply.user_info.email").Add("content", "$$reply.content").Add("created_at", "$$reply.created_at").Build(),
			),
		}...)).
		Unwind("$combined", nil).
		ReplaceWith("$combined").
		Sort(bsonx.M("created_at", -1)).
		Limit(int64(cnt)).Build()
	//pipeline := mongo.Pipeline{
	//	{primitive.E{Key: "$match", Value: bson.M{"approval_status": CommentStatusApproved}}},
	//	{primitive.E{Key: "$project", Value: bson.M{
	//		"combined": bson.M{
	//			"$concatArrays": []any{
	//				bson.A{bson.M{"post_info": "$post_info", "name": "$user_info.name", "content": "$content", "created_at": "$created_at", "email": "$user_info.email"}},
	//				//"$replies",
	//				bson.M{
	//					"$map": bson.M{
	//						"input": bson.M{
	//							"$filter": bson.M{
	//								"input": "$replies",
	//								"as":    "replyItem",
	//								"cond":  bson.M{"$eq": []interface{}{"$$replyItem.approval_status", CommentStatusApproved}},
	//							},
	//						},
	//						"as": "reply",
	//						"in": bson.M{
	//							"post_info":   "$post_info",
	//							"name":        "$$reply.user_info.name",
	//							"content":     "$$reply.content",
	//							"created_at": "$$reply.created_at",
	// 							"email": "$$reply.user_info.email",
	//						},
	//					},
	//				},
	//			},
	//		},
	//	}}},
	//	{primitive.E{Key: "$unwind", Value: "$combined"}},
	//	{primitive.E{Key: "$replaceRoot", Value: bson.M{"newRoot": "$combined"}}},
	//	{primitive.E{Key: "$sort", Value: bson.M{"created_at": -1}}},
	//	{primitive.E{Key: "$limit", Value: cnt}},
	//}

	var results []LatestComment

	// 执行聚合查询
	err := d.coll.Aggregator().Pipeline(pipeline).AggregateWithParse(ctx, &results)
	if err != nil {
		return nil, errors.Wrapf(err, "Fails to execute aggregation operation, pipeline=%v", pipeline)
	}

	return results, nil
}

func (d *CommentDao) AddCommentReply(ctx context.Context, objectID primitive.ObjectID, commentReply Reply) error {
	// 构建查询条件
	filter := query.Id(objectID)
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

func (d *CommentDao) FindApprovedCommentById(ctx context.Context, objectID primitive.ObjectID) (*Comment, error) {
	comment, err := d.coll.Finder().Filter(
		query.BsonBuilder().Id(objectID).Eq("approval_status", true).Build()).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to find the document from comment, cmtId=%s", objectID.Hex())
	}
	return comment, nil
}

func (d *CommentDao) AddComment(ctx context.Context, comment *Comment) (string, error) {
	result, err := d.coll.Creator().InsertOne(ctx, comment)
	if err != nil {
		return "", errors.Wrapf(err, "fails to insert into comment, comment=%v", comment)
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (d *CommentDao) getBeginSecondsAndEndSeconds() (int64, int64) {
	now := time.Now().Local()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0,
		now.Location()).Unix(), time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()).Unix()
}
