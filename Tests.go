package asciichgolangpublic

import (
	"fmt"
	"strings"
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

	if Structs().IsStructOrPointerToStruct(objectToFormat) {
		values, err := Structs().GetFieldValuesAsString(objectToFormat)
		if err != nil {
			return "", TracedErrorf("Unable to get values of '%v' to format as testname", objectToFormat)
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
		return "", TracedError("testname is empty string after evaluation")
	}

	return testname, nil
}

func (t *TestsService) MustFormatAsTestname(objectToFormat interface{}) (testname string) {
	testname, err := t.FormatAsTestname(objectToFormat)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return testname
}
