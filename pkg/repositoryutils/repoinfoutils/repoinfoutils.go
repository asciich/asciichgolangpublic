package repoinfoutils

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Port used if `gitlab_port` is not set in the repo info file.
const DefaultGitlabPort = 2222

// Returns the value of `gitlab_instance`.
//
// repoInfoPath is either the path to the `.repoinfo.yml` file or the path to the directory containing it.
func GetGitlabInstance(ctx context.Context, repoInfoPath string) (string, error) {
	info, err := loadRepositoryInfo(ctx, repoInfoPath)
	if err != nil {
		return "", err
	}

	if info.GitlabInstance == "" {
		return "", tracederrors.TracedErrorf("'gitlab_instance' is not set in the repo info of '%s'.", repoInfoPath)
	}

	logging.LogInfoByCtxf(ctx, "Gitlab instance of repo info '%s' is '%s'.", repoInfoPath, info.GitlabInstance)

	return info.GitlabInstance, nil
}

// Returns the value of `gitlab_port`.
//
// An error is returned if `gitlab_port` is not set.
// Use GetGitlabPortOrDefaultIfUnset to fall back to the default port.
func GetGitlabPort(ctx context.Context, repoInfoPath string) (int, error) {
	info, err := loadRepositoryInfo(ctx, repoInfoPath)
	if err != nil {
		return 0, err
	}

	if info.GitlabPort == nil {
		return 0, tracederrors.TracedErrorf("'gitlab_port' is not set in the repo info of '%s'.", repoInfoPath)
	}

	err = checkValidPort(*info.GitlabPort)
	if err != nil {
		return 0, err
	}

	logging.LogInfoByCtxf(ctx, "Gitlab port of repo info '%s' is '%d'.", repoInfoPath, *info.GitlabPort)

	return *info.GitlabPort, nil
}

// Returns the value of `gitlab_port` if it is set.
// Otherwise the default port 2222 is returned.
func GetGitlabPortOrDefaultIfUnset(ctx context.Context, repoInfoPath string) (int, error) {
	info, err := loadRepositoryInfo(ctx, repoInfoPath)
	if err != nil {
		return 0, err
	}

	if info.GitlabPort == nil {
		logging.LogInfoByCtxf(
			ctx,
			"'gitlab_port' is not set in the repo info of '%s'. Use default port '%d'.",
			repoInfoPath,
			DefaultGitlabPort,
		)

		return DefaultGitlabPort, nil
	}

	err = checkValidPort(*info.GitlabPort)
	if err != nil {
		return 0, err
	}

	logging.LogInfoByCtxf(ctx, "Gitlab port of repo info '%s' is '%d'.", repoInfoPath, *info.GitlabPort)

	return *info.GitlabPort, nil
}

// Returns the value of `gitlab_path`.
//
// The returned path never includes the `.git` suffix, even if it is added by mistake in the repo info file.
func GetGitlabPath(ctx context.Context, repoInfoPath string) (string, error) {
	info, err := loadRepositoryInfo(ctx, repoInfoPath)
	if err != nil {
		return "", err
	}

	if info.GitlabPath == "" {
		return "", tracederrors.TracedErrorf("'gitlab_path' is not set in the repo info of '%s'.", repoInfoPath)
	}

	ret := strings.TrimSuffix(info.GitlabPath, ".git")
	if ret == "" {
		return "", tracederrors.TracedErrorf("'gitlab_path' is empty after removing the '.git' suffix in the repo info of '%s'.", repoInfoPath)
	}

	logging.LogInfoByCtxf(ctx, "Gitlab path of repo info '%s' is '%s'.", repoInfoPath, ret)

	return ret, nil
}

// Returns the SSH remote URL in the form of `ssh://git@<gitlab_instance>/<gitlab_path>.git`.
func GetSshRemoteUrl(ctx context.Context, repoInfoPath string) (string, error) {
	instance, err := GetGitlabInstance(ctx, repoInfoPath)
	if err != nil {
		return "", err
	}

	gitlabPath, err := GetGitlabPath(ctx, repoInfoPath)
	if err != nil {
		return "", err
	}

	ret := "ssh://git@" + instance + "/" + gitlabPath + ".git"

	logging.LogInfoByCtxf(ctx, "SSH remote URL of repo info '%s' is '%s'.", repoInfoPath, ret)

	return ret, nil
}

func checkValidPort(port int) error {
	if port < 1 || port > 65535 {
		return tracederrors.TracedErrorf("Invalid 'gitlab_port' '%d'. Must be between 1 and 65535.", port)
	}

	return nil
}
