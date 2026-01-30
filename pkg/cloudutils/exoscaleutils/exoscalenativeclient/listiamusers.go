package exoscalenativeclient

import (
	"context"

	v3 "github.com/exoscale/egoscale/v3"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ListIamUsers(ctx context.Context, client *v3.Client) ([]string, error) {
	logging.LogInfoByCtxf(ctx, "List exoscale IAM users started.")

	if client == nil {
		return nil, tracederrors.TracedErrorNil("client")
	}

	resp, err := client.ListUsers(ctx)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to list exoscale IAM users: %w", err)
	}

	users := []string{}
	for _, u := range resp.Users {
		users = append(users, u.Email)
	}

	logging.LogInfoByCtxf(ctx, "List exoscale IAM users finished.")

	return users, nil
}
