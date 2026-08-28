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
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/pkg"
	"github.com/chenmingyong0423/fnote/server/internal/tag"

	"github.com/chenmingyong0423/fnote/server/internal/category"

	"github.com/chenmingyong0423/fnote/server/internal/post"

	"github.com/chenmingyong0423/go-sitemap-generator"

	jsoniter "github.com/json-iterator/go"

	"github.com/google/uuid"

	"github.com/chenmingyong0423/go-eventbus"

	"github.com/chenmingyong0423/fnote/server/internal/file/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/file/internal/repository"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/chenmingyong0423/gkit/uuidx"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/spf13/viper"
)

type IFileService interface {
	Upload(ctx context.Context, fileDTO domain.FileDTO) (*domain.File, error)
	IndexFileMeta(ctx context.Context, fileId []byte, entityId string, entityType string) error
	DeleteIndexFileMeta(ctx context.Context, fileId []byte, entityId string, entityType string) error
	HasOtherUsages(ctx context.Context, fileId []byte, entityId string, entityType string) (bool, error)
	GenerateSitemap(ctx context.Context, postBytes, categoryBytes, tagBytes []byte) error
	GetSitemap(ctx context.Context) (string, bool, error)
	GetRobotsTxt(ctx context.Context) (string, bool, error)
	SaveRobotsTxt(ctx context.Context, content string) error
	GetFiles(ctx context.Context, pageDTO domain.PageDTO) ([]*domain.File, int64, error)
	DeleteFile(ctx context.Context, fileId string) error
	DeleteFiles(ctx context.Context, fileIds []string) error
}

var _ IFileService = (*FileService)(nil)

func NewFileService(repo repository.IFileRepository, eventbus *eventbus.EventBus) *FileService {
	s := &FileService{
		repo:     repo,
		eventBus: eventbus,
	}
	go s.subscribePostEvent()
	go s.subscribePostDraftEvent()
	return s
}

type FileService struct {
	repo     repository.IFileRepository
	eventBus *eventbus.EventBus
}

func (s *FileService) DeleteFile(ctx context.Context, fileId string) error {
	pending, err := s.prepareFileDeletion(ctx, fileId)
	if err != nil {
		return err
	}
	return s.deletePreparedFile(ctx, pending)
}

type pendingFileDeletion struct {
	fileID       []byte
	physicalPath string
}

func (s *FileService) DeleteFiles(ctx context.Context, fileIds []string) error {
	if len(fileIds) == 0 || len(fileIds) > 100 {
		return apiwrap.NewErrorResponseBody(http.StatusBadRequest, "file ids must contain between 1 and 100 items")
	}
	pendingFiles := make([]pendingFileDeletion, 0, len(fileIds))
	seen := make(map[string]struct{}, len(fileIds))
	for _, fileID := range fileIds {
		if _, exists := seen[fileID]; exists {
			continue
		}
		seen[fileID] = struct{}{}
		pending, err := s.prepareFileDeletion(ctx, fileID)
		if err != nil {
			return err
		}
		pendingFiles = append(pendingFiles, pending)
	}
	for _, pending := range pendingFiles {
		if err := s.deletePreparedFile(ctx, pending); err != nil {
			return err
		}
	}
	return nil
}

func (s *FileService) prepareFileDeletion(ctx context.Context, fileId string) (pendingFileDeletion, error) {
	decodedFileID, err := hex.DecodeString(fileId)
	if err != nil || len(decodedFileID) != 16 {
		return pendingFileDeletion{}, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "invalid file id")
	}
	file, err := s.repo.FindByFileId(ctx, decodedFileID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return pendingFileDeletion{}, apiwrap.NewErrorResponseBody(http.StatusNotFound, "file not found")
		}
		return pendingFileDeletion{}, err
	}
	if len(file.UsedIn) > 0 {
		return pendingFileDeletion{}, apiwrap.NewErrorResponseBody(http.StatusConflict, "file is still in use")
	}
	physicalPath, ok := resolveStaticFilePath(file.FileName)
	if !ok {
		return pendingFileDeletion{}, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "invalid file path")
	}
	return pendingFileDeletion{fileID: decodedFileID, physicalPath: physicalPath}, nil
}

