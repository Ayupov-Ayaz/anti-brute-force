package loggercfg

type Logger struct {
	LogLevel string `validate:"oneof=DEBUG INFO WARN ERROR"`
}

func (l Logger) Level() string {
	return l.LogLevel
}
