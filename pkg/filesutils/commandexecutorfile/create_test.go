package commandexecutorfile_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
)

func TestCreateFile(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.CreateFile(ctx, nil, "/tmp/test", &filesoptions.CreateOptions{})
		require.Error(t, err)
	})
	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.CreateFile(ctx, commandexecutorexecoo.Exec(), "", &filesoptions.CreateOptions{})
		require.Error(t, err)
	})
}

func TestCreateDirectory(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.CreateDirectory(ctx, nil, "/tmp/test", &filesoptions.CreateOptions{})
		require.Error(t, err)
	})
	t.Run("empty path returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.CreateDirectory(ctx, commandexecutorexecoo.Exec(), "", &filesoptions.CreateOptions{})
		require.Error(t, err)
	})
}
