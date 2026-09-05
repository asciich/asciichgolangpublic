package packagemanagergeneric

import (
	"context"
	"runtime"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/commandexecutorfile"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/aptget"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/pacman"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/yay"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// PackageManagerType represents the type of package manager detected
type PackageManagerType string

const (
	PackageManagerAptGet  PackageManagerType = "apt-get"
	PackageManagerYay     PackageManagerType = "yay"
	PackageManagerPacman  PackageManagerType = "pacman"
	PackageManagerUnknown PackageManagerType = "unknown"
)

// PackageManagerGeneric provides automatic package manager selection based on OS detection
type PackageManagerGeneric struct {
	commandExecutor commandexecutorinterfaces.CommandExecutor
	packageType     PackageManagerType
}

// NewPackageManagerGeneric creates a new generic package manager that auto-detects the OS
func NewPackageManagerGeneric(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) (*PackageManagerGeneric, error) {
	ret := new(PackageManagerGeneric)

	err := ret.SetCommandExecutor(commandExecutor)
	if err != nil {
		return nil, err
	}

	err = ret.detectPackageManager(ctx)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// SetCommandExecutor sets the command executor
func (p *PackageManagerGeneric) SetCommandExecutor(commandExecutor commandexecutorinterfaces.CommandExecutor) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	p.commandExecutor = commandExecutor
	return nil
}

// GetCommandExecutor returns the command executor
func (p *PackageManagerGeneric) GetCommandExecutor() (commandexecutorinterfaces.CommandExecutor, error) {
	if p.commandExecutor == nil {
		return nil, tracederrors.TracedError("commandExecutor not set")
	}
	return p.commandExecutor, nil
}

// GetPackageManagerType returns the detected package manager type
func (p *PackageManagerGeneric) GetPackageManagerType() PackageManagerType {
	return p.packageType
}

// detectPackageManager detects the appropriate package manager based on the OS
func (p *PackageManagerGeneric) detectPackageManager(ctx context.Context) error {
	if runtime.GOOS != "linux" {
		p.packageType = PackageManagerUnknown
		return tracederrors.TracedErrorf("unsupported OS: %s", runtime.GOOS)
	}

	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return err
	}

	// Read /etc/os-release to detect the distribution
	osRelease, err := commandexecutorfile.ReadAsBytes(commandExecutor, "/etc/os-release")
	if err != nil {
		return tracederrors.TracedErrorf("failed to read /etc/os-release: %w", err)
	}

	osReleaseContent := string(osRelease)

	// Check for Ubuntu or Debian-based systems
	if strings.Contains(osReleaseContent, `ID=ubuntu`) || strings.Contains(osReleaseContent, `ID=debian`) ||
		strings.Contains(osReleaseContent, `ID_LIKE="ubuntu"`) || strings.Contains(osReleaseContent, `ID_LIKE="debian"`) {
		p.packageType = PackageManagerAptGet
		return nil
	}

	// Check for Arch Linux or Manjaro
	if strings.Contains(osReleaseContent, `ID=arch`) || strings.Contains(osReleaseContent, `ID=manjaro`) ||
		strings.Contains(osReleaseContent, `ID_LIKE="arch"`) {
		// Prefer yay if available, fallback to pacman
		if isCommandAvailable(ctx, p.commandExecutor, "yay") {
			p.packageType = PackageManagerYay
		} else {
			p.packageType = PackageManagerPacman
		}
		return nil
	}

	p.packageType = PackageManagerUnknown
	return tracederrors.TracedErrorf("unsupported Linux distribution")
}

// isCommandAvailable checks if a command is available in the system
func isCommandAvailable(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, command string) bool {
	_, err := commandExecutor.RunCommandAndGetStdoutAsString(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"which", command},
		},
	)
	return err == nil
}

// InstallPackages installs packages using the detected package manager
func (p *PackageManagerGeneric) InstallPackages(ctx context.Context, packageNames []string, options *packagemanageroptions.InstallPackageOptions) error {
	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return err
	}

	switch p.packageType {
	case PackageManagerAptGet:
		return aptget.InstallPackages(ctx, commandExecutor, packageNames, options)
	case PackageManagerYay:
		return yay.InstallPackages(ctx, commandExecutor, packageNames, options)
	case PackageManagerPacman:
		return pacman.InstallPackages(ctx, commandExecutor, packageNames, options)
	default:
		return tracederrors.TracedErrorf("unsupported package manager: %s", p.packageType)
	}
}

