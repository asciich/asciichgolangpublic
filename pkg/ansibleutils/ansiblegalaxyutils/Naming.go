package ansiblegalaxyutils

import (
	"regexp"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

var regexCollectionName = regexp.MustCompile(`^[a-z][a-z_]*[a-z]$`)

func IsValidCollectionName(name string) bool {
	return regexCollectionName.Match([]byte(name))
}

func CheckValidCollectionName(name string) error {
	if IsValidCollectionName(name) {
		return nil
	}

	return tracederrors.TracedErrorf("Invalid collection name: '%s'. Both collection name and namespace may only consist of lowercase letters and an underscore, have at least two chars and does start with a letter.", name)
}
