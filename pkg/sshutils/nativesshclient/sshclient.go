package nativesshclient

import (
	"context"
	"io"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (n *SshClient) GetCPUArchitecture(ctx context.Context) (string, error) {
	return "", tracederrors.TracedErrorNotImplemented()
}

func (n *SshClient) GetDeepCopy() commandexecutorinterfaces.CommandExecutor {
	panic("Not implemented")
}

func (n *SshClient) GetHostDescription() (string, error) {
	return "localhost", nil
}

func (n *SshClient) IsRunningOnLocalhost() (bool, error) {
	return true, nil
}

func (n *SshClient) RunCommandAndGetStdinAsIoWriteCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.WriteCloser, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}

func (n *SshClient) RunCommandAndGetStdoutAsIoReadCloser(ctx context.Context, options *parameteroptions.RunCommandOptions) (io.ReadCloser, error) {
	return nil, tracederrors.TracedErrorNotImplemented()
}
