package repository

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {
	return fx.Module(
		"repository",

		fx.Provide(
			NewPlayerRepo,
			NewLeaderboard,
			NewHistoryRepo,
			NewStatsRepo,
		),

		fx.Invoke(func(logger *zap.Logger) *zap.Logger {
			return logger.With(zap.String("module", "repository"))
		}),
	)
}
