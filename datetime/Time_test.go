package datetime

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTimeGetCurrentTimeAsVersionStringString(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"only a string"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				generatedString := Time().GetCurrentTimeAsSortableString()
				require.Len(generatedString, len("YYYYmmdd_HHMMSS"))
			},
		)
	}
}
