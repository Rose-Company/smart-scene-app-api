package server

import (
	"context"
	"smart-scene-app-api/pkg"
	"smart-scene-app-api/pkg/rest_service"
	"smart-scene-app-api/services/logger"

	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
)

type ServerContext interface {
	GetService(prefix string) interface{}
	GetLogger() logger.Loggers
	GetContext() context.Context
	SetUser(value interface{})
	GetUser() interface{}
	GetLoggerWithPrefix(prefix string) logger.Loggers
	GetRedisRedsync(prefix string) redsync.Redsync
	InitAuthorizationData()
	SetTelegramService(service rest_service.RestInterface)
	GetTelegramService() rest_service.RestInterface
	GetAwsSes() *pkg.AWSSesClient
	SetAwsSes(service *pkg.AWSSesClient)
	DB() *gorm.DB
	Ctx() context.Context
}
