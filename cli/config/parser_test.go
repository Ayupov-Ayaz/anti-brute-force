package config

import (
	"testing"

	servercfg "github.com/ayupov-ayaz/anti-brute-force/cli/config/server"
	"github.com/stretchr/testify/require"
)

func Test_mergeConfigs(t *testing.T) {
	table := []struct {
		name       string
		cfg        servercfg.Server
		argPort    int
		argUseGRPC bool
		check      func(t *testing.T, cfg servercfg.Server)
	}{
		{
			name: "empty",
			check: func(t *testing.T, cfg servercfg.Server) {
				require.Empty(t, cfg.Port)
				require.False(t, cfg.UseGRPC)
			},
		},
		{
			name: "default is set",
			cfg: servercfg.Server{
				Port:    8080,
				UseGRPC: true,
			},
			check: func(t *testing.T, cfg servercfg.Server) {
				require.Equal(t, 8080, cfg.Port)
				require.True(t, cfg.UseGRPC)
			},
		},
		{
			name:       "default is not set, arg is set",
			argPort:    9090,
			argUseGRPC: true,
			check: func(t *testing.T, cfg servercfg.Server) {
				require.Equal(t, 9090, cfg.Port)
				require.True(t, cfg.UseGRPC)
			},
		},
		{
			name: "default is set, arg is set (use grpcArg is true)",
			cfg: servercfg.Server{
				Port:    8080,
				UseGRPC: false,
			},
			argPort:    9090,
			argUseGRPC: true,
			check: func(t *testing.T, cfg servercfg.Server) {
				require.Equal(t, 9090, cfg.Port)
				require.True(t, cfg.UseGRPC)
			},
		},
		{
			name: "default is set, arg is set (use grpcArg is false)",
			cfg: servercfg.Server{
				Port:    8080,
				UseGRPC: true,
			},
			argPort:    9090,
			argUseGRPC: false,
			check: func(t *testing.T, cfg servercfg.Server) {
				require.Equal(t, 9090, cfg.Port)
				require.True(t, cfg.UseGRPC)
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Server: tt.cfg}
			mergeConfigs(cfg, tt.argPort, tt.argUseGRPC)
			tt.check(t, cfg.Server)
		})
	}
}
