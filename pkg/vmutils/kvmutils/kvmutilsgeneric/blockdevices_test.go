package kvmutilsgeneric_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/vmutils/kvmutils/kvmutilsoptions"
	libvirtxml "libvirt.org/libvirt-go-xml"
)

func getTestCreateOptionsWithBlockDevicePath(t *testing.T, blockDevices []kvmutilsoptions.AdditionalBlockDevice) *kvmutilsoptions.KvmCreateVmOptions {
	t.Helper()

	ret := kvmutilsoptions.NewKvmCreateVmOptions()

	err := ret.SetVmName("test-vm")
	require.NoError(t, err)

	err = ret.SetDiskImageByPath("/var/lib/libvirt/images/test-vm.qcow2")
	require.NoError(t, err)

	ret.AdditionalBlockDevices = blockDevices

	return ret
}

func getBlockDisksFromXml(t *testing.T, libvirtXml string) []libvirtxml.DomainDisk {
	t.Helper()

	domcfg := &libvirtxml.Domain{}
	err := domcfg.Unmarshal(libvirtXml)
	require.NoError(t, err)
	require.NotNil(t, domcfg.Devices)

	ret := []libvirtxml.DomainDisk{}
	for _, disk := range domcfg.Devices.Disks {
		if disk.Source != nil && disk.Source.Block != nil {
			ret = append(ret, disk)
		}
	}

	return ret
}

func Test_GetVirtioTargetDevName(t *testing.T) {
	tests := []struct {
		index    int
		expected string
	}{
		{0, "vdb"},
		{1, "vdc"},
		{24, "vdz"},
		{25, "vdaa"},
		{26, "vdab"},
		{50, "vdaz"},
		{51, "vdba"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			targetDev, err := kvmutilsgeneric.GetVirtioTargetDevName(tt.index)
			require.NoError(t, err)
			require.Equal(t, tt.expected, targetDev)
		})
	}

	t.Run("negative index", func(t *testing.T) {
		_, err := kvmutilsgeneric.GetVirtioTargetDevName(-1)
		require.Error(t, err)
	})
}

func Test_CheckBlockDeviceSourcePath(t *testing.T) {
	validPaths := []string{
		"/dev/sda",
		"/dev/sda1",
		"/dev/nvme0n1p2",
		"/dev/mapper/abc",
		"/dev/vg0/lv_data",
		"/dev/disk/by-id/ata-Samsung_SSD_860_EVO_500GB_S3Z1NB0K",
		"/dev/disk/by-path/pci-0000:00:17.0-ata-1",
	}

	for _, path := range validPaths {
		t.Run("valid "+path, func(t *testing.T) {
			err := kvmutilsgeneric.CheckBlockDeviceSourcePath(path)
			require.NoError(t, err)
		})
	}

	invalidPaths := []string{
		"",
		"sda",
		"/tmp/disk.img",
		"/dev/",
		"/dev//sda",
		"/dev/sda/",
		"/dev/../etc/passwd",
		"/dev/sda'",
		"/dev/sd a",
		"/dev/sda<",
		"/dev/sda&b",
	}

	for _, path := range invalidPaths {
		t.Run("invalid '"+path+"'", func(t *testing.T) {
			err := kvmutilsgeneric.CheckBlockDeviceSourcePath(path)
			require.Error(t, err)
		})
	}
}

func Test_CheckVirtioBlkSerial(t *testing.T) {
	t.Run("empty is valid", func(t *testing.T) {
		require.NoError(t, kvmutilsgeneric.CheckVirtioBlkSerial(""))
	})

	t.Run("valid serial", func(t *testing.T) {
		require.NoError(t, kvmutilsgeneric.CheckVirtioBlkSerial("data_01-backup"))
	})

	t.Run("max length is valid", func(t *testing.T) {
		require.NoError(t, kvmutilsgeneric.CheckVirtioBlkSerial("abcdefghijklmnopqrst"))
	})

	t.Run("too long", func(t *testing.T) {
		require.Error(t, kvmutilsgeneric.CheckVirtioBlkSerial("abcdefghijklmnopqrstu"))
	})

	t.Run("invalid characters", func(t *testing.T) {
		require.Error(t, kvmutilsgeneric.CheckVirtioBlkSerial("data'01"))
	})
}

