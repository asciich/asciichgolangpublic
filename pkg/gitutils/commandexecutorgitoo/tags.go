package commandexecutorgitoo

import (
	"context"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (g *GitRepository) CreateTag(ctx context.Context, options *gitparameteroptions.GitRepositoryCreateTagOptions) (createdTag gitinterfaces.GitTag, err error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	tagName, err := options.GetTagName()
	if err != nil {
		return nil, err
	}

	tagMessage := tagName
	if options.IsTagCommentSet() {
		tagMessage, err = options.GetTagComment()
		if err != nil {
			return nil, err
		}
	}

	hashToTag := ""
	if options.IsCommitHashSet() {
		hashToTag, err = options.GetCommitHash()
		if err != nil {
			return nil, err
		}
	} else {
		hashToTag, err = g.GetCurrentCommitHash(ctx)
		if err != nil {
			return nil, err
		}
	}

	_, err = g.RunGitCommand(ctx, []string{"tag", "-a", tagName, hashToTag, "-m", tagMessage})
	if err != nil {
		return nil, err
	}

	path, hostDescription, err := g.GetPathAndHostDescription()
	if err != nil {
		return nil, err
	}

	createdTag, err = g.GetTagByName(tagName)
	if err != nil {
		return nil, err
	}

	logging.LogChangedByCtxf(ctx, "Created tag '%s' for commit '%s' in git repository '%s' on host '%s'.", tagName, hashToTag, path, hostDescription)

	return createdTag, nil
}

func (g *GitRepository) GetTagByName(name string) (tag gitinterfaces.GitTag, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	toReturn := gitgeneric.NewGitRepositoryTag()

	err = toReturn.SetName(name)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetName(name)
	if err != nil {
		return nil, err
	}

	err = toReturn.SetGitRepository(g)
	if err != nil {
		return nil, err
	}

	return toReturn, nil
}

func (g *GitRepository) GetHashByTagName(tagName string) (hash string, err error) {
	if tagName == "" {
		return "", tracederrors.TracedErrorEmptyString("tagName")
	}

	stdoutLines, err := g.RunGitCommandAndGetStdoutAsLines(
		contextutils.ContextSilent(),
		[]string{"show-ref", "--dereference", tagName},
	)
	if err != nil {
		return "", err
	}

	for _, line := range stdoutLines {
		if strings.HasSuffix(line, "{}") {
			hash = strings.Split(line, " ")[0]
			break
		}
	}

	hash = strings.TrimSpace(hash)

	if hash == "" {
		return "", tracederrors.TracedError("hash is empty string after evaluation")
	}

	return hash, nil
}

func (g *GitRepository) ListTagNames(ctx context.Context) (tagNames []string, err error) {
	return g.RunGitCommandAndGetStdoutAsLines(
		contextutils.WithSilent(ctx), // Do not clutter output by pritning all tags.
		[]string{"tag"},
	)
}

func (g *GitRepository) ListTags(ctx context.Context) (tags []gitinterfaces.GitTag, err error) {
	tagNames, err := g.ListTagNames(ctx)
	if err != nil {
		return nil, err
	}

	tags = []gitinterfaces.GitTag{}
	for _, name := range tagNames {
		toAdd, err := g.GetTagByName(name)
		if err != nil {
			return nil, err
		}

		tags = append(tags, toAdd)
	}

	return tags, nil
}

func (g *GitRepository) ListTagsForCommitHash(ctx context.Context, hash string) (tags []gitinterfaces.GitTag, err error) {
	if hash == "" {
		return nil, tracederrors.TracedErrorEmptyString("hash")
	}

	tagNames, err := g.RunGitCommandAndGetStdoutAsLines(
		contextutils.ContextSilent(),
		[]string{"tag", "--points-at", "HEAD"},
	)
	if err != nil {
		return nil, err
	}

	tags = []gitinterfaces.GitTag{}
	for _, name := range tagNames {
		toAdd, err := g.GetTagByName(name)
		if err != nil {
			return nil, err
		}

		tags = append(tags, toAdd)
	}

	return tags, nil
}
