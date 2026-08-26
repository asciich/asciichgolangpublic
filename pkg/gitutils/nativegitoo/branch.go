package nativegitoo

import (
	"context"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (n *NativeGitRepository) CreateBranch(ctx context.Context, createOptions *parameteroptions.CreateBranchOptions) (err error) {
	if createOptions == nil {
		return tracederrors.TracedErrorNil("createOptions")
	}

	branchName, err := createOptions.GetName()
	if err != nil {
		return err
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Create branch '%s' in git repository '%s' on '%s' started.", branchName, path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	headRef, err := repo.Head()
	if err != nil {
		return tracederrors.TracedErrorf("Get HEAD for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	refName := plumbing.NewBranchReferenceName(branchName)
	ref := plumbing.NewHashReference(refName, headRef.Hash())

	err = repo.Storer.SetReference(ref)
	if err != nil {
		return tracederrors.TracedErrorf("Create branch '%s' in git repository '%s' on '%s' failed: %w", branchName, path, hostDescription, err)
	}

	logging.LogInfoByCtxf(ctx, "Create branch '%s' in git repository '%s' on '%s' finished.", branchName, path, hostDescription)

	return nil
}

func (n *NativeGitRepository) CheckoutBranchByName(ctx context.Context, name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Checkout branch '%s' in git repository '%s' on '%s' started.", name, path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return tracederrors.TracedErrorf("Get worktree for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name),
	})
	if err != nil {
		return tracederrors.TracedErrorf("Checkout branch '%s' in git repository '%s' on '%s' failed: %w", name, path, hostDescription, err)
	}

	logging.LogInfoByCtxf(ctx, "Checkout branch '%s' in git repository '%s' on '%s' finished.", name, path, hostDescription)

	return nil
}

func (n *NativeGitRepository) GetCurrentBranchName(ctx context.Context) (branchName string, err error) {
	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return "", err
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	headRef, err := repo.Head()
	if err != nil {
		return "", tracederrors.TracedErrorf("Get HEAD for git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	if !headRef.Name().IsBranch() {
		return "", tracederrors.TracedErrorf("HEAD in git repository '%s' on '%s' is not pointing to a branch.", path, hostDescription)
	}

	branchName = headRef.Name().Short()

	return branchName, nil
}

func (n *NativeGitRepository) DeleteBranchByName(ctx context.Context, name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	path, hostDescription, err := n.GetPathAndHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete branch '%s' in git repository '%s' on '%s' started.", name, path, hostDescription)

	repo, err := git.PlainOpen(path)
	if err != nil {
		return tracederrors.TracedErrorf("Open git repository '%s' on '%s' failed: %w", path, hostDescription, err)
	}

	refName := plumbing.NewBranchReferenceName(name)

	err = repo.Storer.RemoveReference(refName)
	if err != nil {
		return tracederrors.TracedErrorf("Delete branch '%s' in git repository '%s' on '%s' failed: %w", name, path, hostDescription, err)
	}

	logging.LogInfoByCtxf(ctx, "Delete branch '%s' in git repository '%s' on '%s' finished.", name, path, hostDescription)

	return nil
}
