package asciichgolangpublic

import (
	"fmt"
	"reflect"
)

type TypesServices struct{}

func NewTypesServices() (t *TypesServices) {
	return new(TypesServices)
}

func Types() (t *TypesServices) {
	return NewTypesServices()
}

func (t *TypesServices) GetTypeName(input interface{}) (typeName string, err error) {
	if input == nil {
		return "", TracedErrorNil("input")
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

		asTracedError, err := Errors().GetAsTracedError(inputAsError)
		if err == nil {
			tracedErrorMessage, err := asTracedError.GetErrorMessage()
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

func (t *TypesServices) MustGetTypeName(input interface{}) (typeName string) {
	typeName, err := t.GetTypeName(input)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return typeName
}
