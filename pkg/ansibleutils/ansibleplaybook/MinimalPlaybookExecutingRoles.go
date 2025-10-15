package ansibleplaybook

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func NewMinimalPlaybookExecutingRoles(ctx context.Context, options *MinimalPlaybookOptions) (*Playbook, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	hostname, err := options.GetHostName()
	if err != nil {
		return nil, err
	}

	if len(options.Roles) <= 0 {
		return nil, tracederrors.TracedError("roles has no elements")
	}

	roles := slicesutils.GetDeepCopyOfStringsSlice(options.Roles)

	playbook := &Playbook{
		Plays: []*Play{
			{
				Name:       "Minimal playbook executing roles",
				Hosts:      []string{hostname},
				Roles:      roles,
				RemoteUser: options.RemoteUser,
			},
		},
	}

	return playbook, nil
}

func WriteTemporaryMinimalPlaybookExecutingRoles(ctx context.Context, options *MinimalPlaybookOptions) (string, error) {
	playbook, err := NewMinimalPlaybookExecutingRoles(ctx, options)
	if err != nil {
		return "", err
	}

	path, err := tempfiles.CreateTemporaryFile(ctx)
	if err != nil {
		return "", err
	}

	return path, WritePlaybook(ctx, playbook, path)
}
