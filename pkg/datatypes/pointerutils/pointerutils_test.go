package pointerutils_test

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/pointerutils"
)

var constIntForTesting int = 10
var testVal1 = 5
var testVal2 = 5

type testStruct struct {
	Value int
}

func TestIsPointer(t *testing.T) {
	var intPtr *int = &constIntForTesting
	var nilIntPtr *int

	tests := []struct {
		name              string
		pointerToCheck    interface{}
		expectedIsPointer bool
	}{
		{"int", 5, false},
		{"string", "hello", false},
		{"struct value", testStruct{Value: 1}, false},
		{"slice", []int{1, 2}, false},
		{"map", map[string]int{"a": 1}, false},
		{"func", func() {}, false},
		{"untyped nil", nil, false},
		{"pointer to int", &constIntForTesting, true},
		{"pointer to struct", &testStruct{Value: 1}, true},
		{"pointer to pointer", &intPtr, true},
		{"typed nil pointer", nilIntPtr, true},
		// unsafe.Pointer has reflect.Kind UnsafePointer, not Pointer.
		{"unsafe.Pointer", unsafe.Pointer(&constIntForTesting), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedIsPointer, pointerutils.IsPointer(tt.pointerToCheck))
		})
	}
}

func TestCheckIsPointer(t *testing.T) {
	t.Run("not a pointer", func(t *testing.T) {
		err := pointerutils.CheckIsPointer(5)
		require.Error(t, err)
		// Regression test: with '%s' this would read "%!s(int=5)".
		require.ErrorContains(t, err, "'5' is not a pointer")
	})

	t.Run("struct value", func(t *testing.T) {
		err := pointerutils.CheckIsPointer(testStruct{Value: 1})
		require.Error(t, err)
		require.ErrorContains(t, err, "is not a pointer")
	})

	t.Run("is a pointer", func(t *testing.T) {
		err := pointerutils.CheckIsPointer(&constIntForTesting)
		require.NoError(t, err)
	})

	t.Run("typed nil pointer is still a pointer", func(t *testing.T) {
		var nilIntPtr *int
		err := pointerutils.CheckIsPointer(nilIntPtr)
		require.NoError(t, err)
	})

	t.Run("nil", func(t *testing.T) {
		err := pointerutils.CheckIsPointer(nil)
		require.Error(t, err)
		// Regression test: with '%s' this would read "%!s(<nil>)".
		require.ErrorContains(t, err, "'<nil>' is not a pointer")
	})
}

func TestGetMemoryAddressAsUIntPtr(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		addr, err := pointerutils.GetMemoryAddressAsUIntPtr(nil)
		require.Error(t, err)
		require.EqualValues(t, 0, addr)
	})

	t.Run("not a pointer", func(t *testing.T) {
		addr, err := pointerutils.GetMemoryAddressAsUIntPtr(5)
		require.Error(t, err)
		require.ErrorContains(t, err, "input is not a pointer")
		require.EqualValues(t, 0, addr)
	})

	t.Run("valid pointer", func(t *testing.T) {
		addr, err := pointerutils.GetMemoryAddressAsUIntPtr(&constIntForTesting)
		require.NoError(t, err)
		require.NotZero(t, addr)
	})

	t.Run("matches the real address", func(t *testing.T) {
		// Pointer -> uintptr is a legal conversion and serves as an independent oracle.
		expected := uintptr(unsafe.Pointer(&testVal1))

		addr, err := pointerutils.GetMemoryAddressAsUIntPtr(&testVal1)
		require.NoError(t, err)
		require.Equal(t, expected, addr)
	})

	t.Run("same pointer gives same address", func(t *testing.T) {
		addr1, err := pointerutils.GetMemoryAddressAsUIntPtr(&testVal1)
		require.NoError(t, err)

		addr2, err := pointerutils.GetMemoryAddressAsUIntPtr(&testVal1)
		require.NoError(t, err)

		require.Equal(t, addr1, addr2)
	})

	t.Run("different variables give different addresses", func(t *testing.T) {
		addr1, err := pointerutils.GetMemoryAddressAsUIntPtr(&testVal1)
		require.NoError(t, err)

		addr2, err := pointerutils.GetMemoryAddressAsUIntPtr(&testVal2)
		require.NoError(t, err)

		require.NotEqual(t, addr1, addr2)
	})

	t.Run("typed nil pointer returns address 0", func(t *testing.T) {
		// Documents current behaviour: a typed nil pointer is a pointer, its address is 0.
		var nilIntPtr *int

		addr, err := pointerutils.GetMemoryAddressAsUIntPtr(nilIntPtr)
		require.NoError(t, err)
		require.EqualValues(t, 0, addr)
	})
}

func TestGetMemoryAddressAsUInt64(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		addr, err := pointerutils.GetMemoryAddressAsUInt64(nil)
		require.Error(t, err)
		require.EqualValues(t, 0, addr)
	})

	t.Run("not a pointer", func(t *testing.T) {
		addr, err := pointerutils.GetMemoryAddressAsUInt64(5)
		require.Error(t, err)
		require.ErrorContains(t, err, "input is not a pointer")
		require.EqualValues(t, 0, addr)
	})

	t.Run("valid pointer", func(t *testing.T) {
		addr, err := pointerutils.GetMemoryAddressAsUInt64(&constIntForTesting)
		require.NoError(t, err)
		require.NotZero(t, addr)
	})

	t.Run("consistent with uintptr variant", func(t *testing.T) {
		addrPtr, err := pointerutils.GetMemoryAddressAsUIntPtr(&testVal1)
		require.NoError(t, err)

		addr64, err := pointerutils.GetMemoryAddressAsUInt64(&testVal1)
		require.NoError(t, err)

		require.Equal(t, uint64(addrPtr), addr64)
	})
}

