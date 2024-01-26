package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				isStruct := Structs().IsStruct(tt.objectToTest)
				assert.EqualValues(tt.isStruct, isStruct)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				isStruct := Structs().IsPointerToStruct(tt.objectToTest)
				assert.EqualValues(tt.isPointerToStruct, isStruct)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				isStruct := Structs().IsStructOrPointerToStruct(tt.objectToTest)
				assert.EqualValues(tt.isStructOrPointerToStruct, isStruct)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				fieldValues := Structs().MustGetFieldValuesAsString(tt.objectToTest)

				assert.Len(fieldValues, 0)
			},
		)
	}
}
