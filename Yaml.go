package asciichgolangpublic

import "gopkg.in/yaml.v3"

type YamlService struct{}

func NewYamlService() (y *YamlService) {
	return new(YamlService)
}

func Yaml() (yaml *YamlService) {
	return new(YamlService)
}

func (y *YamlService) DataToYamlBytes(input interface{}) (yamlBytes []byte, err error) {
	yamlBytes, err = yaml.Marshal(input)
	if err != nil {
		return nil, TracedErrorf("Failed to marshal data to yaml: '%w'", err)
	}

	yamlBytes = append([]byte("---\n"), yamlBytes...)

	return yamlBytes, nil
}

func (y *YamlService) DataToYamlFile(jsonData interface{}, outputFile File, verbose bool) (err error) {
	if outputFile == nil {
		return TracedErrorNil("outputFile")
	}

	yamlString, err := y.DataToYamlString(jsonData)
	if err != nil {
		return err
	}

	err = outputFile.WriteString(yamlString, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (y *YamlService) DataToYamlString(input interface{}) (yamlString string, err error) {
	yamlBytes, err := y.DataToYamlBytes(input)
	if err != nil {
		return "", err
	}

	yamlString = string(yamlBytes)

	return yamlString, err
}

func (y *YamlService) MustDataToYamlBytes(input interface{}) (yamlBytes []byte) {
	yamlBytes, err := y.DataToYamlBytes(input)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return yamlBytes
}

func (y *YamlService) MustDataToYamlFile(jsonData interface{}, outputFile File, verbose bool) {
	err := y.DataToYamlFile(jsonData, outputFile, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (y *YamlService) MustDataToYamlString(input interface{}) (yamlString string) {
	yamlString, err := y.DataToYamlString(input)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return yamlString
}
