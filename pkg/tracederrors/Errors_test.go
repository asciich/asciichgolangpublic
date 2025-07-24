package tracederrors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func TestErrorsIsTracedError(t *testing.T) {
	tests := []struct {
		err                   error
		expectedIsTracedError bool
	}{
		{fmt.Errorf("an error"), false},
		{tracederrors.TracedError("an error"), true},
		{nil, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				require.EqualValues(
					tt.expectedIsTracedError,
					tracederrors.IsTracedError(tt.err),
				)
			},
		)
	}
}
