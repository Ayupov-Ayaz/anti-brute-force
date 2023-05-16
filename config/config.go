package config

import (
	listcfg "github.com/ayupov-ayaz/anti-brute-force/config/list"
	loggercfg "github.com/ayupov-ayaz/anti-brute-force/config/logger"
	rediscfg "github.com/ayupov-ayaz/anti-brute-force/config/redis"
	servercfg "github.com/ayupov-ayaz/anti-brute-force/config/server"
)

type Config struct {
	Server servercfg.Server `mapstructure:"server"`
	Redis  rediscfg.Redis   `mapstructure:"redis"`
	Logger loggercfg.Logger `mapstructure:"logger"`
	IPList listcfg.IPList   `mapstructure:"ip_list"`
}
