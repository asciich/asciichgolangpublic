package commandexecutornicutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/commandexecutornicutils"
)

func TestGetFirstActiveNicName(t *testing.T) {
	t.Run("returns a non-empty NIC name", func(t *testing.T) {
		ctx := contextutils.GetVerbosityContextByBool(true)

		nicName, err := commandexecutornicutils.GetFirstActiveNicName(ctx, commandexecutorexecoo.Exec())

		require.NoError(t, err)
		require.NotEmpty(t, nicName)
	})
}
