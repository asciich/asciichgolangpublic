//go:build !linux

package nativehost

import (
	"context"
	"runtime"

	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsoptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (n *NativeHost) SetHostName(ctx context.Context, hostname string, options *hostsutilsoptions.SetHostnameOptions) error {
	if hostname == "" {
		return tracederrors.TracedErrorEmptyString("hostname")
	}
	return tracederrors.TracedErrorf("SetHostName is not supported on '%s'.", runtime.GOOS)
}
