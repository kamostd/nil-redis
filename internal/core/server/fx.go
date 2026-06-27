package server

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"server",

		fx.Provide(
			NewRouterGroup,
			NewGinEngine,
			NewHTTPServer,
		),
	)
}