func Test_GetBlockDeviceAttachments(t *testing.T) {
	t.Run("nil returns empty slice", func(t *testing.T) {
		attachments, err := kvmutilsgeneric.GetBlockDeviceAttachments(nil)
		require.NoError(t, err)
		require.Empty(t, attachments)
	})

	t.Run("targets assigned in order", func(t *testing.T) {
		attachments, err := kvmutilsgeneric.GetBlockDeviceAttachments([]kvmutilsoptions.AdditionalBlockDevice{
			{SourcePath: "/dev/sda"},
			{SourcePath: "/dev/sdb1", Serial: "data01"},
			{SourcePath: "/dev/mapper/abc", ReadOnly: true},
		})
		require.NoError(t, err)
		require.Equal(t, []kvmutilsgeneric.BlockDeviceAttachment{
			{SourcePath: "/dev/sda", TargetDev: "vdb"},
			{SourcePath: "/dev/sdb1", TargetDev: "vdc", Serial: "data01"},
			{SourcePath: "/dev/mapper/abc", TargetDev: "vdd", ReadOnly: true},
		}, attachments)
	})

	t.Run("duplicate source path", func(t *testing.T) {
		_, err := kvmutilsgeneric.GetBlockDeviceAttachments([]kvmutilsoptions.AdditionalBlockDevice{
			{SourcePath: "/dev/sda"},
			{SourcePath: "/dev/sda"},
		})
		require.Error(t, err)
	})

	t.Run("duplicate serial", func(t *testing.T) {
		_, err := kvmutilsgeneric.GetBlockDeviceAttachments([]kvmutilsoptions.AdditionalBlockDevice{
			{SourcePath: "/dev/sda", Serial: "data"},
			{SourcePath: "/dev/sdb", Serial: "data"},
		})
		require.Error(t, err)
	})

	t.Run("invalid source path", func(t *testing.T) {
		_, err := kvmutilsgeneric.GetBlockDeviceAttachments([]kvmutilsoptions.AdditionalBlockDevice{
			{SourcePath: "/tmp/disk.img"},
		})
		require.Error(t, err)
	})
}

func Test_CreateXmlForVmAsString_AdditionalBlockDevices(t *testing.T) {
	t.Run("no additional block devices", func(t *testing.T) {
		libvirtXml, err := kvmutilsgeneric.CreateXmlForVmAsString(getTestCreateOptionsWithBlockDevicePath(t, nil))
		require.NoError(t, err)

		attachments, err := kvmutilsgeneric.GetAdditionalBlockDevicesFromXmlString(libvirtXml)
		require.NoError(t, err)
		require.Empty(t, attachments)
	})

	t.Run("block devices round trip", func(t *testing.T) {
		libvirtXml, err := kvmutilsgeneric.CreateXmlForVmAsString(getTestCreateOptionsWithBlockDevicePath(t, []kvmutilsoptions.AdditionalBlockDevice{
			{SourcePath: "/dev/sda"},
			{SourcePath: "/dev/sdb1", Serial: "data01"},
			{SourcePath: "/dev/mapper/abc", ReadOnly: true},
		}))
		require.NoError(t, err)

		attachments, err := kvmutilsgeneric.GetAdditionalBlockDevicesFromXmlString(libvirtXml)
		require.NoError(t, err)
		require.Equal(t, []kvmutilsgeneric.BlockDeviceAttachment{
			{SourcePath: "/dev/sda", TargetDev: "vdb"},
			{SourcePath: "/dev/sdb1", TargetDev: "vdc", Serial: "data01"},
			{SourcePath: "/dev/mapper/abc", TargetDev: "vdd", ReadOnly: true},
		}, attachments)
	})

	t.Run("block devices use raw virtio driver without host cache", func(t *testing.T) {
		libvirtXml, err := kvmutilsgeneric.CreateXmlForVmAsString(getTestCreateOptionsWithBlockDevicePath(t, []kvmutilsoptions.AdditionalBlockDevice{
			{SourcePath: "/dev/sda1"},
		}))
		require.NoError(t, err)

		blockDisks := getBlockDisksFromXml(t, libvirtXml)
		require.Len(t, blockDisks, 1)

		disk := blockDisks[0]
		require.Equal(t, "disk", disk.Device)
		require.NotNil(t, disk.Driver)
		require.Equal(t, "raw", disk.Driver.Type)
		require.Equal(t, "none", disk.Driver.Cache)
		require.Equal(t, "native", disk.Driver.IO)
		require.NotNil(t, disk.Target)
		require.Equal(t, "virtio", disk.Target.Bus)
	})

	t.Run("root disk stays vda and is the only boot device", func(t *testing.T) {
		libvirtXml, err := kvmutilsgeneric.CreateXmlForVmAsString(getTestCreateOptionsWithBlockDevicePath(t, []kvmutilsoptions.AdditionalBlockDevice{
			{SourcePath: "/dev/sda"},
		}))
		require.NoError(t, err)

		domcfg := &libvirtxml.Domain{}
		err = domcfg.Unmarshal(libvirtXml)
		require.NoError(t, err)
		require.NotNil(t, domcfg.Devices)
		require.Len(t, domcfg.Devices.Disks, 2)

		for _, disk := range domcfg.Devices.Disks {
			require.NotNil(t, disk.Target)

			if disk.Target.Dev == "vda" {
				require.NotNil(t, disk.Boot)
				require.EqualValues(t, 1, disk.Boot.Order)
			} else {
				require.Nil(t, disk.Boot)
			}
		}
	})

	t.Run("invalid block device fails", func(t *testing.T) {
		libvirtXml, err := kvmutilsgeneric.CreateXmlForVmAsString(getTestCreateOptionsWithBlockDevicePath(t, []kvmutilsoptions.AdditionalBlockDevice{
			{SourcePath: "/dev/sda'/><disk"},
		}))
		require.Error(t, err)
		require.Empty(t, libvirtXml)
	})
}

