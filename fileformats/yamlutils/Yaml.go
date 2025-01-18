package yamlutils

import (
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
	"gopkg.in/yaml.v3"
)


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
