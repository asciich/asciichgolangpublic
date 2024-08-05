package asciichgolangpublic

import (
	"fmt"
	"reflect"
	"unsafe"
)

type PointersService struct{}

func NewPointersService() (p *PointersService) {
	return new(PointersService)
}

func Pointers() (pointers *PointersService) {
	return new(PointersService)
}

func (p *PointersService) CheckIsPointer(objectToTest interface{}) (err error) {
	isPointer := p.IsPointer(objectToTest)
	if !isPointer {
		return TracedErrorf("'%s' is not a pointer", objectToTest)
	}

	return nil
}

func (p *PointersService) GetMemoryAddressAsHexString(input interface{}) (memoryAddress string, err error) {
	memoryAddressUint64, err := p.GetMemoryAddressAsUInt64(input)
	if err != nil {
		return "", err
	}

	memoryAddress = fmt.Sprintf("0x%x", memoryAddressUint64)

	return memoryAddress, nil
}

func (p *PointersService) GetMemoryAddressAsUInt64(input interface{}) (memoryAddress uint64, err error) {
	if input == nil {
		return 0, TracedErrorNil("input")
	}

	memoryAddressUIntPtr, err := p.GetMemoryAddressAsUIntPtr(input)
	if err != nil {
		return 0, err
	}

	memoryAddress = uint64(memoryAddressUIntPtr)

	return memoryAddress, nil
}

func (p *PointersService) GetMemoryAddressAsUIntPtr(input interface{}) (memoryAddress uintptr, err error) {
	if input == nil {
		return 0, TracedErrorNil("input")
	}

	if !p.IsPointer(input) {
		return 0, TracedErrorf("input is not a pointer: '%v'", input)
	}

	pointer := reflect.ValueOf(input).Pointer()

	var unsafePtr unsafe.Pointer = unsafe.Pointer(pointer)

	memoryAddress = uintptr(unsafePtr)

	return memoryAddress, nil
}

func (p *PointersService) IsPointer(objectToTest interface{}) (isPointer bool) {
	if objectToTest == nil {
		return false
	}

	isPointer = reflect.ValueOf(objectToTest).Kind() == reflect.Ptr
	return isPointer
}

func (p *PointersService) MustCheckIsPointer(objectToTest interface{}) {
	err := p.CheckIsPointer(objectToTest)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (p *PointersService) MustGetMemoryAddressAsHexString(input interface{}) (memoryAddress string) {
	memoryAddress, err := p.GetMemoryAddressAsHexString(input)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return memoryAddress
}

func (p *PointersService) MustGetMemoryAddressAsUInt64(input interface{}) (memoryAddress uint64) {
	memoryAddress, err := p.GetMemoryAddressAsUInt64(input)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return memoryAddress
}

func (p *PointersService) MustGetMemoryAddressAsUIntPtr(input interface{}) (memoryAddress uintptr) {
	memoryAddress, err := p.GetMemoryAddressAsUIntPtr(input)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return memoryAddress
}

func (p *PointersService) MustPointersEqual(ptr1 interface{}, ptr2 interface{}) (addressEqual bool) {
	addressEqual, err := p.PointersEqual(ptr1, ptr2)
	if err != nil {
		LogGoErrorFatal(err)
	}

	return addressEqual
}

func (p *PointersService) PointersEqual(ptr1 interface{}, ptr2 interface{}) (addressEqual bool, err error) {
	if ptr1 == nil && ptr2 == nil {
		return true, nil
	}

	if ptr1 == nil && p.IsPointer(ptr2) {
		return false, nil
	}

	if p.IsPointer(ptr1) && ptr2 == nil {
		return false, nil
	}

	err = p.CheckIsPointer(ptr1)
	if err != nil {
		return false, err
	}

	err = p.CheckIsPointer(ptr2)
	if err != nil {
		return false, err
	}

	addrPtr1, err := p.GetMemoryAddressAsUIntPtr(ptr1)
	if err != nil {
		return false, err
	}

	addrPtr2, err := p.GetMemoryAddressAsUIntPtr(ptr2)
	if err != nil {
		return false, err
	}

	addressEqual = addrPtr1 == addrPtr2

	return addressEqual, nil
}