func Test_ParseAdditionalBlockDeviceString(t *testing.T) {
	validTests := []struct {
		input    string
		expected kvmutilsoptions.AdditionalBlockDevice
	}{
		{"/dev/sdb", kvmutilsoptions.AdditionalBlockDevice{SourcePath: "/dev/sdb"}},
		{"/dev/sdb1,ro", kvmutilsoptions.AdditionalBlockDevice{SourcePath: "/dev/sdb1", ReadOnly: true}},
		{"/dev/mapper/abc,serial=data01", kvmutilsoptions.AdditionalBlockDevice{SourcePath: "/dev/mapper/abc", Serial: "data01"}},
		{"/dev/mapper/abc,serial=data01,ro", kvmutilsoptions.AdditionalBlockDevice{SourcePath: "/dev/mapper/abc", Serial: "data01", ReadOnly: true}},
		{"/dev/mapper/abc,ro,serial=data01", kvmutilsoptions.AdditionalBlockDevice{SourcePath: "/dev/mapper/abc", Serial: "data01", ReadOnly: true}},
		{"/dev/disk/by-path/pci-0000:00:17.0-ata-1,serial=data01", kvmutilsoptions.AdditionalBlockDevice{SourcePath: "/dev/disk/by-path/pci-0000:00:17.0-ata-1", Serial: "data01"}},
	}

	for _, tt := range validTests {
		t.Run("valid "+tt.input, func(t *testing.T) {
			blockDevice, err := kvmutilsgeneric.ParseAdditionalBlockDeviceString(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.expected, blockDevice)
		})
	}

	invalidInputs := []string{
		"",
		"/tmp/disk.img",
		"/dev/sdb,",
		"/dev/sdb,rw",
		"/dev/sdb,ro=true",
		"/dev/sdb,ro,ro",
		"/dev/sdb,serial",
		"/dev/sdb,serial=",
		"/dev/sdb,serial=a,serial=b",
		"/dev/sdb,serial=data'01",
		"/dev/sdb,serial=abcdefghijklmnopqrstu",
	}

	for _, input := range invalidInputs {
		t.Run("invalid '"+input+"'", func(t *testing.T) {
			_, err := kvmutilsgeneric.ParseAdditionalBlockDeviceString(input)
			require.Error(t, err)
		})
	}
}
