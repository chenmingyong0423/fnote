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
	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/service"
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	"github.com/chenmingyong0423/gkit"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func NewAssetHandler(assetFolderServ service.IAssetService) *AssetHandler {
	return &AssetHandler{
		assetServ: assetFolderServ,
	}
}

type AssetHandler struct {
	assetServ service.IAssetService
}

func (h *AssetHandler) RegisterGinRoutes(engine *gin.Engine) {
	group := engine.Group("/admin-api/assets/folders")
	baseFolderIdGroup := group.Group("/:folderId")

	group.GET("", apiwrap.Wrap(h.GetAssetFolders))
	group.POST("", apiwrap.WrapWithBody(h.AddAssetFolder))
	baseFolderIdGroup.PUT("", apiwrap.WrapWithBody(h.ModifyAssetFolder))
	baseFolderIdGroup.PUT("/name", apiwrap.WrapWithBody(h.ModifyAssetFolderName))
	baseFolderIdGroup.DELETE("", apiwrap.Wrap(h.DeleteAssetFolder))

	// 子文件夹 API
	baseFolderIdGroup.POST("/subfolders", apiwrap.WrapWithBody(h.AddSubAssetFolder))
	baseFolderIdGroup.PUT("/subfolders/:subId", apiwrap.WrapWithBody(h.ModifySubAssetFolder))
	baseFolderIdGroup.DELETE("/subfolders/:subId", apiwrap.Wrap(h.DeleteSubAssetFolder))

	// 文件 API
	baseFolderIdGroup.GET("/assets", apiwrap.Wrap(h.GetAssetsByFolderId))
	baseFolderIdGroup.POST("/assets", apiwrap.WrapWithBody(h.AddAsset))
	// 删除文件
	baseFolderIdGroup.DELETE("/assets/:assetId", apiwrap.Wrap(h.DeleteAsset))
}

func (h *AssetHandler) GetAssetFolders(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[*AssetFolderVO]], error) {
	assertType := ctx.Query("assetType")
	if assertType == "" {
		return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "assetType is required")
	}
	typ := ctx.Query("type")
	if typ == "" {
		return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "type is required")
	}
	assetFolders, err := h.assetServ.GetFoldersByAssetTypeAndType(ctx, assertType, typ)
	if err != nil {
		return nil, err
	}
	if assetFolders == nil {
		return nil, apiwrap.NewErrorResponseBody(http.StatusNotFound, "asset folder not found")
	}

	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(h.toVOs(assetFolders))), nil
}

func (h *AssetHandler) toVOs(assetFolders []*domain.AssetFolder) []*AssetFolderVO {
	result := make([]*AssetFolderVO, 0, len(assetFolders))

	for _, assetFolder := range assetFolders {
		assetFolderVO := toVO(assetFolder)
		assetFolderVO.Children = h.toVOs(assetFolder.ChildFolders)
		result = append(result, assetFolderVO)
	}

	return result
}

