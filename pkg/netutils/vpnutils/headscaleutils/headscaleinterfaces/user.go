package headscaleinterfaces

import "time"

// User represents a single user as listed by `headscale users list`.
type User interface {
	// GetID returns the unique ID of the user (e.g. 1).
	GetID() (uint64, error)

	// GetName returns the display name of the user (column "Name").
	// Returns an empty string if no display name is set.
	GetName() (string, error)

	// GetUsername returns the username (column "Username", e.g. "reto").
	GetUsername() (string, error)

	// GetEmail returns the email address of the user.
	// Returns an empty string if no email is set.
	GetEmail() (string, error)

	// GetCreated returns the time the user was created.
	GetCreated() (time.Time, error)
}
