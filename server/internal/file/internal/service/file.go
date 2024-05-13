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
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

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
	go s.SubscribePostDeletedEvent()
	go s.SubscribePostAddedEvent()
	go s.SubscribePostUpdatedEvent()
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

func (s *FileService) SubscribePostDeletedEvent() {
	eventChan := s.eventBus.Subscribe("post-delete")
	for event := range eventChan {
		rid := uuid.NewString()
		ctx := context.WithValue(context.Background(), "X-Request-ID", rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "post-delete", "payload", string(event.Payload))
		var postEvent domain.PostEvent
		err := json.Unmarshal(event.Payload, &postEvent)
		if err != nil {
			l.ErrorContext(ctx, "post-delete: failed to json.Unmarshal", "err", err)
			continue
		}
		fid, err := hex.DecodeString(postEvent.FileId)
		if err != nil {
			l.ErrorContext(ctx, "post-delete: failed to hex.DecodeString", "fileId", postEvent.FileId, "err", err)
			continue
		}
		err = s.DeleteIndexFileMeta(ctx, fid, postEvent.PostId, "post")
		if err != nil {
			l.ErrorContext(ctx, "post-delete: failed to delete the file-meta index", "fileId", postEvent.FileId, "postId", postEvent.PostId, "err", err)
			continue
		}
		l.InfoContext(ctx, "post-delete: successfully delete the file-meta index", "fileId", postEvent.FileId, "postId", postEvent.PostId)
	}
}

func (s *FileService) SubscribePostAddedEvent() {
	eventChan := s.eventBus.Subscribe("post-addition")
	for event := range eventChan {
		rid := uuid.NewString()
		ctx := context.WithValue(context.Background(), "X-Request-ID", rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "post-addition", "payload", string(event.Payload))
		var postEvent domain.PostEvent
		err := json.Unmarshal(event.Payload, &postEvent)
		if err != nil {
			l.ErrorContext(ctx, "post-addition: failed to json.Unmarshal", "err", err)
			continue
		}
		fid, err := hex.DecodeString(postEvent.FileId)
		if err != nil {
			l.ErrorContext(ctx, "post-addition: failed to hex.DecodeString", "fileId", postEvent.FileId, "err", err)
			continue
		}
		err = s.IndexFileMeta(ctx, fid, postEvent.PostId, "post")
		if err != nil {
			l.ErrorContext(ctx, "post-addition: failed to index the file-meta ", "fileId", postEvent.FileId, "postId", postEvent.PostId, "err", err)
			continue
		}
		l.InfoContext(ctx, "post-addition: successfully index the file-meta ", "fileId", postEvent.FileId, "postId", postEvent.PostId)
	}

}

func (s *FileService) SubscribePostUpdatedEvent() {
	eventChan := s.eventBus.Subscribe("post-update")
	for event := range eventChan {
		rid := uuid.NewString()
		ctx := context.WithValue(context.Background(), "X-Request-ID", rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "post-update", "payload", string(event.Payload))
		var postEvent domain.UpdatedPostEvent
		err := json.Unmarshal(event.Payload, &postEvent)
		if err != nil {
			l.ErrorContext(ctx, "post-update: failed to json.Unmarshal", "err", err)
			continue
		}
		if postEvent.NewFileId != postEvent.OldFileId {
			oldFid, err2 := hex.DecodeString(postEvent.OldFileId)
			if err != nil {
				l.ErrorContext(ctx, "post-update: failed to hex.DecodeString", "fileId", postEvent.OldFileId, "err", err2)
				continue
			}
			err2 = s.DeleteIndexFileMeta(ctx, oldFid, postEvent.PostId, "post")
			if err2 != nil {
				l.ErrorContext(ctx, "post-update: failed to delete the file-meta index", "fileId", postEvent.OldFileId, "postId", postEvent.PostId, "err", err2)
				continue
			}
			newFid, err2 := hex.DecodeString(postEvent.NewFileId)
			if err2 != nil {
				l.ErrorContext(ctx, "post-update: failed to hex.DecodeString", "fileId", postEvent.NewFileId, "err", err2)
				continue
			}
			err2 = s.IndexFileMeta(ctx, newFid, postEvent.PostId, "post")
			if err2 != nil {
				l.ErrorContext(ctx, "post-update: failed to index the file-meta ", "fileId", postEvent.NewFileId, "postId", postEvent.PostId, "err", err2)
				continue
			}
			l.InfoContext(ctx, "post-update: successfully update the index of file-meta ", "newFileId", postEvent.NewFileId, "oldFileId", postEvent.OldFileId, "postId", postEvent.PostId)
		} else {
			l.InfoContext(ctx, "post-update: file not changed", "fileId", postEvent.OldFileId, "postId", postEvent.PostId)
		}
	}
}
