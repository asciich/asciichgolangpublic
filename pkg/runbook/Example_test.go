package runbook_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/runbook"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// This is a simple example to show the ideas behind the runbook package using an example.
//
// For simplicity only a nummerical calculation is performed to keep the example small and focused on the runbook package.
func Test_Example_test(t *testing.T) {
	// We use a verbose output for this example:
	ctx := contextutils.ContextVerbose()

	// For simplicity we just touch this integer:
	var number int

	// lets define our process with 3 simple steps:
	process := &runbook.RunBook{
		Name:        "Simple example",
		Description: "A simple example process to show the 'runbook' package.",
		Steps: []runbook.Runnable{
			&runbook.Step{
				Name: "Set initial values",
				Description: `In this step we set the initial values.

For this example we just set the 'number' to 5.`,
				Run: func(ctx context.Context) error {
					number = 5
					return nil
				},
			},
			&runbook.Step{
				Name: "Perform calculation",
				Description: `In this step we perform the actual calculation.

For this example we just multiply the 'number' by 2`,
				Run: func(ctx context.Context) error {
					number *= 2
					return nil
				},
			},
			&runbook.Step{
				Name: "validate",
				Description: `As a last step we do some sanity checks and validation.

For this example we just check the value of 'number'.`,
				Run: func(ctx context.Context) error {
					if number == 10 {
						logging.LogInfoByCtxf(ctx, "the 'number' has an expected value of %d.", number)
					} else {
						// By using a traced error we include the stack trace in the error message.
						// This allows us to easily find the position in the code where the error occured:
						return tracederrors.TracedErrorf("unexpected 'number' value %d.", number)
					}

					return nil
				},
			},
		},
	}

	// So our process consists of 3 steps:
	require.EqualValues(t, 3, process.GetNSteps(ctx))

	// For documentation purposes we can get the name and description of the process:
	name, err := process.GetName()
	require.NoError(t, err)
	require.EqualValues(t, "Simple example", name)

	description, err := process.GetDescription()
	require.NoError(t, err)
	require.EqualValues(t, "A simple example process to show the 'runbook' package.", description)

	// And we can get the steps documentation describing what is done step by step:
	expected := `1: Set initial values
    In this step we set the initial values.

    For this example we just set the 'number' to 5.
2: Perform calculation
    In this step we perform the actual calculation.

    For this example we just multiply the 'number' by 2
3: validate
    As a last step we do some sanity checks and validation.

    For this example we just check the value of 'number'.
`
	stepDocumentation, err := process.DocumentSteps()
	require.NoError(t, err)
	assert.EqualValues(t, expected, stepDocumentation)

	// To run all steps conscutively:
	err = process.Execute(ctx)
	require.NoError(t, err)
}
