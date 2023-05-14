package rediscfg

type Redis struct {
	Addr string `validate:"required"`
	Pass string
	User string
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
