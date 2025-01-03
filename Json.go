package asciichgolangpublic

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"
)

type JsonService struct {
}

func Json() (jsonService *JsonService) {
	return new(JsonService)
}

func NewJsonService() (jsonService *JsonService) {
	return new(JsonService)
}

func (j *JsonService) DataToJsonBytes(data interface{}) (jsonBytes []byte, err error) {
	jsonBytes, err = json.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, TracedErrorf("Marshal as json failed: '%w', data='%v'", err, data)
	}

	return jsonBytes, nil
}

func (j *JsonService) DataToJsonString(data interface{}) (jsonString string, err error) {
	jsonBytes, err := j.DataToJsonBytes(data)
	if err != nil {
		return "", err
	}

	jsonString = string(jsonBytes)

	return jsonString, nil
}

func (j *JsonService) JsonFileByPathHas(jsonFilePath string, query string, keyToCheck string) (has bool, err error) {
	if jsonFilePath == "" {
		return false, TracedErrorEmptyString("jsonFilePath")
	}

	if query == "" {
		return false, TracedErrorEmptyString("query")
	}

	if keyToCheck == "" {
		return false, TracedErrorEmptyString("keyToCheck")
	}

	jsonFile, err := GetLocalFileByPath(jsonFilePath)
	if err != nil {
		return false, err
	}

	has, err = j.JsonFileHas(jsonFile, query, keyToCheck)
	if err != nil {
		return false, err
	}

	return has, nil
}

func (j *JsonService) JsonFileHas(jsonFile File, query string, keyToCheck string) (has bool, err error) {
	if jsonFile == nil {
		return false, TracedErrorNil("jsonFile")
	}

	if query == "" {
		return false, TracedErrorEmptyString("query")
	}

	if keyToCheck == "" {
		return false, TracedErrorEmptyString("keyToCheck")
	}

	content, err := jsonFile.ReadAsString()
	if err != nil {
		return false, err
	}

	has, err = j.JsonStringHas(content, query, keyToCheck)
	if err != nil {
		return false, err
	}

	return has, nil
}

func (j *JsonService) JsonStringHas(jsonString string, query string, keyToCheck string) (has bool, err error) {
	if query == "" {
		return false, TracedErrorEmptyString("query")
	}

	if keyToCheck == "" {
		return false, TracedErrorEmptyString("keyToCheck")
	}

	has, err = j.RunJqAgainstJsonStringAsBool(jsonString, query+" | has(\""+keyToCheck+"\")")
	if err != nil {
		return false, err
	}

	return has, nil
}

func (j *JsonService) JsonStringToYamlFile(jsonString string, outputFile File, verbose bool) (err error) {
	if outputFile == nil {
		return TracedErrorNil("outputFile")
	}

	jsonData, err := j.ParseJsonString(jsonString)
	if err != nil {
		return err
	}

	err = Yaml().DataToYamlFile(jsonData, outputFile, verbose)
	if err != nil {
		return err
	}

	return nil
}

func (j *JsonService) JsonStringToYamlFileByPath(jsonString string, outputFilePath string, verbose bool) (outputFile File, err error) {
	if outputFilePath == "" {
		return nil, TracedErrorEmptyString("outputFilePath")
	}

	outputFile, err = GetLocalFileByPath(outputFilePath)
	if err != nil {
		return nil, err
	}

	err = j.JsonStringToYamlFile(jsonString, outputFile, verbose)
	if err != nil {
		return nil, err
	}

	return outputFile, nil
}

func (j *JsonService) JsonStringToYamlString(jsonString string) (yamlString string, err error) {
	jsonData, err := j.ParseJsonString(jsonString)
	if err != nil {
		return "", err
	}

	yamlString, err = Yaml().DataToYamlString(jsonData)
	if err != nil {
		return "", err
	}

	return yamlString, nil
}

func (j *JsonService) LoadKeyValueInterfaceDictFromJsonFile(jsonFile File) (keyValues map[string]interface{}, err error) {
	if jsonFile == nil {
		return nil, TracedError("jsonFile is nil")
	}

	jsonContent, err := jsonFile.ReadAsString()
	if err != nil {
		return nil, err
	}

	keyValues, err = j.LoadKeyValueInterfaceDictFromJsonString(jsonContent)
	if err != nil {
		return nil, err
	}

	return keyValues, nil
}

