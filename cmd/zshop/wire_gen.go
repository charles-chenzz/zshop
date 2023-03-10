// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"zshop/internal/biz"
	"zshop/internal/conf"
	"zshop/internal/data"
	"zshop/internal/server"
	"zshop/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	db, err := data.NewMysql(confData)
	if err != nil {
		return nil, nil, err
	}
	dataData, cleanup, err := data.NewData(confData, logger, db)
	if err != nil {
		return nil, nil, err
	}
	greeterRepo := data.NewGreeterRepo(dataData, logger)
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
	greeterService := service.NewGreeterService(greeterUsecase)
	orderRepo := data.NewOrderRepo(dataData, logger)
	orderUserCase := biz.NewOrderUserCase(orderRepo, logger)
	orderService := service.NewOrderService(orderUserCase)
	grpcServer := server.NewGRPCServer(confServer, greeterService, orderService, logger)
	httpServer := server.NewHTTPServer(confServer, greeterService, orderService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
