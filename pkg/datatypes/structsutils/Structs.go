package structsutils

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sort"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/pointersutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

var ErrStructHasNoField = errors.New("struct has no matching field")

func GetFieldValueAsString(structToGetValueFrom any, fieldName string) (string, error) {
	if structToGetValueFrom == nil {
		return "", tracederrors.TracedErrorNil("structToGetValueFrom")
	}

	if fieldName == "" {
		return "", tracederrors.TracedErrorEmptyString("fieldName")
	}

	if !IsStructOrPointerToStruct(structToGetValueFrom) {
		return "", tracederrors.TracedErrorf("'%v' is not as struct", structToGetValueFrom)
	}

	hasField, err := HasField(structToGetValueFrom, fieldName)
	if err != nil {
		return "", err
	}

	if !hasField {
		return "", tracederrors.TracedErrorf("%w: '%s'. Struct is '%s'", ErrStructHasNoField, fieldName, structToGetValueFrom)
	}

	var structWithoutPointer reflect.Value
	if IsPointerToStruct(structToGetValueFrom) {
		structWithoutPointer = reflect.Indirect(reflect.ValueOf(structToGetValueFrom))
	} else {
		structWithoutPointer = reflect.ValueOf(structToGetValueFrom)
	}

	field := structWithoutPointer.FieldByName(fieldName)
	return field.String(), nil
}

func GetFieldValuesAsString(structToGetFieldsFrom any) (values []string, err error) {
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

func ListFieldNames(structToList any) ([]string, error) {
	if !IsStructOrPointerToStruct(structToList) {
		return nil, tracederrors.TracedErrorf("'%v' is not as struct", structToList)
	}

	var structWithoutPointer = reflect.TypeOf(structToList)
	if structWithoutPointer.Kind() == reflect.Ptr {
		structWithoutPointer = structWithoutPointer.Elem()
	}

	numberOfFields := structWithoutPointer.NumField()
	if numberOfFields == 0 {
		return []string{}, nil
	}

	fields := make([]string, 0, numberOfFields)
	for i := 0; i < numberOfFields; i++ {
		fields = append(fields, structWithoutPointer.Field(i).Name)
	}

	sort.Strings(fields)

	return fields, nil
}

func HasField(structToTest any, fieldName string) (bool, error) {
	if structToTest == nil {
		return false, tracederrors.TracedErrorNil("structToTest")
	}

	if fieldName == "" {
		return false, tracederrors.TracedErrorEmptyString("fieldName")
	}

	fields, err := ListFieldNames(structToTest)
	if err != nil {
		return false, err
	}

	return slices.Contains(fields, fieldName), nil
}
