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

package handler

import (
	"log/slog"
	"os"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/chenmingyong0423/fnote/server/internal/backup/service"
	"github.com/gin-gonic/gin"
)

func NewBackupHandler(serv service.IBackupService) *BackupHandler {
	return &BackupHandler{
		serv: serv,
	}
}

type BackupHandler struct {
	serv service.IBackupService
}

func (h *BackupHandler) RegisterGinRoutes(engine *gin.Engine) {
	adminGroup := engine.Group("/admin")

	adminGroup.GET("/backup", h.GetBackups)
	adminGroup.POST("/recovery", apiwrap.Wrap(h.Recovery))
}

func (h *BackupHandler) GetBackups(ctx *gin.Context) {

	zipFileName, err := h.serv.GetBackups(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}
	defer func() {
		if zipFileName != "" {
			gErr := os.Remove(zipFileName)
			if err != nil {
				slog.Error("remote file: %s failed, error: %v", zipFileName, gErr)
			}
		}
	}()
	if zipFileName == "" {
		ctx.JSON(404, gin.H{"message": "empty data"})
		return
	}
	ctx.File(zipFileName)
}

func (h *BackupHandler) Recovery(ctx *gin.Context) (any, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, err
	}
	err = h.serv.Recovery(ctx, file)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
