package slicesutils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlicesContainsInt(t *testing.T) {
	tests := []struct {
		intToSearch      int
		inputSlice       []int
		expectedContains bool
	}{
		{0, []int{}, false},
		{0, []int{1}, false},
		{0, []int{1, 2}, false},
		{0, []int{0, 1, 2}, true},
		{0, []int{0, 1, 2, 0}, true},
		{0, []int{1, 2, 0}, true},
		{1, []int{1, 2, 0}, true},
		{2, []int{1, 2, 0}, true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				require.EqualValues(tt.expectedContains, ContainsInt(tt.inputSlice, tt.intToSearch))
			},
		)
	}
}

func TestSlicesContainsString(t *testing.T) {
	tests := []struct {
		stringToSearch   string
		inputSlice       []string
		expectedContains bool
	}{
		{"0", []string{}, false},
		{"0", []string{"1"}, false},
		{"0", []string{"1", "2"}, false},
		{"0", []string{"0", "1", "2"}, true},
		{"0", []string{"0", "1", "2", "0"}, true},
		{"0", []string{"1", "2", "0"}, true},
		{"1", []string{"1", "2", "0"}, true},
		{"2", []string{"1", "2", "0"}, true},
		{"hello", []string{"hello", "Hello", "world", "World"}, true},
		{"Hello", []string{"hello", "Hello", "world", "World"}, true},
		{"HellO", []string{"hello", "Hello", "world", "World"}, false},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				require.EqualValues(tt.expectedContains, ContainsString(tt.inputSlice, tt.stringToSearch))
			},
		)
	}
}

func TestSlicesContainsStringIgnoreCase(t *testing.T) {
	tests := []struct {
		stringToSearch   string
		inputSlice       []string
		expectedContains bool
	}{
		{"0", []string{}, false},
		{"0", []string{"1"}, false},
		{"0", []string{"1", "2"}, false},
		{"0", []string{"0", "1", "2"}, true},
		{"0", []string{"0", "1", "2", "0"}, true},
		{"0", []string{"1", "2", "0"}, true},
		{"1", []string{"1", "2", "0"}, true},
		{"2", []string{"1", "2", "0"}, true},
		{"hello", []string{"hello", "Hello", "world", "World"}, true},
		{"Hello", []string{"hello", "Hello", "world", "World"}, true},
		{"HellO", []string{"hello", "Hello", "world", "World"}, true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				require.EqualValues(tt.expectedContains, ContainsStringIgnoreCase(tt.inputSlice, tt.stringToSearch))
			},
		)
	}
}

func TestSlicesTrimSpace(t *testing.T) {
	tests := []struct {
		input          []string
		expectedOutput []string
	}{
		{[]string{}, []string{}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", "b"}, []string{"a", "b"}},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{" a", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{" a\t", "b", "c"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b", "c"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b  ", "c"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b  ", "\nc"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b  ", " \nc"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b  ", " \n c"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b  ", " \n \tc"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b  ", " \n \tc\n"}, []string{"a", "b", "c"}},
		{[]string{" a\t", " b  ", " \n \tc\n\n"}, []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				trimmed := TrimSpace(tt.input)
				require.EqualValues(tt.expectedOutput, trimmed)
			},
		)
	}
}

func TestSlicesRemoveMatchingStrings(t *testing.T) {
	tests := []struct {
		input          []string
		removeMatching string
		expectedOutput []string
	}{
		{[]string{}, "", []string{}},
		{[]string{}, "a", []string{}},
		{[]string{""}, "", []string{}},
		{[]string{"a"}, "a", []string{}},
		{[]string{"a", "b"}, "a", []string{"b"}},
		{[]string{"a", "b", "a"}, "a", []string{"b"}},
		{[]string{"a", "b", "c"}, "a", []string{"b", "c"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				removedMatching := RemoveMatchingStrings(tt.input, tt.removeMatching)
				require.EqualValues(tt.expectedOutput, removedMatching)
			},
		)
	}
}

