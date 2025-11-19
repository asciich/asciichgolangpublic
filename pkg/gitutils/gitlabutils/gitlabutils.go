package gitlabutils

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const DEFAULT_TOKEN_NAME = "GITLAB_TOKEN"

func NewAuthenticatedGitlab(ctx context.Context, gitlabUrl string) (*asciichgolangpublic.GitlabInstance, error) {
	if gitlabUrl == "" {
		return nil, tracederrors.TracedErrorEmptyString("gitlabUrl")
	}

	gitlab, err := asciichgolangpublic.GetGitlabByFQDN(gitlabUrl)
	if err != nil {
		return nil, err
	}

	token := os.Getenv(DEFAULT_TOKEN_NAME)
	if token == "" {
		return nil, tracederrors.TracedErrorf("No token to authenticate gitlab. Please set the env var '%s' with the gitlab token to use.", DEFAULT_TOKEN_NAME)
	}

	err = gitlab.Authenticate(ctx, &asciichgolangpublic.GitlabAuthenticationOptions{
		AccessToken: token,
	})
	if err != nil {
		return nil, err
	}

	return gitlab, nil
}
