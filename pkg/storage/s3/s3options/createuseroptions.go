package s3options

type CreateUserOptions struct {

	// If the user already exists
	KeepCurrentPasswordIfUserExists bool

	// Give the user read only permissions.
	// This user can read all buckets and all objects in it.
	ReadOnly bool
}
