package tmux

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestTemux_IsTmuxKey(t *testing.T) {

	tests := []struct {
		input         string
		expectedIsKey bool
	}{
		{"testcase", false},
		{"", false},
		{"enter", true},
		{"C-c", true},
		{"C-d", true},
		{"C-l", true},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require.EqualValues(t, tt.expectedIsKey, IsTmuxKey(tt.input))
			},
		)
	}
}
