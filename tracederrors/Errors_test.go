package tracederrors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorsIsTracedError(t *testing.T) {
	tests := []struct {
		err                   error
		expectedIsTracedError bool
	}{
		{fmt.Errorf("an error"), false},
		{TracedError("an error"), true},
		{nil, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				require.EqualValues(
					tt.expectedIsTracedError,
					IsTracedError(tt.err),
				)
			},
		)
	}
}
