package kubernetesimplementationindependend

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
)

func GetResourcePlural(resourceName string) (string, error) {
	ret := strings.ToLower(resourceName)
	ret = stringsutils.EnsureSuffix(ret, "s")
	return ret, nil
}
