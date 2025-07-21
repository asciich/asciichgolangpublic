package jsonutils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/itchyny/gojq"

	"github.com/asciich/asciichgolangpublic/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/pkg/fileformats/yamlutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func DataToJsonBytes(data interface{}) (jsonBytes []byte, err error) {
	jsonBytes, err = json.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, tracederrors.TracedErrorf("Marshal as json failed: '%w', data='%v'", err, data)
	}

	return jsonBytes, nil
}

func DataToJsonString(data interface{}) (jsonString string, err error) {
	jsonBytes, err := DataToJsonBytes(data)
	if err != nil {
		return "", err
	}

	jsonString = string(jsonBytes)

	return jsonString, nil
}

func JsonFileByPathHas(jsonFilePath string, query string, keyToCheck string) (has bool, err error) {
	if jsonFilePath == "" {
		return false, tracederrors.TracedErrorEmptyString("jsonFilePath")
	}

	if query == "" {
		return false, tracederrors.TracedErrorEmptyString("query")
	}

	if keyToCheck == "" {
		return false, tracederrors.TracedErrorEmptyString("keyToCheck")
	}

	jsonFile, err := files.GetLocalFileByPath(jsonFilePath)
	if err != nil {
		return false, err
	}

	has, err = JsonFileHas(jsonFile, query, keyToCheck)
	if err != nil {
		return false, err
	}

	return has, nil
}

func JsonFileHas(jsonFile files.File, query string, keyToCheck string) (has bool, err error) {
	if jsonFile == nil {
		return false, tracederrors.TracedErrorNil("jsonFile")
	}

	if query == "" {
		return false, tracederrors.TracedErrorEmptyString("query")
	}

	if keyToCheck == "" {
		return false, tracederrors.TracedErrorEmptyString("keyToCheck")
	}

	content, err := jsonFile.ReadAsString()
	if err != nil {
		return false, err
	}

	has, err = JsonStringHas(content, query, keyToCheck)
	if err != nil {
		return false, err
	}

	return has, nil
}

func JsonStringHas(jsonString string, query string, keyToCheck string) (has bool, err error) {
	if query == "" {
		return false, tracederrors.TracedErrorEmptyString("query")
	}

	if keyToCheck == "" {
		return false, tracederrors.TracedErrorEmptyString("keyToCheck")
	}

	has, err = RunJqAgainstJsonStringAsBool(jsonString, query+" | has(\""+keyToCheck+"\")")
	if err != nil {
		return false, err
	}

	return has, nil
}

func JsonStringToYamlFile(jsonString string, outputFile files.File, verbose bool) (err error) {
	if outputFile == nil {
		return tracederrors.TracedErrorNil("outputFile")
	}

	jsonData, err := ParseJsonString(jsonString)
	if err != nil {
		return err
	}

	err = yamlutils.DataToYamlFile(jsonData, outputFile, verbose)
	if err != nil {
		return err
	}

	return nil
}

func JsonStringToYamlFileByPath(jsonString string, outputFilePath string, verbose bool) (outputFile files.File, err error) {
	if outputFilePath == "" {
		return nil, tracederrors.TracedErrorEmptyString("outputFilePath")
	}

	outputFile, err = files.GetLocalFileByPath(outputFilePath)
	if err != nil {
		return nil, err
	}

	err = JsonStringToYamlFile(jsonString, outputFile, verbose)
	if err != nil {
		return nil, err
	}

	return outputFile, nil
}

func JsonStringToYamlString(jsonString string) (yamlString string, err error) {
	jsonData, err := ParseJsonString(jsonString)
	if err != nil {
		return "", err
	}

	yamlString, err = yamlutils.DataToYamlString(jsonData)
	if err != nil {
		return "", err
	}

	return yamlString, nil
}

