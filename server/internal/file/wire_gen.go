// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package file

import (
	"github.com/chenmingyong0423/fnote/server/internal/file/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/file/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/file/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/file/internal/web"
	"github.com/chenmingyong0423/go-eventbus"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

// Injectors from wire.go:

func InitFileModule(mongoDB *mongo.Database, eventBus *eventbus.EventBus) *Module {
	fileDao := dao.NewFileDao(mongoDB)
	fileRepository := repository.NewFileRepository(fileDao)
	fileService := service.NewFileService(fileRepository, eventBus)
	fileHandler := web.NewFileHandler(fileService)
	module := &Module{
		Svc: fileService,
		Hdl: fileHandler,
	}
	return module
}

// wire.go:

var FileProviders = wire.NewSet(web.NewFileHandler, service.NewFileService, repository.NewFileRepository, dao.NewFileDao, wire.Bind(new(service.IFileService), new(*service.FileService)), wire.Bind(new(repository.IFileRepository), new(*repository.FileRepository)), wire.Bind(new(dao.IFileDao), new(*dao.FileDao)))