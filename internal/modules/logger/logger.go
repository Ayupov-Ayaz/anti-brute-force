package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

func New(logLevel string) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("parse log level: %w", err)
	}

	return zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger(), nil
}
