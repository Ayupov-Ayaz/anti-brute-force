package loggercfg

type Logger struct {
	Level string `envconfig:"level" validate:"oneof=DEBUG INFO WARN ERROR"`
}
