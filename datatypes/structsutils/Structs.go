package structsutils

import (
	"fmt"
	"reflect"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/datatypes/pointersutils"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

func GetFieldValuesAsString(structToGetFieldsFrom interface{}) (values []string, err error) {
	if !IsStructOrPointerToStruct(structToGetFieldsFrom) {
		return nil, tracederrors.TracedErrorf("'%v' is not as struct", structToGetFieldsFrom)
	}

	var structWithoutPointer reflect.Value
	if IsPointerToStruct(structToGetFieldsFrom) {
		structWithoutPointer = reflect.Indirect(reflect.ValueOf(structToGetFieldsFrom))
	} else {
		structWithoutPointer = reflect.ValueOf(structToGetFieldsFrom)
	}

	numberOfFields := structWithoutPointer.NumField()
	if numberOfFields == 0 {
		return []string{}, nil
	}

	values = []string{}
	for i := 0; i < numberOfFields; i++ {
		values = append(values, fmt.Sprintf("%v", structWithoutPointer.Field(i)))
	}

	return values, nil
}

func IsPointerToStruct(objectToTest interface{}) (isStruct bool) {
	if objectToTest == nil {
		return false
	}

	if !pointersutils.IsPointer(objectToTest) {
		return false
	}

	isStruct = reflect.Indirect(reflect.ValueOf(objectToTest)).Kind() == reflect.Struct
	return isStruct
}

func IsStruct(objectToTest interface{}) (isStruct bool) {
	if objectToTest == nil {
		return false
	}

	isStruct = reflect.ValueOf(objectToTest).Kind() == reflect.Struct
	return isStruct
}

func IsStructOrPointerToStruct(objectToTest interface{}) (isStruct bool) {
	if IsStruct(objectToTest) {
		return true
	}

	if IsPointerToStruct(objectToTest) {
		return true
	}

	return false
}

func MustGetFieldValuesAsString(structToGetFieldsFrom interface{}) (values []string) {
	values, err := GetFieldValuesAsString(structToGetFieldsFrom)
	if err != nil {
		logging.LogFatalf("structs.GetFieldValuesAsString failed: '%v'", err)
	}
	return values
}