func (h *AssetHandler) AddAssetFolder(ctx *gin.Context, req AssetFolderRequest) (*apiwrap.ResponseBody[any], error) {
	_, err := h.assetServ.AddFolder(ctx, &domain.AssetFolder{
		Name:          req.Name,
		AssetType:     req.AssetType,
		Type:          req.Type,
		SupportDelete: gkit.GetValueOrDefault(req.SupportDelete),
		SupportEdit:   gkit.GetValueOrDefault(req.SupportEdit),
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, apiwrap.NewErrorResponseBody(http.StatusConflict, "asset folder already exists")
		}
		return nil, err
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *AssetHandler) ModifyAssetFolder(ctx *gin.Context, req AssetFolderRequest) (*apiwrap.ResponseBody[any], error) {
	folderId := ctx.Param("folderId")
	modifyCnt, err := h.assetServ.ModifyFolderById(ctx, &domain.AssetFolder{
		Id:            folderId,
		Name:          req.Name,
		AssetType:     req.AssetType,
		Type:          req.Type,
		SupportDelete: gkit.GetValueOrDefault(req.SupportDelete),
		SupportEdit:   gkit.GetValueOrDefault(req.SupportEdit),
	})
	if err != nil {
		return nil, err
	}
	if modifyCnt == 0 {
		return nil, apiwrap.NewErrorResponseBody(http.StatusNotFound, "asset folder not found")
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *AssetHandler) DeleteAssetFolder(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	folderId := ctx.Param("folderId")
	deleteCnt, err := h.assetServ.DeleteFolderById(ctx, folderId)
	if err != nil {
		return nil, err
	}
	if deleteCnt == 0 {
		return nil, apiwrap.NewErrorResponseBody(http.StatusNotFound, "asset folder not found")
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *AssetHandler) AddSubAssetFolder(ctx *gin.Context, req AssetFolderRequest) (*apiwrap.ResponseBody[any], error) {
	folderId := ctx.Param("folderId")
	modifyCnt, _, err := h.assetServ.AddSubFolder(ctx, folderId, &domain.AssetFolder{
		Name:          req.Name,
		AssetType:     req.Type,
		Type:          req.Type,
		SupportDelete: gkit.GetValueOrDefault(req.SupportDelete),
		SupportEdit:   gkit.GetValueOrDefault(req.SupportEdit),
		SupportAdd:    gkit.GetValueOrDefault(req.SupportAdd),
	})
	if err != nil {
		return nil, err
	}
	if modifyCnt == 0 {
		return nil, apiwrap.NewErrorResponseBody(http.StatusNotFound, "asset folder not found")
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *AssetHandler) ModifySubAssetFolder(ctx *gin.Context, req AssetFolderRequest) (*apiwrap.ResponseBody[any], error) {
	folderId := ctx.Param("folderId")
	subId := ctx.Param("subId")
	modifyCnt, err := h.assetServ.ModifySubFolderById(ctx, folderId, &domain.AssetFolder{
		Id:            subId,
		Name:          req.Name,
		AssetType:     req.Type,
		Type:          req.Type,
		SupportDelete: gkit.GetValueOrDefault(req.SupportDelete),
		SupportEdit:   gkit.GetValueOrDefault(req.SupportEdit),
		SupportAdd:    gkit.GetValueOrDefault(req.SupportAdd),
	})
	if err != nil {
		return nil, err
	}
	if modifyCnt == 0 {
		return nil, apiwrap.NewErrorResponseBody(http.StatusNotFound, "asset folder not found")
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *AssetHandler) DeleteSubAssetFolder(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	folderId := ctx.Param("folderId")
	subId := ctx.Param("subId")
	deleteCnt, err := h.assetServ.DeleteSubFolderById(ctx, folderId, subId)
	if err != nil {
		return nil, err
	}
	if deleteCnt == 0 {
		return nil, apiwrap.NewErrorResponseBody(http.StatusNotFound, "asset folder not found")
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *AssetHandler) ModifyAssetFolderName(ctx *gin.Context, req ModifyFolderNameRequest) (*apiwrap.ResponseBody[any], error) {
	folderId := ctx.Param("folderId")
	modifyCnt, err := h.assetServ.ModifyFolderNameById(ctx, folderId, req.Name)
	if err != nil {
		return nil, err
	}
	if modifyCnt == 0 {
		return nil, apiwrap.NewErrorResponseBody(http.StatusNotFound, "asset folder not found")
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *AssetHandler) GetAssetsByFolderId(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[AssetVO]], error) {
	folderId := ctx.Param("folderId")
	assets, err := h.assetServ.GetAssetsByIDs(ctx, folderId)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData[apiwrap.ListVO[AssetVO]](apiwrap.NewListVO[AssetVO](h.toAssetVOs(assets))), nil
}

func (h *AssetHandler) toAssetVOs(assets []*domain.Asset) []AssetVO {
	assetVOs := make([]AssetVO, 0, len(assets))
	for _, asset := range assets {
		assetVOs = append(assetVOs, toAssetVO(asset))
	}
	return assetVOs
}

func (h *AssetHandler) AddAsset(ctx *gin.Context, req PostAssetRequest) (*apiwrap.ResponseBody[any], error) {
	folderId := ctx.Param("folderId")
	_, err := h.assetServ.AddAsset(ctx, folderId, &domain.Asset{
		Title:       req.Title,
		Content:     req.Content,
		Description: req.Description,
		AssetType:   req.AssetType,
		Type:        req.Type,
		Metadata:    req.Metadata,
	})
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *AssetHandler) DeleteAsset(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	folderId := ctx.Param("folderId")
	assetId := ctx.Param("assetId")
	err := h.assetServ.DeleteAsset(ctx, folderId, assetId)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponse(), nil
}

func toAssetVO(asset *domain.Asset) AssetVO {
	return AssetVO{
		Id:          asset.Id,
		Title:       asset.Title,
		Content:     asset.Content,
		Description: asset.Description,
		AssetType:   asset.AssetType,
		Type:        asset.Type,
		Metadata:    asset.Metadata,
	}
}

func toVO(assetFolder *domain.AssetFolder) *AssetFolderVO {
	return &AssetFolderVO{
		Id:            assetFolder.Id,
		Name:          assetFolder.Name,
		SupportDelete: assetFolder.SupportDelete,
		SupportEdit:   assetFolder.SupportEdit,
		SupportAdd:    assetFolder.SupportAdd,
	}
}
