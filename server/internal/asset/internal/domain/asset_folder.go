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

package domain

type AssetFolder struct {
	Id string
	// 文件夹名称
	Name string
	// 文件夹归属的素材类型，image ···
	AssetType string
	// 文件夹类型，post-editor ···
	Type          string
	Assets        []string
	ChildFolders  []*AssetFolder
	SupportDelete bool
	SupportEdit   bool
	SupportAdd    bool
}
