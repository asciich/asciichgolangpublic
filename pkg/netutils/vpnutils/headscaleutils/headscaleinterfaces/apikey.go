package headscaleinterfaces

import "time"

// ApiKey represents a single API key as listed by `headscale apikeys list`.
//
// Only the prefix of an API key is available after creation.
// The full key is shown once when it is created and can not be retrieved later.
type ApiKey interface {
	// GetID returns the unique ID of the API key.
	GetID() (uint64, error)

	// GetPrefix returns the prefix of the API key.
	GetPrefix() (string, error)

	// GetExpiration returns the time the API key expires.
	// Returns nil if no expiration is set (headscale zero value 0001-01-01 00:00:00).
	GetExpiration() (*time.Time, error)

	// GetCreated returns the time the API key was created.
	GetCreated() (time.Time, error)
}