func TestSlicesRemoveStringsWhichContains(t *testing.T) {
	tests := []struct {
		input          []string
		searchString   string
		expectedOutput []string
	}{
		{[]string{}, "a", []string{}},
		{[]string{"a"}, "a", []string{}},
		{[]string{"a", "b"}, "a", []string{"b"}},
		{[]string{"a", "b", "a"}, "a", []string{"b"}},
		{[]string{"a", "b", "c"}, "a", []string{"b", "c"}},
		{[]string{"a", "b", "ca"}, "a", []string{"b"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				removedContains := MustRemoveStringsWhichContains(tt.input, tt.searchString)
				require.EqualValues(tt.expectedOutput, removedContains)
			},
		)
	}
}

func TestSlicesMaxIntValuePerIndex(t *testing.T) {
	tests := []struct {
		input1         []int
		input2         []int
		expectedOutput []int
	}{
		{nil, nil, []int{}},
		{[]int{0}, []int{1}, []int{1}},
		{[]int{0}, []int{-1}, []int{0}},
		{[]int{-10}, []int{-1}, []int{-1}},
		{[]int{-10}, []int{-1, 1, 2, 3}, []int{-1, 1, 2, 3}},
		{[]int{-10, 0, 0, 0, 4}, []int{-1, 1, 2, 3}, []int{-1, 1, 2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				maxValues := MaxIntValuePerIndex(tt.input1, tt.input2)
				require.EqualValues(tt.expectedOutput, maxValues)
			},
		)
	}
}

func TestSlicesRemoveLastElementIfEmptyString(t *testing.T) {
	tests := []struct {
		input          []string
		expectedOutput []string
	}{
		{[]string{}, []string{}},
		{nil, []string{}},
		{[]string{""}, []string{}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", ""}, []string{"a"}},
		{[]string{"a", "b", ""}, []string{"a", "b"}},
		{[]string{"a", "b", "", "c"}, []string{"a", "b", "", "c"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				output := RemoveLastElementIfEmptyString(tt.input)
				require.EqualValues(tt.expectedOutput, output)
			},
		)
	}
}

func TestSlicesRemoveDuplicatedEntries(t *testing.T) {
	tests := []struct {
		input          []string
		expectedOutput []string
	}{
		{[]string{}, []string{}},
		{nil, []string{}},
		{[]string{""}, []string{""}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", ""}, []string{"a", ""}},
		{[]string{"a", "b", ""}, []string{"a", "b", ""}},
		{[]string{"a", "b", "", "c"}, []string{"a", "b", "", "c"}},
		{[]string{"a", "a", "", "c"}, []string{"a", "", "c"}},
		{[]string{"a", "a", "a", "c"}, []string{"a", "c"}},
		{[]string{"a", "a", "a", "a"}, []string{"a"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				output := RemoveDuplicatedStrings(tt.input)
				require.EqualValues(tt.expectedOutput, output)
			},
		)
	}
}

func TestSlicesStringSlicesEqual(t *testing.T) {
	tests := []struct {
		input1        []string
		input2        []string
		expectedEqual bool
	}{
		{nil, nil, false},
		{nil, []string{}, false},
		{[]string{}, nil, false},
		{[]string{}, []string{}, true},
		{[]string{}, []string{"a"}, false},
		{[]string{"a"}, []string{}, false},
		{[]string{"a"}, []string{"A"}, false},
		{[]string{"a"}, []string{"a"}, true},
		{[]string{"a", "b"}, []string{"a"}, false},
		{[]string{"a"}, []string{"a", "b"}, false},
		{[]string{"a", "b"}, []string{"a", "b"}, true},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				require.EqualValues(
					tt.expectedEqual,
					StringSlicesEqual(tt.input1, tt.input2),
				)
			},
		)
	}
}

