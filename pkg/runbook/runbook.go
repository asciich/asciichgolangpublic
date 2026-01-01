package runbook

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Defines a runbook
type RunBook struct {
	Name        string
	Description string
	Steps       []Runnable
}

func (r *RunBook) GetDescription() (string, error) {
	if r.Description == "" {
		return "", tracederrors.TracedError("Description not set")
	}

	return r.Description, nil
}

func (r *RunBook) GetName() (string, error) {
	if r.Name == "" {
		return "", tracederrors.TracedError("Name not set")
	}

	return r.Name, nil
}

func (r *RunBook) GetNSteps(ctx context.Context) int {
	nsteps := len(r.Steps)

	logging.LogInfoByCtxf(ctx, "Process '%s' has %d steps in total.", r.Name, nsteps)

	return nsteps
}

func (r *RunBook) DocumentSteps() (string, error) {
	if len(r.Steps) == 0 {
		return "", tracederrors.TracedError("Process has no steps to document.")
	}

	var documentation strings.Builder

	for i, step := range r.Steps {
		name, err := step.GetName()
		if err != nil {
			return "", err
		}

		description, err := step.GetDescription()
		if err != nil {
			return "", err
		}

		fmt.Fprintf(&documentation, "%d: %s\n", i+1, name)

		toWrite := stringsutils.EnsureEndsWithExactlyOneLineBreak(stringsutils.AddIndent(description, "    "))
		documentation.WriteString(toWrite)
	}

	return documentation.String(), nil
}

func (r *RunBook) Validate(ctx context.Context) error {
	name, err := r.GetName()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Validate runbook '%s' started.", name)

	if len(r.Steps) == 0 {
		return tracederrors.TracedErrorf("There are no steps defined in the runbook '%s'.", name)
	}

	for _, s := range r.Steps {
		err := s.Validate(ctx)
		if err != nil {
			return err
		}

	}

	logging.LogInfoByCtxf(ctx, "Validate runbook '%s' finished.", name)

	return nil
}

func (r *RunBook) Execute(ctx context.Context) error {
	name, err := r.GetName()
	if err != nil {
		return err
	}

	tStart := time.Now()

	logging.LogInfoByCtxf(ctx, "Runbook '%s' started.", name)

	err = r.Validate(ctx)
	if err != nil {
		return err
	}

	for _, s := range r.Steps {
		err := s.Execute(ctx)
		if err != nil {
			return err
		}
	}

	duration := time.Since(tStart)

	logging.LogInfoByCtxf(ctx, "Runbook '%s' finished (took %s).", name, duration)

	return nil
}
