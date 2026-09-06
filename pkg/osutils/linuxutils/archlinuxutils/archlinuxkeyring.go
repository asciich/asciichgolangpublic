package archlinuxutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/pacman"
)

// Updates the pacman database and then the 'archlinux-keyring' package.
//
// The keyring-package contains the signing keys.
// If an archlinux was not updated for a long time it makes sense to update 'archlinux-keyring' first before doing a system update to ensure all new signing keys are present.
func UpdateArchLinuxKeyringPackage(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, useSudo bool) error {
	return pacman.UpdateArchLinuxKeyringPackage(ctx, commandExecutor, useSudo)
}
