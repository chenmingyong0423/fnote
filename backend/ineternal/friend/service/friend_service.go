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
	"fmt"
	configServ "github.com/chenmingyong0423/fnote/backend/ineternal/config/service"
	emailConfig "github.com/chenmingyong0423/fnote/backend/ineternal/email/service"
	"github.com/chenmingyong0423/fnote/backend/ineternal/friend/repository"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/domain"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
	"net/http"
)

type IFriendService interface {
	GetFriends(ctx context.Context) (api.ListVO[*domain.FriendVO], error)
	ApplyForFriend(ctx context.Context, friend domain.Friend) error
}

var _ IFriendService = (*FriendService)(nil)

func NewFriendService(repo repository.IFriendRepository, emailServ emailConfig.IEmailService, configServ configServ.IConfigService) *FriendService {
	return &FriendService{
		repo:       repo,
		emailServ:  emailServ,
		configServ: configServ,
	}
}

type FriendService struct {
	repo       repository.IFriendRepository
	emailServ  emailConfig.IEmailService
	configServ configServ.IConfigService
}

func (s *FriendService) ApplyForFriend(ctx context.Context, friend domain.Friend) error {
	switchConfig, err := s.configServ.GetSwitchStatusByTyp(ctx, "friend")
	if err != nil {
		return err
	}
	if !switchConfig.Status {
		return api.NewHttpCodeError(http.StatusForbidden)
	}
	if f, err := s.repo.FindByUrl(ctx, friend.Url); !errors.Is(err, mongo.ErrNoDocuments) {
		if err != nil {
			return err
		}
		return fmt.Errorf("the friend had already applied for, friend=%v", f)
	}
	err = s.repo.Add(ctx, friend)
	if err != nil {
		return errors.WithMessage(err, "s.repo.Add failed")
	}
	// 发送邮件
	go func() {
		emailCfg, gErr := s.configServ.GetEmailConfig(ctx)
		if gErr != nil {
			slog.ErrorContext(ctx, "emailConfig", gErr)
			slog.ErrorContext(ctx, "friend", "Fails to send email message.")
			return
		}
		webNMasterCfg, gErr := s.configServ.GetWebmasterInfo(ctx, "webmaster")
		if gErr != nil {
			slog.ErrorContext(ctx, "webNMasterCfg", gErr)
			slog.ErrorContext(ctx, "friend", "Fails to send email message.")
			return
		}
		// todo 后面标题内容弄成动态的形式
		gErr = s.emailServ.SendEmail(ctx, domain.Email{
			Host:        emailCfg.Host,
			Port:        emailCfg.Port,
			Account:     emailCfg.Account,
			Password:    emailCfg.Password,
			Name:        webNMasterCfg.Name,
			To:          []string{friend.Email},
			Subject:     "友链申请通知",
			Body:        fmt.Sprintf("您好，您在《%s》网站中提交的友链申请已通过，详情请前往<a href='https://%s/friends'>友链</a>进行查看。", webNMasterCfg.Name, webNMasterCfg.Domain),
			ContentType: "text/plain",
		})
		if err != nil {
			slog.ErrorContext(ctx, "friend", errors.WithMessage(gErr, "Fails to send email message."))
		}
	}()
	return nil
}

func (s *FriendService) GetFriends(ctx context.Context) (api.ListVO[*domain.FriendVO], error) {
	vo := api.ListVO[*domain.FriendVO]{}
	friends, err := s.repo.FindDisplaying(ctx)
	if err != nil {
		return vo, errors.WithMessage(err, "s.repo.FindDisplaying failed")
	}
	vo.List = s.toFriendVOs(friends)
	return vo, nil
}

func (s *FriendService) toFriendVOs(friends []*domain.Friend) []*domain.FriendVO {
	result := make([]*domain.FriendVO, 0, len(friends))
	for _, friend := range friends {
		result = append(result, s.toFriendVO(friend))
	}
	return result
}

func (s *FriendService) toFriendVO(friend *domain.Friend) *domain.FriendVO {
	return &domain.FriendVO{
		Name:        friend.Name,
		Url:         friend.Url,
		Logo:        friend.Logo,
		Description: friend.Description,
		Status:      friend.Status,
		Priority:    friend.Priority,
	}
}
