package environmentvariables

import (
	"context"
	"os"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// GetEnvValueAsStringOrDefault returns the value of the environment variable specified by envName.
// If the environment variable is not set or empty, the provided defaultValue is returned instead.
func GetEnvValueAsStringOrDefault(ctx context.Context, envName string, defaultValue string) (string, error) {
	if envName == "" {
		return "", tracederrors.TracedErrorEmptyString("envName")
	}

	envValue := os.Getenv(envName)

	if envValue == "" {
		logging.LogInfoByCtxf(ctx, "Environment variable '%s' is not set or empty value. Using default value.", envName)
		return defaultValue, nil
	}

	logging.LogInfoByCtxf(ctx, "Environment variable '%s' was read.", envName)

	return envValue, nil
}

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
