package ansibleplaybook

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"gopkg.in/yaml.v3"
)

type Playbook struct {
	Plays []*Play
}

func ReadPlaybook(ctx context.Context, path string) (*Playbook, error) {
	if path == "" {
		return nil, tracederrors.TracedErrorEmptyString("path")
	}

	logging.LogInfoByCtxf(ctx, "Read ansible playbook '%s' started.", path)

	in, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parsed := new([]*Play)

	err = yaml.Unmarshal(in, parsed)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to read ansible playbook '%s': %w", path, err)
	}

	logging.LogInfoByCtxf(ctx, "Read ansible playbook '%s' finished.", path)

	return &Playbook{
		Plays: *parsed,
	}, nil
}

func WritePlaybook(ctx context.Context, playbook *Playbook, path string) error {
	if playbook == nil {
		return tracederrors.TracedErrorNil("playbook")
	}

	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	toWrite, err := yaml.Marshal(playbook.Plays)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to marshal playbook: %w", err)
	}

	err = os.WriteFile(path, toWrite, 0644)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to write playbook to file '%s': %w", path, err)
	}

	logging.LogInfoByCtxf(ctx, "Wrote ansible playbook to '%s'.", path)

	return nil
}