func LoadKeyValueInterfaceDictFromJsonFile(jsonFile files.File) (keyValues map[string]interface{}, err error) {
	if jsonFile == nil {
		return nil, tracederrors.TracedError("jsonFile is nil")
	}

	jsonContent, err := jsonFile.ReadAsString()
	if err != nil {
		return nil, err
	}

	keyValues, err = LoadKeyValueInterfaceDictFromJsonString(jsonContent)
	if err != nil {
		return nil, err
	}

	return keyValues, nil
}

func LoadKeyValueInterfaceDictFromJsonString(jsonString string) (keyValues map[string]interface{}, err error) {
	jsonString = strings.TrimSpace(jsonString)

	if jsonString == "" {
		return nil, tracederrors.TracedError("jsonString is empty string")
	}

	keyValues = map[string]interface{}{}
	err = json.Unmarshal([]byte(jsonString), &keyValues)
	if err != nil {
		return nil, tracederrors.TracedError(err.Error())
	}

	return keyValues, nil

}

func LoadKeyValueStringDictFromJsonString(jsonString string) (keyValues map[string]string, err error) {
	keyValuesInterface, err := LoadKeyValueInterfaceDictFromJsonString(jsonString)
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

func MustDataToJsonBytes(data interface{}) (jsonBytes []byte) {
	jsonBytes, err := DataToJsonBytes(data)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return jsonBytes
}

func MustDataToJsonString(data interface{}) (jsonString string) {
	jsonString, err := DataToJsonString(data)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return jsonString
}

func MustJsonFileByPathHas(jsonFilePath string, query string, keyToCheck string) (has bool) {
	has, err := JsonFileByPathHas(jsonFilePath, query, keyToCheck)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return has
}

func MustJsonFileHas(jsonFile files.File, query string, keyToCheck string) (has bool) {
	has, err := JsonFileHas(jsonFile, query, keyToCheck)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return has
}

func MustJsonStringHas(jsonString string, query string, keyToCheck string) (has bool) {
	has, err := JsonStringHas(jsonString, query, keyToCheck)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return has
}

func MustJsonStringToYamlFile(jsonString string, outputFile files.File, verbose bool) {
	err := JsonStringToYamlFile(jsonString, outputFile, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func MustJsonStringToYamlFileByPath(jsonString string, outputFilePath string, verbose bool) (outputFile files.File) {
	outputFile, err := JsonStringToYamlFileByPath(jsonString, outputFilePath, verbose)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return outputFile
}

func MustJsonStringToYamlString(jsonString string) (yamlString string) {
	yamlString, err := JsonStringToYamlString(jsonString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return yamlString
}

func MustLoadKeyValueInterfaceDictFromJsonFile(jsonFile files.File) (keyValues map[string]interface{}) {
	keyValues, err := LoadKeyValueInterfaceDictFromJsonFile(jsonFile)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return keyValues
}

func MustLoadKeyValueInterfaceDictFromJsonString(jsonString string) (keyValues map[string]interface{}) {
	keyValues, err := LoadKeyValueInterfaceDictFromJsonString(jsonString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return keyValues
}

func MustLoadKeyValueStringDictFromJsonString(jsonString string) (keyValues map[string]string) {
	keyValues, err := LoadKeyValueStringDictFromJsonString(jsonString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return keyValues
}

func MustParseJsonString(jsonString string) (data interface{}) {
	data, err := ParseJsonString(jsonString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return data
}

func MustPrettyFormatJsonString(jsonString string) (formatted string) {
	formatted, err := PrettyFormatJsonString(jsonString)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return formatted
}

func MustRunJqAgainstJsonFileAsString(jsonFile files.File, query string) (result string) {
	result, err := RunJqAgainstJsonFileAsString(jsonFile, query)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return result
}

func MustRunJqAgainstJsonStringAsBool(jsonString string, query string) (result bool) {
	result, err := RunJqAgainstJsonStringAsBool(jsonString, query)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return result
}

func MustRunJqAgainstJsonStringAsInt(jsonString string, query string) (result int) {
	result, err := RunJqAgainstJsonStringAsInt(jsonString, query)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return result
}

func MustRunJqAgainstJsonStringAsString(jsonString string, query string) (result string) {
	result, err := RunJqAgainstJsonStringAsString(jsonString, query)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return result
}

func ParseJsonString(jsonString string) (data interface{}, err error) {
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return nil, tracederrors.TracedError(err.Error())
	}

	return data, err
}

func PrettyFormatJsonString(jsonString string) (formatted string, err error) {
	data, err := ParseJsonString(jsonString)
	if err != nil {
		return "", err
	}

	formatted, err = DataToJsonString(data)
	if err != nil {
		return "", err
	}

	formatted = stringsutils.EnsureEndsWithExactlyOneLineBreak(formatted)

	return formatted, nil
}

func RunJqAgainstJsonFileAsString(jsonFile files.File, query string) (result string, err error) {
	if jsonFile == nil {
		return "", tracederrors.TracedErrorNil("jsonFile")
	}

	jsonString, err := jsonFile.ReadAsString()
	if err != nil {
		return "", err
	}

	result, err = RunJqAgainstJsonStringAsString(jsonString, query)
	if err != nil {
		return "", err
	}

	return result, nil
}

func RunJqAgainstJsonStringAsBool(jsonString string, query string) (result bool, err error) {
	if len(jsonString) <= 0 {
		return false, tracederrors.TracedError("jsonString is empty string")
	}

	if len(query) <= 0 {
		return false, tracederrors.TracedError("query is empty string")
	}

	resultString, err := RunJqAgainstJsonStringAsString(jsonString, query)
	if err != nil {
		return false, err
	}

	resultString = strings.TrimSpace(resultString)

	result, err = strconv.ParseBool(resultString)
	if err != nil {
		return false, tracederrors.TracedError(err.Error())
	}

	return result, nil
}

func RunJqAgainstJsonStringAsInt(jsonString string, query string) (result int, err error) {
	if len(jsonString) <= 0 {
		return -1, tracederrors.TracedError("jsonString is empty string")
	}

	if len(query) <= 0 {
		return -1, tracederrors.TracedError("query is empty string")
	}

	resultString, err := RunJqAgainstJsonStringAsString(jsonString, query)
	if err != nil {
		return -1, err
	}

	resultString = strings.TrimSpace(resultString)

	result, err = strconv.Atoi(resultString)
	if err != nil {
		return -1, tracederrors.TracedError(err.Error())
	}

	return result, nil
}

func RunJqAgainstJsonStringAsString(jsonString string, query string) (result string, err error) {
	if len(jsonString) <= 0 {
		return "", tracederrors.TracedError("json is empty string")
	}

	if len(query) <= 0 {
		return "", tracederrors.TracedError("query is empty string")
	}

	jsonData, err := ParseJsonString(jsonString)
	if err != nil {
		return "", err
	}

	jqQuery, err := gojq.Parse(query)
	if err != nil {
		return "", tracederrors.TracedError(err.Error())
	}
	iter := jqQuery.Run(jsonData)

	result = ""
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return "", tracederrors.TracedError(err.Error())
		}
		switch v := v.(type) {
		case int:
			result += strconv.Itoa(v) + "\n"
		case int64:
			result += strconv.FormatInt(v, 10) + "\n"
		case string:
			result += v + "\n"
		case map[string]interface{}:
			toAdd, err := DataToJsonString(v)
			if err != nil {
				return "", tracederrors.TracedErrorf("Failed to marshal map[string]interface{}")
			}

			result += toAdd + "\n"
		case []interface{}:
			toAdd, err := DataToJsonString(v)
			if err != nil {
				return "", tracederrors.TracedErrorf("Failed to marshal []interface{}")
			}

			result += toAdd + "\n"
		default:
			result += fmt.Sprintf("%#v\n", v)
		}
	}

	result = strings.TrimSpace(result)

	return result, nil
}
