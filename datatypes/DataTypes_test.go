package datatypes

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var ErrGlobalForTesting = errors.New("this global error is for testing")

func TestTypes_GetTypeName(t *testing.T) {
	var ErrForTesting = errors.New("this error is for testing")

	tests := []struct {
		input        interface{}
		expectedName string
	}{
		{5, "int"},
		{"", "string"},
		{float32(5), "float32"},
		{float64(5), "float64"},

		// Instead of directly access the variables a string identifier is used.
		// This allows single selection in test runs while addresses can change and are not suitable for test selection.
		{"&ErrForTesting", "&error{message='this error is for testing'}"},
		{"ErrForTesting", "error{message='this error is for testing'}"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				input := tt.input

				if input == "&ErrForTesting" {
					input = &ErrForTesting
				}

				if input == "ErrForTesting" {
					input = ErrForTesting
				}

				if input == "&ErrGlobalForTesting" {
					input = &ErrGlobalForTesting
				}

				if input == "ErrGlobalForTesting" {
					input = ErrGlobalForTesting
				}

				require.EqualValues(
					t,
					tt.expectedName,
					MustGetTypeName(input),
				)
			},
		)
	}
}
