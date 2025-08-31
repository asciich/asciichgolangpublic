package gitparameteroptions

import (
	"fmt"

	"github.com/asciich/asciichgolangpublic/pkg/binaryinfo"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type GitRepositoryCreateTagOptions struct {
	// Commit hash to tag.
	// If not set the currently checked out commit is tagged (depends on the implementation if supported)
	CommitHash string

	// Name and comment/ message of the tag:
	TagName    string
	TagComment string

	PushTagsToAllRemotes bool
}

func NewGitRepositoryCreateTagOptions() (g *GitRepositoryCreateTagOptions) {
	return new(GitRepositoryCreateTagOptions)
}

func (g *GitRepositoryCreateTagOptions) GetCommitHash() (commitHash string, err error) {
	if g.CommitHash == "" {
		return "", tracederrors.TracedErrorf("CommitHash not set")
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
		return "", tracederrors.TracedErrorf("TagComment not set")
	}

	return g.TagComment, nil
}

func (g *GitRepositoryCreateTagOptions) GetTagCommentOrDefaultIfUnset() (tagComment string) {
	if g.TagComment == "" {
		return fmt.Sprintf(
			"Create tag '%s' by '%s' version '%s'.",
			g.GetTagNameOrEmptyStringIfUnset(),
			binaryinfo.GetSoftwareNameString(),
			binaryinfo.GetSoftwareVersionString(),
		)
	}

	return g.TagComment
}

func (g *GitRepositoryCreateTagOptions) GetTagName() (tagName string, err error) {
	if g.TagName == "" {
		return "", tracederrors.TracedErrorf("TagName not set")
	}

	return g.TagName, nil
}

func (g *GitRepositoryCreateTagOptions) GetTagNameOrEmptyStringIfUnset() (tagName string) {
	return g.TagName
}

func (g *GitRepositoryCreateTagOptions) IsCommitHashSet() (isSet bool) {
	return g.CommitHash != ""
}

func (g *GitRepositoryCreateTagOptions) IsTagCommentSet() (isSet bool) {
	return g.TagComment != ""
}

func (g *GitRepositoryCreateTagOptions) SetCommitHash(commitHash string) (err error) {
	if commitHash == "" {
		return tracederrors.TracedErrorf("commitHash is empty string")
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
		return tracederrors.TracedErrorf("tagComment is empty string")
	}

	g.TagComment = tagComment

	return nil
}

func (g *GitRepositoryCreateTagOptions) SetTagName(tagName string) (err error) {
	if tagName == "" {
		return tracederrors.TracedErrorf("tagName is empty string")
	}

	g.TagName = tagName

	return nil
}
