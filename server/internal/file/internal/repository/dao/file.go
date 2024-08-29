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
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/chenmingyong0423/go-mongox/builder/query"

	"github.com/pkg/errors"

	"github.com/chenmingyong0423/go-mongox/bsonx"

	"github.com/chenmingyong0423/go-mongox/builder/update"

	"github.com/chenmingyong0423/go-mongox"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	EntityTypePost EntityType = "post"
)

type IFileDao interface {
	Save(ctx context.Context, file *File) (string, error)
	PushIntoUsedIn(ctx context.Context, fileId []byte, fileUsage FileUsage) error
	PullUsedIn(ctx context.Context, fileId []byte, fileUsage FileUsage) error
	FindByFileName(ctx context.Context, filename string) (*File, error)
	FindPageByFileType(ctx context.Context, pageNum int64, pageSize int64, fileType []string) ([]*File, int64, error)
}

var _ IFileDao = (*FileDao)(nil)

func NewFileDao(db *mongo.Database) *FileDao {
	return &FileDao{coll: mongox.NewCollection[File](db.Collection("file_meta"))}
}

type FileDao struct {
	coll *mongox.Collection[File]
}

func (d *FileDao) FindPageByFileType(ctx context.Context, pageNum int64, pageSize int64, fileType []string) ([]*File, int64, error) {
	filter := query.In("file_type", fileType...)
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
	if updateOne.ModifiedCount == 0 {
		return fmt.Errorf("pull used in error, file id: %s, file usage: %+v", fileId, fileUsage)
	}
	return nil
}

func (d *FileDao) PushIntoUsedIn(ctx context.Context, fileId []byte, fileUsage FileUsage) error {
	updateOne, err := d.coll.Updater().Filter(bsonx.M("file_id", fileId)).Updates(update.NewBuilder().Push("used_in", fileUsage).Set("updated_at", time.Now().Local()).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "push into used in error, file id: %s, file usage: %+v", fileId, fileUsage)
	}
	if updateOne.ModifiedCount == 0 {
		return fmt.Errorf("push into used in error, file id: %s, file usage: %+v", fileId, fileUsage)
	}
	return nil
}

func (d *FileDao) Save(ctx context.Context, file *File) (string, error) {
	oneResult, err := d.coll.Creator().InsertOne(ctx, file)
	if err != nil {
		return "", err
	}
	return oneResult.InsertedID.(primitive.ObjectID).Hex(), nil
}
