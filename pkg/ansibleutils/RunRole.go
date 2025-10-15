package ansibleutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansibleparemeteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/ansibleutils/ansibleplaybook"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func RunRoles(ctx context.Context, roles []string, options *ansibleparemeteroptions.RunOptions) error {
	if len(roles) <= 0 {
		return tracederrors.TracedError("roles has no elements")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	hostname, err := options.GetLimit()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Running ansible roles '%v' against host '%s' started.", roles, hostname)

	tempPlaybook, err := ansibleplaybook.WriteTemporaryMinimalPlaybookExecutingRoles(
		ctx,
		&ansibleplaybook.MinimalPlaybookOptions{
			Hostname:   hostname,
			Roles:      roles,
			RemoteUser: options.RemoteUser,
		},
	)
	if err != nil {
		return err
	}
	if !options.KeepTemporaryPlaybook {
		defer func() { _ = nativefiles.Delete(ctx, tempPlaybook, &filesoptions.DeleteOptions{}) }()
	}

	optionsToUse := options.DeepCopy()
	optionsToUse.PlaybookPath = tempPlaybook
	err = RunPlaybook(ctx, optionsToUse)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Running ansible roles '%v' against host '%s' finished.", roles, hostname)

	return nil
}
