package config

import (
	httpcfg "github.com/ayupov-ayaz/anti-brute-force/config/http"
	loggercfg "github.com/ayupov-ayaz/anti-brute-force/config/logger"
	rediscfg "github.com/ayupov-ayaz/anti-brute-force/config/redis"
)

type GRPC struct {
	Port int
}

type IPList struct {
	BlackListAddr string
	WhiteListAddr string
}

type Config struct {
	HTTP   httpcfg.HTTP
	GRPC   GRPC
	Redis  rediscfg.Redis
	Logger loggercfg.Logger
	IPList IPList
}

func (c Config) UseHTTP() bool {
	return c.HTTP.Port > 0
}

func (c Config) UseGRPC() bool {
	return c.GRPC.Port > 0
}
