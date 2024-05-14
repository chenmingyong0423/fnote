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
	"archive/tar"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type IBackupService interface {
	GetBackups(ctx context.Context) (string, error)
	Recovery(ctx context.Context, file *multipart.FileHeader) error
}

var _ IBackupService = (*BackupService)(nil)

func NewBackupService(db *mongo.Database) *BackupService {
	return &BackupService{db: db}
}

type BackupService struct {
	db *mongo.Database
}

func (s *BackupService) Recovery(ctx context.Context, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	// 创建一个tar的读取器
	tarReader := tar.NewReader(src)
	for {
		var (
			header      *tar.Header
			fileContent []byte
		)
		header, err = tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// 读取文件内容
		fileContent, err = io.ReadAll(tarReader)
		if err != nil {
			return err
		}
		// 将文件内容写入到 MongoDB 的集合中
		// header.Name 为文件名.后缀，获取到文件名，去除 fnote_ 和后缀
		// 例如：fnote_config.json，获取到 config
		split := strings.Split(header.Name, ".")
		colName := ""
		if len(split) != 2 {
			return fmt.Errorf("file name error: %s", header.Name)
		}
		colName = strings.Replace(split[0], "fnote_", "", 1)
		err = s.DeleteAndInsertCollectionDoc(ctx, colName, fileContent)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BackupService) GetBackups(ctx context.Context) (zipFileName string, err error) {
	var files []string
	defer func() {
		for _, file := range files {
			fErr := os.Remove(file)
			if fErr != nil {
				slog.Error("remove file: %s failed, error: %v", file, fErr)
			}
		}
	}()
	dbName := s.db.Name()
	// 获取到所有的 db 集合
	collections, err := s.db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return "", err
	}
	for _, collectionName := range collections {
		var (
			cur         *mongo.Cursor
			documents   []bson.M
			fileContent []byte
		)
		cur, err = s.db.Collection(collectionName).Find(ctx, bson.M{})
		if err != nil {
			return "", err
		}

		if err = cur.All(ctx, &documents); err != nil {
			return "", err
		}
		if len(documents) > 0 {
			filename := fmt.Sprintf("%s%s_%s.json", viper.GetString("system.static_path"), dbName, collectionName)
			if collectionName == "configs" {
				for _, document := range documents {
					if typ, ok := document["typ"]; ok && typ == "social" {
						props := document["props"].(primitive.M)
						socialList := props["social_info_list"].(primitive.A)
						for _, m := range socialList {
							obj := m.(primitive.M)
							obj["id"] = hex.EncodeToString(obj["id"].(primitive.Binary).Data)
						}
					}
				}
			}
			fileContent, err = json.MarshalIndent(documents, "", "    ")
			if err != nil {
				return "", err
			}
			if err = os.WriteFile(filename, fileContent, 0644); err != nil {
				return "", err
			}
			files = append(files, filename)
		}
	}
	if len(files) == 0 {
		return "", nil
	}
	zipFileName = fmt.Sprintf("%sbackup_%s.zip", viper.GetString("system.static_path"), time.Now().Local().Format(time.DateOnly))
	if err = s.createTar(zipFileName, files); err != nil {
		return "", err
	}
	return zipFileName, nil
}

// createTar 创建一个tar文件，并将files列表中的文件添加到这个tar文件中。
func (s *BackupService) createTar(tarFileName string, files []string) error {
	createSuccess := false
	// 创建tar文件
	newTarFile, err := os.Create(tarFileName)
	if err != nil {
		return err
	}
	defer func() {
		fErr := newTarFile.Close()
		if fErr != nil {
			slog.Error("close file: %s failed, error: %v", tarFileName, fErr)
		}
		if err != nil && createSuccess {
			fErr = os.Remove(tarFileName)
			if fErr != nil {
				slog.Error("remove file: %s failed, error: %v", tarFileName, fErr)
			}
		}
	}()
	createSuccess = true
	// 创建一个tar的写入器
	tarWriter := tar.NewWriter(newTarFile)
	defer tarWriter.Close()

	// 遍历files列表，将每个文件添加到tar中
	for _, file := range files {
		if err := s.addFileToTar(tarWriter, file); err != nil {
			return err
		}
	}
	return nil
}

