//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"muxi-thief/api/route"
	"muxi-thief/controller"
	"muxi-thief/middleware"
)

func WireApp() route.App {
	panic(wire.Build(
		route.ProviderSet,
		controller.ProviderSet,
		middleware.ProviderSet,
		wire.Bind(new(controller.GenerateJWTer), new(*middleware.JWTClient)),
		wire.Bind(new(middleware.ParTokener), new(*middleware.JWTClient)),
		wire.Bind(new(route.ControllerProxy), new(*controller.Controller)),
	))
}
