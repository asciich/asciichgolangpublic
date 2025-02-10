package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestTcpPortsIsPortOpen(t *testing.T) {
	tests := []struct {
		hostname       string
		portNumber     int
		expectedIsOpen bool
	}{
		{"google.ch", 80, true},
		{"google.ch", 443, true},
		{"google.ch", 442, false},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				require.EqualValues(
					tt.expectedIsOpen,
					TcpPorts().MustIsPortOpen(tt.hostname, tt.portNumber, verbose),
				)
			},
		)
	}
}
