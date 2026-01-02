package packagemanageroptions

type InstallPackageOptions struct {
	// Set this to true will update the package database before
	// The installation is done.
	UpdateDatabaseFirst bool

	// Also check for updates and install a package update.
	UpdatePackage bool

	// Forces an installation
	Force bool

	// Use sudo to perform installation as root:
	UseSudo bool
}
