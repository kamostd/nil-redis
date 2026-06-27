package repository

import (
	"context"
	"errors"
	"fmt"
	"nil-redis/internal/domain"
	"nil-redis/internal/features/game/storage/mapper"
	"nil-redis/internal/features/game/storage/model"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const playerHmapName = "player"

type HmapClient interface {
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd

	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd

	Exists(ctx context.Context, keys ...string) *redis.IntCmd
}

type Game struct {
	client HmapClient
	logger *zap.Logger
}

func NewPlayerRepo(client HmapClient, logger *zap.Logger) *Game {
	return &Game{
		client: client,
		logger: logger,
	}
}

func (g Game) WriteRegistrationDate(
	ctx context.Context,
	name string,
	registredAt time.Time,
) error {
	register := map[string]string{"registered_at": registredAt.Format(time.DateOnly)}

	ctxSave, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()
	if _, err := g.client.HSet(
		ctxSave,
		fmt.Sprintf("%s:%s", playerHmapName, name),
		register,
	).Result(); err != nil {
		g.logger.Error(
			"SavePlayer failed",
			zap.String("redis request type", "HSet"),
			zap.Error(err),
			zap.Duration("duration", time.Since(prep)),
		)

		return err
	}

	g.logger.Info(
		"SavePlayer success",
		zap.String("redis request type", "HSet"),
		zap.Duration("duration", time.Since(prep)),
	)

	return nil
}

func (g Game) UpdateGamesPlayed(ctx context.Context, name string) error {
	ctxUpdate, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()
	if _, err := g.client.HIncrBy(
		ctxUpdate,
		fmt.Sprintf("%s:%s", playerHmapName, name),
		"games_played",
		1,
	).Result(); err != nil {
		g.logger.Error(
			"UpdateGamesPlayed failed",
			zap.String("redis request type", "HincrBy"),
			zap.Duration("duration", time.Since(prep)),
			zap.Error(err),
		)

		return err
	}

	g.logger.Info(
		"UpdateGamesPlayed failed",
		zap.String("redis request type", "HincrBy"),
		zap.Duration("duration", time.Since(prep)),
	)

	return nil
}

func (g Game) GetBestScoreByName(ctx context.Context, name string) (int, error) {
	ctxUpdate, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()

	result, err := g.client.HGet(
		ctxUpdate,
		fmt.Sprintf("%s:%s", playerHmapName, name),
		"best_score",
	).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return -1, nil
		}

		g.logger.Error(
			"GetBestScoreByName failed",
			zap.String("redis request type", "HGet"),
			zap.Duration("duration", time.Since(prep)),
			zap.Error(err),
		)

		return -1, err
	}

	g.logger.Info(
		"GetBestScoreByName success",
		zap.String("redis request type", "HGet"),
		zap.Duration("duration", time.Since(prep)),
	)

	bestScore, err := strconv.Atoi(result)
	if err != nil {
		return -1, err
	}

	return bestScore, nil
}

func (g Game) UpdateBestScore(ctx context.Context, name string, score int) error {
	ctxUpdate, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()
	if _, err := g.client.HSet(
		ctxUpdate,
		fmt.Sprintf("%s:%s", playerHmapName, name),
		"best_score",
		score,
	).Result(); err != nil {
		g.logger.Error(
			"UpdateBestScore failed",
			zap.String("redis request type", "HSet"),
			zap.Duration("duration", time.Since(prep)),
			zap.Error(err),
		)

		return err
	}

	g.logger.Info(
		"UpdateBestScore failed",
		zap.String("redis request type", "HSet"),
		zap.Duration("duration", time.Since(prep)),
	)

	return nil
}

func (g Game) IsExists(ctx context.Context, name string) (bool, error) {
	ctxGet, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()
	count, err := g.client.Exists(ctxGet, fmt.Sprintf("%s:%s", playerHmapName, name)).Result()
	after := time.Since(prep)
	if err != nil {
		g.logger.Error(
			"IsExists failed",
			zap.String("redis request type", "HGetAll"),
			zap.Duration("duration", after),
			zap.Error(err),
		)

		return false, err
	}

	g.logger.Info(
		"IsExists success",
		zap.String("redis request type", "HGetAll"),
		zap.Duration("duration", after),
	)

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (g Game) GetPlayerByName(ctx context.Context, name string) (domain.Player, error) {
	ctxGet, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	var player model.Player
	prep := time.Now()
	if err := g.client.HGetAll(ctxGet, fmt.Sprintf("%s:%s", playerHmapName, name)).Scan(&player); err != nil {
		g.logger.Error(
			"GetPlayerByName failed",
			zap.String("redis request type", "HGetAll"),
			zap.Duration("duration", time.Since(prep)),
		)

		return domain.Player{}, err
	}

	g.logger.Debug("redis player", zap.Any("player", player))

	g.logger.Info(
		"GetPlayerByName success",
		zap.String("redis request type", "HGetAll"),
		zap.Duration("duration", time.Since(prep)),
	)

	g.logger.Debug("GetPlayerByName storage player", zap.Any("player", player))

	parsedPlayer, err := mapper.ModelPlayerToDomain(player, name)
	if err != nil {
		return domain.Player{}, err
	}

	g.logger.Debug("GetPlayerByName mapping to domain", zap.Any("time", parsedPlayer.RegisteredAt()))

	return parsedPlayer, nil
}
