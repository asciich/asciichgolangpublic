package ansibleutils

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

var ErrUnknwnAnsibleCliOutput = errors.New("unknown ansible CLI output")

var regexListHostsOutput = regexp.MustCompile(`^\s*hosts \(\d+\):`)

func isListHostsOutput(toCheck string) (isOutput bool) {
	return regexListHostsOutput.Match([]byte(toCheck))
}

func parseListHostsCliOutput(ctx context.Context, cliOutput string) (ansibleOutput *AnsibleCliOuput, err error) {
	ansibleOutput = NewAnsibleCliOutput()
	inventory := ansibleOutput.CreateInventory()

	addCounter := 0
	for i, line := range stringsutils.SplitLines(cliOutput, true) {
		if i == 0 {
			if !isListHostsOutput(line) {
				return nil, tracederrors.TracedErrorf("%w, Unknown first line to parse as ansible --list-hosts output: '%s'", ErrUnknwnAnsibleCliOutput, line)
			}
			continue
		}

		toAdd := strings.TrimSpace(line)
		if toAdd == "" {
			continue
		}

		_, err = inventory.CreateHostByName(ctx, toAdd)
		if err != nil {
			return nil, err
		}

		addCounter += 1
	}

	if addCounter > 0 {
		logging.LogChangedByCtxf(ctx, "'%d' hosts added to '%s' from parsed CLI output.", addCounter, ansibleOutput.Name())
	} else {
		logging.LogInfoByCtxf(ctx, "No hosts added to '%s' from parsed CLI output.", ansibleOutput.Name())
	}

	return ansibleOutput, nil
}

func ParseCliOutput(ctx context.Context, cliOutput string) (ansibleOutput *AnsibleCliOuput, err error) {
	logging.LogInfoByCtxf(ctx, "Parse ansible output started.")

	if cliOutput == "" {
		return nil, tracederrors.TracedErrorEmptyString("cliOutput")
	}

	if isListHostsOutput(cliOutput) {
		ansibleOutput, err = parseListHostsCliOutput(ctx, cliOutput)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, ErrUnknwnAnsibleCliOutput
	}

	logging.LogInfoByCtxf(ctx, "Parse ansible output finished.")

	return ansibleOutput, nil
}
