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

package domain

type File struct {
	Id               string
	FileId           string
	FileName         string
	OriginalFileName string
	FileType         string
	FileSize         int64
	FilePath         string
	Url              string
	UsedIn           []FileUsage
	CreatedAt        int64
	UpdatedAt        int64
}

type FileUsage struct {
	EntityId   string
	EntityType string
}

type FileDTO struct {
	FileName       string `json:"file_name"`
	FileSize       int64  `json:"file_size"`
	Content        []byte `json:"content"`
	FileType       string `json:"file_type"`
	FileExt        string `json:"file_ext"`
	CustomFileName string
}

type PageDTO struct {
	PageNum  int64
	PageSize int64
	FileType []string
}
