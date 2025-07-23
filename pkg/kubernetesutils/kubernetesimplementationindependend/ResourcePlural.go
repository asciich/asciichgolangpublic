package kubernetesimplementationindependend

import (
	"strings"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/datatypes/stringsutils"
)

func GetObjectPlural(objectName string) (string, error) {
	ret := strings.ToLower(objectName)

	if ret == "gitrepository" {
		return "gitrepositories", nil
	}

	ret = stringsutils.EnsureSuffix(ret, "s")
	return ret, nil
}
