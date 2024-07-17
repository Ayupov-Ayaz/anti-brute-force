package limitercfg

import "time"

type ValueLimiter struct {
	Count    int64         `envconfig:"COUNT" validate:"required"`
	Interval time.Duration `envconfig:"INTERVAL" validate:"required"`
}

type Limiter struct {
	Login    ValueLimiter `envconfig:"LOGIN" validate:"required"`
	Password ValueLimiter `envconfig:"PASSWORD" validate:"required"`
	IP       ValueLimiter `envconfig:"IP" validate:"required"`
}
