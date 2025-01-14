package asciichgolangpublic

import (
	"fmt"
	"reflect"

	"github.com/asciich/asciichgolangpublic/datatypes/pointers"
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type StructsService struct{}

func NewStructsService() (s *StructsService) {
	return new(StructsService)
}

func Structs() (structs *StructsService) {
	return new(StructsService)
}

func (s *StructsService) GetFieldValuesAsString(structToGetFieldsFrom interface{}) (values []string, err error) {
	if !s.IsStructOrPointerToStruct(structToGetFieldsFrom) {
		return nil, errors.TracedErrorf("'%v' is not as struct", structToGetFieldsFrom)
	}

	var structWithoutPointer reflect.Value
	if s.IsPointerToStruct(structToGetFieldsFrom) {
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

func (s *StructsService) IsPointerToStruct(objectToTest interface{}) (isStruct bool) {
	if objectToTest == nil {
		return false
	}

	if !pointers.IsPointer(objectToTest) {
		return false
	}

	isStruct = reflect.Indirect(reflect.ValueOf(objectToTest)).Kind() == reflect.Struct
	return isStruct
}

func (s *StructsService) IsStruct(objectToTest interface{}) (isStruct bool) {
	if objectToTest == nil {
		return false
	}

	isStruct = reflect.ValueOf(objectToTest).Kind() == reflect.Struct
	return isStruct
}

func (s *StructsService) IsStructOrPointerToStruct(objectToTest interface{}) (isStruct bool) {
	if s.IsStruct(objectToTest) {
		return true
	}

	if s.IsPointerToStruct(objectToTest) {
		return true
	}

	return false
}

func (s *StructsService) MustGetFieldValuesAsString(structToGetFieldsFrom interface{}) (values []string) {
	values, err := s.GetFieldValuesAsString(structToGetFieldsFrom)
	if err != nil {
		logging.LogFatalf("structs.GetFieldValuesAsString failed: '%v'", err)
	}
	return values
}
