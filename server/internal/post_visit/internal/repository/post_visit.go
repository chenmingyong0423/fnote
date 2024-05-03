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
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/chenmingyong0423/fnote/server/internal/post_visit/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post_visit/internal/repository/dao"
)

type IPostVisitRepository interface {
	Insert(ctx context.Context, postVisit domain.PostVisit) error
}

var _ IPostVisitRepository = (*PostVisitRepository)(nil)

func NewPostVisitRepository(dao dao.IPostVisitDao) *PostVisitRepository {
	return &PostVisitRepository{dao: dao}
}

type PostVisitRepository struct {
	dao dao.IPostVisitDao
}

func (r *PostVisitRepository) Insert(ctx context.Context, postVisit domain.PostVisit) error {
	_, err := r.dao.Insert(ctx, &dao.PostVisit{
		Id:        primitive.NewObjectID(),
		PostId:    postVisit.PostId,
		Ip:        postVisit.Ip,
		UserAgent: postVisit.UserAgent,
		Origin:    postVisit.Origin,
		Referer:   postVisit.Referer,
		StayTime:  postVisit.StayTime,
		VisitAt:   time.UnixMilli(postVisit.VisitAt).Local(),
	})
	if err != nil {
		return err
	}
	return nil
}
