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

func CheckSignatureValid(ctx context.Context, signatureFile filesinterfaces.File, toValidateFile filesinterfaces.File) error {
	if signatureFile == nil {
		return tracederrors.TracedErrorNil("signatureFile")
	}

	if toValidateFile == nil {
		return tracederrors.TracedErrorNil("toValidateFile")
	}

	commandExecutor := commandexecutorexecoo.Exec()
	return commandexecutorgnupg.CheckSignatureValid(ctx, commandExecutor, signatureFile, toValidateFile)
}

func CheckSingnatureByPathValid(ctx context.Context, signaturePath string, toValidatePath string) error {
	return commandexecutorgnupg.CheckSingnatureByPathValid(ctx, commandexecutorexecoo.Exec(), signaturePath, toValidatePath)
}

func SignContentString(ctx context.Context, content string, options *gnupgoptions.SignOption) (string, error) {
	if content == "" {
		return "", tracederrors.TracedErrorEmptyString("content")
	}

	commandExecutor := commandexecutorexecoo.Exec()
	return commandexecutorgnupg.SignContentString(ctx, commandExecutor, content, options)
}

func SignContentBytes(ctx context.Context, content []byte, options *gnupgoptions.SignOption) ([]byte, error) {
	if content == nil {
		return nil, tracederrors.TracedErrorNil("content")
	}

	commandExecutor := commandexecutorexecoo.Exec()
	return commandexecutorgnupg.SignContentBytes(ctx, commandExecutor, content, options)
}

func CheckSignatureValidForContentString(ctx context.Context, signature string, content string) error {
	if signature == "" {
		return tracederrors.TracedErrorEmptyString("signature")
	}

	if content == "" {
		return tracederrors.TracedErrorEmptyString("content")
	}

	commandExecutor := commandexecutorexecoo.Exec()
	return commandexecutorgnupg.CheckSignatureValidForContentString(ctx, commandExecutor, signature, content)
}

func CheckSignatureValidForContentBytes(ctx context.Context, signature []byte, content []byte) error {
	if signature == nil {
		return tracederrors.TracedErrorNil("signature")
	}

	if content == nil {
		return tracederrors.TracedErrorNil("content")
	}

	commandExecutor := commandexecutorexecoo.Exec()
	return commandexecutorgnupg.CheckSignatureValidForContentBytes(ctx, commandExecutor, signature, content)
}
