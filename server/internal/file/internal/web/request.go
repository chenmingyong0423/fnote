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

type FileRequest struct{}

type BatchDeleteFileRequest struct {
	FileIds []string `json:"file_ids" binding:"required,min=1,max=100,dive,required"`
}

type PageRequest struct {
	PageNum  int64    `form:"pageNum" binding:"required,min=1"`
	PageSize int64    `form:"pageSize" binding:"required,min=1,max=100"`
	FileType []string `form:"fileType"`
	Unused   bool     `form:"unused"`
}
