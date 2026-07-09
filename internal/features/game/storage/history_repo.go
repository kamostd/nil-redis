package repository

import (
	"context"
	txmanager "nil-redis/internal/core/tools/tx_manager"
	"nil-redis/internal/domain"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const historyTableName = "games:history"

type ListClient interface {
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd

	LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
}

func ctxListExtractor(ctx context.Context, client ListClient) ListClient {
	pipe, ok := ctx.Value(txmanager.PipelineKey).(redis.Pipeline)
	if !ok {
		return client
	}

	return pipe
}

type HistoryRepo struct {
	client ListClient
	logger *zap.Logger
}

func NewHistoryRepo(client ListClient, logger *zap.Logger) *HistoryRepo {
	return &HistoryRepo{
		client: client,
		logger: logger,
	}
}

// TODO: вынести trim в отдельную функцию чтобы задавать range отделения
func (h HistoryRepo) SaveGame(ctx context.Context, historyNote string) error {
	ctxSave, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	client := ctxListExtractor(ctx, h.client)

	prepLpush := time.Now()
	if _, err := client.LPush(ctxSave, historyTableName, historyNote).Result(); err != nil {
		h.logger.Error(
			"SaveGame failed",
			zap.String("redis request type", "LPush"),
			zap.Duration("duration", time.Since(prepLpush)),
			zap.Error(err),
		)

		return err
	}

	h.logger.Info(
		"SaveGame success",
		zap.String("redis request type", "LPush"),
		zap.Duration("duration", time.Since(prepLpush)),
	)

	ctxTrim, cancelTrim := context.WithTimeout(ctx, time.Second)
	defer cancelTrim()

	prepTirm := time.Now()
	if _, err := client.LTrim(ctxTrim, historyTableName, 0, 9).Result(); err != nil {
		h.logger.Error(
			"SaveGame failed",
			zap.String("redis request type", "LTrim"),
			zap.Duration("duration", time.Since(prepTirm)),
			zap.Error(err),
		)

		return err
	}

	h.logger.Info(
		"SaveGame success",
		zap.String("redis request type", "LTrim"),
		zap.Duration("duration", time.Since(prepLpush)),
	)

	return nil
}

func (h HistoryRepo) ListHistory(ctx context.Context) ([]domain.HistoryNote, error) {
	ctxSave, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	prepLpush := time.Now()

	storageNotes, err := h.client.LRange(ctxSave, historyTableName, 0, 9).Result()
	if err != nil {
		h.logger.Error(
			"ListHistory failed",
			zap.String("redis request type", "Range"),
			zap.Duration("duration", time.Since(prepLpush)),
			zap.Error(err),
		)

		return nil, err
	}

	h.logger.Info(
		"ListHistory success",
		zap.String("redis request type", "Range"),
		zap.Duration("duration", time.Since(prepLpush)),
	)

	notes := make([]domain.HistoryNote, 0, len(storageNotes))

	for _, storageNote := range storageNotes {
		note, err := domain.NewNoteByString(storageNote)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, nil

}
