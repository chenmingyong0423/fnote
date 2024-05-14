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
	"encoding/hex"
	"log/slog"
	"net/http"
	"os"

	jsoniter "github.com/json-iterator/go"

	"github.com/google/uuid"

	"github.com/chenmingyong0423/go-eventbus"

	"github.com/chenmingyong0423/fnote/server/internal/file/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/file/internal/repository"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/chenmingyong0423/gkit/uuidx"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/spf13/viper"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"
)

type IFileService interface {
	Upload(ctx context.Context, fileDTO dto.FileDTO) (*domain.File, error)
	IndexFileMeta(ctx context.Context, fileId []byte, entityId string, entityType string) error
	DeleteIndexFileMeta(ctx context.Context, fileId []byte, entityId string, entityType string) error
}

var _ IFileService = (*FileService)(nil)

func NewFileService(repo repository.IFileRepository, eventbus *eventbus.EventBus) *FileService {
	s := &FileService{
		repo:     repo,
		eventBus: eventbus,
	}
	go s.subscribePostEvent()
	return s
}

type FileService struct {
	repo     repository.IFileRepository
	eventBus *eventbus.EventBus
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

func (s *FileService) subscribePostEvent() {
	eventChan := s.eventBus.Subscribe("post")
	type contextKey string
	for event := range eventChan {
		rid := uuid.NewString()
		var key contextKey = "X-Request-ID"
		ctx := context.WithValue(context.Background(), key, rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "File: post event", "payload", string(event.Payload))
		var e domain.PostEvent
		err := jsoniter.Unmarshal(event.Payload, &e)
		if err != nil {
			l.ErrorContext(ctx, "File: post event: failed to json.Unmarshal", "err", err)
			continue
		}
		switch e.Type {
		case "create":
			s.createIndexFileMeta4PostEvent(ctx, e.NewFileId, e.PostId, l)
		case "update":
			if e.NewFileId != e.OldFileId {
				s.createIndexFileMeta4PostEvent(ctx, e.NewFileId, e.PostId, l)
				s.deleteIndexFileMeta4PostEvent(ctx, e.OldFileId, e.PostId, l)
			}
		case "delete":
			s.deleteIndexFileMeta4PostEvent(ctx, e.OldFileId, e.PostId, l)
		}
		l.InfoContext(ctx, "File: post event: handle successfully")
	}
}

func (s *FileService) deleteIndexFileMeta4PostEvent(ctx context.Context, oldFileId string, postId string, l *slog.Logger) {
	fid, sErr := hex.DecodeString(oldFileId)
	if sErr != nil {
		l.ErrorContext(ctx, "File: post event: failed to hex.DecodeString", "fileId", oldFileId, "err", sErr)
		return
	}
	sErr = s.DeleteIndexFileMeta(ctx, fid, postId, "post")
	if sErr != nil {
		l.ErrorContext(ctx, "File: post event: failed to delete the index of file-meta ", "fileId", oldFileId, "postId", postId, "err", sErr)
	}
}

func (s *FileService) createIndexFileMeta4PostEvent(ctx context.Context, newFileId string, postId string, l *slog.Logger) {
	fid, sErr := hex.DecodeString(newFileId)
	if sErr != nil {
		l.ErrorContext(ctx, "File: post event: failed to hex.DecodeString", "fileId", newFileId, "err", sErr)
		return
	}
	sErr = s.IndexFileMeta(ctx, fid, postId, "post")
	if sErr != nil {
		l.ErrorContext(ctx, "File: post event: failed to index the file-meta ", "fileId", newFileId, "postId", postId, "err", sErr)
	}
}
