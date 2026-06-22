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
	"strings"

	"github.com/chenmingyong0423/fnote/server/internal/file/internal/domain"

	"github.com/chenmingyong0423/fnote/server/internal/file/internal/repository/dao"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IFileRepository interface {
	Save(ctx context.Context, file *domain.File) error
	FindByFileId(ctx context.Context, fileId []byte) (*domain.File, error)
	PushIntoUsedIn(ctx context.Context, fileId []byte, entityId string, entityType string) error
	PullUsedIn(ctx context.Context, fileId []byte, entityId string, entityType string) error
	FindByFileName(ctx context.Context, filename string) (*domain.File, error)
	FindPageFilesByFileType(ctx context.Context, pageDTO domain.PageDTO) ([]*domain.File, int64, error)
	DeleteUnusedByFileId(ctx context.Context, fileId []byte) (int64, error)
}

var _ IFileRepository = (*FileRepository)(nil)

func NewFileRepository(dao dao.IFileDao) *FileRepository {
	return &FileRepository{dao: dao}
}

type FileRepository struct {
	dao dao.IFileDao
}

func (r *FileRepository) DeleteUnusedByFileId(ctx context.Context, fileId []byte) (int64, error) {
	return r.dao.DeleteUnusedByFileId(ctx, fileId)
}

func (r *FileRepository) FindByFileId(ctx context.Context, fileId []byte) (*domain.File, error) {
	file, err := r.dao.FindByFileId(ctx, fileId)
	if err != nil {
		return nil, err
	}
	return r.toDomainFile(file), nil
}

func (r *FileRepository) FindPageFilesByFileType(ctx context.Context, pageDTO domain.PageDTO) ([]*domain.File, int64, error) {
	files, cnt, err := r.dao.FindPageByFileType(ctx, pageDTO.PageNum, pageDTO.PageSize, pageDTO.FileType, pageDTO.Unused, pageDTO.Keyword)
	if err != nil {
		return nil, 0, err
	}
	domainFiles := r.toDomainFiles(files)
	if err = r.resolveFileUsages(ctx, domainFiles); err != nil {
		return nil, 0, err
	}
	return domainFiles, cnt, nil
}

func (r *FileRepository) resolveFileUsages(ctx context.Context, files []*domain.File) error {
	postIDs := make(map[string]struct{})
	postDraftIDs := make(map[string]struct{})
	configIDs := make(map[bson.ObjectID]struct{})
	assetIDs := make(map[bson.ObjectID]struct{})
	for _, file := range files {
		for _, usage := range file.UsedIn {
			switch usage.EntityType {
			case string(dao.EntityTypePost):
				postIDs[usage.EntityId] = struct{}{}
			case string(dao.EntityTypePostDraft):
				postDraftIDs[usage.EntityId] = struct{}{}
			case string(dao.EntityTypeConfig):
				id, err := bson.ObjectIDFromHex(usage.EntityId)
				if err == nil {
					configIDs[id] = struct{}{}
				}
			case string(dao.EntityTypeAsset):
				id, err := bson.ObjectIDFromHex(usage.EntityId)
				if err == nil {
					assetIDs[id] = struct{}{}
				}
			}
		}
	}

	posts, err := r.dao.FindUsagePosts(ctx, mapKeys(postIDs))
	if err != nil {
		return err
	}
	postDrafts, err := r.dao.FindUsagePostDrafts(ctx, mapKeys(postDraftIDs))
	if err != nil {
		return err
	}
	configs, err := r.dao.FindUsageConfigs(ctx, mapKeys(configIDs))
	if err != nil {
		return err
	}
	assets, err := r.dao.FindUsageAssets(ctx, mapKeys(assetIDs))
	if err != nil {
		return err
	}

	postMap := make(map[string]*dao.UsageContent, len(posts))
	for _, item := range posts {
		postMap[item.ID] = item
	}
	postDraftMap := make(map[string]*dao.UsageContent, len(postDrafts))
	for _, item := range postDrafts {
		postDraftMap[item.ID] = item
	}
	configMap := make(map[string]*dao.UsageConfig, len(configs))
	for _, item := range configs {
		configMap[item.ID.Hex()] = item
	}
	assetMap := make(map[string]*dao.UsageAsset, len(assets))
	for _, item := range assets {
		assetMap[item.ID.Hex()] = item
	}

	for _, file := range files {
		for i := range file.UsedIn {
			usage := &file.UsedIn[i]
			switch usage.EntityType {
			case string(dao.EntityTypePost):
				if item := postMap[usage.EntityId]; item != nil {
					usage.Name = item.Title
					usage.Locations = contentLocations(file.Url, item, "文章正文", "文章封面")
				}
			case string(dao.EntityTypePostDraft):
				if item := postDraftMap[usage.EntityId]; item != nil {
					usage.Name = item.Title
					usage.Locations = contentLocations(file.Url, item, "草稿正文", "草稿封面")
				}
			case string(dao.EntityTypeConfig):
				if item := configMap[usage.EntityId]; item != nil {
					usage.Name = configName(item.Typ)
					usage.Locations = configLocations(file.Url, item.Props)
				}
			case string(dao.EntityTypeAsset):
				if item := assetMap[usage.EntityId]; item != nil {
					usage.Name = item.Title
					if usage.Name == "" {
						usage.Name = "素材"
					}
					usage.Locations = []string{assetLocation(item.Type)}
				}
			}
		}
	}
	return nil
}

