package config

import (
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

	if err := parseYaml(cfg); err != nil {
		return nil, err
	}

	mergeConfigs(cfg, argPort, argUseGRPC)

	return cfg, nil
}
