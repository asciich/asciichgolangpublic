package asciichgolangpublic

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ErrGlobalForTesting = errors.New("this global error is for testing")

func TestTypes_GetTypeName(t *testing.T) {
	var ErrForTesting = errors.New("this error is for testing")
	var tracedError = TracedError("this is a traced error")

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
		{"TracedError", "TracedErrorType{message='this is a traced error'}"},
		{"&TracedError", "&TracedErrorType{message='this is a traced error'}"},
		// TODO {"ErrGlobalForTesting", "ErrGlobalForTesting{message='this global error is for testing'}"},
		// TODO {"&ErrGlobalForTesting", "&ErrGlobalForTesting{message='this global error is for testing'}"},

	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				input := tt.input

				if input == "&ErrForTesting" {
					input = &ErrForTesting
				}

				if input == "ErrForTesting" {
					input = ErrForTesting
				}

				if input == "&TracedError" {
					input = &tracedError
				}

				if input == "TracedError" {
					input = tracedError
				}

				if input == "&ErrGlobalForTesting" {
					input = &ErrGlobalForTesting
				}

				if input == "ErrGlobalForTesting" {
					input = ErrGlobalForTesting
				}

				assert.EqualValues(
					tt.expectedName,
					Types().MustGetTypeName(input),
				)
			},
		)
	}
}
