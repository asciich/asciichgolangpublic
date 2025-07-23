package mathutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func TestMathUtils_MaxInt(t *testing.T) {
	tests := []struct {
		i1             int
		i2             int
		expectedResult int
	}{
		{0, 0, 0},
		{0, -1, 0},
		{-1, -1, -1},
		{-1, 0, 0},
		{-1, 5, 5},
		{42, -1, 42},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require.EqualValues(
					t,
					tt.expectedResult,
					MaxInt(tt.i1, tt.i2),
				)
			},
		)
	}
}
