package config

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/joho/godotenv"

	"github.com/kelseyhightower/envconfig"
)

const (
	envPrefix = "ABF"
	envFile   = ".env"
	envPath   = "."
)

func readEnv() error {
	filePath := path.Join(envPath, envFile)
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("os stat : %w", err)
	}

	if err := godotenv.Load(filePath); err != nil {
		return fmt.Errorf("godotenv load : %w", err)
	}

	return nil
}

func parseEnv(cfg *Config) error {
	if err := readEnv(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := envconfig.Process(envPrefix, cfg); err != nil {
		return fmt.Errorf("error parsing env: %w", err)
	}

	return nil
}

func mergeConfigs(cfg *Config, argPort int, argUseGRPC bool) {
	if cfg.Server.Port != argPort && argPort != 0 {
		cfg.Server.Port = argPort
	}

	if cfg.Server.UseGRPC != argUseGRPC && cfg.Server.UseGRPC != true {
		cfg.Server.UseGRPC = argUseGRPC
	}
}

func ParseConfig(argPort int, argUseGRPC bool) (*Config, error) {
	cfg := &Config{}

	if err := parseEnv(cfg); err != nil {
		return nil, err
	}

	mergeConfigs(cfg, argPort, argUseGRPC)

	return cfg, nil
}
