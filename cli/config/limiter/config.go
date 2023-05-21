package limitercfg

import "time"

type ValueLimiter struct {
	Count    int64         `mapstructure:"count" validate:"required"`
	Interval time.Duration `mapstructure:"interval" validate:"required"`
}

type Limiter struct {
	Login    ValueLimiter `mapstructure:"login" validate:"required"`
	Password ValueLimiter `mapstructure:"password" validate:"required"`
	IP       ValueLimiter `mapstructure:"ip" validate:"required"`
}
