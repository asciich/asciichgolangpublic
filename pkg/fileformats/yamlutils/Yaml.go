package yamlutils

import (
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
	gologging "gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v3"
)

var ErrInvalidYaml = errors.New("invalid yaml")

type errTypeInvalidYamlEmptyString struct{}

func (e errTypeInvalidYamlEmptyString) Error() string {
	return "empty string is not a valid yaml"
}

func (e errTypeInvalidYamlEmptyString) Unwrap() error {
	return ErrInvalidYaml
}

var ErrInvalidYamlEmptyString = errTypeInvalidYamlEmptyString{}

func disableYqlibLogging() {
	logger := yqlib.GetLogger()

	backend1 := gologging.NewLogBackend(io.Discard, "", 0)
	backend1Leveled := gologging.AddModuleLevel(backend1)
	logger.SetBackend(backend1Leveled)
}

func init() {
	disableYqlibLogging()
}

func DataToYamlBytes(input interface{}) (yamlBytes []byte, err error) {
	yamlBytes, err = yaml.Marshal(input)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to marshal data to yaml: '%w'", err)
	}

	yamlBytes = append([]byte("---\n"), yamlBytes...)

	return yamlBytes, nil
}

func DataToYamlFile(data interface{}, outputFile files.File, verbose bool) (err error) {
	if outputFile == nil {
		return tracederrors.TracedErrorNil("outputFile")
	}

	yamlString, err := DataToYamlString(data)
	if err != nil {
		return err
	}

	err = outputFile.WriteString(yamlString, verbose)
	if err != nil {
		return err
	}

	return nil
}

func DataToYamlString(input interface{}) (yamlString string, err error) {
	yamlBytes, err := DataToYamlBytes(input)
	if err != nil {
		return "", err
	}

	yamlString = string(yamlBytes)

	return yamlString, err
}

func MustDataToYamlBytes(input interface{}) (yamlBytes []byte) {
	yamlBytes, err := DataToYamlBytes(input)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return yamlBytes
}

func MustDataToYamlFile(jsonData interface{}, outputFile files.File, verbose bool) {
	err := DataToYamlFile(jsonData, outputFile, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func MustDataToYamlString(input interface{}) (yamlString string) {
	yamlString, err := DataToYamlString(input)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return yamlString
}

func MustRunYqQueryAginstYamlStringAsString(yamlString string, query string) (result string) {
	result, err := RunYqQueryAginstYamlStringAsString(yamlString, query)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return result
}

func RunYqQueryAginstYamlStringAsString(yamlString string, query string) (result string, err error) {
	encoder := yqlib.NewYamlEncoder(yqlib.NewDefaultYamlPreferences())
	decoder := yqlib.NewYamlDecoder(yqlib.NewDefaultYamlPreferences())

	result, err = yqlib.NewStringEvaluator().EvaluateAll(
		query,
		yamlString,
		encoder,
		decoder,
	)

	if err != nil {
		return "", tracederrors.TracedErrorf("Failed to evaluate query '%s': %w", query, err)
	}

	result = strings.TrimSuffix(result, "\n")

	return result, nil
}

func SplitMultiYaml(yamlString string) (splitted []string) {
	var toAdd string

	for _, line := range stringsutils.SplitLines(yamlString, true) {
		trimmedLine := stringsutils.TrimSpacesRight(line)
		if trimmedLine == "---" {
			if toAdd == "" {
				continue
			}

			toAdd = "---\n" + stringsutils.EnsureEndsWithExactlyOneLineBreak(toAdd)
			splitted = append(splitted, toAdd)
			toAdd = ""
			continue
		}

		if toAdd == "" {
			toAdd = trimmedLine
		} else {
			toAdd += "\n" + trimmedLine
		}
	}

	if toAdd != "" {
		toAdd = "---\n" + stringsutils.EnsureEndsWithExactlyOneLineBreak(toAdd)
		splitted = append(splitted, toAdd)
	}

	return splitted
}

var regexDocumentStartRemoval = regexp.MustCompile("^---.*\n")

func MergeMultiYaml(yamls []string) (merged string, err error) {
	if yamls == nil {
		return "", tracederrors.TracedError("yamls")
	}

	for _, y := range yamls {
		y = stringsutils.TrimSpacesRight(y)
		y = stringsutils.TrimAllLeadingNewLines(y)

		if strings.HasPrefix(y, "---") {
			y = string(regexDocumentStartRemoval.ReplaceAll([]byte(y), []byte("")))
			y = stringsutils.TrimAllLeadingNewLines(y)
		}

		if y == "" {
			continue
		}

		merged += "---\n" + stringsutils.EnsureEndsWithExactlyOneLineBreak(y)
	}

	return merged, nil
}

func MustMergeMultiYaml(yamls []string) (merged string) {
	merged, err := MergeMultiYaml(yamls)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	merged = stringsutils.EnsureEndsWithExactlyOneLineBreak(merged)

	return merged
}

func MustLoadGeneric(input string) (data interface{}) {
	data, err := LoadGeneric(input)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return data
}

func LoadGeneric(input string) (data interface{}, err error) {
	data = new(interface{})

	err = yaml.Unmarshal([]byte(input), data)
	if err != nil {
		return nil, tracederrors.TracedErrorf("%w: %w", ErrInvalidYaml, err)
	}

	return data, nil
}

func Validate(toValidate string) (err error) {
	trimmed := strings.TrimSpace(toValidate)

	if trimmed == "" {
		return ErrInvalidYamlEmptyString
	}

	_, err = LoadGeneric(toValidate)
	if err != nil {
		return err
	}

	return nil
}

func EnsureDocumentStart(input string) (output string) {
	trimmed := stringsutils.TrimSpacesLeft(input)

	if trimmed == "---" {
		return "---\n"
	}

	if strings.HasPrefix(trimmed, "---\n") {
		return trimmed
	}

	if strings.HasPrefix(trimmed, "#") {
		withoutComments := stringsutils.RemoveComments(input)
		if withoutComments == "---" {
			return trimmed + "\n"
		}

		if strings.HasPrefix(trimmed, "---\n") {
			return trimmed
		}
	}

	return "---\n" + trimmed
}

func EnsureDocumentStartAndEnd(input string) (output string) {
	output = EnsureDocumentStart(input)
	return stringsutils.EnsureEndsWithExactlyOneLineBreak(output)
}
