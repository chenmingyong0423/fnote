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
	"github.com/chenmingyong0423/fnote/backend/ineternal/domain"
	"github.com/chenmingyong0423/fnote/backend/ineternal/visit_log/repository/dao"
	"github.com/pkg/errors"
	"time"
)

type IVisitLogRepository interface {
	Add(ctx context.Context, visitHistory domain.VisitHistory) error
}

var _ IVisitLogRepository = (*VisitLogRepository)(nil)

type VisitLogRepository struct {
	dao dao.IVisitLogDao
}

func (r *VisitLogRepository) Add(ctx context.Context, visitHistory domain.VisitHistory) error {
	err := r.dao.Add(ctx, dao.VisitHistory{Url: visitHistory.Url, Ip: visitHistory.Ip, UserAgent: visitHistory.UserAgent, Origin: visitHistory.Origin, Referer: visitHistory.Referer, CreateTime: time.Now().Unix()})
	if err != nil {
		return errors.WithMessage(err, "r.dao.Add failed")
	}
	return nil
}

func NewVisitLogRepository(dao dao.IVisitLogDao) *VisitLogRepository {
	return &VisitLogRepository{dao: dao}
}
