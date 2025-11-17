package xmlutils

import (
	"strings"

	"github.com/go-xmlfmt/xmlfmt"
)

func PrettyFormat(toFormat string) (string, error) {
	formatted := xmlfmt.FormatXML(toFormat, "", "  ")

	formatted = strings.TrimPrefix(formatted, "\n")

	return formatted, nil
}
