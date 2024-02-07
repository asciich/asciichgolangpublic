package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOsIsRunningOnWindows(t *testing.T) {
	tests := []struct {
		testmessage string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.False(OS().IsRunningOnWindows())
			},
		)
	}
}
