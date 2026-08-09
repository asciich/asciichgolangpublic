package commandexecutorfile_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
)

func TestDelete(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Delete(ctx, nil, "/tmp/test", &filesoptions.DeleteOptions{})
		require.Error(t, err)
	})
	t.Run("empty filePath returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), "", &filesoptions.DeleteOptions{})
		require.Error(t, err)
	})
	t.Run("relative path returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.Delete(ctx, commandexecutorexecoo.Exec(), "relative/path", &filesoptions.DeleteOptions{})
		require.Error(t, err)
	})
}

func TestDeleteDirectory(t *testing.T) {
	t.Run("nil commandExecutor returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.DeleteDirectory(ctx, nil, "/tmp/test", &filesoptions.DeleteOptions{})
		require.Error(t, err)
	})
	t.Run("empty filePath returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.DeleteDirectory(ctx, commandexecutorexecoo.Exec(), "", &filesoptions.DeleteOptions{})
		require.Error(t, err)
	})
	t.Run("relative path returns error", func(t *testing.T) {
		ctx := getCtx()
		err := commandexecutorfile.DeleteDirectory(ctx, commandexecutorexecoo.Exec(), "relative/path", &filesoptions.DeleteOptions{})
		require.Error(t, err)
	})
}
