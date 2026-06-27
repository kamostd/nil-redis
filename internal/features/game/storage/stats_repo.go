package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const statsTotalGamesKeyName = "stats:total_games"

type HincrClient interface {
	Incr(ctx context.Context, key string) *redis.IntCmd

	Get(ctx context.Context, key string) *redis.StringCmd
}

type StatsRepo struct {
	client HincrClient
	logger *zap.Logger
}

func NewStatsRepo(client HincrClient, logger *zap.Logger) *StatsRepo {
	return &StatsRepo{
		client: client,
		logger: logger,
	}
}

func (s StatsRepo) UpdateTotalGame(ctx context.Context) error {
	ctxIncr, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()
	if _, err := s.client.Incr(ctxIncr, statsTotalGamesKeyName).Result(); err != nil {
		s.logger.Error(
			"UpdateTotalGame failed",
			zap.String("redis request type", "HIncrBy"),
			zap.Duration("duration", time.Since(prep)),
			zap.Error(err),
		)

		return err
	}

	s.logger.Info(
		"UpdateTotalGame success",
		zap.String("redis request type", "HIncrBy"),
		zap.Duration("duration", time.Since(prep)),
	)

	return nil

}

func (s StatsRepo) GetTotalGames(ctx context.Context) (int, error) {
	ctxIncr, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()
	gamesS, err := s.client.Get(ctxIncr, statsTotalGamesKeyName).Result()
	if err != nil {
		s.logger.Error(
			"FetchTotalGames failed",
			zap.String("redis request type", "Get"),
			zap.Duration("duration", time.Since(prep)),
			zap.Error(err),
		)

		return 0, err
	}

	s.logger.Info(
		"FetchTotalGames success",
		zap.String("redis request type", "Get"),
		zap.Duration("duration", time.Since(prep)),
	)

	totalGames, err := strconv.Atoi(gamesS)
	if err != nil {
		s.logger.Error(
			"FetchTotalGames success",
			zap.String("convert", "total_games"),
			zap.Duration("duration", time.Since(prep)),
		)

		return 0, err
	}

	return totalGames, nil
}
