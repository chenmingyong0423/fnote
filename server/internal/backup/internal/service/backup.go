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
	"archive/zip"
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	backupDataDir   = "data"
	backupStaticDir = "static"
)

type IBackupService interface {
	GetBackups(ctx context.Context) (string, error)
	Recovery(ctx context.Context, file *multipart.FileHeader) error
}

var _ IBackupService = (*BackupService)(nil)

func NewBackupService(db *mongox.Database) *BackupService {
	return &BackupService{db: db.Database()}
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

	content, err := io.ReadAll(src)
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return fmt.Errorf("backup file is empty")
	}
	if !isZipFile(content) {
		return fmt.Errorf("unsupported backup file format: only zip backup files are supported")
	}

	err = s.recoveryZip(ctx, bytes.NewReader(content), int64(len(content)))
	if err != nil {
		return fmt.Errorf("restore zip backup failed: %w", err)
	}
	return nil
}

func (s *BackupService) recoveryZip(ctx context.Context, src io.ReaderAt, size int64) error {
	zipReader, err := zip.NewReader(src, size)
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		name := cleanArchiveName(file.Name)
		if name == "" {
			continue
		}
		name = normalizeBackupEntryName(name)
		if file.FileInfo().IsDir() {
			if strings.HasPrefix(name, backupStaticDir+"/") {
				if err = mkdirStaticDir(strings.TrimPrefix(name, backupStaticDir+"/")); err != nil {
					return err
				}
			}
			continue
		}

		switch {
		case strings.HasPrefix(name, backupDataDir+"/") && strings.HasSuffix(name, ".json"):
			if err = s.restoreZipCollection(ctx, file); err != nil {
				return err
			}
		case strings.HasPrefix(name, backupStaticDir+"/"):
			if err = restoreZipStaticFile(file, strings.TrimPrefix(name, backupStaticDir+"/")); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *BackupService) restoreZipCollection(ctx context.Context, file *zip.File) error {
	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return s.restoreCollection(ctx, path.Base(cleanArchiveName(file.Name)), content)
}

func (s *BackupService) restoreCollection(ctx context.Context, filename string, content []byte) error {
	colName, err := s.collectionNameFromBackupFile(filename)
	if err != nil {
		return err
	}
	return s.DeleteAndInsertCollectionDoc(ctx, colName, content)
}

func (s *BackupService) collectionNameFromBackupFile(filename string) (string, error) {
	base := filepath.Base(filename)
	if filepath.Ext(base) != ".json" {
		return "", fmt.Errorf("file name error: %s", filename)
	}

	name := strings.TrimSuffix(base, ".json")
	dbPrefix := s.db.Name() + "_"
	if strings.HasPrefix(name, dbPrefix) {
		return strings.TrimPrefix(name, dbPrefix), nil
	}

	if strings.HasPrefix(name, "fnote_") {
		return strings.TrimPrefix(name, "fnote_"), nil
	}

	return "", fmt.Errorf("file name error: %s", filename)
}

func (s *BackupService) GetBackups(ctx context.Context) (zipFileName string, err error) {
	staticPath := viper.GetString("system.static_path")
	if staticPath == "" {
		return "", fmt.Errorf("system.static_path is empty")
	}

	if err = os.MkdirAll(staticPath, os.ModePerm); err != nil {
		return "", err
	}

	tempDir, err := os.MkdirTemp("", "fnote-backup-*")
	if err != nil {
		return "", err
	}
	defer func() {
		if fErr := os.RemoveAll(tempDir); fErr != nil {
			slog.Error("remove backup temp dir failed", "dir", tempDir, "error", fErr)
		}
	}()

	if err = s.exportCollections(ctx, filepath.Join(tempDir, backupDataDir)); err != nil {
		return "", err
	}

	zipFileName = filepath.Join(staticPath, fmt.Sprintf("backup_%s.zip", time.Now().Local().Format("2006-01-02_150405")))
	fileCount, err := s.createZip(zipFileName, filepath.Join(tempDir, backupDataDir), staticPath)
	if err != nil {
		return "", err
	}
	if fileCount == 0 {
		if fErr := os.Remove(zipFileName); fErr != nil && !os.IsNotExist(fErr) {
			slog.Error("remove empty backup file failed", "file", zipFileName, "error", fErr)
		}
		return "", nil
	}
	return zipFileName, nil
}

func (s *BackupService) exportCollections(ctx context.Context, dataDir string) error {
	dbName := s.db.Name()
	collections, err := s.db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return err
	}

	for _, collectionName := range collections {
		cur, err := s.db.Collection(collectionName).Find(ctx, bson.M{})
		if err != nil {
			return err
		}

		var documents []bson.M
		if err = cur.All(ctx, &documents); err != nil {
			return err
		}
		if len(documents) == 0 {
			continue
		}

		if collectionName == "configs" {
			for _, document := range documents {
				if typ, ok := document["typ"]; ok && typ == "social" {
					props := document["props"].(bson.D)
					for _, prop := range props {
						if prop.Key == "social_info_list" {
							socialList := prop.Value.(bson.A)
							for _, m := range socialList {
								encodeSocialID(m)
							}
							break
						}
					}
					break
				}
			}
		}

		fileContent, err := json.MarshalIndent(documents, "", "    ")
		if err != nil {
			return err
		}
		if err = os.MkdirAll(dataDir, os.ModePerm); err != nil {
			return err
		}

		filename := filepath.Join(dataDir, fmt.Sprintf("%s_%s.json", dbName, collectionName))
		if err = os.WriteFile(filename, fileContent, 0644); err != nil {
			return err
		}
	}
	return nil
}

