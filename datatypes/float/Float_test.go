package float

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatToString(t *testing.T) {
	tests := []struct {
		input          float64
		maxDigits      int
		expectedOutput string
	}{
		{0, 0, "0"},
		{0.1, 0, "0"},
		{0.1, 1, "0.1"},
		{0.5, 0, "1"},
		{1.0, 0, "1"},
		{1.0, 1, "1"},
		{1.0, 2, "1"},
		{1.0, 3, "1"},
		{11.0, 0, "11"},
		{11.0, 1, "11"},
		{11.0, 2, "11"},
		{11.0, 3, "11"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				floatAsString := MustToString(tt.input, tt.maxDigits)
				assert.EqualValues(tt.expectedOutput, floatAsString)
			},
		)
	}
}

func TestFloatRound(t *testing.T) {
	tests := []struct {
		input          float64
		digits         int
		expectedOutput float64
	}{
		{0, 0, 0.0},
		{0.1, 0, 0.0},
		{0.1, 1, 0.1},
		{0.11, 1, 0.1},
		{0.1, 2, 0.1},
		{0.11, 2, 0.11},
		{0.49, 0, 0},
		{0.49, 1, 0.5},
		{0.49, 2, 0.49},
		{0.49, 3, 0.49},
		{0.5, 0, 1.0},
		{0.55, 0, 1.0},
		{0.555, 0, 1.0},
		{0.5555, 0, 1.0},
		{0.55555, 0, 1.0},
		{0.55555, 1, 0.6},
		{0.55555, 2, 0.56},
		{0.55555, 3, 0.556},
		{0.55555, 4, 0.5556},
		{0.55555, 5, 0.55555},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				rounded := MustRound(tt.input, tt.digits)
				assert.EqualValues(tt.expectedOutput, rounded)
			},
		)
	}
}
