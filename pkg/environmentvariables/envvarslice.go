package environmentvariables

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Addes (or override if already present) the variable "name" with the given "value" in the environment variable slice.
// The slice must be in the format of 'key=value', same as the return value of os.Environ().
func SetEnvVarInStringSlice(envVarSlice []string, name string, value string) ([]string, error) {
	if name == "" {
		return nil, tracederrors.TracedErrorEmptyString("name")
	}

	if len(envVarSlice) == 0 {
		return []string{name + "=" + value}, nil
	}

	ret := make([]string, 0, len(envVarSlice)+1)

	var overwritten bool
	for _, entry := range envVarSlice {
		splitted := strings.SplitN(entry, "=", 2)
		if len(splitted) != 2 {
			return nil, tracederrors.TracedErrorf("Invalid key=value in envVarSlice: '%s'.", entry)
		}

		if splitted[0] == name {
			ret = append(ret, name+"="+value)
			overwritten = true
		} else {
			ret = append(ret, entry)
		}
	}

	if !overwritten {
		ret = append(ret, name+"="+value)
	}

	return ret, nil
}

// Addes (or override if already present) the variables "key of the map" with the corresponding "value value in the map" in the environment variable slice.
// The slice must be in the format of 'key=value', same as the return value of os.Environ().
func SetEnvVarsInStringSlice(envVarSlice []string, toAdd map[string]string) ([]string, error) {
	if toAdd == nil && envVarSlice == nil {
		return []string{}, nil
	}

	if toAdd == nil {
		return slicesutils.GetDeepCopyOfStringsSlice(envVarSlice), nil
	}

	ret := slicesutils.GetDeepCopyOfStringsSlice(envVarSlice)
	var err error
	for k, v := range toAdd {
		ret, err = SetEnvVarInStringSlice(ret, k, v)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}
