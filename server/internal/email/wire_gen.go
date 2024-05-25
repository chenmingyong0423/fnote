// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package email

import (
	"github.com/chenmingyong0423/fnote/server/internal/email/internal/service"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

// Injectors from wire.go:

func InitEmailModule(mongoDB *mongo.Database) *Module {
	emailService := service.NewEmailService()
	module := &Module{
		Svc: emailService,
	}
	return module
}

// wire.go:

var EmailProviders = wire.NewSet(service.NewEmailService, wire.Bind(new(service.IEmailService), new(*service.EmailService)))