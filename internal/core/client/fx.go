package client

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"client",

		fx.Provide(NewRedis),
	)
}