func (s *FileService) deletePreparedFile(ctx context.Context, pending pendingFileDeletion) error {
	deletedCount, err := s.repo.DeleteUnusedByFileId(ctx, pending.fileID)
	if err != nil {
		return err
	}
	if deletedCount == 0 {
		return apiwrap.NewErrorResponseBody(http.StatusConflict, "file is still in use")
	}
	if err = os.Remove(pending.physicalPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return errors.Wrap(err, "file metadata was deleted but physical file removal failed")
	}
	return nil
}

func resolveStaticFilePath(fileName string) (string, bool) {
	if fileName == "" || fileName == "." || filepath.Base(fileName) != fileName || strings.Contains(fileName, `\`) {
		return "", false
	}
	configuredStaticPath := viper.GetString("system.static_path")
	if configuredStaticPath == "" {
		return "", false
	}
	staticPath, err := filepath.Abs(configuredStaticPath)
	if err != nil {
		return "", false
	}
	absFilePath, err := filepath.Abs(filepath.Join(staticPath, fileName))
	if err != nil {
		return "", false
	}
	relativePath, err := filepath.Rel(staticPath, absFilePath)
	if err != nil || relativePath == ".." || strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
		return "", false
	}
	return absFilePath, true
}

func (s *FileService) HasOtherUsages(ctx context.Context, fileId []byte, entityId string, entityType string) (bool, error) {
	file, err := s.repo.FindByFileId(ctx, fileId)
	if err != nil {
		return false, err
	}
	for _, usage := range file.UsedIn {
		if usage.EntityId != entityId || usage.EntityType != entityType {
			return true, nil
		}
	}
	return false, nil
}

func (s *FileService) GetFiles(ctx context.Context, pageDTO domain.PageDTO) ([]*domain.File, int64, error) {
	return s.repo.FindPageFilesByFileType(ctx, pageDTO)
}

func (s *FileService) GenerateSitemap(_ context.Context, postBytes, categoryBytes, tagBytes []byte) error {
	baseHost := strings.TrimRight(pkg.GetOrDefault4String(os.Getenv("WEBSITE_BASE_HOST"), "http://localhost:3000"), "/")
	uploaderHost := strings.TrimRight(pkg.GetOrDefault4String(os.Getenv("UPLOADER_HOST"), "http://localhost:8080"), "/")
	var aboutMeLastMod string
	var posts []post.Post
	err := jsoniter.Unmarshal(postBytes, &posts)
	if err != nil {
		return err
	}
	var categories []category.Category
	err = jsoniter.Unmarshal(categoryBytes, &categories)
	if err != nil {
		return err
	}
	var tags []tag.Tag
	err = jsoniter.Unmarshal(tagBytes, &tags)
	if err != nil {
		return err
	}
	sitemapBuilder := sitemap.NewSitemap().
		XmlnsImage("https://www.google.com/schemas/sitemap-image/1.1").
		Url(
			baseHost,
			sitemap.WithLastMod(time.Now().Format(time.DateOnly)),
			sitemap.WithChangeFreq("always"),
			sitemap.WithPriority(1.0),
		).Output(filepath.Join(viper.GetString("system.static_path"), "sitemap.xml"))
	for _, p := range posts {
		if p.Id == "about-me" {
			aboutMeLastMod = time.Unix(p.UpdatedAt, 0).Format(time.DateOnly)
			continue
		}
		sitemapBuilder.Url(
			fmt.Sprintf("%s/posts/%s", baseHost, p.Id),
			sitemap.WithLastMod(time.Unix(p.UpdatedAt, 0).Format(time.DateOnly)),
			sitemap.WithChangeFreq("monthly"),
			sitemap.WithPriority(0.9),
			sitemap.WithImage(sitemap.NewUrlImage(fmt.Sprintf("%s%s", uploaderHost, p.CoverImg))),
		)
	}
	for _, c := range categories {
		sitemapBuilder.Url(
			fmt.Sprintf("%s/categories/%s", baseHost, c.Route),
			sitemap.WithLastMod(time.Unix(c.UpdatedAt, 0).Format(time.DateOnly)),
			sitemap.WithChangeFreq("weekly"),
			sitemap.WithPriority(0.8),
		)
	}
	for _, t := range tags {
		sitemapBuilder.Url(
			fmt.Sprintf("%s/tags/%s", baseHost, t.Route),
			sitemap.WithLastMod(time.Unix(t.UpdatedAt, 0).Format(time.DateOnly)),
			sitemap.WithChangeFreq("weekly"),
			sitemap.WithPriority(0.8),
		)
	}
	if aboutMeLastMod != "" {
		sitemapBuilder.Url(
			fmt.Sprintf("%s/about-me", baseHost),
			sitemap.WithLastMod(aboutMeLastMod),
			sitemap.WithChangeFreq("monthly"),
			sitemap.WithPriority(0.9),
		)
	}
	sitemapBuilder.Url(
		fmt.Sprintf("%s/friend", baseHost),
		sitemap.WithChangeFreq("always"),
		sitemap.WithPriority(0.5),
	)
	err = sitemapBuilder.GenerateXml()
	if err != nil {
		return err
	}
	return nil
}

func (s *FileService) GetSitemap(_ context.Context) (string, bool, error) {
	return s.getStaticTextFile("sitemap.xml")
}

func (s *FileService) GetRobotsTxt(_ context.Context) (string, bool, error) {
	return s.getStaticTextFile("robots.txt")
}

func (s *FileService) getStaticTextFile(filename string) (string, bool, error) {
	content, err := os.ReadFile(filepath.Join(viper.GetString("system.static_path"), filename))
	if errors.Is(err, os.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return string(content), true, nil
}

func (s *FileService) SaveRobotsTxt(_ context.Context, content string) error {
	staticPath := viper.GetString("system.static_path")
	if err := os.MkdirAll(staticPath, os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(staticPath, "robots.txt"), []byte(content), 0o644)
}

func (s *FileService) DeleteIndexFileMeta(ctx context.Context, fileId []byte, entityId string, entityType string) error {
	return s.repo.PullUsedIn(ctx, fileId, entityId, entityType)
}

func (s *FileService) IndexFileMeta(ctx context.Context, fileId []byte, entityId string, entityType string) error {
	return s.repo.PushIntoUsedIn(ctx, fileId, entityId, entityType)
}

func (s *FileService) Upload(ctx context.Context, fileDTO domain.FileDTO) (*domain.File, error) {
	var (
		filename string
	)
	fileId := uuidx.RearrangeUUID4()
	if fileDTO.CustomFileName != "" {
		if filepath.Base(fileDTO.CustomFileName) != fileDTO.CustomFileName || fileDTO.CustomFileName == "." {
			return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "invalid custom file name")
		}
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

	configuredStaticPath := viper.GetString("system.static_path")
	if configuredStaticPath == "" {
		return nil, errors.New("system.static_path is empty")
	}
	staticPath, err := filepath.Abs(configuredStaticPath)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(staticPath, os.ModePerm)
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(staticPath, filename)
	create, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	if _, err = create.Write(fileDTO.Content); err != nil {
		_ = create.Close()
		_ = os.Remove(filePath)
		return nil, err
	}
	if err = create.Close(); err != nil {
		_ = os.Remove(filePath)
		return nil, err
	}
	file := &domain.File{
		FileId:           fileId,
		FileName:         filename,
		OriginalFileName: fileDTO.FileName,
		FileType:         fileDTO.FileType,
		FileSize:         fileDTO.FileSize,
		FilePath:         filePath,
		Url:              "/static/" + filename,
	}
	err = s.repo.Save(ctx, file)
	if err != nil {
		if removeErr := os.Remove(filePath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
			return nil, errors.Wrapf(err, "save file metadata failed; rollback file failed: %v", removeErr)
		}
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
			l.ErrorContext(ctx, "File: post event: failed to unmarshal", "error", err)
			continue
		}
		switch e.Type {
		case "create":
			if len(e.AddedFileIds) > 0 {
				s.indexFileUsages(ctx, e.AddedFileIds, e.PostId, "post", l)
			} else {
				s.createIndexFileMeta4PostEvent(ctx, e.NewFileId, e.PostId, l)
			}
		case "update":
			if len(e.AddedFileIds) > 0 || len(e.DeletedFileIds) > 0 {
				s.indexFileUsages(ctx, e.AddedFileIds, e.PostId, "post", l)
				s.deleteFileUsages(ctx, e.DeletedFileIds, e.PostId, "post", l)
			} else if e.NewFileId != e.OldFileId {
				s.createIndexFileMeta4PostEvent(ctx, e.NewFileId, e.PostId, l)
				s.deleteIndexFileMeta4PostEvent(ctx, e.OldFileId, e.PostId, l)
			}
		case "delete":
			if len(e.DeletedFileIds) > 0 {
				s.deleteFileUsages(ctx, e.DeletedFileIds, e.PostId, "post", l)
			} else {
				s.deleteIndexFileMeta4PostEvent(ctx, e.OldFileId, e.PostId, l)
			}
		}
		l.InfoContext(ctx, "File: post event: handle successfully")
	}
}

func (s *FileService) subscribePostDraftEvent() {
	eventChan := s.eventBus.Subscribe("post-draft")
	for event := range eventChan {
		ctx := context.Background()
		l := slog.Default()
		var e domain.ContentFileEvent
		if err := jsoniter.Unmarshal(event.Payload, &e); err != nil {
			l.ErrorContext(ctx, "File: post draft event: failed to unmarshal", "error", err)
			continue
		}
		s.indexFileUsages(ctx, e.AddedFileIds, e.EntityId, "post-draft", l)
		s.deleteFileUsages(ctx, e.DeletedFileIds, e.EntityId, "post-draft", l)
	}
}

func (s *FileService) indexFileUsages(ctx context.Context, fileIDs []string, entityID, entityType string, l *slog.Logger) {
	for _, fileID := range fileIDs {
		fid, err := hex.DecodeString(fileID)
		if err != nil {
			l.ErrorContext(ctx, "failed to decode file id", "file_id", fileID, "error", err)
			continue
		}
		if err = s.IndexFileMeta(ctx, fid, entityID, entityType); err != nil {
			l.ErrorContext(ctx, "failed to index file usage", "file_id", fileID, "entity_id", entityID, "entity_type", entityType, "error", err)
		}
	}
}

func (s *FileService) deleteFileUsages(ctx context.Context, fileIDs []string, entityID, entityType string, l *slog.Logger) {
	for _, fileID := range fileIDs {
		fid, err := hex.DecodeString(fileID)
		if err != nil {
			l.ErrorContext(ctx, "failed to decode file id", "file_id", fileID, "error", err)
			continue
		}
		if err = s.DeleteIndexFileMeta(ctx, fid, entityID, entityType); err != nil {
			l.ErrorContext(ctx, "failed to delete file usage", "file_id", fileID, "entity_id", entityID, "entity_type", entityType, "error", err)
		}
	}
}

func (s *FileService) deleteIndexFileMeta4PostEvent(ctx context.Context, oldFileId string, postId string, l *slog.Logger) {
	fid, sErr := hex.DecodeString(oldFileId)
	if sErr != nil {
		l.ErrorContext(ctx, "File: post event: failed to hex.DecodeString", "fileId", oldFileId, "error", sErr)
		return
	}
	sErr = s.DeleteIndexFileMeta(ctx, fid, postId, "post")
	if sErr != nil {
		l.ErrorContext(ctx, "File: post event: failed to delete the index of file-meta ", "fileId", oldFileId, "postId", postId, "error", sErr)
	}
}

func (s *FileService) createIndexFileMeta4PostEvent(ctx context.Context, newFileId string, postId string, l *slog.Logger) {
	fid, sErr := hex.DecodeString(newFileId)
	if sErr != nil {
		l.ErrorContext(ctx, "File: post event: failed to hex.DecodeString", "fileId", newFileId, "error", sErr)
		return
	}
	sErr = s.IndexFileMeta(ctx, fid, postId, "post")
	if sErr != nil {
		l.ErrorContext(ctx, "File: post event: failed to index the file-meta ", "fileId", newFileId, "postId", postId, "error", sErr)
	}
}
