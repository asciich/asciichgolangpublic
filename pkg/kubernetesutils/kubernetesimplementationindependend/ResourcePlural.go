package kubernetesimplementationindependend

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
)

func GetObjectPlural(objectName string) (string, error) {
	ret := strings.ToLower(objectName)

	if ret == "gitrepository" {
		return "gitrepositories", nil
	}

	ret = stringsutils.EnsureSuffix(ret, "s")
	return ret, nil
}
