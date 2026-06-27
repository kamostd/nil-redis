package repository

import (
	"context"
	"nil-redis/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const leaderboardZName = "leaderboard"

type ClientLeaderboard interface {
	ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd
	ZRevRank(ctx context.Context, key, member string) *redis.IntCmd

	ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd

	ZScore(ctx context.Context, key, member string) *redis.FloatCmd
}

type Leaderboard struct {
	client ClientLeaderboard
	logger *zap.Logger
}

func NewLeaderboard(client ClientLeaderboard, logger *zap.Logger) *Leaderboard {
	return &Leaderboard{
		client: client,
		logger: logger,
	}
}

func (l Leaderboard) SaveOrUpdateMember(ctx context.Context, name string, score int) (int, error) {
	ctxSaveOrUpdate, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()
	result, err := l.client.ZIncrBy(
		ctxSaveOrUpdate,
		leaderboardZName,
		float64(score),
		name,
	).Result()
	if err != nil {
		l.logger.Error(
			"SaveOrUpdateMember failed",
			zap.String("redis operation type", "ZIncrBy"),
			zap.Duration("duration", time.Since(prep)),
			zap.Error(err),
		)

		return -1, err
	}

	l.logger.Info(
		"SaveOrUpdateMember success",
		zap.String("redis operation type", "ZIncrBy"),
		zap.Duration("duration", time.Since(prep)),
	)

	return int(result), nil
}

func (l Leaderboard) RankByName(ctx context.Context, name string) (int, error) {
	ctxRank, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()
	rank, err := l.client.ZRevRank(ctxRank, leaderboardZName, name).Result()
	if err != nil {
		l.logger.Error(
			"RankByName failed",
			zap.String("redis operation type", "ZRank"),
			zap.Duration("duration", time.Since(prep)),
			zap.Error(err),
		)

		return -1, err
	}

	l.logger.Info(
		"RankByName success",
		zap.String("redis operation type", "ZRank"),
		zap.Duration("duration", time.Since(prep)),
	)

	return int(rank) + 1, nil
}

func (l Leaderboard) ListTopPlayerWithScore(ctx context.Context, top int) ([]domain.LeaderboardPlayer, error) {
	ctxList, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()
	notes, err := l.client.ZRevRangeWithScores(ctxList, leaderboardZName, 0, int64(top)).Result()
	if err != nil {
		l.logger.Error(
			"ListTopPlayerWithScore failed",
			zap.String("redis operation type", "ZRevRangeWithScores"),
			zap.Duration("duration", time.Since(prep)),
			zap.Error(err),
		)

		return nil, err
	}

	l.logger.Info(
		"ListTopPlayerWithScore success",
		zap.String("redis operation type", "ZRevRangeWithScores"),
		zap.Duration("duration", time.Since(prep)),
	)

	leaderboard := make([]domain.LeaderboardPlayer, 0, len(notes))

	for i, note := range notes {
		leaderboard = append(leaderboard, domain.NewLeaderboard(
			note.Member.(string),
			i+1,
			int(note.Score),
		))
	}

	return leaderboard, nil
}

func (l Leaderboard) GetTotalScoreByName(ctx context.Context, name string) (int, error) {
	ctxGet, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prep := time.Now()
	result, err := l.client.ZScore(ctxGet, leaderboardZName, name).Result()
	if err != nil {
		l.logger.Error(
			"GetTotalScoreByName failed",
			zap.String("redis operation type", "ZScore"),
			zap.Duration("duration", time.Since(prep)),
			zap.Error(err),
		)

		return -1, err
	}

	l.logger.Info(
		"GetTotalScoreByName success",
		zap.String("redis operation type", "ZScore"),
		zap.Duration("duration", time.Since(prep)),
	)

	return int(result), nil
}
