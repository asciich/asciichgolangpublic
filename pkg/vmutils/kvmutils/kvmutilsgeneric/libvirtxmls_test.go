package kvmutilsgeneric_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsoptions"
	libvirtxml "libvirt.org/libvirt-go-xml"
)

// Returns valid create options with 4 GB RAM and 4 CPUs.
// The disk image is a real (empty) file in a temporary directory, so the
// test does not depend on whether SetDiskImageByPath requires an existing file.
func getTestCreateOptions(t *testing.T) *kvmutilsoptions.KvmCreateVmOptions {
	t.Helper()

	diskPath := filepath.Join(t.TempDir(), "testvm.qcow2")
	require.NoError(t, os.WriteFile(diskPath, []byte{}, 0o644))

	options := &kvmutilsoptions.KvmCreateVmOptions{
		VmName:      "testvm",
		MacAddress:  "96:56:d6:0e:04:11",
		VmRamSizeGb: 4,
		CpuCount:    4,
	}
	require.NoError(t, options.SetDiskImageByPath(diskPath))

	return options
}

// Renders the XML for the given options and parses it back into a libvirt domain.
func renderAndParseDomain(t *testing.T, options *kvmutilsoptions.KvmCreateVmOptions) *libvirtxml.Domain {
	t.Helper()

	xmlString, err := kvmutilsgeneric.CreateXmlForVmAsString(options)
	require.NoError(t, err)

	domain := &libvirtxml.Domain{}
	require.NoError(t, domain.Unmarshal(xmlString))

	require.NotNil(t, domain.Memory)
	require.NotNil(t, domain.CurrentMemory)
	require.NotNil(t, domain.VCPU)

	return domain
}

func TestCreateXmlForVmAsString_RamAndCpu(t *testing.T) {
	t.Run("RAM and CPU are rendered", func(t *testing.T) {
		domain := renderAndParseDomain(t, getTestCreateOptions(t))

		require.EqualValues(t, 4*1024*1024, domain.Memory.Value)
		require.Equal(t, "KiB", domain.Memory.Unit)
		require.EqualValues(t, 4*1024*1024, domain.CurrentMemory.Value)
		require.Equal(t, "KiB", domain.CurrentMemory.Unit)
		require.EqualValues(t, 4, domain.VCPU.Value)
	})

	t.Run("Unset RAM and CPU fall back to defaults", func(t *testing.T) {
		options := getTestCreateOptions(t)
		options.VmRamSizeGb = 0
		options.CpuCount = 0

		domain := renderAndParseDomain(t, options)

		require.EqualValues(t, kvmutilsoptions.DefaultVmRamSizeGb*1024*1024, domain.Memory.Value) // 2097152
		require.EqualValues(t, kvmutilsoptions.DefaultVmRamSizeGb*1024*1024, domain.CurrentMemory.Value)
		require.EqualValues(t, kvmutilsoptions.DefaultCpuCount, domain.VCPU.Value) // 2
	})

	t.Run("Negative RAM size results in error", func(t *testing.T) {
		options := getTestCreateOptions(t)
		options.VmRamSizeGb = -1
		_, err := kvmutilsgeneric.CreateXmlForVmAsString(options)
		require.Error(t, err)
	})

	t.Run("Negative CPU count results in error", func(t *testing.T) {
		options := getTestCreateOptions(t)
		options.CpuCount = -1
		_, err := kvmutilsgeneric.CreateXmlForVmAsString(options)
		require.Error(t, err)
	})
}
