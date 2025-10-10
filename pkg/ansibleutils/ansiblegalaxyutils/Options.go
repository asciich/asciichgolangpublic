package ansiblegalaxyutils

import (
	"errors"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/structsutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetAnsiblePath(options any) (string, error) {
	if options == nil {
		return "", tracederrors.TracedErrorNil("options")
	}

	path, err := structsutils.GetFieldValueAsString(options, "AnsibleVirtualenvPath")
	if err != nil {
		if errors.Is(err, structsutils.ErrStructHasNoField) {
			return "ansible", nil
		}
		return "", err
	}

	if path == "" {
		return "ansible", nil
	}

	return filepath.Join(path, "bin", "ansible"), nil
}

func GetAnsibleGalaxyPath(options any) (string, error) {
	path, err := GetAnsiblePath(options)
	if err != nil {
		return "", err
	}

	return path + "-galaxy", nil
}
