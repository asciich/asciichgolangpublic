package pointerutils

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CheckIsPointer(objectToTest interface{}) (err error) {
	isPointer := IsPointer(objectToTest)
	if !isPointer {
		return tracederrors.TracedErrorf("'%s' is not a pointer", objectToTest)
	}

	return nil
}

func GetMemoryAddressAsHexString(input interface{}) (memoryAddress string, err error) {
	memoryAddressUint64, err := GetMemoryAddressAsUInt64(input)
	if err != nil {
		return "", err
	}

	memoryAddress = fmt.Sprintf("0x%x", memoryAddressUint64)

	return memoryAddress, nil
}

func GetMemoryAddressAsUInt64(input interface{}) (memoryAddress uint64, err error) {
	if input == nil {
		return 0, tracederrors.TracedErrorNil("input")
	}

	memoryAddressUIntPtr, err := GetMemoryAddressAsUIntPtr(input)
	if err != nil {
		return 0, err
	}

	memoryAddress = uint64(memoryAddressUIntPtr)

	return memoryAddress, nil
}

func GetMemoryAddressAsUIntPtr(input interface{}) (memoryAddress uintptr, err error) {
	if input == nil {
		return 0, tracederrors.TracedErrorNil("input")
	}

	if !IsPointer(input) {
		return 0, tracederrors.TracedErrorf("input is not a pointer: '%v'", input)
	}

	pointer := reflect.ValueOf(input).Pointer()

	var unsafePtr unsafe.Pointer = unsafe.Pointer(pointer)

	memoryAddress = uintptr(unsafePtr)

	return memoryAddress, nil
}

func IsPointer(objectToTest interface{}) (isPointer bool) {
	if objectToTest == nil {
		return false
	}

	isPointer = reflect.ValueOf(objectToTest).Kind() == reflect.Ptr
	return isPointer
}

func PointersEqual(ptr1 interface{}, ptr2 interface{}) (addressEqual bool, err error) {
	if ptr1 == nil && ptr2 == nil {
		return true, nil
	}

	if ptr1 == nil && IsPointer(ptr2) {
		return false, nil
	}

	if IsPointer(ptr1) && ptr2 == nil {
		return false, nil
	}

	err = CheckIsPointer(ptr1)
	if err != nil {
		return false, err
	}

	err = CheckIsPointer(ptr2)
	if err != nil {
		return false, err
	}

	addrPtr1, err := GetMemoryAddressAsUIntPtr(ptr1)
	if err != nil {
		return false, err
	}

	addrPtr2, err := GetMemoryAddressAsUIntPtr(ptr2)
	if err != nil {
		return false, err
	}

	addressEqual = addrPtr1 == addrPtr2

	return addressEqual, nil
}

func ToInt64Pointer(data int64) (*int64) {
	return &data
}