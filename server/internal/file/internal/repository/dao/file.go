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

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/chenmingyong0423/go-mongox/v2/builder/query"

	"github.com/pkg/errors"

	"github.com/chenmingyong0423/go-mongox/v2/bsonx"

	"github.com/chenmingyong0423/go-mongox/v2/builder/update"

	"github.com/chenmingyong0423/go-mongox/v2"
)

type File struct {
	mongox.Model     `bson:",inline"`
	FileId           []byte      `bson:"file_id"`
	FileName         string      `bson:"file_name"`
	OriginalFileName string      `bson:"original_file_name"`
	FileType         string      `bson:"file_type"`
	FileSize         int64       `bson:"size"`
	FilePath         string      `bson:"file_path"`
	Url              string      `bson:"url"`
	UsedIn           []FileUsage `bson:"used_in"`
}

type FileUsage struct {
	EntityId   string     `bson:"entity_id"`
	EntityType EntityType `bson:"entity_type"`
}

type EntityType string

const (
	EntityTypePost      EntityType = "post"
	EntityTypePostDraft EntityType = "post-draft"
	EntityTypeAsset     EntityType = "asset"
	EntityTypeConfig    EntityType = "config"
)

type UsageContent struct {
	ID       string `bson:"_id"`
	Title    string `bson:"title"`
	Content  string `bson:"content"`
	CoverImg string `bson:"cover_img"`
}

type UsageAsset struct {
	ID    bson.ObjectID `bson:"_id"`
	Title string        `bson:"title"`
	Type  string        `bson:"type"`
}

type UsageConfig struct {
	ID    bson.ObjectID    `bson:"_id"`
	Typ   string           `bson:"typ"`
	Props UsageConfigProps `bson:"props"`
}

type UsageConfigProps struct {
	WebsiteIcon        string            `bson:"website_icon"`
	WebsiteOwnerAvatar string            `bson:"website_owner_avatar"`
	WebsiteRecords     []string          `bson:"website_records"`
	OgImage            string            `bson:"og_image"`
	Content            string            `bson:"content"`
	Introduction       string            `bson:"introduction"`
	List               []UsageConfigItem `bson:"list"`
	SocialInfoList     []UsageSocialInfo `bson:"social_info_list"`
}

type UsageConfigItem struct {
	Image    string `bson:"image"`
	CoverImg string `bson:"cover_img"`
}

type UsageSocialInfo struct {
	SocialValue string `bson:"social_value"`
}

type IFileDao interface {
	Save(ctx context.Context, file *File) (string, error)
	FindByFileId(ctx context.Context, fileId []byte) (*File, error)
	PushIntoUsedIn(ctx context.Context, fileId []byte, fileUsage FileUsage) error
	PullUsedIn(ctx context.Context, fileId []byte, fileUsage FileUsage) error
	FindByFileName(ctx context.Context, filename string) (*File, error)
	FindPageByFileType(ctx context.Context, pageNum int64, pageSize int64, fileType []string, unused bool) ([]*File, int64, error)
	FindUsagePosts(ctx context.Context, ids []string) ([]*UsageContent, error)
	FindUsagePostDrafts(ctx context.Context, ids []string) ([]*UsageContent, error)
	FindUsageConfigs(ctx context.Context, ids []bson.ObjectID) ([]*UsageConfig, error)
	FindUsageAssets(ctx context.Context, ids []bson.ObjectID) ([]*UsageAsset, error)
	DeleteUnusedByFileId(ctx context.Context, fileId []byte) (int64, error)
}

var _ IFileDao = (*FileDao)(nil)

func NewFileDao(db *mongox.Database) *FileDao {
	return &FileDao{
		coll:          mongox.NewCollection[File](db, "file_meta"),
		postColl:      mongox.NewCollection[UsageContent](db, "posts"),
		postDraftColl: mongox.NewCollection[UsageContent](db, "post_draft"),
		configColl:    mongox.NewCollection[UsageConfig](db, "configs"),
		assetColl:     mongox.NewCollection[UsageAsset](db, "assets"),
	}
}

type FileDao struct {
	coll          *mongox.Collection[File]
	postColl      *mongox.Collection[UsageContent]
	postDraftColl *mongox.Collection[UsageContent]
	configColl    *mongox.Collection[UsageConfig]
	assetColl     *mongox.Collection[UsageAsset]
}

