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
	"github.com/chenmingyong0423/fnote/backend/ineternal/comment/repository/dao"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/domain"
	"github.com/google/uuid"
	"time"
)

type ICommentRepository interface {
	AddComment(ctx context.Context, comment domain.Comment) (any, error)
}

func NewCommentRepository(dao dao.ICommentDao) *CommentRepository {
	return &CommentRepository{
		dao: dao,
	}
}

var _ ICommentRepository = (*CommentRepository)(nil)

type CommentRepository struct {
	dao dao.ICommentDao
}

func (r *CommentRepository) AddComment(ctx context.Context, comment domain.Comment) (any, error) {
	unix := time.Now().Unix()
	return r.dao.AddComment(ctx, dao.Comment{
		Id:         uuid.NewString(),
		Comment:    comment.Comment,
		Replies:    make([]dao.CommentReply, 0),
		Status:     1,
		CreateTime: unix,
		UpdateTime: unix,
	})
}
