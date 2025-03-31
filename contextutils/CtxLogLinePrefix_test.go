package contextutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/contextutils"
)

func Test_LogLinePrefix(t *testing.T) {
	t.Run("nil and empty string", func(t *testing.T) {
		ctx := contextutils.WithLogLinePrefix(nil, "")
		require.NotNil(t, ctx)

		require.EqualValues(t, "", contextutils.GetLogLinePrefixFromCtx(ctx))
	})

	t.Run("nil and empty string", func(t *testing.T) {
		ctx := contextutils.WithLogLinePrefix(nil, "abc")
		require.NotNil(t, ctx)

		require.EqualValues(t, "abc", contextutils.GetLogLinePrefixFromCtx(ctx))
	})
}
