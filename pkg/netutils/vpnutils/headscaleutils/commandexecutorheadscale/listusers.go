package commandexecutorheadscale

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// User represents the main object inside the array
type listUserEntry struct {
	ID        int                    `json:"id"`
	Name      string                 `json:"name"`
	CreatedAt listUserEntryTimestamp `json:"created_at"`
}

// Timestamp matches the nested object structure
type listUserEntryTimestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int64 `json:"nanos"`
}

func ListUsersRaw(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) ([]*listUserEntry, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "List headscale users on '%s' started.", hostDescription)

	output, err := commandExecutor.RunCommandAndGetStdoutAsBytes(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{"headscale", "user", "list", "-o", "json"},
		},
	)
	if err != nil {
		return nil, err
	}

	ret := []*listUserEntry{}
	err = json.Unmarshal(output, &ret)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to unmarshal headscale users: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "List headscale users on '%s' finished. Found %d users.", hostDescription, len(ret))

	return ret, nil
}

func ListUserNames(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) ([]string, error) {
	rawUsers, err := ListUsersRaw(ctx, commandExecutor)
	if err != nil {
		return nil, err
	}

	users := []string{}
	for _, u := range rawUsers {
		users = append(users, u.Name)
	}

	sort.Strings(users)

	return users, nil
}
