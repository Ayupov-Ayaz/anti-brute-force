package config

import (
	limitercfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/limiter"
	listcfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/list"
	loggercfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/logger"
	rediscfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/redis"
	servercfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/server"
)

type Config struct {
	Server  servercfg.Server   `envconfig:"SERVER" validate:"required"`
	Redis   rediscfg.Redis     `envconfig:"REDIS" validate:"required"`
	Logger  loggercfg.Logger   `envconfig:"LOGGER" validate:"required"`
	IPList  listcfg.IPList     `envconfig:"IP_LIST" validate:"required"`
	Limiter limitercfg.Limiter `envconfig:"LIMITER" validate:"required"`
}
