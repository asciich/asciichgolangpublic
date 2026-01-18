package nativedocker

import "github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"

func IsRemovalAlreadyInProgressError(err error) bool {
	if err == nil {
		return false
	}

	return stringsutils.ContainsAllIgnoreCase(err.Error(), []string{"removal of", "is already in progress"})
}
