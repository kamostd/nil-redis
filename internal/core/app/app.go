package app

import (
	"context"
	"errors"
	"net/http"
	"nil-redis/internal/core/client"
	"nil-redis/internal/core/logger"
	"nil-redis/internal/core/server"
	"nil-redis/internal/core/swagger"
	txmanager "nil-redis/internal/core/tools/tx_manager"
	repository "nil-redis/internal/features/game/storage"
	"nil-redis/internal/features/game/transport"
	"nil-redis/internal/features/game/usecase"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func provideListClient(client *redis.Client) repository.ListClient               { return client }
func provideClientLeaderboard(client *redis.Client) repository.ClientLeaderboard { return client }
func provideHmapClient(client *redis.Client) repository.HmapClient               { return client }
func provideHincrClient(client *redis.Client) repository.HincrClient             { return client }

func provideGameRepo(repo *repository.Player) usecase.GameRepo { return repo }
func provideLeaderboardSaverGeter(repo *repository.Leaderboard) usecase.LeaderboardSaverGeter {
	return repo
}
func provideHistorySaver(repo *repository.HistoryRepo) usecase.HistorySaver { return repo }
func provideStatsUpdater(repo *repository.StatsRepo) usecase.StatsUpdater   { return repo }

func provideHistoryLister(repo *repository.HistoryRepo) usecase.HistoryLister         { return repo }
func provideLeaderboardLister(repo *repository.Leaderboard) usecase.LeaderboardLister { return repo }
func providePlayerGeter(repo *repository.Player) usecase.PlayerGeter                  { return repo }
func provideRankScoreGeter(repo *repository.Leaderboard) usecase.RankScoreGeter       { return repo }
func provideStatsGeter(repo *repository.StatsRepo) usecase.StatsGeter                 { return repo }

func provideGameService(service *usecase.GameService) transport.GameService           { return service }
func provideHistoryService(service *usecase.HistoryService) transport.HistoryServiece { return service }
func provideLeaderboardService(service *usecase.LeaderboardService) transport.LeaderboardService {
	return service
}
func provideProfileCollector(service *usecase.ProfileService) transport.ProfileCollector {
	return service
}
func provideStatsService(service *usecase.StatsService) transport.StatsService {
	return service
}

func provideTxManager(trm *txmanager.Trm) usecase.TxManager {
	return trm
}

// func providePipliner(client *redis.Client) txmanager.Pipeliner {
// 	return client
// }

func App() *fx.App {
	return fx.New(
		logger.NewModule(),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{
				Logger: logger,
			}
		}),

		fx.Provide(
			provideTxManager,
			// providePipliner,
		),

		fx.Provide(
			provideListClient,
			provideClientLeaderboard,
			provideHmapClient,
			provideHincrClient,
		),

		fx.Provide(
			provideGameRepo,
			provideLeaderboardSaverGeter,
			provideHistorySaver,
			provideStatsUpdater,
			provideHistoryLister,
			provideLeaderboardLister,
			providePlayerGeter,
			provideRankScoreGeter,
			provideStatsGeter,
		),

		fx.Provide(
			provideGameService,
			provideHistoryService,
			provideLeaderboardService,
			provideProfileCollector,
			provideStatsService,
		),

		txmanager.NewModule(),
		swagger.NewModule(),
		usecase.NewModule(),
		repository.NewModule(),
		transport.NewModule(),
		client.NewModule(),
		server.NewModule(),

		fx.Invoke(func(lc fx.Lifecycle, server *http.Server, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					go func() {
						if err := server.ListenAndServe(); err != nil {
							if !errors.Is(err, http.ErrServerClosed) {
								logger.Error(
									"server.ListenAndServe() failed",
									zap.Error(err),
								)
							}
						}
					}()

					return nil
				},
				OnStop: func(context.Context) error {
					return server.Close()
				},
			})
		}),
	)
}
