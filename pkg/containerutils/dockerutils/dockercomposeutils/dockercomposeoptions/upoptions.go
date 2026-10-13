package dockercomposeoptions

type UpOptions struct {
	// Remove containers of services which belong to the compose project
	// but are no longer defined in the compose file.
	// Equivalent to 'docker compose up --remove-orphans'.
	RemoveOrphans bool
}
