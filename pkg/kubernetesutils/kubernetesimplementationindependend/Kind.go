package kubernetesimplementationindependend

import (
	"strings"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

func SanitizeKindName(name string) (string, error) {
	if name == "" {
		return "", tracederrors.TracedErrorEmptyString("name")
	}

	if len(name) <= 1 {
		return "", tracederrors.TracedErrorf("At least two chars required for a kind name but got '%s'.", name)
	}

	return strings.ToUpper(string(name[0])) + name[1:], nil
}
