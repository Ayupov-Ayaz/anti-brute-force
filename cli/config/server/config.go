package servercfg

type Server struct {
	Port    int  `envconfig:"PORT" validate:"required"`
	UseGRPC bool `envconfig:"USE_GRPC"`
}
