package usecase

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {
	return fx.Module(
		"service",

		fx.Provide(
			NewGameService,
			NewHistoryService,
			NewLeaderboardService,
			NewProfileService,
			NewStatsService,
		),

		fx.Invoke(func(logger *zap.Logger) *zap.Logger {
			return logger.With(zap.String("module", "service"))
		}),
	)
}