func (j *JsonService) LoadKeyValueInterfaceDictFromJsonString(jsonString string) (keyValues map[string]interface{}, err error) {
	jsonString = strings.TrimSpace(jsonString)

	if jsonString == "" {
		return nil, TracedError("jsonString is empty string")
	}

	keyValues = map[string]interface{}{}
	err = json.Unmarshal([]byte(jsonString), &keyValues)
	if err != nil {
		return nil, TracedError(err.Error())
	}

	return keyValues, nil

}

func (j *JsonService) LoadKeyValueStringDictFromJsonString(jsonString string) (keyValues map[string]string, err error) {
	keyValuesInterface, err := j.LoadKeyValueInterfaceDictFromJsonString(jsonString)
	if err != nil {
		return nil, err
	}

	keyValues = map[string]string{}
	for k, v := range keyValuesInterface {
		valueString := fmt.Sprintf("%v", v)
		keyValues[k] = valueString
	}

	return keyValues, nil
}

func (j *JsonService) MustDataToJsonBytes(data interface{}) (jsonBytes []byte) {
	jsonBytes, err := j.DataToJsonBytes(data)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return jsonBytes
}

func (j *JsonService) MustDataToJsonString(data interface{}) (jsonString string) {
	jsonString, err := j.DataToJsonString(data)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return jsonString
}

func (j *JsonService) MustJsonFileByPathHas(jsonFilePath string, query string, keyToCheck string) (has bool) {
	has, err := j.JsonFileByPathHas(jsonFilePath, query, keyToCheck)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return has
}

func (j *JsonService) MustJsonFileHas(jsonFile File, query string, keyToCheck string) (has bool) {
	has, err := j.JsonFileHas(jsonFile, query, keyToCheck)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return has
}

func (j *JsonService) MustJsonStringHas(jsonString string, query string, keyToCheck string) (has bool) {
	has, err := j.JsonStringHas(jsonString, query, keyToCheck)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return has
}

