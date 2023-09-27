// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/chenmingyong0423/fnote/backend/internal/category/handler"
	"github.com/chenmingyong0423/fnote/backend/internal/category/repository"
	"github.com/chenmingyong0423/fnote/backend/internal/category/repository/dao"
	"github.com/chenmingyong0423/fnote/backend/internal/category/service"
	"github.com/chenmingyong0423/fnote/backend/internal/comment/hanlder"
	repository2 "github.com/chenmingyong0423/fnote/backend/internal/comment/repository"
	dao2 "github.com/chenmingyong0423/fnote/backend/internal/comment/repository/dao"
	service2 "github.com/chenmingyong0423/fnote/backend/internal/comment/service"
	handler2 "github.com/chenmingyong0423/fnote/backend/internal/config/handler"
	repository3 "github.com/chenmingyong0423/fnote/backend/internal/config/repository"
	dao3 "github.com/chenmingyong0423/fnote/backend/internal/config/repository/dao"
	service3 "github.com/chenmingyong0423/fnote/backend/internal/config/service"
	service5 "github.com/chenmingyong0423/fnote/backend/internal/email/service"
	hanlder2 "github.com/chenmingyong0423/fnote/backend/internal/friend/hanlder"
	repository5 "github.com/chenmingyong0423/fnote/backend/internal/friend/repository"
	dao5 "github.com/chenmingyong0423/fnote/backend/internal/friend/repository/dao"
	service7 "github.com/chenmingyong0423/fnote/backend/internal/friend/service"
	"github.com/chenmingyong0423/fnote/backend/internal/ioc"
	service6 "github.com/chenmingyong0423/fnote/backend/internal/message/service"
	handler3 "github.com/chenmingyong0423/fnote/backend/internal/post/handler"
	repository4 "github.com/chenmingyong0423/fnote/backend/internal/post/repository"
	dao4 "github.com/chenmingyong0423/fnote/backend/internal/post/repository/dao"
	service4 "github.com/chenmingyong0423/fnote/backend/internal/post/service"
	handler4 "github.com/chenmingyong0423/fnote/backend/internal/visit_log/handler"
	repository6 "github.com/chenmingyong0423/fnote/backend/internal/visit_log/repository"
	dao6 "github.com/chenmingyong0423/fnote/backend/internal/visit_log/repository/dao"
	service8 "github.com/chenmingyong0423/fnote/backend/internal/visit_log/service"
	"github.com/gin-gonic/gin"
)

// Injectors from wire.go:

func initializeApp(username ioc.Username, password ioc.Password) (*gin.Engine, error) {
	database := ioc.NewMongoDB(username, password)
	categoryDao := dao.NewCategoryDao(database)
	categoryRepository := repository.NewCategoryRepository(categoryDao)
	categoryService := service.NewCategoryService(categoryRepository)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	commentDao := dao2.NewCommentDao(database)
	commentRepository := repository2.NewCommentRepository(commentDao)
	commentService := service2.NewCommentService(commentRepository)
	configDao := dao3.NewConfigDao(database)
	configRepository := repository3.NewConfigRepository(configDao)
	configService := service3.NewConfigService(configRepository)
	postDao := dao4.NewPostDao(database)
	postRepository := repository4.NewPostRepository(postDao)
	postService := service4.NewPostService(postRepository)
	emailService := service5.NewEmailService()
	messageService := service6.NewMessageService(configService, emailService)
	commentHandler := hanlder.NewCommentHandler(commentService, configService, postService, messageService)
	configHandler := handler2.NewConfigHandler(configService)
	friendDao := dao5.NewFriendDao(database)
	friendRepository := repository5.NewFriendRepository(friendDao)
	friendService := service7.NewFriendService(friendRepository, emailService, configService)
	friendHandler := hanlder2.NewFriendHandler(friendService)
	postHandler := handler3.NewPostHandler(postService)
	visitLogDao := dao6.NewVisitLogDao(database)
	visitLogRepository := repository6.NewVisitLogRepository(visitLogDao)
	visitLogService := service8.NewVisitLogService(visitLogRepository)
	visitLogHandler := handler4.NewVisitLogHandler(visitLogService, configService)
	engine, err := ioc.NewGinEngine(categoryHandler, commentHandler, configHandler, friendHandler, postHandler, visitLogHandler)
	if err != nil {
		return nil, err
	}
	return engine, nil
}