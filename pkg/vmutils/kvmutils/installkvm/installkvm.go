package installkvm

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxutils/systemdutils"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanagergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/userutils"
)

// Name of the libvirt daemon systemd service.
const libvirtdServiceName = "libvirtd"

// Returns the KVM related package names for the detected package manager.
func getPackageNames(packageManager *packagemanagergeneric.PackageManagerGeneric) ([]string, error) {
	if packageManager == nil {
		return nil, tracederrors.TracedErrorNil("packageManager")
	}

	packageManagerType := packageManager.GetPackageManagerType()

	switch packageManagerType {
	case packagemanagergeneric.PackageManagerAptGet:
		ret := []string{
			"qemu-system-x86",
			"qemu-utils",
			"libvirt-daemon-system",
			"libvirt-clients",
			"virtinst",
			"ovmf",
		}
		return ret, nil
	case packagemanagergeneric.PackageManagerPacman, packagemanagergeneric.PackageManagerYay:
		ret := []string{
			"qemu-system-x86",
			"qemu-img",
			"libvirt",
			"virt-install",
			"dnsmasq",
			"edk2-ovmf",
		}
		return ret, nil
	default:
		return nil, tracederrors.TracedErrorf("Unsupported package manager type '%s' to install KVM packages.", packageManagerType)
	}
}

// Returns an error naming the given operation if the current process is not running as root.
func checkRunningAsRoot(ctx context.Context, operationName string) error {
	if operationName == "" {
		return tracederrors.TracedErrorEmptyString("operationName")
	}

	isRoot, err := userutils.IsRunningAsRoot(ctx)
	if err != nil {
		return err
	}

	if !isRoot {
		return tracederrors.TracedErrorf("%s requires root privileges.", operationName)
	}

	return nil
}

// InstallPackages installs all needed KVM related system packages on localhost.
//
// The function is idempotent: Packages are only installed if not all of them are already present.
//
// Requires root privileges.
func InstallPackages(ctx context.Context) error {
	err := checkRunningAsRoot(ctx, "Install KVM packages")
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Install KVM packages started.")

	packageManager, err := packagemanagergeneric.NewPackageManagerGeneric(ctx, commandexecutorexecoo.Exec())
	if err != nil {
		return err
	}

	packageNames, err := getPackageNames(packageManager)
	if err != nil {
		return err
	}

	isInstalled, err := packageManager.IsPackagesInstalled(ctx, packageNames)
	if err != nil {
		return err
	}

	if isInstalled {
		logging.LogInfoByCtxf(ctx, "All '%d' KVM packages %v are already installed.", len(packageNames), packageNames)
	} else {
		err = packageManager.InstallPackages(ctx, packageNames, &packagemanageroptions.InstallPackageOptions{})
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Installed KVM packages %v.", packageNames)
	}

	logging.LogInfoByCtxf(ctx, "Install KVM packages finished. '%d' KVM packages ensured to be installed.", len(packageNames))

	return nil
}
// Start ensures all needed KVM packages are installed and the libvirtd service is running on localhost.
//
// The libvirtd service is NOT enabled. Use Enable to start libvirtd on boot.
// The function is idempotent: Only missing packages are installed and libvirtd is only started if needed.
//
// Requires root privileges.
func Start(ctx context.Context) error {
	err := checkRunningAsRoot(ctx, "Start KVM")
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Start KVM started.")

	err = InstallPackages(ctx)
	if err != nil {
		return err
	}

	// StartService is idempotent and logs changed/info messages by itself.
	err = systemdutils.StartService(ctx, commandexecutorexecoo.Exec(), libvirtdServiceName)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Start KVM finished. KVM packages are installed and '%s' is running.", libvirtdServiceName)

	return nil
}

// Enable ensures all needed KVM packages are installed and the libvirtd service is enabled (started on boot) on localhost.
//
// The libvirtd service is NOT started. Use Start to ensure libvirtd is running.
// The function is idempotent: Only missing packages are installed and libvirtd is only enabled if needed.
//
// Requires root privileges.
func Enable(ctx context.Context) error {
	err := checkRunningAsRoot(ctx, "Enable KVM")
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Enable KVM started.")

	err = InstallPackages(ctx)
	if err != nil {
		return err
	}

	// EnableService is idempotent and logs changed/info messages by itself.
	err = systemdutils.EnableService(ctx, commandexecutorexecoo.Exec(), libvirtdServiceName)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Enable KVM finished. KVM packages are installed and '%s' is enabled.", libvirtdServiceName)

	return nil
}
