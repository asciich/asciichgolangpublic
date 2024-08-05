package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(tt.expectedContains, Slices().ContainsInt(tt.inputSlice, tt.intToSearch))
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				trimmed := Slices().TrimSpace(tt.input)
				assert.EqualValues(tt.expectedOutput, trimmed)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				removedMatching := Slices().RemoveMatchingStrings(tt.input, tt.removeMatching)
				assert.EqualValues(tt.expectedOutput, removedMatching)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				removedContains := Slices().MustRemoveStringsWhichContains(tt.input, tt.searchString)
				assert.EqualValues(tt.expectedOutput, removedContains)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				maxValues := Slices().MaxIntValuePerIndex(tt.input1, tt.input2)
				assert.EqualValues(tt.expectedOutput, maxValues)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				output := Slices().RemoveLastElementIfEmptyString(tt.input)
				assert.EqualValues(tt.expectedOutput, output)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				output := Slices().RemoveDuplicatedStrings(tt.input)
				assert.EqualValues(tt.expectedOutput, output)
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
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedEqual,
					Slices().StringSlicesEqual(tt.input1, tt.input2),
				)
			},
		)
	}
}