// IsPackagesInstalled returns true when all given packages are installed using the detected package manager
func (p *PackageManagerGeneric) IsPackagesInstalled(ctx context.Context, packageNames []string) (bool, error) {
	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	switch p.packageType {
	case PackageManagerAptGet:
		return aptget.IsPackagesInstalled(ctx, commandExecutor, packageNames)
	case PackageManagerYay:
		return yay.IsPackagesInstalled(ctx, commandExecutor, packageNames)
	case PackageManagerPacman:
		return pacman.IsPackagesInstalled(ctx, commandExecutor, packageNames)
	default:
		return false, tracederrors.TracedErrorf("unsupported package manager: %s", p.packageType)
	}
}

// IsPackageInstalled returns true when the given package is installed using the detected package manager
func (p *PackageManagerGeneric) IsPackageInstalled(ctx context.Context, packageName string) (bool, error) {
	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	switch p.packageType {
	case PackageManagerAptGet:
		return aptget.IsPackageInstalled(ctx, commandExecutor, packageName)
	case PackageManagerYay:
		return yay.IsPackageInstalled(ctx, commandExecutor, packageName)
	case PackageManagerPacman:
		return pacman.IsPackageInstalled(ctx, commandExecutor, packageName)
	default:
		return false, tracederrors.TracedErrorf("unsupported package manager: %s", p.packageType)
	}
}

// IsPackagesUpdateAvailable returns true when an update is available for at least one of the given packages using the detected package manager
func (p *PackageManagerGeneric) IsPackagesUpdateAvailable(ctx context.Context, packageNames []string, options *packagemanageroptions.UpdateDatabaseOptions) (bool, error) {
	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	switch p.packageType {
	case PackageManagerAptGet:
		return aptget.IsPackagesUpdateAvailable(ctx, commandExecutor, packageNames, options)
	case PackageManagerYay:
		return yay.IsPackagesUpdateAvailable(ctx, commandExecutor, packageNames, options)
	case PackageManagerPacman:
		return pacman.IsPackagesUpdateAvailable(ctx, commandExecutor, packageNames, options)
	default:
		return false, tracederrors.TracedErrorf("unsupported package manager: %s", p.packageType)
	}
}

// IsPackageUpdateAvailable returns true when an update is available for the given package using the detected package manager
func (p *PackageManagerGeneric) IsPackageUpdateAvailable(ctx context.Context, packageName string, options *packagemanageroptions.UpdateDatabaseOptions) (bool, error) {
	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return false, err
	}

	switch p.packageType {
	case PackageManagerAptGet:
		return aptget.IsPackageUpdateAvailable(ctx, commandExecutor, packageName, options)
	case PackageManagerYay:
		return yay.IsPackageUpdateAvailable(ctx, commandExecutor, packageName, options)
	case PackageManagerPacman:
		return pacman.IsPackageUpdateAvailable(ctx, commandExecutor, packageName, options)
	default:
		return false, tracederrors.TracedErrorf("unsupported package manager: %s", p.packageType)
	}
}

// UpdateDatabase updates the package database/index using the detected package manager
func (p *PackageManagerGeneric) UpdateDatabase(ctx context.Context, options *packagemanageroptions.UpdateDatabaseOptions) error {
	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return err
	}

	switch p.packageType {
	case PackageManagerAptGet:
		return aptget.UpdateDatabase(ctx, commandExecutor, options)
	case PackageManagerYay:
		return yay.UpdateDatabase(ctx, commandExecutor, options)
	case PackageManagerPacman:
		return pacman.UpdateDatabase(ctx, commandExecutor, options)
	default:
		return tracederrors.TracedErrorf("unsupported package manager: %s", p.packageType)
	}
}

// UpdatePackages updates packages using the detected package manager
func (p *PackageManagerGeneric) UpdatePackages(ctx context.Context, packageNames []string, options *packagemanageroptions.UpdatePackageOptions) error {
	commandExecutor, err := p.GetCommandExecutor()
	if err != nil {
		return err
	}

	switch p.packageType {
	case PackageManagerAptGet:
		return aptget.UpdatePackages(ctx, commandExecutor, packageNames, options)
	case PackageManagerYay:
		return yay.UpdatePackages(ctx, commandExecutor, packageNames, options)
	case PackageManagerPacman:
		return pacman.UpdatePackages(ctx, commandExecutor, packageNames, options)
	default:
		return tracederrors.TracedErrorf("unsupported package manager: %s", p.packageType)
	}
}
