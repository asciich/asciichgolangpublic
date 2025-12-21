package runbook

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type Step struct {
	Name        string
	Description string
	Run         func(context.Context) error
}

func (s *Step) GetDescription() (string, error) {
	if s.Description == "" {
		return "", tracederrors.TracedError("Description not set")
	}

	return s.Description, nil
}

func (s *Step) GetName() (string, error) {
	if s.Name == "" {
		return "", tracederrors.TracedError("Name not set")
	}

	return s.Name, nil
}

func (s *Step) Execute(ctx context.Context) error {
	name, err := s.GetName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Step '%s' started.", name)

	if s.Run == nil {
		return tracederrors.TracedError("Run function not set")
	}

	err = s.Run(ctx)
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Step '%s' finished.", name)

	return nil
}

func (s *Step) Validate(ctx context.Context) error {
	name, err := s.GetName()
	if err != nil {
		return tracederrors.TracedErrorf("The step has no name")
	}

	_, err = s.GetDescription()
	if err != nil {
		return tracederrors.TracedErrorf("The step '%s' has no description.", name)
	}

	if s.Run == nil {
		return tracederrors.TracedErrorf("The Run function is not set for step '%s'.", name)
	}

	logging.LogInfoByCtxf(ctx, "Step '%s' validated successfully.", name)

	return nil
}
