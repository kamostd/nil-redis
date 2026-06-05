package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type HincrClient interface {
	HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd
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
	if _, err := s.client.HIncrBy(ctxIncr, "stats", "total_games", 1).Result(); err != nil {
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