// addFileToTar 向tar文件中添加一个文件。
func (s *BackupService) addFileToTar(tarWriter *tar.Writer, filename string) error {
	// 打开需要添加到tar的文件
	fileToTar, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToTar.Close()

	// 获取文件信息，用于设置tar头信息
	info, err := fileToTar.Stat()
	if err != nil {
		return err
	}

	// 创建一个对应于当前文件的tar头
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	// 确保使用相对路径或文件名，避免在tar中创建不必要的目录结构
	header.Name = filepath.Base(filename)

	// 写入文件头到tar
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}
	// 将文件内容写入tar
	_, err = io.Copy(tarWriter, fileToTar)
	return err
}

func (s *BackupService) DeleteAndInsertCollectionDoc(ctx context.Context, colName string, content []byte) error {
	// content 是一个 json 数组，将其插入到集合中
	var documents []map[string]any
	if err := json.Unmarshal(content, &documents); err != nil {
		return err
	}
	col := s.db.Collection(colName)
	for _, doc := range documents {
		objectID, fErr := primitive.ObjectIDFromHex(doc["_id"].(string))
		// 部分集合的 id 不是 objectID
		if fErr == nil {
			doc["_id"] = objectID
		}

		if colName == "configs" {
			if typ, ok := doc["typ"]; ok && typ == "social" {
				props := doc["props"].(map[string]any)
				socialList := props["social_info_list"].([]any)
				for _, m := range socialList {
					obj := m.(map[string]any)
					var fErr2 error
					obj["id"], fErr2 = hex.DecodeString(obj["id"].(string))
					if fErr2 != nil {
						return fErr2
					}
				}
			}
			if typ, ok := doc["typ"]; ok && typ == "website" {
				props := doc["props"].(map[string]any)
				websiteRunTime := props["website_runtime"].(string)
				parse, fErr2 := time.Parse(time.RFC3339, websiteRunTime)
				if fErr2 != nil {
					return fErr2
				}
				props["website_runtime"] = parse
			}
			if typ, ok := doc["typ"]; ok && typ == "notice" {
				props := doc["props"].(map[string]any)
				publishTime := props["publish_time"].(string)
				parse, fErr2 := time.Parse(time.RFC3339, publishTime)
				if fErr2 != nil {
					return fErr2
				}
				props["publish_time"] = parse
			}
			if typ, ok := doc["typ"]; ok && typ == "carousel" {
				props := doc["props"].(map[string]any)
				list := props["list"].([]any)
				for _, m := range list {
					obj := m.(map[string]any)
					createdAt := obj["created_at"].(string)
					parse, fErr2 := time.Parse(time.RFC3339, createdAt)
					if fErr2 != nil {
						return fErr2
					}
					obj["created_at"] = parse
					updatedAt := obj["updated_at"].(string)
					parse, fErr2 = time.Parse(time.RFC3339, updatedAt)
					if fErr2 != nil {
						return fErr2
					}
					obj["updated_at"] = parse
				}
			}
		}
		if createdAt, ok := doc["created_at"].(string); ok {
			parse, fErr2 := time.Parse(time.RFC3339, createdAt)
			if fErr2 != nil {
				return fErr2
			}
			doc["created_at"] = parse
		}
		if updatedAt, ok := doc["updated_at"].(string); ok {
			parse, fErr2 := time.Parse(time.RFC3339, updatedAt)
			if fErr2 != nil {
				return fErr2
			}
			doc["updated_at"] = parse
		}
	}
	// 先清空数据
	_, err := col.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}
	var documentsAny []any
	for _, doc := range documents {
		documentsAny = append(documentsAny, doc)
	}
	_, err = col.InsertMany(ctx, documentsAny)
	if err != nil {
		return err
	}
	return nil
}
