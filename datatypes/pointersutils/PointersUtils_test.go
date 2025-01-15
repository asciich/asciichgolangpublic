package pointersutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
		tt := tt
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedIsPointer,
					IsPointer(tt.pointerToCheck),
				)
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
				assert := assert.New(t)

				expectedPointerAddress := fmt.Sprintf("%p", tt.pointer)

				assert.EqualValues(
					expectedPointerAddress,
					MustGetMemoryAddressAsHexString(tt.pointer),
				)
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
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedIsEqual,
					MustPointersEqual(tt.pointer1, tt.pointer2),
				)
			},
		)
	}
}
