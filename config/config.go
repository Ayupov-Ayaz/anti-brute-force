package config

import (
	loggercfg "github.com/ayupov-ayaz/anti-brute-force/config/logger"
	rediscfg "github.com/ayupov-ayaz/anti-brute-force/config/redis"
)

type HTTP struct {
	Port int `mapstructure:"port"`
}

type GRPC struct {
	Port int `mapstructure:"port"`
}

type IPList struct {
	BlackListAddr string `mapstructure:"blacklist_addr"`
	WhiteListAddr string `mapstructure:"whitelist_addr"`
}

type Config struct {
	HTTP   HTTP             `mapstructure:"http"`
	GRPC   GRPC             `mapstructure:"grpc"`
	Redis  rediscfg.Redis   `mapstructure:"redis"`
	Logger loggercfg.Logger `mapstructure:"logger"`
	IPList IPList           `mapstructure:"ip_list"`
}

func (c Config) UseHTTP() bool {
	return c.HTTP.Port > 0
}

func (c Config) UseGRPC() bool {
	return c.GRPC.Port > 0
}
