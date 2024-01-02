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

package service

import (
	"context"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/server/internal/visit_log/repository"
	"github.com/pkg/errors"
)

type IVisitLogService interface {
	CollectVisitLog(ctx context.Context, visitHistory domain.VisitHistory) error
}

var _ IVisitLogService = (*VisitLogService)(nil)

type VisitLogService struct {
	repo repository.IVisitLogRepository
}

func (s *VisitLogService) CollectVisitLog(ctx context.Context, visitHistory domain.VisitHistory) error {
	err := s.repo.Add(ctx, visitHistory)
	if err != nil {
		return errors.WithMessage(err, "s.repo.Add failed")
	}
	return nil
}

func NewVisitLogService(repo repository.IVisitLogRepository) *VisitLogService {
	return &VisitLogService{repo: repo}
}
