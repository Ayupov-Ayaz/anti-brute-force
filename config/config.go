package config

type HTTP struct {
	Port int
}

type GRPC struct {
	Port int
}

type IPList struct {
	BlackListAddr string
	WhiteListAddr string
}

type Redis struct {
	Addr     string
	Password string
	User     string
}

type Config struct {
	HTTP   HTTP
	GRPC   GRPC
	Redis  Redis
	IPList IPList
}

func (c Config) UseHTTP() bool {
	return c.HTTP.Port > 0
}

func (c Config) UseGRPC() bool {
	return c.GRPC.Port > 0
}
