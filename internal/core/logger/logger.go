package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.Logger, func() error, error) {
	zapLevel, err := zapcore.ParseLevel("INFO")
	if err != nil {
		return nil, nil, err
	}

	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, nil, err

	}

	filename := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02T15-04-05.000000Z07"))
	filepath := filepath.Join("logs", filename)

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return nil, nil, err
		}
	}

	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000Z07")
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(cfg), zapcore.AddSync(file), zapLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), os.Stdout, zapLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)), file.Close, nil
}
