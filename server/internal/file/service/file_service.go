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

package service

import (
	"context"
	"net/http"
	"os"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/chenmingyong0423/gkit/uuidx"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/spf13/viper"

	"github.com/chenmingyong0423/fnote/server/internal/file/repository"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"
)

type IFileService interface {
	Upload(ctx context.Context, fileDTO dto.FileDTO) (*domain.File, error)
	IndexFileMeta(ctx context.Context, fileId []byte, entityId string, entityType string) error
	DeleteIndexFileMeta(ctx context.Context, fileId []byte, entityId string, entityType string) error
}

var _ IFileService = (*FileService)(nil)

func NewFileService(repo repository.IFileRepository) *FileService {
	return &FileService{
		repo: repo,
	}
}

type FileService struct {
	repo repository.IFileRepository
}

func (s *FileService) DeleteIndexFileMeta(ctx context.Context, fileId []byte, entityId string, entityType string) error {
	return s.repo.PullUsedIn(ctx, fileId, entityId, entityType)
}

func (s *FileService) IndexFileMeta(ctx context.Context, fileId []byte, entityId string, entityType string) error {
	return s.repo.PushIntoUsedIn(ctx, fileId, entityId, entityType)
}

func (s *FileService) Upload(ctx context.Context, fileDTO dto.FileDTO) (*domain.File, error) {
	var (
		filename string
	)
	fileId := uuidx.RearrangeUUID4()
	if fileDTO.CustomFileName != "" {
		filename = fileDTO.CustomFileName + fileDTO.FileExt
		file, err := s.repo.FindByFileName(ctx, filename)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		if file != nil {
			return nil, apiwrap.NewErrorResponseBody(http.StatusConflict, "file already exists")
		}
	} else {
		filename = fileId + fileDTO.FileExt
	}

	staticPath := viper.GetString("system.static_path")
	err := os.MkdirAll(staticPath, os.ModePerm)
	if err != nil {
		return nil, err
	}
	create, err := os.Create(staticPath + filename)
	if err != nil {
		return nil, err
	}
	defer create.Close()
	_, err = create.Write(fileDTO.Content)
	if err != nil {
		return nil, err
	}
	file := &domain.File{
		FileId:           fileId,
		FileName:         filename,
		OriginalFileName: fileDTO.FileName,
		FileType:         fileDTO.FileType,
		FileSize:         fileDTO.FileSize,
		FilePath:         staticPath + filename,
		Url:              "/static/" + filename,
	}
	err = s.repo.Save(ctx, file)
	if err != nil {
		return nil, err
	}
	return file, nil
}
