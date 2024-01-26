package asciichgolangpublic

import "reflect"

type PointersService struct{}

func NewPointersService() (p *PointersService) {
	return new(PointersService)
}

func Pointers() (pointers *PointersService) {
	return new(PointersService)
}

func (p *PointersService) IsPointer(objectToTest interface{}) (isPointer bool) {
	if objectToTest == nil {
		return false
	}

	isPointer = reflect.ValueOf(objectToTest).Kind() == reflect.Ptr
	return isPointer
}
