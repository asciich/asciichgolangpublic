package httpclientcmdoptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httputilsinterfaces"
)

type HttpClientCmdOptions struct {
	GetClient func() httputilsinterfaces.Client
}
