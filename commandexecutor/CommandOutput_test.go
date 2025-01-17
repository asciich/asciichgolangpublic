package commandexecutor

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert := assert.New(t)

				commandOutput := NewCommandOutput()
				commandOutput.MustSetReturnCode(tt.returnCode)
				assert.EqualValues(
					tt.returnCode,
					commandOutput.MustGetReturnCode(),
				)
			},
		)
	}
}
