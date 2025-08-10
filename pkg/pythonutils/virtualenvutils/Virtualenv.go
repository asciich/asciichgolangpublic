package virtualenvutils

import (
	"context"
	"path/filepath"
	"slices"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const PYENV_CONFIG_BASENAME = "pyvenv.cfg"

type VirtualEnv struct {
	Path            string
	CommandExecutor *commandexecutorinterfaces.CommandExecutor
}

func GetVirtualEnv(path string) (*VirtualEnv, error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	return &VirtualEnv{
		Path: path,
	}, nil
}

func IsVirtualEnv(ctx context.Context, path string) (bool, error) {
	if path == "" {
		return false, tracederrors.TracedErrorEmptyString("path")
	}

	if !nativefiles.IsDir(contextutils.WithSilent(ctx), path) {
		logging.LogInfoByCtxf(ctx, "'%s' is not a directory and can therefore not be a virtualenv.", path)
		return false, nil
	}

	configFilePath := filepath.Join(path, PYENV_CONFIG_BASENAME)
	if !nativefiles.IsFile(contextutils.WithSilent(ctx), configFilePath) {
		logging.LogInfoByCtxf(ctx, "'%s' is not a virtualenv since the virtualenv config '%s' is missing.", path, configFilePath)
		return false, nil
	}

	logging.LogInfoByCtxf(ctx, "'%s' is a virtualenv", path)
	return true, nil
}

func (v *VirtualEnv) GetCommandExecutorOrDefaultIfUnset() commandexecutorinterfaces.CommandExecutor {
	if v.CommandExecutor != nil {
		return *v.CommandExecutor
	}

	return commandexecutorexecoo.Exec()
}

func (v *VirtualEnv) Create(ctx context.Context) error {
	vePath, err := v.GetPath()
	if err != nil {
		return err
	}

	isVirtuelEnv, err := IsVirtualEnv(ctx, vePath)
	if err != nil {
		return err
	}

	if isVirtuelEnv {
		logging.LogInfoByCtxf(ctx, "Virtualenv '%s' already exists.", vePath)
	} else {
		err := nativefiles.CreateDirectory(ctx, vePath)
		if err != nil {
			return err
		}

		ce := v.GetCommandExecutorOrDefaultIfUnset()
		cmd := []string{"virtualenv", vePath}
		_, err = ce.RunCommand(
			commandexecutorgeneric.WithLiveOutputOnStdout(ctx),
			&parameteroptions.RunCommandOptions{
				Command: cmd,
			},
		)
		if err != nil {
			return err
		}

		logging.LogInfoByCtxf(ctx, "Created virtualenv '%s'.", vePath)
	}

	return nil
}

// Creates a virtualenv and ensures the packages as specified in the options are present.
func CreateVirtualEnv(ctx context.Context, options *CreateVirtualenvOptions) (*VirtualEnv, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	vePath, err := options.GetPath()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Install virtualenv '%s' started.", options.Path)

	ve, err := GetVirtualEnv(vePath)
	if err != nil {
		return nil, err
	}

	err = ve.Create(ctx)
	if err != nil {
		return nil, err
	}

	if len(options.Packages) <= 0 {
		logging.LogInfoByCtxf(ctx, "No packages to install in virtualenv '%s'.", vePath)
	} else {
		err = ve.InstallPackages(ctx, options.Packages)
		if err != nil {
			return nil, err
		}
	}

	logging.LogInfoByCtxf(ctx, "Install virtualenv '%s' finished.", vePath)

	return ve, nil
}

func (v *VirtualEnv) GetPath() (string, error) {
	if v.Path == "" {
		return "", tracederrors.TracedError("Path is empty string")
	}

	return v.Path, nil
}

func (v *VirtualEnv) IsVirtualEnv(ctx context.Context) (bool, error) {
	path, err := v.GetPath()
	if err != nil {
		return false, err
	}

	return IsVirtualEnv(ctx, path)
}

func (v *VirtualEnv) GetPipPath() (string, error) {
	path, err := v.GetPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(path, "bin", "pip"), nil
}

func (v *VirtualEnv) ListInstalledPackageNames(ctx context.Context) ([]string, error) {
	path, err := v.GetPath()
	if err != nil {
		return nil, err
	}

	isVirtualenv, err := v.IsVirtualEnv(contextutils.WithSilent(ctx))
	if err != nil {
		return nil, err
	}

	if !isVirtualenv {
		logging.LogInfoByCtxf(ctx, "Unable to list installed python packages in virtualenv '%s'. '%s' is not a virtualenv.", path, path)
		return nil, err
	}

	pipPath, err := v.GetPipPath()
	if err != nil {
		return nil, err
	}

	freezeOutput, err := commandexecutorexec.RunCommand(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{pipPath, "freeze"},
	})
	if err != nil {
		return nil, err
	}

	lines, err := freezeOutput.GetStdoutAsLines(true)
	if err != nil {
		return nil, err
	}

	installedPackages := []string{}
	for _, line := range lines {
		toAdd := strings.ToLower(strings.Split(line, "=")[0])
		toAdd = strings.TrimSpace(toAdd)

		if toAdd == "" {
			continue
		}

		installedPackages = append(installedPackages, toAdd)
	}

	return installedPackages, nil
}

