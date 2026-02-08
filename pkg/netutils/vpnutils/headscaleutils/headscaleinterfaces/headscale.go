package headscaleinterfaces

import "context"

type HeadScale interface {
	CreateUser(ctx context.Context, userName string) error
	GeneratePreauthKeyForUser(ctx context.Context, userName string) (string, error)
	GetUserId(ctx context.Context, userName string) (int, error)
	ListUserNames(ctx context.Context) ([]string, error)
}
