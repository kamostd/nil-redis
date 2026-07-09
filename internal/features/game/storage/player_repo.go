package repository

import (
	"context"
	"errors"
	"fmt"
	txmanager "nil-redis/internal/core/tools/tx_manager"
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

func ctxPlayerExtractor(ctx context.Context, client HmapClient) HmapClient {
	pipe, ok := ctx.Value(txmanager.PipelineKey).(redis.Pipeline)

	if !ok {
		return client
	}

	return pipe
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

func (p Player) WriteRegistrationDate(
	ctx context.Context,
	name string,
	registredAt time.Time,
) error {
	register := map[string]string{"registered_at": registredAt.Format(time.DateOnly)}

	client := ctxPlayerExtractor(ctx, p.client)

	ctxSave, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()
	if _, err := client.HSet(
		ctxSave,
		fmt.Sprintf("%s:%s", playerHmapName, name),
		register,
	).Result(); err != nil {
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

func (p Player) UpdateGamesPlayed(ctx context.Context, name string) error {
	ctxUpdate, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	client := ctxPlayerExtractor(ctx, p.client)

	prep := time.Now()
	if _, err := client.HIncrBy(
		ctxUpdate,
		fmt.Sprintf("%s:%s", playerHmapName, name),
		"games_played",
		1,
	).Result(); err != nil {
		p.logger.Error(
			"UpdateGamesPlayed failed",
			zap.String("redis request type", "HincrBy"),
			zap.Duration("duration", time.Since(prep)),
			zap.Error(err),
		)

		return err
	}

	p.logger.Info(
		"UpdateGamesPlayed failed",
		zap.String("redis request type", "HincrBy"),
		zap.Duration("duration", time.Since(prep)),
	)

	return nil
}

func (p Player) GetBestScoreByName(ctx context.Context, name string) (int, error) {
	ctxUpdate, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	client := ctxPlayerExtractor(ctx, p.client)

	prep := time.Now()

	result, err := client.HGet(
		ctxUpdate,
		fmt.Sprintf("%s:%s", playerHmapName, name),
		"best_score",
	).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return -1, nil
		}

		p.logger.Error(
			"GetBestScoreByName failed",
			zap.String("redis request type", "HGet"),
			zap.Duration("duration", time.Since(prep)),
			zap.Error(err),
		)

		return -1, err
	}

	p.logger.Info(
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

func (p Player) UpdateBestScore(ctx context.Context, name string, score int) error {
	ctxUpdate, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	client := ctxPlayerExtractor(ctx, p.client)

	prep := time.Now()
	if _, err := client.HSet(
		ctxUpdate,
		fmt.Sprintf("%s:%s", playerHmapName, name),
		"best_score",
		score,
	).Result(); err != nil {
		p.logger.Error(
			"UpdateBestScore failed",
			zap.String("redis request type", "HSet"),
			zap.Duration("duration", time.Since(prep)),
			zap.Error(err),
		)

		return err
	}

	p.logger.Info(
		"UpdateBestScore failed",
		zap.String("redis request type", "HSet"),
		zap.Duration("duration", time.Since(prep)),
	)

	return nil
}

func (p Player) IsExists(ctx context.Context, name string) (bool, error) {
	ctxGet, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()
	count, err := p.client.Exists(ctxGet, fmt.Sprintf("%s:%s", playerHmapName, name)).Result()
	after := time.Since(prep)
	if err != nil {
		p.logger.Error(
			"IsExists failed",
			zap.String("redis request type", "HGetAll"),
			zap.Duration("duration", after),
			zap.Error(err),
		)

		return false, err
	}

	p.logger.Info(
		"IsExists success",
		zap.String("redis request type", "HGetAll"),
		zap.Duration("duration", after),
	)

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (p Player) GetPlayerByName(ctx context.Context, name string) (domain.Player, error) {
	ctxGet, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	var player model.Player
	prep := time.Now()
	if err := p.client.HGetAll(ctxGet, fmt.Sprintf("%s:%s", playerHmapName, name)).Scan(&player); err != nil {
		p.logger.Error(
			"GetPlayerByName failed",
			zap.String("redis request type", "HGetAll"),
			zap.Duration("duration", time.Since(prep)),
		)

		return domain.Player{}, err
	}

	p.logger.Debug("redis player", zap.Any("player", player))

	p.logger.Info(
		"GetPlayerByName success",
		zap.String("redis request type", "HGetAll"),
		zap.Duration("duration", time.Since(prep)),
	)

	p.logger.Debug("GetPlayerByName storage player", zap.Any("player", player))

	parsedPlayer, err := mapper.ModelPlayerToDomain(player, name)
	if err != nil {
		return domain.Player{}, err
	}

	p.logger.Debug("GetPlayerByName mapping to domain", zap.Any("time", parsedPlayer.RegisteredAt()))

	return parsedPlayer, nil
}
