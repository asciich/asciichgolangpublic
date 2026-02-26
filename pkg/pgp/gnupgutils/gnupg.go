package gnupgutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/pgp/gnupgutils/commandexecutorgnupg"
	"github.com/asciich/asciichgolangpublic/pkg/pgp/gnupgutils/gnupgoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func SignFileByPath(ctx context.Context, toSignPath string, options *gnupgoptions.SignOption) error {
	if toSignPath == "" {
		return tracederrors.TracedErrorEmptyString("toSignPath")
	}

	commandExecutor := commandexecutorexecoo.Exec()
	return commandexecutorgnupg.SignFileByPath(ctx, commandExecutor, toSignPath, options)
}

func SignFile(ctx context.Context, toSign filesinterfaces.File, options *gnupgoptions.SignOption) error {
	if toSign == nil {
		return tracederrors.TracedErrorNil("toSign")
	}

	commandExecutor := commandexecutorexecoo.Exec()
	return commandexecutorgnupg.SignFile(ctx, commandExecutor, toSign, options)
}

func CheckSignatureValid(ctx context.Context, signatureFile filesinterfaces.File) error {
	if signatureFile == nil {
		return tracederrors.TracedErrorNil("signatureFile")
	}

	commandExecutor := commandexecutorexecoo.Exec()
	return commandexecutorgnupg.CheckSignatureValid(ctx, commandExecutor, signatureFile)
}

func CheckSingnatureByPathValid(ctx context.Context, signaturePath string) error {
	return commandexecutorgnupg.CheckSingnatureByPathValid(ctx, commandexecutorexecoo.Exec(), signaturePath)
}
