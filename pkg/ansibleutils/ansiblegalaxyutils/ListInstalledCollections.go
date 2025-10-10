package ansiblegalaxyutils

import (
	"context"
	"encoding/json"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexec"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ListInstalledCollections(ctx context.Context, options *ListInstalledCollectionsOptions) (map[string]string, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	ansibleGalaxyCmd, err := options.GetAnsibleGalaxyPath()
	if err != nil {
		return nil, err
	}

	type Collection struct {
		Version string `json:"version"`
	}

	type InstallPaths map[string]Collection

	type Output map[string]InstallPaths

	output, err := commandexecutorexec.RunCommand(
		ctx,
		&parameteroptions.RunCommandOptions{
			Command: []string{ansibleGalaxyCmd, "collection", "list", "--format=json"},
		},
	)
	if err != nil {
		return nil, err
	}

	outputBytes, err := output.GetStdoutAsBytes()
	if err != nil {
		return nil, err
	}

	var parsed Output

	err = json.Unmarshal(outputBytes, &parsed)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to unmrashal json of ansible-galaxy collection list: %w", err)
	}

	ret := map[string]string{}
	for _, installPath := range parsed {
		for name, collection := range installPath {
			version := collection.Version
			ret[name] = version
		}
	}

	logging.LogInfoByCtxf(ctx, "There are '%d' ansible collections installed.", len(ret))

	return ret, nil
}
