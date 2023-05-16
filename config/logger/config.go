package loggercfg

type Logger struct {
	LogLevel string `mapstructure:"level" validate:"oneof=DEBUG INFO WARN ERROR"`
}

func (l Logger) Level() string {
	return l.LogLevel
}
