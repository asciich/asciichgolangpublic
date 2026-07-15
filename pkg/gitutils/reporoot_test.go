package gitutils_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils"
)

func Test_GetRepositoryRootPathByPath(t *testing.T) {
	absPath, err := filepath.Abs(".")
	require.NoError(t, err)

	tests := []struct {
		Path string
	}{
		{"."},
		{absPath},
	}

	for _, tt := range tests {
		ctx := getCtx()

		repoRootPath, err := gitutils.GetRepositoryRootPathByPath(ctx, tt.Path)
		require.NoError(t, err)
		require.EqualValues(t, "asciichgolangpublic", filepath.Base(repoRootPath))
	}
}
