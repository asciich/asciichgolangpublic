package cobrautils

import (
	"bytes"

	"github.com/gavv/cobradoc"
	"github.com/spf13/cobra"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GenerateMarkdownDocumentation(name string, cmd *cobra.Command) (string, error) {
	if cmd == nil {
		return "", tracederrors.TracedErrorNil("cmd")
	}

	if name == "" {
		return "", tracederrors.TracedErrorEmptyString("name")
	}

	var buf bytes.Buffer

	options := cobradoc.Options{
		Name:             name,
		Header:           "Command Line Reference",
		ShortDescription: "Detailed reference for all available commands.",
	}

	err := cobradoc.WriteDocument(&buf, cmd, cobradoc.Markdown, options)
	if err != nil {
		return "", tracederrors.TracedErrorf("failed to generate documentation: %w", err)
	}

	return buf.String(), nil
}
