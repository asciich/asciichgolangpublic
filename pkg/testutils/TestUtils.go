package testutils

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/structsutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// GetKindClusterNameForTest returns a KinD cluster name based on the test's package name.
// This ensures that parallel running tests in different packages do not interfere with each other.
func GetKindClusterNameForTest(t testing.TB) string {
	// Get the package name from the test's file path
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return continuousintegration.GetDefaultKindClusterName()
	}

	// Extract package directory from file path
	// Example: /path/to/pkg/kubernetesutils/Namespace_test.go -> kubernetesutils
	parts := strings.Split(file, "/")
	if len(parts) >= 2 {
		packageName := parts[len(parts)-2]
		return continuousintegration.GetKindClusterNameByPackageName(packageName)
	}

	return continuousintegration.GetDefaultKindClusterName()
}

func FormatAsTestname(objectToFormat interface{}) (testname string, err error) {
	testname = ""

	if structsutils.IsStructOrPointerToStruct(objectToFormat) {
		values, err := structsutils.GetFieldValuesAsString(objectToFormat)
		if err != nil {
			return "", tracederrors.TracedErrorf("Unable to get values of '%v' to format as testname", objectToFormat)
		}
		testname = strings.Join(values, "-")
	}

	if len(testname) <= 0 {
		testname = fmt.Sprintf("%v", objectToFormat)
	}

	testname = strings.TrimSpace(testname)
	for _, toReplace := range []string{",", "/", "\\", " ", "\n", "\t", "[", "]", "{", "}", "*"} {
		testname = strings.ReplaceAll(testname, toReplace, "_")
	}

	if testname == "" {
		testname = "emptyTestName"
	}

	if testname == "" {
		return "", tracederrors.TracedError("testname is empty string after evaluation")
	}

	return testname, nil
}

func MustFormatAsTestname(objectToFormat interface{}) (testname string) {
	testname, err := FormatAsTestname(objectToFormat)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return testname
}

func SkipIfRunningInContinuousIntegration(t testing.TB) {
	if continuousintegration.IsRunningInContinuousIntegration() {
		t.Skip("Test not available in continuous integration")
	}
}

func SkipIfRunningInGithub(t testing.TB) {
	if continuousintegration.IsRunningInGithub() {
		t.Skip("Test not available in Github continuous integration")
	}
}
