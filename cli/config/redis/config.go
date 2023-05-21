package rediscfg

type Redis struct {
	Addr string `envconfig:"addr" validate:"required"`
	Pass string `envconfig:"password"`
	User string `envconfig:"user"`
}
