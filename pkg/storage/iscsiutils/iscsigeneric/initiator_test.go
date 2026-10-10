package iscsigeneric_test

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/storage/iscsiutils/iscsigeneric"
)

var initiatorNameRegex = regexp.MustCompile(`^iqn\.2016-04\.com\.open-iscsi:[0-9a-f]{12}$`)

func Test_GenerateInitiatorName(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		ctx := context.Background()

		_, err := iscsigeneric.GenerateInitiatorName(ctx)
		require.NoError(t, err)
	})

	t.Run("starts with open-iscsi prefix", func(t *testing.T) {
		ctx := context.Background()

		name, err := iscsigeneric.GenerateInitiatorName(ctx)
		require.NoError(t, err)
		require.True(t, strings.HasPrefix(name, iscsigeneric.DefaultInitiatorNamePrefix+":"))
	})

	t.Run("matches open-iscsi format", func(t *testing.T) {
		ctx := context.Background()

		name, err := iscsigeneric.GenerateInitiatorName(ctx)
		require.NoError(t, err)
		require.Regexp(t, initiatorNameRegex, name)
	})

	t.Run("has expected length", func(t *testing.T) {
		ctx := context.Background()

		name, err := iscsigeneric.GenerateInitiatorName(ctx)
		require.NoError(t, err)
		require.Len(t, name, len(iscsigeneric.DefaultInitiatorNamePrefix)+1+12)
	})

	t.Run("does not exceed max iSCSI name length of 223 bytes", func(t *testing.T) {
		ctx := context.Background()

		name, err := iscsigeneric.GenerateInitiatorName(ctx)
		require.NoError(t, err)
		require.LessOrEqual(t, len(name), 223)
	})

	t.Run("generated names are unique", func(t *testing.T) {
		ctx := context.Background()

		const count = 1000
		seen := make(map[string]struct{}, count)

		for i := 0; i < count; i++ {
			name, err := iscsigeneric.GenerateInitiatorName(ctx)
			require.NoError(t, err)

			_, exists := seen[name]
			require.False(t, exists, "Duplicate initiator name generated: '%s'", name)

			seen[name] = struct{}{}
		}
	})
}
