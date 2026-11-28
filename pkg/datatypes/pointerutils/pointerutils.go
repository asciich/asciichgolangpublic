package pointerutils

import (
	"fmt"
	"reflect"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func CheckIsPointer(objectToTest interface{}) (err error) {
	isPointer := IsPointer(objectToTest)
	if !isPointer {
		// %v instead of %s: %s prints garbage like %!s(<nil>) or %!s(int=5) for non-strings
		return tracederrors.TracedErrorf("'%v' is not a pointer", objectToTest)
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

	// Pointer() already returns a uintptr, no unsafe.Pointer conversion needed.
	memoryAddress = reflect.ValueOf(input).Pointer()

	return memoryAddress, nil
}

func IsPointer(objectToTest interface{}) (isPointer bool) {
	if objectToTest == nil {
		return false
	}

	// reflect.Ptr is the old name; reflect.Pointer exists since Go 1.18
	isPointer = reflect.ValueOf(objectToTest).Kind() == reflect.Pointer
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

func ToInt64Pointer(data int64) *int64 {
	return &data
}
