package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type Config interface {
	Level() string
}

func New(config Config) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(config.Level())
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("parse log level: %w", err)
	}

	return zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger(), nil
}