func (v *VirtualEnv) IsPackageInstalled(ctx context.Context, packageName string) (isInstalled bool, err error) {
	if packageName == "" {
		return false, tracederrors.TracedErrorEmptyString("packageName")
	}

	return v.IsPackagesInstalled(ctx, []string{packageName})
}

func (v *VirtualEnv) IsPackagesInstalled(ctx context.Context, packageNames []string) (isInstalled bool, err error) {
	vePath, err := v.GetPath()
	if err != nil {
		return false, err
	}

	if len(packageNames) <= 0 {
		return false, tracederrors.TracedError("No elements in packageNames.")
	}

	installedPackages, err := v.ListInstalledPackageNames(ctx)
	if err != nil {
		return false, err
	}

	for _, p := range packageNames {
		isInstalled = slices.Contains(installedPackages, p)
		if isInstalled {
			logging.LogInfoByCtxf(ctx, "Package '%s' is installed in virtualenv '%s'.", p, vePath)
		} else {
			logging.LogInfoByCtxf(ctx, "Package '%s' is not installed in virtualenv '%s'.", p, vePath)
			break
		}
	}

	return isInstalled, nil
}

func (v *VirtualEnv) InstallPackages(ctx context.Context, packages []string) error {
	if packages == nil {
		return tracederrors.TracedErrorNil("packages")
	}

	vePath, err := v.GetPath()
	if err != nil {
		return err
	}

	if len(packages) <= 0 {
		return tracederrors.TracedErrorf("No packages to install in virtualenv '%s'.", vePath)
	}

	installed, err := v.IsPackagesInstalled(ctx, packages)
	if err != nil {
		return err
	}

	if installed {
		logging.LogInfoByCtxf(ctx, "All packages '%v' already installed in virtualenv '%s'.", packages, vePath)
	} else {
		for _, p := range packages {
			err = v.InstallPackage(ctx, p)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *VirtualEnv) InstallPackage(ctx context.Context, packageName string) error {
	if packageName == "" {
		return tracederrors.TracedErrorEmptyString("packageName")
	}

	vePath, err := v.GetPath()
	if err != nil {
		return err
	}

	isInstalled, err := v.IsPackageInstalled(ctx, packageName)
	if err != nil {
		return err
	}

	if isInstalled {
		logging.LogInfoByCtxf(ctx, "Package '%s' already installed in virtualenv '%s'.", packageName, vePath)
	} else {
		pipPath, err := v.GetPipPath()
		if err != nil {
			return err
		}

		ce := v.GetCommandExecutorOrDefaultIfUnset()
		_, err = ce.RunCommand(
			commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
			&parameteroptions.RunCommandOptions{
				Command: []string{pipPath, "install", packageName},
			},
		)
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Installed package '%s' in virtualenv '%s'.", packageName, vePath)
	}

	return nil
}
