package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			MustFormatAsTestname(tt),
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
