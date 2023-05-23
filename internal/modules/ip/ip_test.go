package ip

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIPToUint32(t *testing.T) {
	tests := []struct {
		ip       string
		expected uint32
	}{
		{
			ip:       "64.233.187.99",
			expected: 1089059683,
		},
		{
			ip:       "255.255.255.254",
			expected: 4294967294,
		},
	}

	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			ip := net.ParseIP(tt.ip)
			require.NotZero(t, ip)
			require.Equal(t, tt.expected, ipToUint32(ip))
		})
	}
}
