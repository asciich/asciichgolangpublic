package pointerutils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/pointerutils"
)

var constIntForTesting int = 10

func TestPointersIsPointer(t *testing.T) {
	tests := []struct {
		pointerToCheck    interface{}
		expectedIsPointer bool
	}{
		{5, false},
		{&constIntForTesting, true},
		{nil, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require.EqualValues(t, tt.expectedIsPointer, pointerutils.IsPointer(tt.pointerToCheck))
			},
		)
	}
}

var testVal1 = 5
var testVal2 = 5

func TestPointersGetMemoryAddressAsHexString(t *testing.T) {
	tests := []struct {
		pointer interface{}
	}{
		{&testVal1},
		{&testVal2},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				expectedPointerAddress := fmt.Sprintf("%p", tt.pointer)

				address, err := pointerutils.GetMemoryAddressAsHexString(tt.pointer)
				require.NoError(t, err)
				require.EqualValues(t, expectedPointerAddress, address)
			},
		)
	}
}

func TestPointersPointersEqual(t *testing.T) {
	tests := []struct {
		pointer1        interface{}
		pointer2        interface{}
		expectedIsEqual bool
	}{
		{nil, nil, true},
		{nil, &testVal1, false},
		{&testVal1, nil, false},
		{&testVal1, &testVal1, true},
		{&testVal2, &testVal2, true},
		{&testVal1, &testVal2, false},
		{&testVal2, &testVal1, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				equal, err := pointerutils.PointersEqual(tt.pointer1, tt.pointer2)
				require.NoError(t, err)
				require.EqualValues(t, tt.expectedIsEqual, equal)
			},
		)
	}
}

func TestToPointer(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		pointer := pointerutils.ToInt64Pointer(123)
		require.True(t, pointerutils.IsPointer(pointer))
	})
}
