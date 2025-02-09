// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"muxi-thief/api/route"
	"muxi-thief/controller"
	"muxi-thief/middleware"
)

// Injectors from wire.go:

func WireApp() route.App {
	jwtClient := middleware.NewJWTClient()
	controllerController := controller.NewAuthController(jwtClient)
	middlewareMiddleware := middleware.NewMiddleware(jwtClient)
	engine := route.NewRouter(controllerController, middlewareMiddleware)
	app := route.NewApp(engine)
	return app
}
