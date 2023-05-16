package servercfg

type Server struct {
	Port    int  `mapstructure:"port"`
	UseGRPC bool `mapstructure:"use_grpc"`
}
