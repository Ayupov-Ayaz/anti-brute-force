package config

import "github.com/go-playground/validator/v10"

func ParseConfig() (Config, error) {
	cfg := Config{}
	// todo: implement me
	err := validator.New().Struct(cfg)
	return cfg, err
}