func (j *JsonService) MustJsonStringToYamlFile(jsonString string, outputFile File, verbose bool) {
	err := j.JsonStringToYamlFile(jsonString, outputFile, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (j *JsonService) MustJsonStringToYamlFileByPath(jsonString string, outputFilePath string, verbose bool) (outputFile File) {
	outputFile, err := j.JsonStringToYamlFileByPath(jsonString, outputFilePath, verbose)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return outputFile
}

func (j *JsonService) MustJsonStringToYamlString(jsonString string) (yamlString string) {
	yamlString, err := j.JsonStringToYamlString(jsonString)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return yamlString
}

func (j *JsonService) MustLoadKeyValueInterfaceDictFromJsonFile(jsonFile File) (keyValues map[string]interface{}) {
	keyValues, err := j.LoadKeyValueInterfaceDictFromJsonFile(jsonFile)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return keyValues
}

func (j *JsonService) MustLoadKeyValueInterfaceDictFromJsonString(jsonString string) (keyValues map[string]interface{}) {
	keyValues, err := j.LoadKeyValueInterfaceDictFromJsonString(jsonString)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return keyValues
}

func (j *JsonService) MustLoadKeyValueStringDictFromJsonString(jsonString string) (keyValues map[string]string) {
	keyValues, err := j.LoadKeyValueStringDictFromJsonString(jsonString)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return keyValues
}

func (j *JsonService) MustParseJsonString(jsonString string) (data interface{}) {
	data, err := j.ParseJsonString(jsonString)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return data
}

func (j *JsonService) MustPrettyFormatJsonString(jsonString string) (formatted string) {
	formatted, err := j.PrettyFormatJsonString(jsonString)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return formatted
}

func (j *JsonService) MustRunJqAgainstJsonFileAsString(jsonFile File, query string) (result string) {
	result, err := j.RunJqAgainstJsonFileAsString(jsonFile, query)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return result
}

func (j *JsonService) MustRunJqAgainstJsonStringAsBool(jsonString string, query string) (result bool) {
	result, err := j.RunJqAgainstJsonStringAsBool(jsonString, query)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return result
}

func (j *JsonService) MustRunJqAgainstJsonStringAsInt(jsonString string, query string) (result int) {
	result, err := j.RunJqAgainstJsonStringAsInt(jsonString, query)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return result
}

func (j *JsonService) MustRunJqAgainstJsonStringAsString(jsonString string, query string) (result string) {
	result, err := j.RunJqAgainstJsonStringAsString(jsonString, query)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return result
}

func (j *JsonService) ParseJsonString(jsonString string) (data interface{}, err error) {
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return nil, TracedError(err.Error())
	}

	return data, err
}

func (j *JsonService) PrettyFormatJsonString(jsonString string) (formatted string, err error) {
	data, err := j.ParseJsonString(jsonString)
	if err != nil {
		return "", err
	}

	formatted, err = j.DataToJsonString(data)
	if err != nil {
		return "", err
	}

	formatted = Strings().EnsureEndsWithExactlyOneLineBreak(formatted)

	return formatted, nil
}

func (j *JsonService) RunJqAgainstJsonFileAsString(jsonFile File, query string) (result string, err error) {
	if jsonFile == nil {
		return "", TracedErrorNil("jsonFile")
	}

	jsonString, err := jsonFile.ReadAsString()
	if err != nil {
		return "", err
	}

	result, err = j.RunJqAgainstJsonStringAsString(jsonString, query)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (j *JsonService) RunJqAgainstJsonStringAsBool(jsonString string, query string) (result bool, err error) {
	if len(jsonString) <= 0 {
		return false, TracedError("jsonString is empty string")
	}

	if len(query) <= 0 {
		return false, TracedError("query is empty string")
	}

	resultString, err := j.RunJqAgainstJsonStringAsString(jsonString, query)
	if err != nil {
		return false, err
	}

	resultString = strings.TrimSpace(resultString)

	result, err = strconv.ParseBool(resultString)
	if err != nil {
		return false, TracedError(err.Error())
	}

	return result, nil
}

func (j *JsonService) RunJqAgainstJsonStringAsInt(jsonString string, query string) (result int, err error) {
	if len(jsonString) <= 0 {
		return -1, TracedError("jsonString is empty string")
	}

	if len(query) <= 0 {
		return -1, TracedError("query is empty string")
	}

	resultString, err := j.RunJqAgainstJsonStringAsString(jsonString, query)
	if err != nil {
		return -1, err
	}

	resultString = strings.TrimSpace(resultString)

	result, err = strconv.Atoi(resultString)
	if err != nil {
		return -1, TracedError(err.Error())
	}

	return result, nil
}

func (j *JsonService) RunJqAgainstJsonStringAsString(jsonString string, query string) (result string, err error) {
	if len(jsonString) <= 0 {
		return "", TracedError("json is empty string")
	}

	if len(query) <= 0 {
		return "", TracedError("query is empty string")
	}

	jsonData, err := j.ParseJsonString(jsonString)
	if err != nil {
		return "", err
	}

	jqQuery, err := gojq.Parse(query)
	if err != nil {
		return "", TracedError(err.Error())
	}
	iter := jqQuery.Run(jsonData)

	result = ""
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return "", TracedError(err.Error())
		}
		switch v := v.(type) {
		case int:
			result += strconv.Itoa(v) + "\n"
		case int64:
			result += strconv.FormatInt(v, 10) + "\n"
		case string:
			result += v + "\n"
		case map[string]interface{}:
			toAdd, err := Json().DataToJsonString(v)
			if err != nil {
				return "", TracedErrorf("Failed to marshal map[string]interface{}")
			}

			result += toAdd + "\n"
		case []interface{}:
			toAdd, err := Json().DataToJsonString(v)
			if err != nil {
				return "", TracedErrorf("Failed to marshal []interface{}")
			}

			result += toAdd + "\n"
		default:
			result += fmt.Sprintf("%#v\n", v)
		}
	}

	result = strings.TrimSpace(result)

	return result, nil
}
