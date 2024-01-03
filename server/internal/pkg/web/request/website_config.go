// Copyright 2023 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package request

type UpdateWebsiteConfigReq struct {
	Name     string `json:"name" binding:"required"`
	Icon     string `json:"icon" binding:"required"`
	LiveTime int64  `json:"live_time" binding:"required"`
}

type UpdateOwnerConfigReq struct {
	Name    string `json:"name" binding:"required"`
	Profile string `json:"profile" binding:"required"`
	Picture string `json:"picture" binding:"required"`
}

type UpdateSeoMetaConfigReq struct {
	Title                 string `json:"title" binding:"required"`
	Description           string `json:"description" binding:"required"`
	OgTitle               string `json:"og_title" binding:"required"`
	Keywords              string `json:"keywords" binding:"required"`
	Author                string `json:"author" binding:"required"`
	Robots                string `json:"robots" binding:"required"`
	OgImage               string `json:"og_image"`
	BaiduSiteVerification string `json:"baidu_site_verification"`
}

type UpdateCommentConfigReq struct {
	EnableComment *bool `json:"enable_comment" binding:"required"`
}
