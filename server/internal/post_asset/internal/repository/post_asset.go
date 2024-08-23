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

package repository

import (
	"github.com/chenmingyong0423/fnote/server/internal/post_asset/internal/repository/dao"
)

type IPostAssetRepository interface {
}

var _ IPostAssetRepository = (*PostAssetRepository)(nil)

func NewPostAssetRepository(dao dao.IPostAssetDao) *PostAssetRepository {
	return &PostAssetRepository{dao: dao}
}

type PostAssetRepository struct {
	dao dao.IPostAssetDao
}