func TestSlicesDiffStringSlices(t *testing.T) {
	tests := []struct {
		input1          []string
		input2          []string
		expectedANotInB []string
		expectedBNotInA []string
	}{
		{nil, nil, []string{}, []string{}},
		{[]string{}, nil, []string{}, []string{}},
		{nil, []string{}, []string{}, []string{}},
		{[]string{}, []string{}, []string{}, []string{}},
		{[]string{"a"}, []string{}, []string{"a"}, []string{}},
		{[]string{"a"}, []string{"b"}, []string{"a"}, []string{"b"}},
		{[]string{""}, []string{"b"}, []string{""}, []string{"b"}},
		{[]string{""}, []string{"b", "a"}, []string{""}, []string{"a", "b"}},
		{[]string{"c"}, []string{"b", "a"}, []string{"c"}, []string{"a", "b"}},
		{[]string{"a", "c"}, []string{"b", "a"}, []string{"c"}, []string{"b"}},
		{[]string{"a", "c"}, []string{"a"}, []string{"c"}, []string{}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				aNotInB, bNotInA := DiffStringSlices(tt.input1, tt.input2)

				require.EqualValues(
					tt.expectedANotInB,
					aNotInB,
				)
				require.EqualValues(
					tt.expectedBNotInA,
					bNotInA,
				)
			},
		)
	}
}

func TestSlicesGetDeepCopyOfByteSlice(t *testing.T) {
	tests := []struct {
		input           []byte
		expected_output []byte
	}{
		{[]byte{}, []byte{}},
		{nil, nil},
		{[]byte("a"), []byte("a")},
		{[]byte("ab"), []byte("ab")},
		{[]byte("abc"), []byte("abc")},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				copy := GetDeepCopyOfByteSlice(tt.input)
				require.EqualValues(tt.expected_output, copy)

				for i := 0; i < len(tt.input); i++ {
					tt.input[i] = 0x00
				}

				require.EqualValues(tt.expected_output, copy)
			},
		)
	}
}

func TestSlicesGetDeepCopyOfStringSlice(t *testing.T) {
	tests := []struct {
		input           []string
		expected_output []string
	}{
		{nil, []string{}},
		{[]string{"a"}, []string{"a"}},
		{[]string{"a", ""}, []string{"a", ""}},
		{[]string{"a", "b"}, []string{"a", "b"}},
		{[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				copy := GetDeepCopyOfStringsSlice(tt.input)
				require.EqualValues(tt.expected_output, copy)
			},
		)
	}
}

func TestSlices_GetSortedDeepCopyString(t *testing.T) {
	require := require.New(t)
	inputSlice := []string{"c", "b", "a"}

	sorted := GetSortedDeepCopyOfStringsSlice(inputSlice)
	require.EqualValues([]string{"a", "b", "c"}, sorted)
	require.NotEqual(inputSlice, sorted)
	require.EqualValues([]string{"c", "b", "a"}, inputSlice)
}

func TestSlices_RemoveEmptyStringsAtEnd(t *testing.T) {
	tests := []struct {
		input          []string
		expectedOutput []string
	}{
		{[]string{}, []string{}},
		{[]string{""}, []string{}},
		{[]string{"", ""}, []string{}},
		{[]string{"", "", ""}, []string{}},
		{[]string{"a", ""}, []string{"a"}},
		{[]string{"a", "", ""}, []string{"a"}},
		{[]string{"a", "", "", ""}, []string{"a"}},
		{[]string{"a", "b", ""}, []string{"a", "b"}},
		{[]string{"a", "b", "", ""}, []string{"a", "b"}},
		{[]string{"a", "b", "", "", ""}, []string{"a", "b"}},
		{[]string{"a", "", "b", ""}, []string{"a", "", "b"}},
		{[]string{"a", "", "b", "", ""}, []string{"a", "", "b"}},
		{[]string{"a", "", "b", "", "", ""}, []string{"a", "", "b"}},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				copy := RemoveEmptyStringsAtEnd(tt.input)
				require.EqualValues(tt.expectedOutput, copy)
			},
		)
	}
}
