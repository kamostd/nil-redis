package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type ListClient interface {
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd
}

type GameRepo struct {
	client ListClient
	logger *zap.Logger
}

func NewGameRepo(client ListClient, logger *zap.Logger) *GameRepo {
	return &GameRepo{
		client: client,
		logger: logger,
	}
}

func (g GameRepo) SaveGame(
	ctx context.Context,
	historyNote string,
) error {
	ctxSave, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prepLpush := time.Now()
	if _, err := g.client.LPush(ctxSave, "games", historyNote).Result(); err != nil {
		g.logger.Error(
			"SaveGame failed",
			zap.String("redis request type", "LPush"),
			zap.Duration("duration", time.Since(prepLpush)),
			zap.Error(err),
		)

		return err
	}

	g.logger.Info(
		"SaveGame success",
		zap.String("redis request type", "LPush"),
		zap.Duration("duration", time.Since(prepLpush)),
	)

	ctxTrim, cancelTrim := context.WithTimeout(ctx, time.Second)
	defer cancelTrim()

	prepTirm := time.Now()
	if _, err := g.client.LTrim(ctxTrim, "games", 0, 9).Result(); err != nil {
		g.logger.Error(
			"SaveGame failed",
			zap.String("redis request type", "LTrim"),
			zap.Duration("duration", time.Since(prepTirm)),
			zap.Error(err),
		)

		return err
	}

	g.logger.Info(
		"SaveGame success",
		zap.String("redis request type", "LTrim"),
		zap.Duration("duration", time.Since(prepLpush)),
	)

	return nil
}
