package asciichgolangpublic

import (
	aslices "github.com/asciich/asciichgolangpublic/datatypes/slices"
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type GitlabSyncBranchOptions struct {
	// Define the target branch where the files are synced to.
	// Can be specified by either setting the target branch as object or by string.
	TargetBranch     *GitlabBranch
	TargetBranchName string

	PathsToSync []string

	Verbose bool
}

func NewGitlabSyncBranchOptions() (g *GitlabSyncBranchOptions) {
	return new(GitlabSyncBranchOptions)
}

func (g *GitlabSyncBranchOptions) GetDeepCopy() (copy *GitlabSyncBranchOptions) {
	copy = NewGitlabSyncBranchOptions()

	*copy = *g

	if g.TargetBranch != nil {
		copy.TargetBranch = g.TargetBranch.GetDeepCopy()
	}

	if g.PathsToSync != nil {
		copy.PathsToSync = aslices.GetDeepCopyOfStringsSlice(g.PathsToSync)
	}

	return copy
}

func (g *GitlabSyncBranchOptions) GetPathsToSync() (pathsToSync []string, err error) {
	if g.PathsToSync == nil {
		return nil, errors.TracedErrorf("PathsToSync not set")
	}

	if len(g.PathsToSync) <= 0 {
		return nil, errors.TracedErrorf("PathsToSync has no elements")
	}

	return g.PathsToSync, nil
}

func (g *GitlabSyncBranchOptions) GetTargetBranch() (targetBranch *GitlabBranch, err error) {
	if g.TargetBranch == nil {
		return nil, errors.TracedErrorf("TargetBranch not set")
	}

	return g.TargetBranch, nil
}

func (g *GitlabSyncBranchOptions) GetTargetBranchName() (targetBranchName string, err error) {
	if g.TargetBranchName != "" {
		targetBranchName = g.TargetBranchName
	}

	if targetBranchName == "" {
		if g.IsTargetBranchSet() {
			targetBranch, err := g.GetTargetBranch()
			if err != nil {
				return "", err
			}

			targetBranchName, err = targetBranch.GetName()
			if err != nil {
				return "", err
			}
		}
	}

	if targetBranchName == "" {
		return targetBranchName, errors.TracedErrorf("TargetBranchName not set")
	}

	return targetBranchName, nil
}

func (g *GitlabSyncBranchOptions) GetVerbose() (verbose bool) {
	return g.Verbose
}

func (g *GitlabSyncBranchOptions) IsTargetBranchSet() (isSet bool) {
	return g.TargetBranch != nil
}

func (g *GitlabSyncBranchOptions) MustGetPathsToSync() (pathsToSync []string) {
	pathsToSync, err := g.GetPathsToSync()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return pathsToSync
}

func (g *GitlabSyncBranchOptions) MustGetTargetBranch() (targetBranch *GitlabBranch) {
	targetBranch, err := g.GetTargetBranch()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return targetBranch
}

func (g *GitlabSyncBranchOptions) MustGetTargetBranchName() (targetBranchName string) {
	targetBranchName, err := g.GetTargetBranchName()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return targetBranchName
}

func (g *GitlabSyncBranchOptions) MustSetPathsToSync(pathsToSync []string) {
	err := g.SetPathsToSync(pathsToSync)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabSyncBranchOptions) MustSetTargetBranch(targetBranch *GitlabBranch) {
	err := g.SetTargetBranch(targetBranch)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabSyncBranchOptions) MustSetTargetBranchName(targetBranchName string) {
	err := g.SetTargetBranchName(targetBranchName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabSyncBranchOptions) MustSetTargetBranchNameAndUnsetTargetBranchObject(targetBranchName string) {
	err := g.SetTargetBranchNameAndUnsetTargetBranchObject(targetBranchName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabSyncBranchOptions) MustUnsetTargetBranch() {
	err := g.UnsetTargetBranch()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func (g *GitlabSyncBranchOptions) SetPathsToSync(pathsToSync []string) (err error) {
	if pathsToSync == nil {
		return errors.TracedErrorf("pathsToSync is nil")
	}

	if len(pathsToSync) <= 0 {
		return errors.TracedErrorf("pathsToSync has no elements")
	}

	g.PathsToSync = pathsToSync

	return nil
}

func (g *GitlabSyncBranchOptions) SetTargetBranch(targetBranch *GitlabBranch) (err error) {
	if targetBranch == nil {
		return errors.TracedErrorf("targetBranch is nil")
	}

	g.TargetBranch = targetBranch

	return nil
}

func (g *GitlabSyncBranchOptions) SetTargetBranchName(targetBranchName string) (err error) {
	if targetBranchName == "" {
		return errors.TracedErrorf("targetBranchName is empty string")
	}

	g.TargetBranchName = targetBranchName

	return nil
}

func (g *GitlabSyncBranchOptions) SetTargetBranchNameAndUnsetTargetBranchObject(targetBranchName string) (err error) {
	err = g.SetTargetBranchName(targetBranchName)
	if err != nil {
		return err
	}

	err = g.UnsetTargetBranch()
	if err != nil {
		return err
	}

	return err
}

func (g *GitlabSyncBranchOptions) SetVerbose(verbose bool) {
	g.Verbose = verbose
}

func (g *GitlabSyncBranchOptions) UnsetTargetBranch() (err error) {
	g.TargetBranch = nil
	return nil
}