func (s *BackupService) createZip(zipFileName, dataDir, staticPath string) (int, error) {
	newZipFile, err := os.Create(zipFileName)
	if err != nil {
		return 0, err
	}

	createSuccess := false
	defer func() {
		if fErr := newZipFile.Close(); fErr != nil {
			slog.Error("close backup file failed", "file", zipFileName, "error", fErr)
		}
		if !createSuccess {
			if fErr := os.Remove(zipFileName); fErr != nil && !os.IsNotExist(fErr) {
				slog.Error("remove failed backup file failed", "file", zipFileName, "error", fErr)
			}
		}
	}()

	zipWriter := zip.NewWriter(newZipFile)

	fileCount := 0
	if err = filepath.WalkDir(dataDir, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}

		archiveName := path.Join(backupDataDir, filepath.Base(filePath))
		if err = addFileToZip(zipWriter, filePath, archiveName); err != nil {
			return err
		}
		fileCount++
		return nil
	}); err != nil && !os.IsNotExist(err) {
		return 0, err
	}

	zipAbs, err := filepath.Abs(zipFileName)
	if err != nil {
		return 0, err
	}
	if err = filepath.WalkDir(staticPath, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}

		fileAbs, err := filepath.Abs(filePath)
		if err != nil {
			return err
		}
		if fileAbs == zipAbs || isBackupArchive(filePath) {
			return nil
		}

		relPath, err := filepath.Rel(staticPath, filePath)
		if err != nil {
			return err
		}
		archiveName := path.Join(backupStaticDir, filepath.ToSlash(relPath))
		if err = addFileToZip(zipWriter, filePath, archiveName); err != nil {
			return err
		}
		fileCount++
		return nil
	}); err != nil {
		return 0, err
	}

	if err = zipWriter.Close(); err != nil {
		return 0, err
	}
	createSuccess = true
	return fileCount, nil
}

func addFileToZip(zipWriter *zip.Writer, filename, archiveName string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.ToSlash(archiveName)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	return err
}

func restoreZipStaticFile(file *zip.File, relPath string) error {
	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	mode := file.Mode()
	if mode == 0 {
		mode = 0644
	}
	return restoreStaticFile(reader, relPath, mode)
}

