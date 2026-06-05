package repository

import (
	"context"
	"nil-redis/internal/features/game/storage/model"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type HmapClient interface {
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
	HGetDel(ctx context.Context, key string, fields ...string) *redis.StringSliceCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
}

type Player struct {
	client HmapClient
	logger *zap.Logger
}

func NewPlayerRepo(client HmapClient, logger *zap.Logger) *Player {
	return &Player{
		client: client,
		logger: logger,
	}
}

func (p Player) SavePlayer(
	ctx context.Context,
	name string,
	bestScore, totalGames int,
	registredAt time.Time,
) error {
	ctxSaveTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()
	_, err := p.client.HSet(ctxSaveTimeout, name, model.NewPlayer(
		registredAt.Format(time.DateOnly),
		bestScore,
		totalGames,
	)).Result()
	if err != nil {
		p.logger.Error(
			"SavePlayer failed",
			zap.String("redis request type", "HSet"),
			zap.Error(err),
			zap.Duration("duration", time.Since(prep)),
		)

		return err
	}

	p.logger.Info(
		"SavePlayer success",
		zap.String("redis request type", "HSet"),
		zap.Duration("duration", time.Since(prep)),
	)

	return nil
}

func (p Player) FindPlayerByName(ctx context.Context, name string) (model.Player, error) {
	ctxGetTimeout, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()

	player := model.Player{}
	if err := p.client.HGetAll(ctxGetTimeout, name).Scan(&player); err != nil {
		p.logger.Error(
			"FindPlayerByName failed",
			zap.String("redis request type", "HGetAll"),
			zap.Error(err),
			zap.Duration("duration", time.Since(prep)),
		)

		return model.Player{}, err
	}

	p.logger.Info(
		"FindPlayerByName success",
		zap.String("redis request type", "HGetAll"),
		zap.Duration("duration", time.Since(prep)),
	)

	return player, nil
}
