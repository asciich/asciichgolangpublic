package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestMacAddressesIsStringAMacAddress(t *testing.T) {
	tests := []struct {
		input        string
		isMacAddress bool
	}{
		{"testcase", false},
		{"z", false},
		{"52:54:00:b0:90:86", true},
		{"52:54:00:B0:90:86", true},
		{"52:54:00:B0:90:86:", false},
		{":52:54:00:B0:90:86:", false},
		{":52:54:00:B0:90:86", false},
		{"52:54:00:B0:90:8g", false},
		{"52:54:00:B0:90:", false},
		{"52:54:00:B0:90", false},
		{"52:54:00:B0", false},
		{"52:54:00", false},
		{"52:54", false},
		{"52", false},
		{"5", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)
				assert.EqualValues(tt.isMacAddress, MacAddresses().IsStringAMacAddress(tt.input))
			},
		)
	}
}
