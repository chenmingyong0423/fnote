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
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	"github.com/chenmingyong0423/fnote/server/internal/post_asset/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post_asset/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewPostAssetHandler(serv service.IPostAssetService, assetFolderServ service.IAssetFolderService) *PostAssetHandler {
	return &PostAssetHandler{
		serv:            serv,
		assetFolderServ: assetFolderServ,
	}
}

type PostAssetHandler struct {
	serv            service.IPostAssetService
	assetFolderServ service.IAssetFolderService
}

func (h *PostAssetHandler) RegisterGinRoutes(engine *gin.Engine) {
	group := engine.Group("/admin-api/assets")
	group.GET("/folders", apiwrap.Wrap(h.GetFolders))

}

func (h *PostAssetHandler) GetFolders(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[AssetFolderVO]], error) {
	assertType := ctx.Query("assetType")
	if assertType == "" {
		return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "assetType is required")
	}
	typ := ctx.Query("type")
	if typ == "" {
		return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "type is required")
	}
	assetFolder, err := h.assetFolderServ.GetFoldersByAssetTypeAndType(ctx, assertType, typ)
	if err != nil {
		return nil, err
	}
	if assetFolder == nil {
		return nil, apiwrap.NewErrorResponseBody(http.StatusNotFound, "asset folder not found")
	}

	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(h.toVOs(assetFolder))), nil
}

func (h *PostAssetHandler) toVOs(assetFolder *domain.AssetFolder) []AssetFolderVO {
	var result []AssetFolderVO

	result = append(result, toVO(assetFolder))

	for _, child := range assetFolder.ChildFolders {
		result = append(result, h.toVOs(&child)...)
	}

	return result
}

func toVO(assetFolder *domain.AssetFolder) AssetFolderVO {
	return AssetFolderVO{
		Id:            assetFolder.Id,
		Name:          assetFolder.Name,
		SupportDelete: assetFolder.SupportDelete,
		SupportEdit:   assetFolder.SupportEdit,
	}
}