func restoreStaticFile(reader io.Reader, relPath string, mode fs.FileMode) error {
	targetPath, err := staticTargetPath(relPath)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
		return err
	}

	if mode == 0 {
		mode = 0644
	}
	file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

func mkdirStaticDir(relPath string) error {
	targetPath, err := staticTargetPath(relPath)
	if err != nil {
		return err
	}
	return os.MkdirAll(targetPath, os.ModePerm)
}

func staticTargetPath(relPath string) (string, error) {
	staticPath := viper.GetString("system.static_path")
	if staticPath == "" {
		return "", fmt.Errorf("system.static_path is empty")
	}

	cleanRel := filepath.Clean(filepath.FromSlash(relPath))
	if cleanRel == "." || cleanRel == ".." || filepath.IsAbs(cleanRel) || strings.HasPrefix(cleanRel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid static file path: %s", relPath)
	}

	baseAbs, err := filepath.Abs(staticPath)
	if err != nil {
		return "", err
	}
	targetAbs, err := filepath.Abs(filepath.Join(baseAbs, cleanRel))
	if err != nil {
		return "", err
	}

	basePrefix := strings.TrimRight(baseAbs, `\/`) + string(os.PathSeparator)
	if targetAbs != baseAbs && !strings.HasPrefix(targetAbs, basePrefix) {
		return "", fmt.Errorf("invalid static file path: %s", relPath)
	}
	return targetAbs, nil
}

func cleanArchiveName(name string) string {
	cleaned := path.Clean(strings.ReplaceAll(name, "\\", "/"))
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") || strings.HasPrefix(cleaned, "/") {
		return ""
	}
	return cleaned
}

func normalizeBackupEntryName(name string) string {
	if strings.HasPrefix(name, backupDataDir+"/") || strings.HasPrefix(name, backupStaticDir+"/") {
		return name
	}
	parts := strings.SplitN(name, "/", 2)
	if len(parts) != 2 {
		return name
	}
	if strings.HasPrefix(parts[1], backupDataDir+"/") || strings.HasPrefix(parts[1], backupStaticDir+"/") {
		return parts[1]
	}
	return name
}

func isZipFile(content []byte) bool {
	return len(content) >= 2 && bytes.Equal(content[:2], []byte{'P', 'K'})
}

func isBackupArchive(filename string) bool {
	base := filepath.Base(filename)
	matched, err := filepath.Match("backup_*.zip", base)
	return err == nil && matched
}

func (s *BackupService) DeleteAndInsertCollectionDoc(ctx context.Context, colName string, content []byte) error {
	var documents []map[string]any
	if err := json.Unmarshal(content, &documents); err != nil {
		return err
	}
	col := s.db.Collection(colName)
	for _, doc := range documents {
		objectID, fErr := bson.ObjectIDFromHex(doc["_id"].(string))
		// Some collections use custom string IDs instead of ObjectIDs.
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
					obj["id"], fErr2 = decodeSocialID(obj["id"])
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

func encodeSocialID(value any) {
	switch obj := value.(type) {
	case bson.M:
		if binary, ok := obj["id"].(bson.Binary); ok {
			obj["id"] = hex.EncodeToString(binary.Data)
		}
	case bson.D:
		for i := range obj {
			if obj[i].Key != "id" {
				continue
			}
			if binary, ok := obj[i].Value.(bson.Binary); ok {
				obj[i].Value = hex.EncodeToString(binary.Data)
			}
			return
		}
	}
}

func decodeSocialID(value any) ([]byte, error) {
	switch v := value.(type) {
	case string:
		return hex.DecodeString(v)
	case map[string]any:
		if data, ok := v["Data"].(string); ok {
			return hex.DecodeString(data)
		}
	}
	return nil, fmt.Errorf("invalid social id format")
}
