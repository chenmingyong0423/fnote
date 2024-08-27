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

type PostAssetRequest struct {
	// 素材标题
	Title string `json:"title"`
	// 素材内容
	Content string `json:"content" binding:"required"`
	// 素材描述
	Description string `json:"description"`
	// 素材类型，image ···
	AssetType string `json:"asset_type" binding:"required"`
	// 文件夹类型，post-editor ···
	Type string `json:"type" binding:"required"`
	// 元数据
	Metadata map[string]any `json:"metadata"`
}

type AssetFolderRequest struct {
	Name string `json:"name" binding:"required"`
	// 文件夹归属的素材类型，image ···
	AssetType string `json:"asset_type" binding:"required"`
	// 文件夹类型，post-editor ···
	Type          string `json:"type" binding:"required"`
	SupportDelete *bool  `json:"support_delete" binding:"required"`
	SupportEdit   *bool  `json:"support_edit" binding:"required"`
	SupportAdd    *bool  `json:"support_add" binding:"required"`
}

type ModifyFolderNameRequest struct {
	Name string `json:"name" binding:"required"`
}
