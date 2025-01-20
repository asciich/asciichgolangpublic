package yamlutils

import (
	"io"
	"strings"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
	gologging "gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v3"
)

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

func DataToYamlFile(jsonData interface{}, outputFile files.File, verbose bool) (err error) {
	if outputFile == nil {
		return tracederrors.TracedErrorNil("outputFile")
	}

	yamlString, err := DataToYamlString(jsonData)
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
