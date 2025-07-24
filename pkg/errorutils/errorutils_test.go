package errorutils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/errorutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func Test_AppendToErrorMessage(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := errorutils.AppendToErrorMessage(nil, "appendix")
		require.EqualValues(t, "appendix", err.Error())
	})

	t.Run("fmt.Errorf", func(t *testing.T) {
		err := errorutils.AppendToErrorMessage(fmt.Errorf("abc"), "appendix")
		require.EqualValues(t, "abc appendix", err.Error())
	})

	t.Run("tracedError", func(t *testing.T) {
		err := tracederrors.TracedError("abc")
		err = errorutils.AppendToErrorMessage(err, "appendix")
		require.NotEqual(t, "abc appendix", err.Error())
		require.Equal(t, "abc appendix", errorutils.GetErrorMessage(err))
	})
}

func Test_GetErrorMessage(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		require.EqualValues(t, "", errorutils.GetErrorMessage(nil))
	})

	t.Run("empty string fmt.Errorf", func(t *testing.T) {
		err := fmt.Errorf("")
		require.EqualValues(t, "", errorutils.GetErrorMessage(err))
	})

	t.Run("empty traced error", func(t *testing.T) {
		err := tracederrors.TracedError("")
		require.EqualValues(t, "", errorutils.GetErrorMessage(err))
	})

	t.Run("non empty string fmt.Errorf", func(t *testing.T) {
		err := fmt.Errorf("abc")
		require.EqualValues(t, "abc", errorutils.GetErrorMessage(err))
	})

	t.Run("non empty traced error", func(t *testing.T) {
		err := tracederrors.TracedError("abc")
		require.EqualValues(t, "abc", errorutils.GetErrorMessage(err))
	})
}
