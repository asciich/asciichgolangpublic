package tracederrors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func inThisFunctionSomethingGoesWrong() (err error) {
	// Use TracedErrors when an error occures.
	return tracederrors.TracedError("This is an error message")
}

func Test_ExampleUsage(t *testing.T) {
	// This function calls emualtes an error.
	err := inThisFunctionSomethingGoesWrong()

	// The retured error is a TracedError:
	require.True(t, tracederrors.IsTracedError(err)) // returns true for all TracedErrors.
	require.False(t, tracederrors.IsTracedError(fmt.Errorf("another error"))) // returns false for all non TracedErrors.

	// Get the error message
	msg := err.Error() // includes the error message and the stack trace as human readable text.

	// The error message contains the stack trace:
	require.Contains(t, msg, "inThisFunctionSomethingGoesWrong")
	require.Contains(t, msg, "Test_ExampleUsage")

	// The error message contains the file names in the stack trace as well:
	require.Contains(t, msg, "/Example_usage_test.go")

	// The error message itself is of corse as well part of the tracederror error message:
	require.Contains(t, msg, "This is an error message")
}
