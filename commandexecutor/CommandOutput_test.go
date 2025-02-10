package commandexecutor

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func TestCommandOutputGetAndSetReturnCode(t *testing.T) {

	tests := []struct {
		returnCode int
	}{
		{-1},
		{1},
		{2},
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				commandOutput := NewCommandOutput()
				commandOutput.MustSetReturnCode(tt.returnCode)
				require.EqualValues(
					tt.returnCode,
					commandOutput.MustGetReturnCode(),
				)
			},
		)
	}
}