func TestGetMemoryAddressAsHexString(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		address, err := pointerutils.GetMemoryAddressAsHexString(nil)
		require.Error(t, err)
		require.Empty(t, address)
	})

	t.Run("not a pointer", func(t *testing.T) {
		address, err := pointerutils.GetMemoryAddressAsHexString(5)
		require.Error(t, err)
		require.Empty(t, address)
	})

	tests := []struct {
		name    string
		pointer interface{}
	}{
		{"testVal1", &testVal1},
		{"testVal2", &testVal2},
		{"constIntForTesting", &constIntForTesting},
	}

	for _, tt := range tests {
		t.Run("valid pointer "+tt.name, func(t *testing.T) {
			expectedPointerAddress := fmt.Sprintf("%p", tt.pointer)

			address, err := pointerutils.GetMemoryAddressAsHexString(tt.pointer)
			require.NoError(t, err)
			require.Equal(t, expectedPointerAddress, address)
			require.Regexp(t, `^0x[0-9a-f]+$`, address)
		})
	}

	t.Run("consistent with uint64 variant", func(t *testing.T) {
		addr64, err := pointerutils.GetMemoryAddressAsUInt64(&testVal1)
		require.NoError(t, err)

		address, err := pointerutils.GetMemoryAddressAsHexString(&testVal1)
		require.NoError(t, err)

		require.Equal(t, fmt.Sprintf("0x%x", addr64), address)
	})

	t.Run("different variables give different strings", func(t *testing.T) {
		address1, err := pointerutils.GetMemoryAddressAsHexString(&testVal1)
		require.NoError(t, err)

		address2, err := pointerutils.GetMemoryAddressAsHexString(&testVal2)
		require.NoError(t, err)

		require.NotEqual(t, address1, address2)
	})
}

func TestPointersEqual(t *testing.T) {
	var nilIntPtr1 *int
	var nilIntPtr2 *int
	copyOfPtr1 := &testVal1

	tests := []struct {
		name            string
		pointer1        interface{}
		pointer2        interface{}
		expectedIsEqual bool
		expectError     bool
	}{
		// Valid comparisons
		{"nil and nil", nil, nil, true, false},
		{"nil and pointer", nil, &testVal1, false, false},
		{"pointer and nil", &testVal1, nil, false, false},
		{"same pointer testVal1", &testVal1, &testVal1, true, false},
		{"same pointer testVal2", &testVal2, &testVal2, true, false},
		{"copy of pointer", &testVal1, copyOfPtr1, true, false},
		{"different pointers, same value (1,2)", &testVal1, &testVal2, false, false},
		{"different pointers, same value (2,1)", &testVal2, &testVal1, false, false},

		// Typed nil pointers
		{"two typed nil pointers", nilIntPtr1, nilIntPtr2, true, false},
		{"typed nil and valid pointer", nilIntPtr1, &testVal1, false, false},
		{"valid pointer and typed nil", &testVal1, nilIntPtr1, false, false},
		{"untyped nil and typed nil", nil, nilIntPtr1, false, false},
		{"typed nil and untyped nil", nilIntPtr1, nil, false, false},

		// Error cases
		{"both not pointers", 5, 5, false, true},
		{"pointer and non-pointer", &testVal1, 5, false, true},
		{"non-pointer and pointer", 5, &testVal1, false, true},
		{"nil and non-pointer", nil, 5, false, true},
		{"non-pointer and nil", 5, nil, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equal, err := pointerutils.PointersEqual(tt.pointer1, tt.pointer2)

			if tt.expectError {
				require.Error(t, err)
				require.False(t, equal)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedIsEqual, equal)
		})
	}
}

func TestToInt64Pointer(t *testing.T) {
	t.Run("is a pointer", func(t *testing.T) {
		pointer := pointerutils.ToInt64Pointer(123)
		require.True(t, pointerutils.IsPointer(pointer))
	})

	values := []int64{0, 1, -1, 123, -123, 9223372036854775807, -9223372036854775808}
	for _, value := range values {
		t.Run(fmt.Sprintf("value %d", value), func(t *testing.T) {
			pointer := pointerutils.ToInt64Pointer(value)
			require.NotNil(t, pointer)
			require.Equal(t, value, *pointer)
		})
	}

	t.Run("each call returns a new pointer", func(t *testing.T) {
		p1 := pointerutils.ToInt64Pointer(42)
		p2 := pointerutils.ToInt64Pointer(42)

		equal, err := pointerutils.PointersEqual(p1, p2)
		require.NoError(t, err)
		require.False(t, equal)
	})

	t.Run("modifying the result does not affect other results", func(t *testing.T) {
		p1 := pointerutils.ToInt64Pointer(42)
		p2 := pointerutils.ToInt64Pointer(42)

		*p1 = 99

		require.EqualValues(t, 99, *p1)
		require.EqualValues(t, 42, *p2)
	})
}
