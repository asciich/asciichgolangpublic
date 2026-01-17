package archlinuxutils

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/pacman"
	"github.com/asciich/asciichgolangpublic/pkg/runbook"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func NewUpdateArchLinuxKeyringPackageRunbook(commandExecutor commandexecutorinterfaces.CommandExecutor, useSudo bool) *runbook.RunBook {
	return &runbook.RunBook{
		Name: "Update the archlinux-keyring package.",
		Description: `Update the 'archlinux-keyring' package containing all singing keys.
Running this before an system update makes sense to ensure all signing keys are present and up to date.`,
		Steps: []runbook.Runnable{
			&runbook.Step{
				Name:        "Preflight check",
				Description: "Check if this runbook can be run.",
				Run: func(ctx context.Context) error {
					if commandExecutor == nil {
						return tracederrors.TracedErrorEmptyString("commandExecutor is not set.")
					}

					return nil
				},
			},
			&runbook.Step{
				Name:        "Update pacman database",
				Description: "Update the pacman database before checking for archlinux-keyring updates.",
				Run: func(ctx context.Context) error {
					pacman, err := pacman.NewPacman(commandExecutor)
					if err != nil {
						return err
					}

					return pacman.UpdateDatabase(
						ctx,
						&packagemanageroptions.UpdateDatabaseOptions{
							UseSudo: useSudo,
						},
					)
				},
			},
			&runbook.Step{
				Name:        "Install or update archlinux-keyring package.",
				Description: "Install or update archlinx-keyring package.",
				Run: func(ctx context.Context) error {
					pacman, err := pacman.NewPacman(commandExecutor)
					if err != nil {
						return err
					}

					return pacman.InstallPackages(
						ctx,
						[]string{"archlinux-keyring"},
						&packagemanageroptions.InstallPackageOptions{
							UpdatePackage:       true,
							UpdateDatabaseFirst: false, // Update was already done in previous step.
							Force:               true,
							UseSudo:             useSudo,
						},
					)
				},
			},
		},
	}
}

// Updates the pacman database and then the 'archlinux-keyring' package.
//
// The keyring-package contains the signing keys.
// If an archlinux was not updated for a long time it makes sense to update 'archlinux-keyring' first before doing a system update to ensure all new signing keys are present.
func UpdateArchLinuxKeyringPackage(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, useSudo bool) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	runBook := NewUpdateArchLinuxKeyringPackageRunbook(commandExecutor, useSudo)

	return runBook.Execute(ctx)
}