func (d *FileDao) FindUsagePosts(ctx context.Context, ids []string) ([]*UsageContent, error) {
	if len(ids) == 0 {
		return []*UsageContent{}, nil
	}
	return d.postColl.Finder().Filter(query.In("_id", ids...)).Find(ctx)
}

func (d *FileDao) FindUsagePostDrafts(ctx context.Context, ids []string) ([]*UsageContent, error) {
	if len(ids) == 0 {
		return []*UsageContent{}, nil
	}
	return d.postDraftColl.Finder().Filter(query.In("_id", ids...)).Find(ctx)
}

func (d *FileDao) FindUsageConfigs(ctx context.Context, ids []bson.ObjectID) ([]*UsageConfig, error) {
	if len(ids) == 0 {
		return []*UsageConfig{}, nil
	}
	return d.configColl.Finder().Filter(query.In("_id", ids...)).Find(ctx)
}

func (d *FileDao) FindUsageAssets(ctx context.Context, ids []bson.ObjectID) ([]*UsageAsset, error) {
	if len(ids) == 0 {
		return []*UsageAsset{}, nil
	}
	return d.assetColl.Finder().Filter(query.In("_id", ids...)).Find(ctx)
}

func (d *FileDao) FindByFileId(ctx context.Context, fileId []byte) (*File, error) {
	return d.coll.Finder().Filter(query.Eq("file_id", fileId)).FindOne(ctx)
}

func (d *FileDao) DeleteUnusedByFileId(ctx context.Context, fileId []byte) (int64, error) {
	result, err := d.coll.Deleter().Filter(bson.M{
		"file_id": fileId,
		"$or": bson.A{
			bson.M{"used_in": bson.M{"$exists": false}},
			bson.M{"used_in": bson.M{"$size": 0}},
		},
	}).DeleteOne(ctx)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (d *FileDao) FindPageByFileType(ctx context.Context, pageNum int64, pageSize int64, fileType []string, unused bool) ([]*File, int64, error) {
	filter := bson.D{}
	if len(fileType) > 0 {
		filter = append(filter, query.In("file_type", fileType...)...)
	}
	if unused {
		filter = append(filter, bson.E{Key: "$or", Value: bson.A{
			bson.M{"used_in": bson.M{"$exists": false}},
			bson.M{"used_in": bson.M{"$size": 0}},
		}})
	}
	count, err := d.coll.Finder().Filter(filter).Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	findOptions := options.Find().
		SetSkip((pageNum - 1) * pageSize).
		SetLimit(pageSize).
		SetSort(bsonx.M("created_at", -1))
	files, err := d.coll.Finder().Filter(filter).Find(ctx, findOptions)
	if err != nil {
		return nil, 0, err
	}
	return files, count, nil
}

func (d *FileDao) FindByFileName(ctx context.Context, filename string) (*File, error) {
	return d.coll.Finder().Filter(query.Eq("file_name", filename)).FindOne(ctx)
}

func (d *FileDao) PullUsedIn(ctx context.Context, fileId []byte, fileUsage FileUsage) error {
	updateOne, err := d.coll.Updater().Filter(bsonx.M("file_id", fileId)).Updates(update.NewBuilder().Pull("used_in", fileUsage).Set("updated_at", time.Now().Local()).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "pull used in error, file id: %s, file usage: %+v", fileId, fileUsage)
	}
	if updateOne.MatchedCount == 0 {
		return fmt.Errorf("pull used in error, file id: %s, file usage: %+v", fileId, fileUsage)
	}
	return nil
}

func (d *FileDao) PushIntoUsedIn(ctx context.Context, fileId []byte, fileUsage FileUsage) error {
	updateOne, err := d.coll.Updater().Filter(bsonx.M("file_id", fileId)).Updates(update.NewBuilder().AddToSet("used_in", fileUsage).Set("updated_at", time.Now().Local()).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "push into used in error, file id: %s, file usage: %+v", fileId, fileUsage)
	}
	if updateOne.MatchedCount == 0 {
		return fmt.Errorf("push into used in error, file id: %s, file usage: %+v", fileId, fileUsage)
	}
	return nil
}

func (d *FileDao) Save(ctx context.Context, file *File) (string, error) {
	oneResult, err := d.coll.Creator().InsertOne(ctx, file)
	if err != nil {
		return "", err
	}
	return oneResult.InsertedID.(bson.ObjectID).Hex(), nil
}
