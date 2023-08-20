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

package repository

import (
	"context"

	"github.com/chenmingyong0423/fnote/backend/ineternal/config/repository/dao"
	"github.com/chenmingyong0423/fnote/backend/ineternal/domain"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type IConfigRepository interface {
	FindByTyp(ctx context.Context, typ string) (*domain.WebMasterConfig, error)
}

func NewConfigRepository(dao dao.IConfigDao) *ConfigRepository {
	return &ConfigRepository{
		dao: dao,
	}
}

var _ IConfigRepository = (*ConfigRepository)(nil)

type ConfigRepository struct {
	dao dao.IConfigDao
}

func (r *ConfigRepository) FindByTyp(ctx context.Context, typ string) (*domain.WebMasterConfig, error) {
	config, err := r.dao.FindByTyp(ctx, typ)
	if err != nil {
		return nil, errors.WithMessage(err, "r.dao.FindByTyp failed")
	}
	marshal, err := bson.Marshal(config.Props)
	if err != nil {
		return nil, err
	}
	webMasterConfig := &domain.WebMasterConfig{}
	err = bson.Unmarshal(marshal, webMasterConfig)
	if err != nil {
		return nil, err
	}
	return webMasterConfig, nil
}
