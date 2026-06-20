package main

import (
	"context"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type fileUsage struct {
	EntityID   string `bson:"entity_id" json:"entity_id"`
	EntityType string `bson:"entity_type" json:"entity_type"`
}

type fileRecord struct {
	ID       bson.ObjectID `bson:"_id"`
	FileID   []byte        `bson:"file_id"`
	FileName string        `bson:"file_name"`
	FilePath string        `bson:"file_path"`
	URL      string        `bson:"url"`
	UsedIn   []fileUsage   `bson:"used_in"`
}

type missingFile struct {
	Record       fileRecord
	ExpectedPath string
	Reason       string
}

func main() {
	mongoURI := flag.String("mongo-uri", envOrDefault("MONGODB_URI", ""), "MongoDB URI，或设置 MONGODB_URI")
	database := flag.String("database", envOrDefault("MONGODB_DATABASE", "fnote"), "数据库名称")
	username := flag.String("username", envOrDefault("MONGODB_USERNAME", ""), "MongoDB 用户名")
	password := flag.String("password", envOrDefault("MONGODB_PASSWORD", ""), "MongoDB 密码")
	authSource := flag.String("auth-source", envOrDefault("MONGODB_AUTH_SOURCE", "fnote"), "MongoDB 认证数据库")
	staticPath := flag.String("static-path", envOrDefault("STATIC_PATH", ""), "服务器当前静态文件目录，例如 /fnote/static")
	output := flag.String("output", "", "报告路径，默认在当前目录生成 TSV 文件")
	flag.Parse()

	if *mongoURI == "" {
		exitWithError("必须通过 -mongo-uri 或 MONGODB_URI 指定 MongoDB 地址")
	}
	if *staticPath == "" {
		fmt.Fprintln(os.Stderr, "警告：未指定 -static-path，将使用 MongoDB 中保存的 file_path")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(*mongoURI)
	if *username != "" {
		clientOptions.SetAuth(options.Credential{
			Username:   *username,
			Password:   *password,
			AuthSource: *authSource,
		})
	}
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		exitWithError("连接 MongoDB 失败：%v", err)
	}
	defer client.Disconnect(context.Background())
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		exitWithError("MongoDB Ping 失败：%v", err)
	}

	cursor, err := client.Database(*database).Collection("file_meta").Find(ctx, bson.D{}, options.Find().SetProjection(bson.M{
		"file_id": 1, "file_name": 1, "file_path": 1, "url": 1, "used_in": 1,
	}))
	if err != nil {
		exitWithError("查询 file_meta 失败：%v", err)
	}
	defer cursor.Close(context.Background())

	missing := make([]missingFile, 0)
	total := 0
	for cursor.Next(ctx) {
		var record fileRecord
		if err = cursor.Decode(&record); err != nil {
			exitWithError("解析 file_meta 失败：%v", err)
		}
		total++
		expectedPath := resolveFilePath(record, *staticPath)
		reason := missingReason(expectedPath)
		if reason != "" {
			missing = append(missing, missingFile{Record: record, ExpectedPath: expectedPath, Reason: reason})
		}
	}
	if err = cursor.Err(); err != nil {
		exitWithError("遍历 file_meta 失败：%v", err)
	}

	if *output == "" {
		*output = fmt.Sprintf("missing-file-report-%s.tsv", time.Now().Format("20060102-150405"))
	}
	if err = writeReport(*output, missing); err != nil {
		exitWithError("写入报告失败：%v", err)
	}

	absOutput, _ := filepath.Abs(*output)
	fmt.Printf("检查完成：MongoDB 文件记录 %d 条，磁盘缺失 %d 条\n", total, len(missing))
	fmt.Printf("报告位置：%s\n", absOutput)
}

func resolveFilePath(record fileRecord, staticPath string) string {
	if staticPath != "" {
		fileName := record.FileName
		if fileName == "" {
			fileName = filepath.Base(record.URL)
		}
		return filepath.Join(staticPath, fileName)
	}
	return record.FilePath
}

func missingReason(path string) string {
	if path == "" {
		return "没有可检查的文件路径"
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "文件不存在"
		}
		if os.IsPermission(err) {
			return "没有文件访问权限"
		}
		return err.Error()
	}
	if !info.Mode().IsRegular() {
		return "路径存在但不是普通文件"
	}
	return ""
}

func writeReport(path string, missing []missingFile) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = '\t'
	defer writer.Flush()
	if err = writer.Write([]string{
		"mongo_id", "file_id", "file_name", "expected_path", "url", "used_in_count", "used_in", "reason",
	}); err != nil {
		return err
	}
	for _, item := range missing {
		usedIn, marshalErr := json.Marshal(item.Record.UsedIn)
		if marshalErr != nil {
			return marshalErr
		}
		if err = writer.Write([]string{
			item.Record.ID.Hex(),
			hex.EncodeToString(item.Record.FileID),
			item.Record.FileName,
			item.ExpectedPath,
			item.Record.URL,
			strconv.Itoa(len(item.Record.UsedIn)),
			string(usedIn),
			item.Reason,
		}); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func exitWithError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
