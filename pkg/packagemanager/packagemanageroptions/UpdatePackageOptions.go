package packagemanageroptions

type UpdatePackageOptions struct {
	// Set this to true will update the package database before
	// The update is done.
	UpdateDatabaseFirst bool

	// Force package update:
	Force bool

	// Use sudo to perform package update:
	UseSudo bool
}
