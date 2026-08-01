package txmanager

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const PipelineKey ctxPipelineKey = "pipline"

type ctxPipelineKey string

func ctxInjector(ctx context.Context, p redis.Pipeliner) context.Context {
	return context.WithValue(ctx, PipelineKey, p)
}

type Pipeliner interface {
	TxPipeline() redis.Pipeliner
}

type Trm struct {
	pipe   Pipeliner
	logger *zap.Logger
}

func NewTrm(
	client *redis.Client,
	logger *zap.Logger,
) *Trm {
	return &Trm{
		pipe:   client,
		logger: logger,
	}
}

func (t Trm) Do(ctx context.Context, tx func(context.Context) error) error {
	pipe := t.pipe.TxPipeline()

	ctxWithPipline := ctxInjector(ctx, pipe)

	if err := tx(ctxWithPipline); err != nil {
		return err
	}

	if _, err := pipe.Exec(ctxWithPipline); err != nil {
		t.logger.Error("Do", zap.Error(err))

		return err
	}

	return nil
}
