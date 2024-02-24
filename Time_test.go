package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeGetCurrentTimeAsVersionStringString(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"only a string"},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				generatedString := Time().GetCurrentTimeAsSortableString()
				assert.Len(generatedString, len("YYYYmmdd_HHMMSS"))
			},
		)
	}
}