func assetLocation(typ string) string {
	if typ == "post-editor" {
		return "文章编辑器素材库"
	}
	return "素材库"
}

func mapKeys[K comparable](values map[K]struct{}) []K {
	keys := make([]K, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}

func contentLocations(fileURL string, content *dao.UsageContent, contentLabel, coverLabel string) []string {
	locations := make([]string, 0, 2)
	if fileURL == "" {
		return locations
	}
	if strings.Contains(content.Content, fileURL) {
		locations = append(locations, contentLabel)
	}
	if strings.Contains(content.CoverImg, fileURL) {
		locations = append(locations, coverLabel)
	}
	return locations
}

func configName(typ string) string {
	names := map[string]string{
		"website":  "网站配置",
		"seo meta": "SEO 配置",
		"pay":      "收款配置",
		"carousel": "轮播图配置",
		"notice":   "公告配置",
		"friend":   "友链配置",
		"social":   "社交配置",
	}
	if name := names[typ]; name != "" {
		return name
	}
	return typ
}

func configLocations(fileURL string, props dao.UsageConfigProps) []string {
	locations := make([]string, 0)
	if fileURL == "" {
		return locations
	}
	appendIfContains := func(value, label string) {
		if strings.Contains(value, fileURL) {
			locations = appendUnique(locations, label)
		}
	}
	appendIfContains(props.WebsiteIcon, "站点 Logo")
	appendIfContains(props.WebsiteOwnerAvatar, "站长头像")
	appendIfContains(props.OgImage, "SEO 分享图片")
	appendIfContains(props.Content, "公告内容")
	appendIfContains(props.Introduction, "友链介绍")
	for _, record := range props.WebsiteRecords {
		appendIfContains(record, "备案信息")
	}
	for _, item := range props.List {
		appendIfContains(item.Image, "收款二维码")
		appendIfContains(item.CoverImg, "轮播图")
	}
	for _, item := range props.SocialInfoList {
		appendIfContains(item.SocialValue, "社交配置")
	}
	return locations
}

func appendUnique(values []string, value string) []string {
	for _, current := range values {
		if current == value {
			return values
		}
	}
	return append(values, value)
}

func (r *FileRepository) toDomainFiles(files []*dao.File) []*domain.File {
	result := make([]*domain.File, 0)
	for _, file := range files {
		result = append(result, r.toDomainFile(file))
	}
	return result
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
	fileId, err := hex.DecodeString(file.FileId)
	if err != nil {
		return err
	}
	_, err = r.dao.Save(ctx, &dao.File{
		FileId:           fileId,
		FileName:         file.FileName,
		OriginalFileName: file.OriginalFileName,
		FileType:         file.FileType,
		FileSize:         file.FileSize,
		FilePath:         file.FilePath,
		Url:              file.Url,
		UsedIn:           make([]dao.FileUsage, 0),
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *FileRepository) toDomainFile(file *dao.File) *domain.File {
	return &domain.File{
		Id:               file.ID.Hex(),
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
					Locations:  make([]string, 0),
				})
			}
			return usedIn
		}(),
		CreatedAt: file.CreatedAt.Unix(),
		UpdatedAt: file.UpdatedAt.Unix(),
	}
}
