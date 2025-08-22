package gettypename

import (
	"fmt"
	"reflect"
)

func GetTypeName(input interface{}) (typeName string, err error) {
	if input == nil {
		return "", fmt.Errorf("input is nil")
	}

	reflectType := reflect.TypeOf(input)

	typeName = reflectType.Name()

	var inputAsError error
	var ptrPrefix = ""

	inputAsError, ok := input.(error)
	if !ok {
		inputAsError = nil
	}

	if inputAsError == nil {
		inputAsErrorPtr, ok := input.(*error)
		if ok {
			inputAsError = *inputAsErrorPtr
			ptrPrefix = "&"
		}
	}

	if inputAsError != nil {
		errorReflectType := reflect.TypeOf(inputAsError)
		errorTypeName := errorReflectType.Name()

		var message = inputAsError.Error()

		withErrorMessage, ok := err.(interface{ GetErrorMessage() (string, error) })
		if ok {
			tracedErrorMessage, err := withErrorMessage.GetErrorMessage()
			if err == nil {
				message = tracedErrorMessage
			}
		}

		if errorTypeName == "" {
			errorTypeName = "error"
		}

		typeName = ptrPrefix + fmt.Sprintf(
			"%s{message='%s'}",
			errorTypeName,
			message,
		)
	}

	return typeName, nil
}
