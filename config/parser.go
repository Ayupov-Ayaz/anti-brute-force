package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func parseYaml(cfg *Config) error {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(".config.yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(cfg)
}

func ParseConfig() (*Config, error) {
	cfg := &Config{}

	if err := parseYaml(cfg); err != nil {
		return nil, err
	}

	if err := validator.New().Struct(cfg); err != nil {
		return nil, fmt.Errorf("validate configs failed: %w", err)
	}

	return cfg, nil
}
