package headscaleinterfaces

import "context"

// Headscale describes the operations available on a headscale server.
//
// Implementations are provided by the native... and commandexecutor... packages.
type Headscale interface {
	// CreateUser ensures the user with the given name exists.
	CreateUser(ctx context.Context, userName string) error

	// GeneratePreauthKeyForUser generates a new preauth key for the given user.
	GeneratePreauthKeyForUser(ctx context.Context, userName string) (string, error)

	// GetUserId returns the ID of the given user.
	GetUserId(ctx context.Context, userName string) (int, error)

	// DeleteNodeByName ensures the node with the given name is deleted.
	//
	// The name is the node name as returned by ListNodeNames.
	// The function is idempotent: If no node with this name exists, nothing is changed and no error is returned.
	// Returns an error if the name is empty or if more than one node has this name,
	// since it is then unclear which node should be deleted. In this case use DeleteNodeById.
	DeleteNodeByName(ctx context.Context, nodeName string) error

	// DeleteNodeById ensures the node with the given ID is deleted.
	//
	// The ID is the unique node ID as returned by Node.GetID.
	// The function is idempotent: If no node with this ID exists, nothing is changed and no error is returned.
	// Returns an error if the ID is 0, since headscale node IDs start at 1.
	DeleteNodeById(ctx context.Context, id uint64) error

	// ListNodes returns the nodes with all details.
	ListNodes(ctx context.Context) ([]Node, error)

	// ListNodeNames returns only the names of the nodes.
	// The names must be sorted alphabetically.
	ListNodeNames(ctx context.Context) ([]string, error)

	// ListUsers returns the users with all details.
	//
	// Returns an empty slice if there are no users.
	ListUsers(ctx context.Context) ([]User, error)

	// ListUserNames returns only the names of the users.
	// The names are the usernames as returned by User.GetUsername.
	// The names must be sorted alphabetically.
	ListUserNames(ctx context.Context) ([]string, error)

	// ListApiKeys returns the API keys with all details.
	//
	// Returns an empty slice if there are no API keys.
	// Only the prefix of an API key is available, never the full key.
	ListApiKeys(ctx context.Context) ([]ApiKey, error)
}
