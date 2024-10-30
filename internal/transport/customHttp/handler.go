package customHttp

import (
	"SimpleShop/internal/service"
	"SimpleShop/pkg/logger"
	"log"
)

type HandlerHttp struct {
	Service  service.HttpModule
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	DebugLog *log.Logger
}

func NewTransportHttpHandler(ServiceObject service.HttpModule, logger *logger.CustomLogger) *HandlerHttp {
	handlerObject := &HandlerHttp{
		Service:  ServiceObject,
		ErrorLog: logger.ErrorLogger,
		InfoLog:  logger.InfoLogger,
		DebugLog: logger.DebugLogger,
	}

	return handlerObject
}
