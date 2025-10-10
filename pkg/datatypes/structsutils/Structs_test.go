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
		objectToTest any
	}{
		{StructWithNoFields{}},
		{&StructWithNoFields{}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				fieldValues, err := GetFieldValuesAsString(tt.objectToTest)
				require.NoError(t, err)

				require.Len(t, fieldValues, 0)
			},
		)
	}
}

func TestStructsGetFieldValues_Values(t *testing.T) {
	type StructWithFields struct {
		A string
	}

	t.Run("no values set", func(t *testing.T) {
		fieldValues, err := GetFieldValuesAsString(StructWithFields{})
		require.NoError(t, err)

		require.Len(t, fieldValues, 1)
	})
}

func Test_GetFieldValueAsString(t *testing.T) {
	type StructWithFields struct {
		A string
	}

	t.Run("no values set", func(t *testing.T) {
		fieldValue, err := GetFieldValueAsString(StructWithFields{}, "A")
		require.NoError(t, err)

		require.EqualValues(t, fieldValue, "")
	})

	t.Run("no values set ptr", func(t *testing.T) {
		fieldValue, err := GetFieldValueAsString(&StructWithFields{}, "A")
		require.NoError(t, err)

		require.EqualValues(t, fieldValue, "")
	})

	t.Run("no values set and invalid Field requested", func(t *testing.T) {
		fieldValue, err := GetFieldValueAsString(&StructWithFields{}, "B")
		require.Errorf(t, err, "'B' does not exists and therefore an error is expected.")

		require.EqualValues(t, fieldValue, "")
	})
}

func Test_ListFieldNames(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		names, err := ListFieldNames(nil)
		require.Error(t, err)
		require.Nil(t, names)
	})

	t.Run("struct with no fields", func(t *testing.T) {
		type Emtpy struct{}

		names, err := ListFieldNames(Emtpy{})
		require.NoError(t, err)
		require.Len(t, names, 0)
	})

	t.Run("struct with no fields ptr", func(t *testing.T) {
		type Emtpy struct{}

		names, err := ListFieldNames(&Emtpy{})
		require.NoError(t, err)
		require.Len(t, names, 0)
	})

	t.Run("struct with one field", func(t *testing.T) {
		type OneField struct {
			A string
		}

		names, err := ListFieldNames(OneField{})
		require.NoError(t, err)
		require.EqualValues(t, []string{"A"}, names)
	})

	t.Run("struct with one field ptr", func(t *testing.T) {
		type OneField struct {
			A string
		}

		names, err := ListFieldNames(&OneField{})
		require.NoError(t, err)
		require.EqualValues(t, []string{"A"}, names)
	})

	t.Run("struct with two fields", func(t *testing.T) {
		type TwoFields struct {
			A string
			B string
		}

		names, err := ListFieldNames(TwoFields{})
		require.NoError(t, err)
		require.EqualValues(t, []string{"A", "B"}, names)
	})

	t.Run("struct with two fields ptr", func(t *testing.T) {
		type TwoFields struct {
			A string
			B string
		}

		names, err := ListFieldNames(&TwoFields{})
		require.NoError(t, err)
		require.EqualValues(t, []string{"A", "B"}, names)
	})
}
