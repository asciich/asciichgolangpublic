package commandexecutorheadscale

import (
	"context"
	"encoding/json"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type ListNodeTimestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int64 `json:"nanos,omitempty"`
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
	Expiration ListNodeTimestamp `json:"expiration"`
	CreatedAt  ListNodeTimestamp `json:"created_at"`
}

type ListNodeEntry struct {
	ID             int                `json:"id"`
	MachineKey     string             `json:"machine_key"`
	NodeKey        string             `json:"node_key"`
	DiscoKey       string             `json:"disco_key"`
	IPAddresses    []string           `json:"ip_addresses"`
	Name           string             `json:"name"`
	User           ListNodeUser       `json:"user"`
	LastSeen       ListNodeTimestamp  `json:"last_seen"`
	Expiry         ListNodeTimestamp  `json:"expiry"`
	PreAuthKey     ListNodePreAuthKey `json:"pre_auth_key"`
	CreatedAt      ListNodeTimestamp  `json:"created_at"`
	RegisterMethod int                `json:"register_method"`
	GivenName      string             `json:"given_name"`
	Online         bool               `json:"online"`
}

func (l *ListNodeEntry) GetName() (string, error) {
	if l.Name == "" {
		return "", tracederrors.TracedError("Name not set")
	}

	return l.Name, nil
}

func ListNodes(ctx context.Context, commandExectuor commandexecutorinterfaces.CommandExecutor) ([]ListNodeEntry, error) {
	if commandExectuor == nil {
		return nil, tracederrors.TracedErrorNil("commandExecutor")
	}

	hostDescription, err := commandExectuor.GetHostDescription()
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "List headscale nodes on '%s' started.", hostDescription)

	output, err := commandExectuor.RunCommandAndGetStdoutAsBytes(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"headscale", "nodes", "list", "-ojson"},
	})
	if err != nil {
		return nil, err
	}

	nodes := []ListNodeEntry{}
	err = json.Unmarshal(output, &nodes)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to parse JSON of listed headscale nodes: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "List headscale nodes on '%s' finished.", hostDescription)

	return nodes, nil
}

func ListNodeNames(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor) ([]string, error) {
	nodes, err := ListNodes(ctx, commandExecutor)
	if err != nil {
		return nil, err
	}

	nodeNames := make([]string, 0, len(nodes))
	for _, node := range nodes {
		toAdd, err := node.GetName()
		if err != nil {
			return nil, err
		}

		nodeNames = append(nodeNames, toAdd)
	}

	return nodeNames, nil
}
