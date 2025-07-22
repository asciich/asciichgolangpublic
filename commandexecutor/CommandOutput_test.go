package commandexecutor_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
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
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				commandOutput := commandexecutor.NewCommandOutput()
				err := commandOutput.SetReturnCode(tt.returnCode)
				require.NoError(t, err)

				returnCode, err := commandOutput.GetReturnCode()
				require.NoError(t, err)
				require.EqualValues(t, tt.returnCode, returnCode)
			},
		)
	}
}
