package repository

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {
	return fx.Module(
		"player repository",

		fx.Provide(NewPlayerRepo),

		fx.Invoke(func(logger *zap.Logger) *zap.Logger {
			return logger.With(zap.String("module", "player repo"))
		}),
	)
}
