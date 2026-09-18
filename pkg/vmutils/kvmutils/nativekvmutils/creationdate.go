package nativekvmutils

import (
	"context"
	"encoding/xml"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// ---------------------------------------------------------------------------
// VM creation date
// ---------------------------------------------------------------------------

// domainDiskXml is a minimal representation of the libvirt domain XML used to
// extract the primary disk image path. Only the fields required to locate the
// backing file are modelled.
type domainDiskXml struct {
	Devices struct {
		Disks []struct {
			Device string `xml:"device,attr"`
			Source struct {
				File string `xml:"file,attr"`
			} `xml:"source"`
		} `xml:"disk"`
	} `xml:"devices"`
}

// getVmDiskImagePath returns the path of the primary disk image backing the
// given VM by parsing its domain XML.
func (k *NativeKvmHypervisor) getVmDiskImagePath(ctx context.Context, vmName string) (diskImagePath string, err error) {
	if vmName == "" {
		return "", tracederrors.TracedErrorEmptyString("vmName")
	}

	domainXml, err := k.GetDomainXmlAsString(ctx, vmName)
	if err != nil {
		return "", err
	}

	parsed := &domainDiskXml{}
	err = xml.Unmarshal([]byte(domainXml), parsed)
	if err != nil {
		return "", tracederrors.TracedErrorf("failed to parse domain XML of VM '%s' to get disk image path: %w", vmName, err)
	}

	for _, disk := range parsed.Devices.Disks {
		// Only consider actual disk devices with a file backing (skip cdrom,
		// floppy and network/block backed disks without a 'file' source).
		if disk.Device != "disk" {
			continue
		}

		if disk.Source.File == "" {
			continue
		}

		return disk.Source.File, nil
	}

	return "", tracederrors.TracedErrorf("No file backed disk image found in domain XML of VM '%s'.", vmName)
}

// GetVmCreationDate returns the creation date of the VM derived from the birth
// date of its primary disk image file. libvirt itself does not store a VM
// creation timestamp, so the disk image's birth date is used as the best
// available proxy.
func (k *NativeKvmHypervisor) GetVmCreationDate(ctx context.Context, vmName string) (creationDate time.Time, err error) {
	if vmName == "" {
		return time.Time{}, tracederrors.TracedErrorEmptyString("vmName")
	}

	diskImagePath, err := k.getVmDiskImagePath(ctx, vmName)
	if err != nil {
		return time.Time{}, err
	}

	creationDate, err = nativefiles.GetBirthDate(ctx, diskImagePath)
	if err != nil {
		return time.Time{}, err
	}

	logging.LogInfoByCtxf(ctx, "Creation date of VM '%s' is '%v'.", vmName, creationDate)

	return creationDate, nil
}

// GetVmCreationDateRFC3339 returns the creation date of the VM as an RFC3339
// formatted string. See GetVmCreationDate for details on how the date is
// derived.
func (k *NativeKvmHypervisor) GetVmCreationDateRFC3339(ctx context.Context, vmName string) (creationDate string, err error) {
	if vmName == "" {
		return "", tracederrors.TracedErrorEmptyString("vmName")
	}

	diskImagePath, err := k.getVmDiskImagePath(ctx, vmName)
	if err != nil {
		return "", err
	}

	creationDate, err = nativefiles.GetBirthDateRFC3339(ctx, diskImagePath)
	if err != nil {
		return "", err
	}

	logging.LogInfoByCtxf(ctx, "Creation date of VM '%s' in RFC3339 is '%s'.", vmName, creationDate)

	return creationDate, nil
}
