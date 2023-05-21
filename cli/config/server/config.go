package servercfg

type Server struct {
	Port    int  `mapstructure:"port" validate:"required"`
	UseGRPC bool `mapstructure:"use_grpc"`
}
