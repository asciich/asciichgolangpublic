package osutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxutils"
)

// TurnSwapOff turns off swap on the local system if needed.
//
// It is a convenience wrapper that delegates to
// commandexecutorlinuxuserutils.TurnSwapOff using the local exec based
// command executor.
func TurnSwapOff(ctx context.Context) error {
	return linuxutils.TurnSwapOff(ctx, commandexecutorexecoo.Exec())
}
