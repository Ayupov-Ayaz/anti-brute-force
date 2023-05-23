package ip

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJoinIPMask(t *testing.T) {
	tests := []struct {
		ip          string
		mask        string
		expected    string
		expectError bool
	}{
		{"192.168.0.1", "255.255.255.0", "192.168.0.0", false},
		{"10.0.0.1", "255.0.0.0", "10.0.0.0", false},
		{"192.168.0.1", "invalid-mask", "<nil>", true},
		{"invalid-ip", "255.255.255.0", "<nil>", true},
	}

	for _, tt := range tests {
		t.Run(tt.ip, func(t *testing.T) {
			ipMasked, err := joinIPMask(tt.ip, tt.mask)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tt.expected, ipMasked.String())
		})
	}
}

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
