// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	handler2 "github.com/chenmingyong0423/fnote/server/internal/category/handler"
	repository2 "github.com/chenmingyong0423/fnote/server/internal/category/repository"
	dao2 "github.com/chenmingyong0423/fnote/server/internal/category/repository/dao"
	service4 "github.com/chenmingyong0423/fnote/server/internal/category/service"
	"github.com/chenmingyong0423/fnote/server/internal/comment/hanlder"
	repository5 "github.com/chenmingyong0423/fnote/server/internal/comment/repository"
	dao5 "github.com/chenmingyong0423/fnote/server/internal/comment/repository/dao"
	service5 "github.com/chenmingyong0423/fnote/server/internal/comment/service"
	repository3 "github.com/chenmingyong0423/fnote/server/internal/count_stats/repository"
	dao3 "github.com/chenmingyong0423/fnote/server/internal/count_stats/repository/dao"
	service2 "github.com/chenmingyong0423/fnote/server/internal/count_stats/service"
	handler8 "github.com/chenmingyong0423/fnote/server/internal/data_analysis/handler"
	service7 "github.com/chenmingyong0423/fnote/server/internal/email/service"
	"github.com/chenmingyong0423/fnote/server/internal/file/handler"
	"github.com/chenmingyong0423/fnote/server/internal/file/repository"
	"github.com/chenmingyong0423/fnote/server/internal/file/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/file/service"
	hanlder2 "github.com/chenmingyong0423/fnote/server/internal/friend/hanlder"
	repository8 "github.com/chenmingyong0423/fnote/server/internal/friend/repository"
	dao8 "github.com/chenmingyong0423/fnote/server/internal/friend/repository/dao"
	service10 "github.com/chenmingyong0423/fnote/server/internal/friend/service"
	"github.com/chenmingyong0423/fnote/server/internal/ioc"
	service9 "github.com/chenmingyong0423/fnote/server/internal/message/service"
	handler6 "github.com/chenmingyong0423/fnote/server/internal/message_template/handler"
	repository7 "github.com/chenmingyong0423/fnote/server/internal/message_template/repository"
	dao7 "github.com/chenmingyong0423/fnote/server/internal/message_template/repository/dao"
	service8 "github.com/chenmingyong0423/fnote/server/internal/message_template/service"
	handler4 "github.com/chenmingyong0423/fnote/server/internal/post/handler"
	repository6 "github.com/chenmingyong0423/fnote/server/internal/post/repository"
	dao6 "github.com/chenmingyong0423/fnote/server/internal/post/repository/dao"
	service6 "github.com/chenmingyong0423/fnote/server/internal/post/service"
	handler7 "github.com/chenmingyong0423/fnote/server/internal/tag/handler"
	repository10 "github.com/chenmingyong0423/fnote/server/internal/tag/repository"
	dao10 "github.com/chenmingyong0423/fnote/server/internal/tag/repository/dao"
	service12 "github.com/chenmingyong0423/fnote/server/internal/tag/service"
	handler5 "github.com/chenmingyong0423/fnote/server/internal/visit_log/handler"
	repository9 "github.com/chenmingyong0423/fnote/server/internal/visit_log/repository"
	dao9 "github.com/chenmingyong0423/fnote/server/internal/visit_log/repository/dao"
	service11 "github.com/chenmingyong0423/fnote/server/internal/visit_log/service"
	handler3 "github.com/chenmingyong0423/fnote/server/internal/website_config/handler"
	repository4 "github.com/chenmingyong0423/fnote/server/internal/website_config/repository"
	dao4 "github.com/chenmingyong0423/fnote/server/internal/website_config/repository/dao"
	service3 "github.com/chenmingyong0423/fnote/server/internal/website_config/service"
	"github.com/gin-gonic/gin"
)

// Injectors from wire.go:

func initializeApp() (*gin.Engine, error) {
	database := ioc.NewMongoDB()
	fileDao := dao.NewFileDao(database)
	fileRepository := repository.NewFileRepository(fileDao)
	fileService := service.NewFileService(fileRepository)
	fileHandler := handler.NewFileHandler(fileService)
	categoryDao := dao2.NewCategoryDao(database)
	categoryRepository := repository2.NewCategoryRepository(categoryDao)
	countStatsDao := dao3.NewCountStatsDao(database)
	countStatsRepository := repository3.NewCountStatsRepository(countStatsDao)
	countStatsService := service2.NewCountStatsService(countStatsRepository)
	websiteConfigDao := dao4.NewWebsiteConfigDao(database)
	websiteConfigRepository := repository4.NewWebsiteConfigRepository(websiteConfigDao)
	websiteConfigService := service3.NewWebsiteConfigService(websiteConfigRepository)
	categoryService := service4.NewCategoryService(categoryRepository, countStatsService, websiteConfigService)
	categoryHandler := handler2.NewCategoryHandler(categoryService)
	commentDao := dao5.NewCommentDao(database)
	commentRepository := repository5.NewCommentRepository(commentDao)
	commentService := service5.NewCommentService(commentRepository)
	postDao := dao6.NewPostDao(database)
	postRepository := repository6.NewPostRepository(postDao)
	postService := service6.NewPostService(postRepository, websiteConfigService, countStatsService, fileService)
	emailService := service7.NewEmailService()
	msgTplDao := dao7.NewMsgTplDao(database)
	msgTplRepository := repository7.NewMsgTplRepository(msgTplDao)
	msgTplService := service8.NewMsgTplService(msgTplRepository)
	messageService := service9.NewMessageService(websiteConfigService, emailService, msgTplService)
	commentHandler := hanlder.NewCommentHandler(commentService, websiteConfigService, postService, messageService)
	websiteConfigHandler := handler3.NewWebsiteConfigHandler(websiteConfigService)
	friendDao := dao8.NewFriendDao(database)
	friendRepository := repository8.NewFriendRepository(friendDao)
	friendService := service10.NewFriendService(friendRepository)
	friendHandler := hanlder2.NewFriendHandler(friendService, messageService, websiteConfigService)
	postHandler := handler4.NewPostHandler(postService, websiteConfigService)
	visitLogDao := dao9.NewVisitLogDao(database)
	visitLogRepository := repository9.NewVisitLogRepository(visitLogDao)
	visitLogService := service11.NewVisitLogService(visitLogRepository)
	visitLogHandler := handler5.NewVisitLogHandler(visitLogService, websiteConfigService)
	msgTplHandler := handler6.NewMsgTplHandler(msgTplService)
	tagDao := dao10.NewTagDao(database)
	tagRepository := repository10.NewTagRepository(tagDao)
	tagService := service12.NewTagService(tagRepository, countStatsService)
	tagHandler := handler7.NewTagHandler(tagService)
	dataAnalysisHandler := handler8.NewDataAnalysisHandler(visitLogService, websiteConfigService)
	writer := ioc.InitLogger()
	v := ioc.InitMiddlewares(writer)
	validators := ioc.InitGinValidators()
	engine, err := ioc.NewGinEngine(fileHandler, categoryHandler, commentHandler, websiteConfigHandler, friendHandler, postHandler, visitLogHandler, msgTplHandler, tagHandler, dataAnalysisHandler, v, validators)
	if err != nil {
		return nil, err
	}
	return engine, nil
}
