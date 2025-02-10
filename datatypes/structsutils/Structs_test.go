package structsutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type testStruct struct{}

func TestStructsIsStruct(t *testing.T) {

	var intToTest int = 1
	var floatTotTest float64 = 1.1

	tests := []struct {
		objectToTest interface{}
		isStruct     bool
	}{
		{"stringObject", false},
		{1, false},
		{&intToTest, false},
		{&floatTotTest, false},
		{1.1, false},
		{false, false},
		{true, false},
		{nil, false},
		{testStruct{}, true},
		{&testStruct{}, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				isStruct := IsStruct(tt.objectToTest)
				require.EqualValues(tt.isStruct, isStruct)
			},
		)
	}
}

func TestStructsIsPointerToStruct(t *testing.T) {

	var intToTest int = 1
	var floatTotTest float64 = 1.1

	tests := []struct {
		objectToTest      interface{}
		isPointerToStruct bool
	}{
		{"stringObject", false},
		{1, false},
		{&intToTest, false},
		{&floatTotTest, false},
		{1.1, false},
		{false, false},
		{true, false},
		{nil, false},
		{testStruct{}, false},
		{&testStruct{}, true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				isStruct := IsPointerToStruct(tt.objectToTest)
				require.EqualValues(tt.isPointerToStruct, isStruct)
			},
		)
	}
}

func TestStructsIsStructOrPointerToStruct(t *testing.T) {

	var intToTest int = 1
	var floatTotTest float64 = 1.1

	tests := []struct {
		objectToTest              interface{}
		isStructOrPointerToStruct bool
	}{
		{"stringObject", false},
		{1, false},
		{&intToTest, false},
		{&floatTotTest, false},
		{1.1, false},
		{false, false},
		{true, false},
		{nil, false},
		{testStruct{}, true},
		{&testStruct{}, true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				isStruct := IsStructOrPointerToStruct(tt.objectToTest)
				require.EqualValues(tt.isStructOrPointerToStruct, isStruct)
			},
		)
	}
}

func TestStructsGetFieldValues_NoValues(t *testing.T) {
	type StructWithNoFields struct{}

	tests := []struct {
		objectToTest interface{}
	}{
		{StructWithNoFields{}},
		{&StructWithNoFields{}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				fieldValues := MustGetFieldValuesAsString(tt.objectToTest)

				require.Len(fieldValues, 0)
			},
		)
	}
}
