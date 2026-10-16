package hostreconcileroptions

// RunOptions configures a single host reconciliation run, for example
// one triggered by 'hostreconciler reconcile-once'.
//
// The zero value is valid and gives the default behavior: idempotent
// reconciliation, which changes only what differs from the desired state.
type RunOptions struct {
	// Recreate all docker compose containers even if their configuration
	// and image haven't changed.
	// Equivalent to 'docker compose up --force-recreate' and set by the
	// CLI flag '--docker-compose-force-recreate'.
	//
	// Note: If set, the reconciliation is intentionally NOT idempotent.
	// Every run recreates the docker compose containers.
	DockerComposeForceRecreate bool
}
