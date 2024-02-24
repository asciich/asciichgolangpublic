package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				assert.EqualValues(
					tt.expectedIsOpen,
					TcpPorts().MustIsPortOpen(tt.hostname, tt.portNumber, verbose),
				)
			},
		)
	}
}
