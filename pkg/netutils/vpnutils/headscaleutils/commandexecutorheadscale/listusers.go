package commandexecutorheadscale

import (
	"context"
	"encoding/json"
	"sort"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalegeneric"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// ListUsersEntry is one entry of `headscale users list -o json`.
//
// Mapping to the columns of `headscale users list`:
//   - Name     = `display_name`
//   - Username = `name`
type ListUsersEntry struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	DisplayName string            `json:"display_name"`
	Email       string            `json:"email"`
	CreatedAt   ListNodeTimestamp `json:"created_at"`
}

// ToUser converts the entry into a headscalegeneric.User.
func (l *ListUsersEntry) ToUser() (*headscalegeneric.User, error) {
	if l == nil {
		return nil, tracederrors.TracedErrorNil("listUsersEntry")
	}

	if l.ID <= 0 {
		return nil, tracederrors.TracedErrorf("Invalid user id '%d'.", l.ID)
	}

	created := l.CreatedAt.ToTime()
	if created == nil {
		return nil, tracederrors.TracedErrorf("Created time of user with id '%d' is not set.", l.ID)
	}

	ret := headscalegeneric.NewUser()

	steps := []error{
		ret.SetID(uint64(l.ID)),
		ret.SetName(l.DisplayName),
		ret.SetUsername(l.Name),
		ret.SetEmail(l.Email),
		ret.SetCreated(*created),
	}

	for _, err := range steps {
		if err != nil {
			return nil, tracederrors.TracedErrorf("Unable to convert headscale user with id '%d': %w", l.ID, err)
		}
	}

	return ret, nil
}

// ListUsers returns all headscale users with details.
func ListUsers(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) ([]headscaleinterfaces.User, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "List headscale users on '%s' started.", hostDescription)

	output, err := commandExecutor.RunCommandAndGetStdoutAsBytes(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"headscale", "users", "list", "-ojson"},
	})
	if err != nil {
		return nil, err
	}

	entries := []ListUsersEntry{}
	err = json.Unmarshal(output, &entries) // "null" leaves entries nil, which is fine.
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse JSON of listed headscale users: %w", err)
	}

	ret := make([]headscaleinterfaces.User, 0, len(entries))
	for i := range entries {
		user, err := entries[i].ToUser()
		if err != nil {
			return nil, err
		}

		ret = append(ret, user)
	}

	logging.LogInfoByCtxf(ctx, "List headscale users on '%s' finished. Found '%d' users.", hostDescription, len(ret))

	return ret, nil
}

// ListUserNames returns the usernames of all headscale users, sorted alphabetically.
func ListUserNames(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) ([]string, error) {
	users, err := ListUsers(ctx, commandExecutor)
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(users))
	for _, user := range users {
		username, err := user.GetUsername()
		if err != nil {
			return nil, err
		}

		ret = append(ret, username)
	}

	sort.Strings(ret)

	return ret, nil
}
