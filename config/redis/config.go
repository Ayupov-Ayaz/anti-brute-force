package rediscfg

type Redis struct {
	Addr string `mapstructure:"addr" validate:"required"`
	Pass string `mapstructure:"password"`
	User string `mapstructure:"user"`
}

func (r Redis) Address() string {
	return r.Addr
}

func (r Redis) Username() string {
	return r.User
}

func (r Redis) Password() string {
	return r.Pass
}
