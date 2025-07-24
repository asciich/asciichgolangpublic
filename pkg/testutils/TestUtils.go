package testutils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/asciich/asciichgolangpublic/datatypes/structsutils"
	"github.com/asciich/asciichgolangpublic/pkg/continuousintegration"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

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
