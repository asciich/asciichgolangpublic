package testutils

import (
	"fmt"
	"strings"

	"github.com/asciich/asciichgolangpublic/datatypes/structsutils"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

type TestsService struct{}

func MustFormatAsTestname(objectToFormat interface{}) (testname string) {
	return Tests().MustFormatAsTestname(objectToFormat)
}

func NewTestsService() (t *TestsService) {
	return new(TestsService)
}

func Tests() (tests *TestsService) {
	return NewTestsService()
}

func (t *TestsService) FormatAsTestname(objectToFormat interface{}) (testname string, err error) {
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

func (t *TestsService) MustFormatAsTestname(objectToFormat interface{}) (testname string) {
	testname, err := t.FormatAsTestname(objectToFormat)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return testname
}