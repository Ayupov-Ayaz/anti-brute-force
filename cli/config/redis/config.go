package rediscfg

type Redis struct {
	Addr string `mapstructure:"addr" validate:"required"`
	Pass string `mapstructure:"password"`
	User string `mapstructure:"user"`
}
