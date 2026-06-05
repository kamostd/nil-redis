package logger

import (
	"context"
	"errors"
	"syscall"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func ProvideLogger(lc fx.Lifecycle) (*zap.Logger, error) {
	logger, close, err := New()
	if err != nil {
		return nil, err
	}

	lc.Append(
		fx.Hook{
			OnStop: func(ctx context.Context) error {
				if err := logger.Sync(); err != nil {
					if err != nil && !errors.Is(err, syscall.EINVAL) {
						return err
					}
				}

				return close()
			},
		},
	)

	return logger, nil
}

func NewModule() fx.Option {
	return fx.Module(
		"logger",
		fx.Provide(
			ProvideLogger,
		),
	)
}
