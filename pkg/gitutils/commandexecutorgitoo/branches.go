package commandexecutorgitoo

import (
	"context"
	"sort"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (g *GitRepository) CheckoutBranchByName(ctx context.Context, name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	currentBranchName, err := g.GetCurrentBranchName(ctx)
	if err != nil {
		return err
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if currentBranchName == name {
		logging.LogInfoByCtxf(ctx, "Git repository '%s' on host '%s' is already checked out on branch '%s'.", path, hostDescription, name)
	} else {
		_, err := g.RunGitCommand(ctx, []string{"checkout", name})
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Git repository '%s' on host '%s' checked out on branch '%s'.", path, hostDescription, name)
	}

	return nil
}

func (g *GitRepository) GetCurrentBranchName(ctx context.Context) (branchName string, err error) {
	stdout, err := g.RunGitCommandAndGetStdoutAsString(ctx, []string{"rev-parse", "--abbrev-ref", "HEAD"})
	if err != nil {
		return "", err
	}

	branchName = strings.TrimSpace(stdout)

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return "", err
	}

	if branchName == "" {
		return "", tracederrors.TracedErrorf(
			"Unable to get branch name for git repository '%s' on host '%s'. branchName is empty string after evaluation.",
			path,
			hostDescription,
		)
	}

	logging.LogInfoByCtxf(ctx, "Branch '%s' is currently checked out in git repository '%s' on host '%s'.", branchName, path, hostDescription)

	return branchName, nil
}

func (g *GitRepository) CreateBranch(ctx context.Context, createOptions *parameteroptions.CreateBranchOptions) (err error) {
	if createOptions == nil {
		return tracederrors.TracedErrorNil("createOptions")
	}

	name, err := createOptions.GetName()
	if err != nil {
		return err
	}

	branchExists, err := g.BranchByNameExists(ctx, name)
	if err != nil {
		return err
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if branchExists {
		logging.LogInfof(
			"Branch '%s' already exists in git repository '%s' on host '%s'.",
			name,
			path,
			hostDescription,
		)
	} else {
		_, err = g.RunGitCommand(ctx, []string{"checkout", "-b", name})
		if err != nil {
			return err
		}

		logging.LogChangedf(
			"Branch '%s' in git repository '%s' on host '%s' created.",
			name,
			path,
			hostDescription,
		)
	}

	return nil
}

func (g *GitRepository) DeleteBranchByName(ctx context.Context, name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	branchExists, err := g.BranchByNameExists(ctx, name)
	if err != nil {
		return err
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	if branchExists {
		_, err := g.RunGitCommand(ctx, []string{"branch", "-D", name})
		if err != nil {
			return err
		}

		logging.LogChangedByCtxf(ctx, "Branch '%s' in git repository '%s' on host '%s' deleted.", name, path, hostDescription)
	} else {
		logging.LogInfoByCtxf(ctx, "Branch '%s' in git repository '%s' on host '%s' is already absent. Skip delete.", name, path, hostDescription)
	}

	return nil
}

func (g *GitRepository) ListBranchNames(ctx context.Context) (branchNames []string, err error) {
	lines, err := g.RunGitCommandAndGetStdoutAsLines(ctx, []string{"branch"})
	if err != nil {
		return nil, err
	}

	branchNames = []string{}

	for _, line := range lines {
		line = strings.TrimPrefix(line, "* ")
		line = strings.TrimSpace(line)

		branchNames = append(branchNames, line)
	}

	sort.Strings(branchNames)

	path, hostDescripton, err := g.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Found '%d' branches in git repository '%s' on host '%s'.", len(branchNames), path, hostDescripton)

	return branchNames, nil
}
