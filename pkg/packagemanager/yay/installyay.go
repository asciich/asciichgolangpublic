package yay

import (
	"context"
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/commandexecutorgitoo"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxuserutils/commandexecutorlinuxuserutils"
	"github.com/asciich/asciichgolangpublic/pkg/osutils/linuxuserutils/linuxuseroptions"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/packagemanageroptions"
	"github.com/asciich/asciichgolangpublic/pkg/packagemanager/pacman"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/shellutils/shelllinehandler"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const YAY_GIT_REPO_URL = "https://aur.archlinux.org/yay.git"
const YAY_INSTALLATION_USER = "yay_installation_user"

func CreateYayInstallationUser(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, useSudo bool) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Create yay installation user '%s' on '%s' started.", YAY_INSTALLATION_USER, hostDescription)

	err = commandexecutorlinuxuserutils.Create(ctx, commandExecutor, &linuxuseroptions.CreateOptions{
		UserName:            YAY_INSTALLATION_USER,
		UseSudo:             useSudo,
		CreateHomeDirectory: true,
	})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Create yay installation user '%s' on '%s' finished.", YAY_INSTALLATION_USER, hostDescription)

	return nil
}

func DeleteYayInstallationUser(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, useSudo bool) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete yay installation user '%s' on '%s' started.", YAY_INSTALLATION_USER, hostDescription)

	err = commandexecutorlinuxuserutils.Delete(ctx, commandExecutor, &linuxuseroptions.DeleteOptions{
		UserName: YAY_INSTALLATION_USER,
		UseSudo:  useSudo,
	})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete yay installation user '%s' on '%s' finished.", YAY_INSTALLATION_USER, hostDescription)

	return nil
}

func InstallYay(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, options *packagemanageroptions.InstallPackageOptions) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if options == nil {
		options = new(packagemanageroptions.InstallPackageOptions)
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return nil
	}

	logging.LogInfoByCtxf(ctx, "Install yay on '%s' started.", hostDescription)

	isInstalled, err := IsInstalled(ctx, commandExecutor)
	if err != nil {
		return err
	}

	if isInstalled {
		logging.LogInfoByCtxf(ctx, "Yay is already installed on '%s'. Skip installation.", hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Yay is not yet installed on '%s'. Going to install.", hostDescription)

		neededPackages := []string{"base-devel", "git", "go"}
		logging.LogInfoByCtxf(ctx, "Install needed packages '%v' to install yay started.", neededPackages)
		err := pacman.InstallPackages(
			ctx,
			commandExecutor,
			neededPackages,
			options,
		)
		if err != nil {
			return err
		}
		logging.LogInfoByCtxf(ctx, "Install needed packages '%v' to install yay finished.", neededPackages)

		err = CreateYayInstallationUser(ctx, commandExecutor, options.UseSudo)
		if err != nil {
			return err
		}

		defer func() {
			err = DeleteYayInstallationUser(ctx, commandExecutor, options.UseSudo)
			if err != nil {
				logging.LogGoError(err)
			}
		}()

		logging.LogInfoByCtxf(ctx, "Clone yay repository from '%s' started.", YAY_GIT_REPO_URL)
		yayRepo, err := commandexecutorgitoo.CloneToTemporaryRepository(ctx, commandExecutor, YAY_GIT_REPO_URL)
		if err != nil {
			return err
		}
		defer yayRepo.Delete(ctx, &filesoptions.DeleteOptions{UseSudo: true})

		yayRepoPath, err := yayRepo.GetPath()
		if err != nil {
			return err
		}

		_, err = commandExecutor.RunCommand(ctx,
			&parameteroptions.RunCommandOptions{
				Command:   []string{"chown", "-R", YAY_INSTALLATION_USER, yayRepoPath},
				RunAsRoot: options.UseSudo,
			},
		)
		if err != nil {
			return err
		}

		logging.LogInfoByCtxf(ctx, "Clone yay repository from '%s' to temporary directory '%s' finished.", YAY_GIT_REPO_URL, yayRepoPath)

		cmd := []string{"bash", "-c", fmt.Sprintf("cd '%s' && makepkg -sf", yayRepoPath)}
		cmdJoined, err := shelllinehandler.Join(cmd)
		if err != nil {
			return err
		}
		logging.LogInfoByCtxf(ctx, "Build yay package using '%s' started.", cmdJoined)
		_, err = commandExecutor.RunCommand(
			commandexecutorgeneric.WithLiveOutputOnStdout(contextutils.WithSilent(ctx)),
			&parameteroptions.RunCommandOptions{
				Command:            cmd,
				RunAsUser:          YAY_INSTALLATION_USER,
				UseSudoToRunAsUser: options.UseSudo,
			})
		if err != nil {
			return err
		}
		logging.LogInfoByCtxf(ctx, "Build yay package using '%s' finished.", cmdJoined)

		cmd = []string{"bash", "-c", fmt.Sprintf("cd '%s' && pacman -U yay-*.pkg.tar.*z* --noconfirm 2>&1", yayRepoPath)}
		cmdJoined, err = shelllinehandler.Join(cmd)
		if err != nil {
			return err
		}
		logging.LogInfoByCtxf(ctx, "Install yay using '%s' started.", cmdJoined)
		_, err = commandExecutor.RunCommand(
			commandexecutorgeneric.WithLiveOutputOnStdout(contextutils.WithSilent(ctx)),
			&parameteroptions.RunCommandOptions{
				Command:   cmd,
				RunAsRoot: options.UseSudo,
			})
		if err != nil {
			return err
		}
		logging.LogInfoByCtxf(ctx, "Install yay using '%s' finished.", cmd)

		logging.LogChangedByCtxf(ctx, "Installed yay on '%s'.", hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Install yay on '%s' finished.", hostDescription)

	return nil
}

func (y *Yay) InstallYay(ctx context.Context, options *packagemanageroptions.InstallPackageOptions) error {
	commaandExecutor, err := y.GetCommandExecutor()
	if err != nil {
		return err
	}

	return InstallYay(ctx, commaandExecutor, options)
}
