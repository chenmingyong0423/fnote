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
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/visit_log/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/visit_log/internal/repository/dao"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type IVisitLogRepository interface {
	Add(ctx context.Context, visitHistory domain.VisitHistory) error
	CountOfToday(ctx context.Context) (int64, error)
	CountOfTodayByIp(ctx context.Context) (int64, error)
	GetViewTendencyStats4PV(ctx context.Context, days int) ([]domain.TendencyData, error)
	GetViewTendencyStats4UV(ctx context.Context, days int) ([]domain.TendencyData, error)
	GetByDate(ctx context.Context, start time.Time, end time.Time) ([]domain.VisitHistory, error)
}

var _ IVisitLogRepository = (*VisitLogRepository)(nil)

type VisitLogRepository struct {
	dao dao.IVisitLogDao
}

func (r *VisitLogRepository) GetByDate(ctx context.Context, start time.Time, end time.Time) ([]domain.VisitHistory, error) {
	visitHistories, err := r.dao.GetByDate(ctx, start, end)
	if err != nil {
		return nil, err
	}
	return r.toDomains(visitHistories), nil
}

func (r *VisitLogRepository) GetViewTendencyStats4UV(ctx context.Context, days int) ([]domain.TendencyData, error) {
	tendencyData, err := r.dao.GetViewTendencyStats4UV(ctx, days)
	if err != nil {
		return nil, err
	}
	return r.tdToDomain(tendencyData), nil
}

func (r *VisitLogRepository) GetViewTendencyStats4PV(ctx context.Context, days int) ([]domain.TendencyData, error) {
	tendencyData, err := r.dao.GetViewTendencyStats4PV(ctx, days)
	if err != nil {
		return nil, err
	}
	return r.tdToDomain(tendencyData), nil
}

func (r *VisitLogRepository) CountOfTodayByIp(ctx context.Context) (int64, error) {
	return r.dao.CountOfTodayByIp(ctx)
}

func (r *VisitLogRepository) CountOfToday(ctx context.Context) (int64, error) {
	return r.dao.CountOfToday(ctx)
}

func (r *VisitLogRepository) Add(ctx context.Context, visitHistory domain.VisitHistory) error {
	err := r.dao.Add(ctx, &dao.VisitHistory{Id: uuid.NewString(), Url: visitHistory.Url, Ip: visitHistory.Ip, UserAgent: visitHistory.UserAgent, Origin: visitHistory.Origin, Referer: visitHistory.Referer, CreatedAt: time.Now().Local()})
	if err != nil {
		return errors.WithMessage(err, "r.dao.Add failed")
	}
	return nil
}

func (r *VisitLogRepository) tdToDomain(data []*dao.TendencyData) []domain.TendencyData {
	var result []domain.TendencyData
	for _, d := range data {
		result = append(result, domain.TendencyData{Timestamp: d.Date.Unix(), ViewCount: d.ViewCount})
	}
	return result
}

func (r *VisitLogRepository) toDomains(visitHistories []*dao.VisitHistory) []domain.VisitHistory {
	result := make([]domain.VisitHistory, 0, len(visitHistories))
	for _, vh := range visitHistories {
		result = append(result, r.toDomain(vh))
	}
	return result
}

func (r *VisitLogRepository) toDomain(vh *dao.VisitHistory) domain.VisitHistory {
	return domain.VisitHistory{
		Url:       vh.Url,
		Ip:        vh.Ip,
		UserAgent: vh.UserAgent,
		Origin:    vh.Origin,
		Type:      vh.UserAgent,
		Referer:   vh.Referer,
	}
}

func NewVisitLogRepository(dao dao.IVisitLogDao) *VisitLogRepository {
	return &VisitLogRepository{dao: dao}
}
