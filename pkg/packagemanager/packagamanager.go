package packagemanager

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanagergeneric"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
)

// InstallPackages installs the given packages inside the given container.
func InstallPackages(ctx context.Context, container containerinterfaces.Container, packageNames []string, options *packagemanageroptions.InstallPackageOptions) error {
	packageManager, err := packagemanagergeneric.NewPackageManagerGeneric(ctx, container)
	if err != nil {
		return err
	}

	return packageManager.InstallPackages(ctx, packageNames, options)
}

// IsPackagesInstalled returns true when all given packages are installed in the given container.
func IsPackagesInstalled(ctx context.Context, container containerinterfaces.Container, packageNames []string) (bool, error) {
	packageManager, err := packagemanagergeneric.NewPackageManagerGeneric(ctx, container)
	if err != nil {
		return false, err
	}

	return packageManager.IsPackagesInstalled(ctx, packageNames)
}

// IsPackageInstalled returns true when the given package is installed in the given container.
func IsPackageInstalled(ctx context.Context, container containerinterfaces.Container, packageName string) (bool, error) {
	packageManager, err := packagemanagergeneric.NewPackageManagerGeneric(ctx, container)
	if err != nil {
		return false, err
	}

	return packageManager.IsPackageInstalled(ctx, packageName)
}

// IsPackagesUpdateAvailable returns true when an update is available for at least one of the given packages in the given container.
func IsPackagesUpdateAvailable(ctx context.Context, container containerinterfaces.Container, packageNames []string, options *packagemanageroptions.UpdateDatabaseOptions) (bool, error) {
	packageManager, err := packagemanagergeneric.NewPackageManagerGeneric(ctx, container)
	if err != nil {
		return false, err
	}

	return packageManager.IsPackagesUpdateAvailable(ctx, packageNames, options)
}

// IsPackageUpdateAvailable returns true when an update is available for the given package in the given container.
func IsPackageUpdateAvailable(ctx context.Context, container containerinterfaces.Container, packageName string, options *packagemanageroptions.UpdateDatabaseOptions) (bool, error) {
	packageManager, err := packagemanagergeneric.NewPackageManagerGeneric(ctx, container)
	if err != nil {
		return false, err
	}

	return packageManager.IsPackageUpdateAvailable(ctx, packageName, options)
}

// UpdateDatabase updates the package database/index in the given container.
func UpdateDatabase(ctx context.Context, container containerinterfaces.Container, options *packagemanageroptions.UpdateDatabaseOptions) error {
	packageManager, err := packagemanagergeneric.NewPackageManagerGeneric(ctx, container)
	if err != nil {
		return err
	}

	return packageManager.UpdateDatabase(ctx, options)
}

// UpdatePackages updates the given packages in the given container.
func UpdatePackages(ctx context.Context, container containerinterfaces.Container, packageNames []string, options *packagemanageroptions.UpdatePackageOptions) error {
	packageManager, err := packagemanagergeneric.NewPackageManagerGeneric(ctx, container)
	if err != nil {
		return err
	}

	return packageManager.UpdatePackages(ctx, packageNames, options)
}
