package commandexecutorfileoo_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/checksumutils"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfileoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func Test_CommandExecutorFileOo_GetSha256Sum(t *testing.T) {
	ctx := getCtx()

	temporaryFilePath, err := tempfiles.CreateTemporaryFile(ctx)
	require.NoError(t, err)
	defer nativefiles.Delete(ctx, temporaryFilePath, &filesoptions.DeleteOptions{})

	const content = "hello world"
	err = nativefiles.WriteString(ctx, temporaryFilePath, content)
	require.NoError(t, err)

	expectedSha256 := checksumutils.GetSha256SumFromString(content)
	file, err := commandexecutorfileoo.New(commandexecutorexecoo.Exec(), temporaryFilePath)
	require.NoError(t, err)
	sha256sum, err := file.GetSha256Sum(ctx)
	require.NoError(t, err)
	require.EqualValues(t, expectedSha256, sha256sum)
}
