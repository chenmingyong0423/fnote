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

package web

import (
	"io"
	"path/filepath"

	"github.com/chenmingyong0423/fnote/server/internal/file/internal/domain"

	"github.com/chenmingyong0423/fnote/server/internal/file/internal/service"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/gin-gonic/gin"
)

func NewFileHandler(serv service.IFileService) *FileHandler {
	return &FileHandler{
		serv: serv,
	}
}

type FileHandler struct {
	serv service.IFileService
}

func (h *FileHandler) RegisterGinRoutes(engine *gin.Engine) {
	adminGroup := engine.Group("/admin-api/files")
	adminGroup.POST("/upload", apiwrap.Wrap(h.UploadFile))
	adminGroup.GET("", apiwrap.WrapWithBody(h.GetFiles))
}

func (h *FileHandler) UploadFile(ctx *gin.Context) (*apiwrap.ResponseBody[FileVO], error) {
	fileName := ctx.PostForm("file_name")
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, err
	}
	openedFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer openedFile.Close()
	content, err := io.ReadAll(openedFile)
	if err != nil {
		return nil, err
	}
	fileDto := domain.FileDTO{
		FileName:       file.Filename,
		FileSize:       file.Size,
		Content:        content,
		FileType:       file.Header.Get("Content-Type"),
		FileExt:        filepath.Ext(file.Filename),
		CustomFileName: fileName,
	}
	fileInfo, err := h.serv.Upload(ctx, fileDto)
	if err != nil {
		return nil, err

	}
	return apiwrap.SuccessResponseWithData(FileVO{
		FileId:   fileInfo.FileId,
		FileName: fileInfo.FileName,
		Url:      fileInfo.Url,
	}), nil
}

func (h *FileHandler) GetFiles(ctx *gin.Context, req PageRequest) (*apiwrap.ResponseBody[apiwrap.PageVO[FileVO]], error) {
	files, total, err := h.serv.GetFiles(ctx, domain.PageDTO{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		FileType: req.FileType,
	})
	if err != nil {
		return nil, err
	}
	pageVO := apiwrap.PageVO[FileVO]{}
	pageVO.PageNo = req.PageNum
	pageVO.PageSize = req.PageSize
	pageVO.List = h.toVOs(files)
	pageVO.SetTotalCountAndCalculateTotalPages(total)
	return apiwrap.SuccessResponseWithData(pageVO), nil
}

func (h *FileHandler) toVOs(files []*domain.File) []FileVO {
	voList := make([]FileVO, len(files))
	for i, file := range files {
		voList[i] = FileVO{
			FileId:   file.FileId,
			FileName: file.FileName,
			Url:      file.Url,
		}
	}
	return voList
}
