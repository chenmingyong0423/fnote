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
		if err == nil {
			for _, file := range files {
				err2 := os.Remove(file)
				if err2 != nil {
					slog.Error("remote file: %s failed, error: %v", file, err2)
				}
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
	zipFileName = fmt.Sprintf("%sbackup_%s.zip", viper.GetString("system.static_path"), time.Now().Format(time.DateOnly))
	if err = s.createTar(zipFileName, files); err != nil {
		return "", err
	}
	return zipFileName, nil
}

// createTar 创建一个tar文件，并将files列表中的文件添加到这个tar文件中。
func (s *BackupService) createTar(tarFileName string, files []string) error {
	// 创建tar文件
	newTarFile, err := os.Create(tarFileName)
	if err != nil {
		return err
	}
	defer newTarFile.Close()

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
	var documents []bson.M
	if err := json.Unmarshal(content, &documents); err != nil {
		return err
	}
	col := s.db.Collection(colName)
	// 挨个插入到集合中，如果有 dup key，忽略
	for _, doc := range documents {
		objectID, err2 := primitive.ObjectIDFromHex(doc["_id"].(string))
		if err2 != nil {
			return err2
		}
		doc["_id"] = objectID
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
