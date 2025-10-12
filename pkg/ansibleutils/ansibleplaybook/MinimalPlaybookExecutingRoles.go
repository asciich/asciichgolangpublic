package ansibleplaybook

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func NewMinimalPlaybookExecutingRoles(ctx context.Context, hostname string, roles []string) (*Playbook, error) {
	if hostname == "" {
		return nil, tracederrors.TracedErrorEmptyString("hostname")
	}

	if len(roles) <= 0 {
		return nil, tracederrors.TracedError("roles has no elements")
	}

	roles = slicesutils.GetDeepCopyOfStringsSlice(roles)

	playbook := &Playbook{
		Plays: []*Play{
			{
				Name:  "Minimal playbook executing roles",
				Hosts: []string{hostname},
				Roles: roles,
			},
		},
	}

	return playbook, nil
}

func WriteTemporaryMinimalPlaybookExecutingRoles(ctx context.Context, hostname string, roles []string) (string, error) {
	playbook, err := NewMinimalPlaybookExecutingRoles(ctx, hostname, roles)
	if err != nil {
		return "", err
	}

	path, err := tempfiles.CreateTemporaryFile(ctx)
	if err != nil {
		return "", err
	}

	return path, WritePlaybook(ctx, playbook, path)
}
