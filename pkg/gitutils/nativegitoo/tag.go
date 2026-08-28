package nativegitoo

import (
	"context"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/gitutils/gitparameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"github.com/asciich/asciichgolangpublic/pkg/versionutils"
)

type NativeGitTag struct {
	name       string
	repository *NativeGitRepository
}

func NewNativeGitTag(name string, repository *NativeGitRepository) *NativeGitTag {
	return &NativeGitTag{
		name:       name,
		repository: repository,
	}
}

func (t *NativeGitTag) GetHash(ctx context.Context) (hash string, err error) {
	path, err := t.repository.GetPath()
	if err != nil {
		return "", err
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", tracederrors.TracedErrorf("Open git repository '%s' failed: %w", path, err)
	}

	tagRef, err := repo.Tag(t.name)
	if err != nil {
		return "", tracederrors.TracedErrorf("Get tag '%s' failed: %w", t.name, err)
	}

	tagObj, err := repo.TagObject(tagRef.Hash())
	if err == nil {
		// Annotated tag - return the target commit hash
		return tagObj.Target.String(), nil
	}

	// Lightweight tag - return the tag hash directly
	return tagRef.Hash().String(), nil
}

func (t *NativeGitTag) GetName() (name string, err error) {
	if t.name == "" {
		return "", tracederrors.TracedErrorEmptyString("name")
	}
	return t.name, nil
}

func (t *NativeGitTag) GetGitRepository() (repo gitinterfaces.GitRepository, err error) {
	return t.repository, nil
}

func (t *NativeGitTag) IsVersionTag() (isVersionTag bool, err error) {
	return versionutils.IsVersionString(t.name), nil
}

func (t *NativeGitTag) SetName(name string) (err error) {
	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}
	t.name = name
	return nil
}

func (t *NativeGitTag) GetVersion() (version versionutils.Version, err error) {
	return versionutils.NewFromString(t.name)
}

func (n *NativeGitRepository) GetTagByName(name string) (tag gitinterfaces.GitTag, err error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	return NewNativeGitTag(name, n), nil
}

func (n *NativeGitRepository) ListTagNames(ctx context.Context) (tagNames []string, err error) {
	path, err := n.GetPath()
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Open git repository '%s' failed: %w", path, err)
	}

	tagsIter, err := repo.Tags()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Get tags for git repository '%s' failed: %w", path, err)
	}

	tagNames = []string{}
	err = tagsIter.ForEach(func(ref *plumbing.Reference) error {
		tagNames = append(tagNames, ref.Name().Short())
		return nil
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Iterate tags for git repository '%s' failed: %w", path, err)
	}

	return tagNames, nil
}

func (n *NativeGitRepository) ListTags(ctx context.Context) (tags []gitinterfaces.GitTag, err error) {
	tagNames, err := n.ListTagNames(ctx)
	if err != nil {
		return nil, err
	}

	tags = []gitinterfaces.GitTag{}
	for _, name := range tagNames {
		toAdd, err := n.GetTagByName(name)
		if err != nil {
			return nil, err
		}

		tags = append(tags, toAdd)
	}

	return tags, nil
}

func (n *NativeGitRepository) ListTagsForCommitHash(ctx context.Context, hash string) (tags []gitinterfaces.GitTag, err error) {
	if hash == "" {
		return nil, tracederrors.TracedErrorEmptyString("hash")
	}

	path, err := n.GetPath()
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Open git repository '%s' failed: %w", path, err)
	}

	hashObj := plumbing.NewHash(hash)

	tagsIter, err := repo.Tags()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Get tags for git repository '%s' failed: %w", path, err)
	}

	tags = []gitinterfaces.GitTag{}
	err = tagsIter.ForEach(func(ref *plumbing.Reference) error {
		tagObj, err := repo.TagObject(ref.Hash())
		if err == nil {
			// Annotated tag - check if target matches
			if tagObj.Target == hashObj {
				tag := gitgeneric.NewGitRepositoryTag()
				err := tag.SetName(ref.Name().Short())
				if err != nil {
					return err
				}
				err = tag.SetGitRepository(n)
				if err != nil {
					return err
				}
				tags = append(tags, tag)
			}
			return nil
		}

		// Lightweight tag - ref hash is the commit hash
		if ref.Hash() == hashObj {
			tag := gitgeneric.NewGitRepositoryTag()
			err := tag.SetName(ref.Name().Short())
			if err != nil {
				return err
			}
			err = tag.SetGitRepository(n)
			if err != nil {
				return err
			}
			tags = append(tags, tag)
		}
		return nil
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Iterate tags for git repository '%s' failed: %w", path, err)
	}

	return tags, nil
}

func (n *NativeGitRepository) CreateTag(ctx context.Context, createOptions *gitparameteroptions.GitRepositoryCreateTagOptions) (createdTag gitinterfaces.GitTag, err error) {
	if createOptions == nil {
		return nil, tracederrors.TracedErrorNil("createOptions")
	}

	tagName, err := createOptions.GetTagName()
	if err != nil {
		return nil, err
	}

	path, err := n.GetPath()
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Open git repository '%s' failed: %w", path, err)
	}

	hashToTag := ""
	if createOptions.IsCommitHashSet() {
		hashToTag, err = createOptions.GetCommitHash()
		if err != nil {
			return nil, err
		}
	} else {
		head, err := repo.Head()
		if err != nil {
			return nil, tracederrors.TracedErrorf("Get HEAD for git repository '%s' failed: %w", path, err)
		}
		hashToTag = head.Hash().String()
	}

	tagMessage := tagName
	if createOptions.IsTagCommentSet() {
		tagMessage, err = createOptions.GetTagComment()
		if err != nil {
			return nil, err
		}
	}

	hashObj := plumbing.NewHash(hashToTag)

	// Get config for tagger info
	config, err := repo.Config()
	if err != nil {
		return nil, tracederrors.TracedErrorf("Get config for git repository '%s' failed: %w", path, err)
	}

	// Create annotated tag
	_, err = repo.CreateTag(tagName, hashObj, &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  config.User.Name,
			Email: config.User.Email,
			When:  time.Now(),
		},
		Message: tagMessage,
	})
	if err != nil {
		return nil, tracederrors.TracedErrorf("Create tag '%s' for commit '%s' failed: %w", tagName, hashToTag, err)
	}

	createdTag = NewNativeGitTag(tagName, n)

	return createdTag, nil
}
