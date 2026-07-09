package txmanager

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {
	return fx.Module(
		"tx manager",

		fx.Provide(NewTrm),

		fx.Invoke(func(logger *zap.Logger) *zap.Logger {
			return logger.With(zap.String("module", "tx manager"))
		}),
	)
}
