package ip

import (
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
