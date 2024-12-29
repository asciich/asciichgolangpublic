package asciichgolangpublic

import (
	"fmt"
)

type GitRepositoryCreateTagOptions struct {
	// Commit hash to tag.
	// If not set the currently checked out commit is tagged (depends on the implementation if supported)
	CommitHash string

	// Name and comment/ message of the tag:
	TagName    string
	TagComment string

	Verbose              bool
	PushTagsToAllRemotes bool
}

func NewGitRepositoryCreateTagOptions() (g *GitRepositoryCreateTagOptions) {
	return new(GitRepositoryCreateTagOptions)
}

func (g *GitRepositoryCreateTagOptions) IsCommitHashSet() (isSet bool) {
	return g.CommitHash != ""
}

func (g *GitRepositoryCreateTagOptions) GetCommitHash() (commitHash string, err error) {
	if g.CommitHash == "" {
		return "", TracedErrorf("CommitHash not set")
	}

	return g.CommitHash, nil
}

func (g *GitRepositoryCreateTagOptions) GetDeepCopy() (copy *GitRepositoryCreateTagOptions) {
	copy = NewGitRepositoryCreateTagOptions()

	*copy = *g

	return copy
}

func (g *GitRepositoryCreateTagOptions) GetPushTagsToAllRemotes() (pushTagsToAllRemotes bool, err error) {

	return g.PushTagsToAllRemotes, nil
}

func (g *GitRepositoryCreateTagOptions) GetTagComment() (tagComment string, err error) {
	if g.TagComment == "" {
		return "", TracedErrorf("TagComment not set")
	}

	return g.TagComment, nil
}

func (g *GitRepositoryCreateTagOptions) GetTagCommentOrDefaultIfUnset() (tagComment string) {
	if g.TagComment == "" {
		return fmt.Sprintf(
			"Create tag '%s' by '%s' version '%s'.",
			g.GetTagNameOrEmptyStringIfUnset(),
			GetSoftwareNameString(),
			GetSoftwareVersionString(),
		)
	}

	return g.TagComment
}

func (g *GitRepositoryCreateTagOptions) GetTagName() (tagName string, err error) {
	if g.TagName == "" {
		return "", TracedErrorf("TagName not set")
	}

	return g.TagName, nil
}

func (g *GitRepositoryCreateTagOptions) GetTagNameOrEmptyStringIfUnset() (tagName string) {
	return g.TagName
}

func (g *GitRepositoryCreateTagOptions) GetVerbose() (verbose bool, err error) {

	return g.Verbose, nil
}

func (g *GitRepositoryCreateTagOptions) IsTagCommentSet() (isSet bool) {
	return g.TagComment != ""
}

func (g *GitRepositoryCreateTagOptions) MustGetCommitHash() (commitHash string) {
	commitHash, err := g.GetCommitHash()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return commitHash
}

func (g *GitRepositoryCreateTagOptions) MustGetPushTagsToAllRemotes() (pushTagsToAllRemotes bool) {
	pushTagsToAllRemotes, err := g.GetPushTagsToAllRemotes()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return pushTagsToAllRemotes
}

func (g *GitRepositoryCreateTagOptions) MustGetTagComment() (tagComment string) {
	tagComment, err := g.GetTagComment()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tagComment
}

func (g *GitRepositoryCreateTagOptions) MustGetTagName() (tagName string) {
	tagName, err := g.GetTagName()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return tagName
}

func (g *GitRepositoryCreateTagOptions) MustGetVerbose() (verbose bool) {
	verbose, err := g.GetVerbose()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return verbose
}

func (g *GitRepositoryCreateTagOptions) MustSetCommitHash(commitHash string) {
	err := g.SetCommitHash(commitHash)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryCreateTagOptions) MustSetPushTagsToAllRemotes(pushTagsToAllRemotes bool) {
	err := g.SetPushTagsToAllRemotes(pushTagsToAllRemotes)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryCreateTagOptions) MustSetTagComment(tagComment string) {
	err := g.SetTagComment(tagComment)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryCreateTagOptions) MustSetTagName(tagName string) {
	err := g.SetTagName(tagName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryCreateTagOptions) MustSetVerbose(verbose bool) {
	err := g.SetVerbose(verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (g *GitRepositoryCreateTagOptions) SetCommitHash(commitHash string) (err error) {
	if commitHash == "" {
		return TracedErrorf("commitHash is empty string")
	}

	g.CommitHash = commitHash

	return nil
}

func (g *GitRepositoryCreateTagOptions) SetPushTagsToAllRemotes(pushTagsToAllRemotes bool) (err error) {
	g.PushTagsToAllRemotes = pushTagsToAllRemotes

	return nil
}

func (g *GitRepositoryCreateTagOptions) SetTagComment(tagComment string) (err error) {
	if tagComment == "" {
		return TracedErrorf("tagComment is empty string")
	}

	g.TagComment = tagComment

	return nil
}

func (g *GitRepositoryCreateTagOptions) SetTagName(tagName string) (err error) {
	if tagName == "" {
		return TracedErrorf("tagName is empty string")
	}

	g.TagName = tagName

	return nil
}

func (g *GitRepositoryCreateTagOptions) SetVerbose(verbose bool) (err error) {
	g.Verbose = verbose

	return nil
}
