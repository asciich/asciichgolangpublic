package changesummary

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChangeSummryIsChanged(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				changeSummary := NewChangeSummary()
				require.False(changeSummary.IsChanged())

				for i := 0; i < 2; i++ {
					changeSummary.SetChanged(false)
					require.False(changeSummary.IsChanged())
				}

				for i := 0; i < 2; i++ {
					changeSummary.SetChanged(true)
					require.True(changeSummary.IsChanged())
				}

				for i := 0; i < 2; i++ {
					changeSummary.SetChanged(false)
					require.False(changeSummary.IsChanged())
				}
			},
		)
	}
}

func TestChangeSummryNumberOfChnages(t *testing.T) {
	tests := []struct {
		testcase string
	}{
		{"testcase"},
	}

	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%v", tt),
			func(t *testing.T) {
				require := require.New(t)

				changeSummary := NewChangeSummary()
				require.False(changeSummary.IsChanged())
				require.EqualValues(0, changeSummary.GetNumberOfChanges())

				nLoops := 3

				for i := 0; i < nLoops; i++ {
					changeSummary.IncrementNumberOfChanges()
					require.True(changeSummary.IsChanged())
					require.EqualValues(i+1, changeSummary.GetNumberOfChanges())
				}

				for i := 0; i < 2; i++ {
					changeSummary.SetChanged(true)
					require.True(changeSummary.IsChanged())
					require.EqualValues(nLoops, changeSummary.GetNumberOfChanges())
				}

				for i := 0; i < 2; i++ {
					changeSummary.SetChanged(false)
					require.False(changeSummary.IsChanged())
					require.EqualValues(0, changeSummary.GetNumberOfChanges())
				}
			},
		)
	}
}
