package config

import (
	limitercfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/limiter"
	listcfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/list"
	loggercfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/logger"
	rediscfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/redis"
	servercfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/server"
)

type Config struct {
	Server  servercfg.Server   `mapstructure:"server" validate:"required"`
	Redis   rediscfg.Redis     `mapstructure:"redis" validate:"required"`
	Logger  loggercfg.Logger   `mapstructure:"logger" validate:"required"`
	IPList  listcfg.IPList     `mapstructure:"ip_list" validate:"required"`
	Limiter limitercfg.Limiter `mapstructure:"limiter" validate:"required"`
}
