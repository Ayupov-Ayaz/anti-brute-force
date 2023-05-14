package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Config interface {
	Level() string
}

func New(config Config) (*zap.Logger, error) {
	logCfg := zap.NewProductionConfig()
	if err := logCfg.Level.UnmarshalText([]byte(config.Level())); err != nil {
		return nil, fmt.Errorf("unmarshal log level failed: %w, level=%s", err, config.Level())
	}

	l, err := logCfg.Build()
	if err != nil {
		return nil, err
	}

	return l, nil
}
