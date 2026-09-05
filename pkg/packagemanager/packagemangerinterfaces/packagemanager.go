package packagemangerinterfaces

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
)

type PackageManager interface {
	// Underlying command executor used by this package manager.
	GetCommandExecutor() (commandexecutorinterfaces.CommandExecutor, error)

	// Install the given packages.
	InstallPackages(ctx context.Context, packageNames []string, options *packagemanageroptions.InstallPackageOptions) error

	// Returns true when all given packages are installed.
	IsPackagesInstalled(ctx context.Context, packageNames []string) (bool, error)

	// Returns true when the given package is installed.
	IsPackageInstalled(ctx context.Context, packageName string) (bool, error)

	// Returns true when an update is available for at least one of the given packages.
	IsPackagesUpdateAvailable(ctx context.Context, packageNames []string, options *packagemanageroptions.UpdateDatabaseOptions) (bool, error)

	// Returns true when an update is available for the given package.
	IsPackageUpdateAvailable(ctx context.Context, packageName string, options *packagemanageroptions.UpdateDatabaseOptions) (bool, error)

	// Update the package database/index.
	UpdateDatabase(ctx context.Context, options *packagemanageroptions.UpdateDatabaseOptions) error

	// Update the given packages.
	UpdatePackages(ctx context.Context, packageNames []string, options *packagemanageroptions.UpdatePackageOptions) error
}
