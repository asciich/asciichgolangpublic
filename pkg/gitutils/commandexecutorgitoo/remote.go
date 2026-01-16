package commandexecutorgitoo

import (
	"context"
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (g *GitRepository) AddRemote(ctx context.Context, remoteOptions *gitparameteroptions.GitRemoteAddOptions) (err error) {
	if remoteOptions == nil {
		return tracederrors.TracedError("remoteOptions is nil")
	}

	remoteName, err := remoteOptions.GetRemoteName()
	if err != nil {
		return err
	}

	remoteUrl, err := remoteOptions.GetRemoteUrl()
	if err != nil {
		return err
	}

	repoPath, err := g.GetPath()
	if err != nil {
		return err
	}

	remoteExists, err := g.RemoteConfigurationExists(
		ctx,
		&gitgeneric.GenericGitRemoteConfig{
			RemoteName: remoteName,
			UrlFetch:   remoteUrl,
			UrlPush:    remoteUrl,
		},
	)
	if err != nil {
		return err
	}

	if remoteExists {
		logging.LogInfoByCtxf(ctx, "Remote '%s' as '%s' to repository '%s' already exists.", remoteUrl, remoteName, repoPath)
	} else {
		err = g.RemoveRemoteByName(ctx, remoteName)
		if err != nil {
			return err
		}

		_, err = g.RunGitCommand(ctx, []string{"remote", "add", remoteName, remoteUrl})
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Added remote '%s' as '%s' to repository '%s'.", remoteUrl, remoteName, repoPath)
	}

	return nil
}

func (g *GitRepository) RemoteConfigurationExists(ctx context.Context, config gitinterfaces.GitRemoteConfig) (exists bool, err error) {
	if config == nil {
		return false, tracederrors.TracedError("config is nil")
	}

	remoteConfigs, err := g.GetRemoteConfigs(ctx)
	if err != nil {
		return false, err
	}

	for _, toCheck := range remoteConfigs {
		if config.Equals(toCheck) {
			return true, nil
		}
	}

	return false, nil
}

func (g *GitRepository) RemoveRemoteByName(ctx context.Context, remoteNameToRemove string) (err error) {
	if len(remoteNameToRemove) <= 0 {
		return tracederrors.TracedError("remoteNameToRemove is empty string")
	}

	remoteExists, err := g.RemoteByNameExists(ctx, remoteNameToRemove)
	if err != nil {
		return err
	}

	repoDirPath, err := g.GetPath()
	if err != nil {
		return err
	}

	if remoteExists {
		_, err := g.RunGitCommand(ctx, []string{"remote", "remove", remoteNameToRemove})
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Remote '%s' for repository '%s' removed.", remoteNameToRemove, repoDirPath)
	} else {
		logging.LogInfoByCtxf(ctx, "Remote '%s' for repository '%s' already deleted.", remoteNameToRemove, repoDirPath)
	}

	return nil
}

func (g *GitRepository) RemoteByNameExists(ctx context.Context, remoteName string) (remoteExists bool, err error) {
	if len(remoteName) <= 0 {
		return false, fmt.Errorf("remoteName is empty string")
	}

	remoteConfigs, err := g.GetRemoteConfigs(ctx)
	if err != nil {
		return false, err
	}

	for _, toCheck := range remoteConfigs {
		toCheckRemoteName, _ := toCheck.GetRemoteName()
		if toCheckRemoteName == remoteName {
			return true, nil
		}
	}

	return false, nil
}

func (g *GitRepository) GetRemoteConfigs(ctx context.Context) (remoteConfigs []gitinterfaces.GitRemoteConfig, err error) {
	output, err := g.RunGitCommand(ctx, []string{"remote", "-v"})
	if err != nil {
		return nil, err
	}

	outputLines, err := output.GetStdoutAsLines(false)
	if err != nil {
		return nil, err
	}

	remoteConfigs = []gitinterfaces.GitRemoteConfig{}
	for _, line := range outputLines {
		line = strings.TrimSpace(line)
		if len(line) <= 0 {
			continue
		}

		lineCleaned := strings.ReplaceAll(line, "\t", " ")

		splitted := stringsutils.SplitAtSpacesAndRemoveEmptyStrings(lineCleaned)
		if len(splitted) != 3 {
			return nil, tracederrors.TracedErrorf("Unable to parse '%s' as remote. splitted is '%v'", line, splitted)
		}

		remoteName := splitted[0]
		remoteUrl := splitted[1]
		remoteDirection := splitted[2]

		var remoteToModify gitinterfaces.GitRemoteConfig = nil
		for _, existingRemote := range remoteConfigs {
			existingRemoteName, _ := existingRemote.GetRemoteName()
			if existingRemoteName == remoteName {
				remoteToModify = existingRemote
			}
		}

		if remoteToModify == nil {
			remoteToAdd := &gitgeneric.GenericGitRemoteConfig{}
			remoteToAdd.RemoteName = remoteName
			remoteConfigs = append(remoteConfigs, remoteToAdd)
			remoteToModify = remoteToAdd
		}

		if remoteDirection == "(fetch)" {
			err = remoteToModify.SetUrlFetch(remoteUrl)
			if err != nil {
				return nil, err
			}
		} else if remoteDirection == "(push)" {
			err = remoteToModify.SetUrlPush(remoteUrl)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, tracederrors.TracedErrorf("Unknown remoteDirection='%s'", remoteDirection)
		}
	}

	return remoteConfigs, nil
}

func (g *GitRepository) Fetch(ctx context.Context) (err error) {
	_, err = g.RunGitCommand(ctx, []string{"fetch"})
	if err != nil {
		return err
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Fetched git repository '%s' on host '%s'", path, hostDescription)

	return nil
}

func (g *GitRepository) PullFromRemote(ctx context.Context, pullOptions *gitparameteroptions.GitPullFromRemoteOptions) (err error) {
	if pullOptions == nil {
		return tracederrors.TracedError("pullOptions not set")
	}

	remoteName, err := pullOptions.GetRemoteName()
	if err != nil {
		return err
	}

	branchName, err := pullOptions.GetBranchName()
	if err != nil {
		return err
	}

	if len(remoteName) <= 0 {
		return tracederrors.TracedError("remoteName is empty string")
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	_, err = g.RunGitCommand(ctx, []string{"pull", remoteName, branchName})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Pulled git repository '%s' on host '%s' from remote '%s'.", path, hostDescription, remoteName)

	return nil
}

func (g *GitRepository) Push(ctx context.Context) (err error) {
	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Push git repository '%s' on '%s' started.", path, hostDescription)

	_, err = g.RunGitCommand(ctx, []string{"push"})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Push git repository '%s' on '%s' finished.", path, hostDescription)

	return
}

func (g *GitRepository) PushTagsToRemote(ctx context.Context, remoteName string) (err error) {
	if len(remoteName) <= 0 {
		return tracederrors.TracedError("remoteName is empty string")
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	_, err = g.RunGitCommand(ctx, []string{"push", remoteName, "--tags"})
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Pushed tags of git repository '%s' on host '%s' to remote '%s'.", path, hostDescription, remoteName)

	return nil
}

func (g *GitRepository) PushToRemote(ctx context.Context, remoteName string) (err error) {
	if len(remoteName) <= 0 {
		return tracederrors.TracedError("remoteName is empty string")
	}

	_, err = g.RunGitCommand(ctx, []string{"push", remoteName})
	if err != nil {
		return err
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Pushed git repository '%s' on host '%s' to remote '%s'.", path, hostDescription, remoteName)

	return nil
}
