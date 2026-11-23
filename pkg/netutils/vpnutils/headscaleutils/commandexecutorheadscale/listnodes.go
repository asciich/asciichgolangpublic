package commandexecutorheadscale

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscalegeneric"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ListNodeTimestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int64 `json:"nanos,omitempty"`
}

// ToTime converts the timestamp to a time.
// Returns nil if the timestamp is not set (missing in the JSON, 0 or the Go zero time 0001-01-01 00:00:00).
func (l ListNodeTimestamp) ToTime() *time.Time {
	if l.Seconds == 0 && l.Nanos == 0 {
		return nil
	}

	ret := time.Unix(l.Seconds, l.Nanos).UTC()
	if ret.IsZero() {
		return nil
	}

	return &ret
}

type ListNodeUser struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	CreatedAt ListNodeTimestamp `json:"created_at"`
}

type ListNodePreAuthKey struct {
	User       ListNodeUser      `json:"user"`
	ID         int               `json:"id"`
	Key        string            `json:"key"`
	Used       bool              `json:"used"`
	Ephemeral  bool              `json:"ephemeral"`
	Expiration ListNodeTimestamp `json:"expiration"`
	CreatedAt  ListNodeTimestamp `json:"created_at"`
}

type ListNodeEntry struct {
	ID             int                `json:"id"`
	MachineKey     string             `json:"machine_key"`
	NodeKey        string             `json:"node_key"`
	DiscoKey       string             `json:"disco_key"`
	IPAddresses    []string           `json:"ip_addresses"`
	Name           string             `json:"name"` // The hostname reported by the client.
	User           ListNodeUser       `json:"user"`
	LastSeen       ListNodeTimestamp  `json:"last_seen"`
	Expiry         ListNodeTimestamp  `json:"expiry"`
	PreAuthKey     ListNodePreAuthKey `json:"pre_auth_key"`
	CreatedAt      ListNodeTimestamp  `json:"created_at"`
	RegisterMethod int                `json:"register_method"`
	GivenName      string             `json:"given_name"` // The name used for MagicDNS.
	Online         bool               `json:"online"`
	ForcedTags     []string           `json:"forced_tags"`
	Tags           []string           `json:"tags"`
}

// ToNode converts the entry into a headscalegeneric.Node.
//
// Mapping to the columns of `headscale nodes list`:
//   - Hostname = `name`
//   - Name     = `given_name`
func (l *ListNodeEntry) ToNode(ctx context.Context) (*headscalegeneric.Node, error) {
	if l == nil {
		return nil, tracederrors.TracedErrorNil("listNodeEntry")
	}

	if l.ID <= 0 {
		return nil, tracederrors.TracedErrorf("Invalid node id '%d'.", l.ID)
	}

	// Union of all tag fields without duplicates.
	tagSet := map[string]struct{}{}
	for _, tag := range append(append([]string{}, l.ForcedTags...), l.Tags...) {
		tagSet[tag] = struct{}{}
	}
	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	ret := headscalegeneric.NewNode()

	steps := []error{
		ret.SetID(uint64(l.ID)),
		ret.SetHostname(l.Name),
		ret.SetName(l.GivenName),
		ret.SetMachineKey(l.MachineKey),
		ret.SetNodeKey(l.NodeKey),
		ret.SetUser(l.User.Name),
		ret.SetTags(tags),
		ret.SetIPAddresses(ctx, l.IPAddresses),
		ret.SetEphemeral(l.PreAuthKey.Ephemeral),
		ret.SetLastSeen(l.LastSeen.ToTime()),
		ret.SetExpiry(l.Expiry.ToTime()),
		ret.SetOnline(l.Online),
	}

	for _, err := range steps {
		if err != nil {
			return nil, tracederrors.TracedErrorf("Unable to convert headscale node with id '%d': %w", l.ID, err)
		}
	}

	return ret, nil
}

func listNodeEntries(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) ([]ListNodeEntry, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	output, err := commandExecutor.RunCommandAndGetStdoutAsBytes(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"headscale", "nodes", "list", "-ojson"},
	})
	if err != nil {
		return nil, err
	}

	ret := []ListNodeEntry{}
	err = json.Unmarshal(output, &ret)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse JSON of listed headscale nodes: %w", err)
	}

	return ret, nil
}

// ListNodes returns all headscale nodes with details.
func ListNodes(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) ([]headscaleinterfaces.Node, error) {
	if commandExecutor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "List headscale nodes on '%s' started.", hostDescription)

	entries, err := listNodeEntries(ctx, commandExecutor)
	if err != nil {
		return nil, err
	}

	ret := make([]headscaleinterfaces.Node, 0, len(entries))
	for i := range entries {
		node, err := entries[i].ToNode(ctx)
		if err != nil {
			return nil, err
		}

		ret = append(ret, node)
	}

	logging.LogInfoByCtxf(ctx, "List headscale nodes on '%s' finished. Found '%d' nodes.", hostDescription, len(ret))

	return ret, nil
}

// ListNodeNames returns the names of all headscale nodes, sorted alphabetically.
func ListNodeNames(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) ([]string, error) {
	nodes, err := ListNodes(ctx, commandExecutor)
	if err != nil {
		return nil, err
	}

	ret := make([]string, 0, len(nodes))
	for _, node := range nodes {
		name, err := node.GetName()
		if err != nil {
			return nil, err
		}

		ret = append(ret, name)
	}

	sort.Strings(ret)

	return ret, nil
}
