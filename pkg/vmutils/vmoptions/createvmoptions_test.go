package vmoptions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsoptions"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/vmoptions"
)

func Test_CreateVmOptions_AdditionalBlockDevices(t *testing.T) {
	t.Run("unset returns empty slice", func(t *testing.T) {
		options := vmoptions.NewCreateVmOptions()
		blockDevices := options.GetAdditionalBlockDevices()
		require.NotNil(t, blockDevices)
		require.Empty(t, blockDevices)
	})

	t.Run("set and get round trip", func(t *testing.T) {
		blockDevices := []kvmutilsoptions.AdditionalBlockDevice{
			{SourcePath: "/dev/sda"},
			{SourcePath: "/dev/mapper/abc", Serial: "data01", ReadOnly: true},
		}

		options := vmoptions.NewCreateVmOptions()
		err := options.SetAdditionalBlockDevices(blockDevices)
		require.NoError(t, err)
		require.Equal(t, blockDevices, options.GetAdditionalBlockDevices())
	})

	t.Run("modifying the input slice after set does not change the options", func(t *testing.T) {
		blockDevices := []kvmutilsoptions.AdditionalBlockDevice{{SourcePath: "/dev/sda"}}

		options := vmoptions.NewCreateVmOptions()
		err := options.SetAdditionalBlockDevices(blockDevices)
		require.NoError(t, err)

		blockDevices[0].SourcePath = "/dev/sdb"

		require.Equal(t, []kvmutilsoptions.AdditionalBlockDevice{{SourcePath: "/dev/sda"}}, options.GetAdditionalBlockDevices())
	})

	t.Run("modifying the returned slice does not change the options", func(t *testing.T) {
		options := vmoptions.NewCreateVmOptions()
		err := options.SetAdditionalBlockDevices([]kvmutilsoptions.AdditionalBlockDevice{{SourcePath: "/dev/sda"}})
		require.NoError(t, err)

		returned := options.GetAdditionalBlockDevices()
		returned[0].SourcePath = "/dev/sdb"

		require.Equal(t, []kvmutilsoptions.AdditionalBlockDevice{{SourcePath: "/dev/sda"}}, options.GetAdditionalBlockDevices())
	})

	t.Run("invalid block device is rejected and options stay unchanged", func(t *testing.T) {
		options := vmoptions.NewCreateVmOptions()
		err := options.SetAdditionalBlockDevices([]kvmutilsoptions.AdditionalBlockDevice{{SourcePath: "/dev/sda"}})
		require.NoError(t, err)

		err = options.SetAdditionalBlockDevices([]kvmutilsoptions.AdditionalBlockDevice{{SourcePath: "/tmp/disk.img"}})
		require.Error(t, err)
		require.Equal(t, []kvmutilsoptions.AdditionalBlockDevice{{SourcePath: "/dev/sda"}}, options.GetAdditionalBlockDevices())
	})

	t.Run("nil clears block devices", func(t *testing.T) {
		options := vmoptions.NewCreateVmOptions()
		err := options.SetAdditionalBlockDevices([]kvmutilsoptions.AdditionalBlockDevice{{SourcePath: "/dev/sda"}})
		require.NoError(t, err)

		err = options.SetAdditionalBlockDevices(nil)
		require.NoError(t, err)
		require.Empty(t, options.GetAdditionalBlockDevices())
	})

	t.Run("add appends in order", func(t *testing.T) {
		options := vmoptions.NewCreateVmOptions()
		require.NoError(t, options.AddAdditionalBlockDevice(kvmutilsoptions.AdditionalBlockDevice{SourcePath: "/dev/sda"}))
		require.NoError(t, options.AddAdditionalBlockDevice(kvmutilsoptions.AdditionalBlockDevice{SourcePath: "/dev/sdb1", Serial: "data01"}))

		require.Equal(t, []kvmutilsoptions.AdditionalBlockDevice{
			{SourcePath: "/dev/sda"},
			{SourcePath: "/dev/sdb1", Serial: "data01"},
		}, options.GetAdditionalBlockDevices())
	})

	t.Run("add rejects duplicate source path and keeps existing devices", func(t *testing.T) {
		options := vmoptions.NewCreateVmOptions()
		require.NoError(t, options.AddAdditionalBlockDevice(kvmutilsoptions.AdditionalBlockDevice{SourcePath: "/dev/sda"}))

		err := options.AddAdditionalBlockDevice(kvmutilsoptions.AdditionalBlockDevice{SourcePath: "/dev/sda"})
		require.Error(t, err)
		require.Len(t, options.GetAdditionalBlockDevices(), 1)
	})
}
