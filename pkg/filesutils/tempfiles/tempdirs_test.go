package tempfiles_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func Test_CreateTempDir(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		tempDirPath, err := tempfiles.CreateTempDir(context.TODO())
		require.NoError(t, err)
		require.True(t, strings.HasPrefix(tempDirPath, "/tmp/"))
	})
}
