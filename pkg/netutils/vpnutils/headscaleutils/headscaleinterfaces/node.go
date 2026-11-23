package headscaleinterfaces

import "time"

// Node represents a single node as listed by `headscale nodes list`.
//
// Implementations must return IP addresses which are validated using the iputils package.
type Node interface {
	// GetID returns the unique ID of the node (e.g. 1).
	GetID() (uint64, error)

	// GetHostname returns the hostname reported by the client.
	GetHostname() (string, error)

	// GetName returns the name used for MagicDNS (given name).
	GetName() (string, error)

	// GetMachineKey returns the machine key.
	GetMachineKey() (string, error)

	// GetNodeKey returns the node key.
	GetNodeKey() (string, error)

	// GetUser returns the name of the user owning the node.
	GetUser() (string, error)

	// GetTags returns the ACL tags assigned to the node.
	// Returns an empty slice if there are no tags.
	GetTags() ([]string, error)

	// GetIPAddresses returns the assigned IPv4 and IPv6 addresses.
	GetIPAddresses() ([]string, error)

	// IsEphemeral returns true if the node is ephemeral.
	IsEphemeral() (bool, error)

	// GetLastSeen returns the time the node was last seen.
	// Returns nil if the node was never seen.
	GetLastSeen() (*time.Time, error)

	// GetExpiry returns the key expiration of the node.
	// Returns nil if the key never expires (headscale zero value 0001-01-01 00:00:00).
	GetExpiry() (*time.Time, error)

	// IsOnline returns true if the node is currently connected.
	IsOnline() (bool, error)

	// IsExpired returns true if the node key is expired.
	IsExpired() (bool, error)
}
