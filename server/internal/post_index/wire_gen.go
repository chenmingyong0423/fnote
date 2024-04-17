// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package post_index

import (
	service2 "github.com/chenmingyong0423/fnote/server/internal/post_index/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/post_index/internal/web"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitPostIndexModule(cfgServ website_config.Service) Model {
	baiduService := service2.NewBaiduService()
	postIndexService := service2.NewPostIndexService(baiduService, cfgServ)
	postIndexHandler := web.NewPostIndexHandler(postIndexService)
	model := Model{
		Svc: postIndexService,
		Hdl: postIndexHandler,
	}
	return model
}

// wire.go:

var ConfigProviders = wire.NewSet(web.NewPostIndexHandler, service2.NewPostIndexService, service2.NewBaiduService, wire.Bind(new(service2.IPostIndexService), new(*service2.PostIndexService)))