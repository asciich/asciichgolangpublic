package fstabconfigfile_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/fstabconfigfile"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
)

func getCtx() context.Context {
	return contextutils.ContextVerbose()
}

func Test_ReadFromFile(t *testing.T) {
	t.Run("example1", func(t *testing.T) {
		ctx := getCtx()

		content := `# <device>                                <dir> <type> <options>                                        <dump> <fsck>
UUID=0a3407de-014b-458b-b5c1-848e92a327a3 /     ext4 defaults                                           0      1
UUID=CBB6-24F2                            /boot vfat defaults,nodev,nosuid,noexec,fmask=0177,dmask=0077 0      2
UUID=f9fe0b69-a280-415d-a03a-a32752370dee none  swap defaults                                           0      0
UUID=b411dc99-f0a0-4c87-9e05-184977be8539 /home ext4 defaults                                           0      2
`
		tmpPath, err := tempfiles.CreateTemporaryFileFromContentString(ctx, content)
		require.NoError(t, err)

		defer func() { _ = nativefiles.Delete(ctx, tmpPath, &filesoptions.DeleteOptions{}) }()

		entries, err := fstabconfigfile.ReadFromFile(ctx, tmpPath)
		require.NoError(t, err)

		require.Len(t, entries, 4)
		
		require.EqualValues(t, "UUID=0a3407de-014b-458b-b5c1-848e92a327a3", entries[0].Device)
		require.EqualValues(t, "/", entries[0].Dir)
		require.EqualValues(t, "ext4", entries[0].Type)
		require.EqualValues(t, "defaults", entries[0].Options)
		require.EqualValues(t, "0", entries[0].Dump)
		require.EqualValues(t, "1", entries[0].Fsck)
				
		require.EqualValues(t, "UUID=CBB6-24F2", entries[1].Device)
		require.EqualValues(t, "/boot", entries[1].Dir)
		require.EqualValues(t, "vfat", entries[1].Type)
		require.EqualValues(t, "defaults,nodev,nosuid,noexec,fmask=0177,dmask=0077", entries[1].Options)
		require.EqualValues(t, "0", entries[1].Dump)
		require.EqualValues(t, "2", entries[1].Fsck)
		
		require.EqualValues(t, "UUID=f9fe0b69-a280-415d-a03a-a32752370dee", entries[2].Device)
		require.EqualValues(t, "none", entries[2].Dir)
		require.EqualValues(t, "swap", entries[2].Type)
		require.EqualValues(t, "defaults", entries[2].Options)
		require.EqualValues(t, "0", entries[2].Dump)
		require.EqualValues(t, "0", entries[2].Fsck)

		require.EqualValues(t, "UUID=b411dc99-f0a0-4c87-9e05-184977be8539", entries[3].Device)
		require.EqualValues(t, "/home", entries[3].Dir)
		require.EqualValues(t, "ext4", entries[3].Type)
		require.EqualValues(t, "defaults", entries[3].Options)
		require.EqualValues(t, "0", entries[3].Dump)
		require.EqualValues(t, "2", entries[3].Fsck)
	})
}
