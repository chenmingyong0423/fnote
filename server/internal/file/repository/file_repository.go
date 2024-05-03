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

package repository

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/file/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IFileRepository interface {
	Save(ctx context.Context, file *domain.File) error
	PushIntoUsedIn(ctx context.Context, fileId []byte, entityId string, entityType string) error
	PullUsedIn(ctx context.Context, fileId []byte, entityId string, entityType string) error
	FindByFileName(ctx context.Context, filename string) (*domain.File, error)
}

var _ IFileRepository = (*FileRepository)(nil)

func NewFileRepository(dao dao.IFileDao) *FileRepository {
	return &FileRepository{dao: dao}
}

type FileRepository struct {
	dao dao.IFileDao
}

func (r *FileRepository) FindByFileName(ctx context.Context, filename string) (*domain.File, error) {
	file, err := r.dao.FindByFileName(ctx, filename)
	if err != nil {
		return nil, err
	}
	return r.toDomainFile(file), nil
}

func (r *FileRepository) PullUsedIn(ctx context.Context, fileId []byte, entityId string, entityType string) error {
	return r.dao.PullUsedIn(ctx, fileId, dao.FileUsage{
		EntityId:   entityId,
		EntityType: dao.EntityType(entityType),
	})
}

func (r *FileRepository) PushIntoUsedIn(ctx context.Context, fileId []byte, entityId string, entityType string) error {
	return r.dao.PushIntoUsedIn(ctx, fileId, dao.FileUsage{
		EntityId:   entityId,
		EntityType: dao.EntityType(entityType),
	})
}

func (r *FileRepository) Save(ctx context.Context, file *domain.File) error {
	unix := time.Now().Local().Unix()
	fileId, err := hex.DecodeString(file.FileId)
	if err != nil {
		return err
	}
	_, err = r.dao.Save(ctx, &dao.File{
		Id:               primitive.ObjectID{},
		FileId:           fileId,
		FileName:         file.FileName,
		OriginalFileName: file.OriginalFileName,
		FileType:         file.FileType,
		FileSize:         file.FileSize,
		FilePath:         file.FilePath,
		Url:              file.Url,
		UsedIn:           make([]dao.FileUsage, 0),
		CreateTime:       unix,
		UpdateTime:       unix,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *FileRepository) toDomainFile(file *dao.File) *domain.File {
	return &domain.File{
		Id:               file.Id.Hex(),
		FileId:           hex.EncodeToString(file.FileId),
		FileName:         file.FileName,
		OriginalFileName: file.OriginalFileName,
		FileType:         file.FileType,
		FileSize:         file.FileSize,
		FilePath:         file.FilePath,
		Url:              file.Url,
		UsedIn: func() []domain.FileUsage {
			usedIn := make([]domain.FileUsage, 0)
			for _, usage := range file.UsedIn {
				usedIn = append(usedIn, domain.FileUsage{
					EntityId:   usage.EntityId,
					EntityType: string(usage.EntityType),
				})
			}
			return usedIn
		}(),
		CreateTime: file.CreateTime,
		UpdateTime: file.UpdateTime,
	}
}
