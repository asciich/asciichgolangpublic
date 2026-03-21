package environmentvariables

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func GetEnvValueAsString(ctx context.Context, envName string, allowEmpty bool) (string, error) {
	if envName == "" {
		return "", tracederrors.TracedErrorEmptyString("envName")
	}

	envValue := os.Getenv(envName)

	if envValue == "" {
		if allowEmpty {
			logging.LogInfoByCtxf(ctx, "Environment variable '%s' is not set or empty value.", envName)
		} else {
			return "", tracederrors.TracedErrorf("Environment variable '%s' is not set or emtpy value", envName)
		}
	}

	logging.LogInfoByCtxf(ctx, "Environment variable '%s' was read.", envName)

	return envValue, nil
}

func SetEnvVar(ctx context.Context, envName string, value string, logValue bool) error {
	if envName == "" {
		return tracederrors.TracedErrorEmptyString("envName")
	}

	os.Setenv(envName, value)

	if logValue {
		logging.LogInfoByCtxf(ctx, "Set environment variable '%s' to '%s'.", envName, value)
	} else {
		logging.LogInfoByCtxf(ctx, "Set environment variable '%s'.", envName)
	}

	return nil
}
