package ports

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInferTunnelScheme(t *testing.T) {
	tests := []struct {
		name        string
		port        uint16
		serviceName string
		want        string
	}{
		{
			name:        "service rtsps takes precedence",
			port:        554,
			serviceName: "rtsps",
			want:        "rtsps",
		},
		{
			name:        "service https takes precedence",
			port:        80,
			serviceName: "https",
			want:        "https",
		},
		{
			name:        "service http takes precedence",
			port:        443,
			serviceName: "http",
			want:        "http",
		},
		{
			name: "fallback rtsps on port 322",
			port: 322,
			want: "rtsps",
		},
		{
			name: "fallback rtsps on port 8322",
			port: 8322,
			want: "rtsps",
		},
		{
			name: "fallback https on port 443",
			port: 443,
			want: "https",
		},
		{
			name: "fallback https on port 8443",
			port: 8443,
			want: "https",
		},
		{
			name: "fallback http on port 80",
			port: 80,
			want: "http",
		},
		{
			name: "fallback http on port 8080",
			port: 8080,
			want: "http",
		},
		{
			name: "unknown port without service",
			port: 554,
			want: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, InferTunnelScheme(test.port, test.serviceName))
		})
	}
}
