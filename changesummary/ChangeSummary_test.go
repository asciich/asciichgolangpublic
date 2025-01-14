package changesummary

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
				assert := assert.New(t)

				changeSummary := NewChangeSummary()
				assert.False(changeSummary.IsChanged())

				for i := 0; i < 2; i++ {
					changeSummary.SetChanged(false)
					assert.False(changeSummary.IsChanged())
				}

				for i := 0; i < 2; i++ {
					changeSummary.SetChanged(true)
					assert.True(changeSummary.IsChanged())
				}

				for i := 0; i < 2; i++ {
					changeSummary.SetChanged(false)
					assert.False(changeSummary.IsChanged())
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
				assert := assert.New(t)

				changeSummary := NewChangeSummary()
				assert.False(changeSummary.IsChanged())
				assert.EqualValues(0, changeSummary.GetNumberOfChanges())

				nLoops := 3

				for i := 0; i < nLoops; i++ {
					changeSummary.IncrementNumberOfChanges()
					assert.True(changeSummary.IsChanged())
					assert.EqualValues(i+1, changeSummary.GetNumberOfChanges())
				}

				for i := 0; i < 2; i++ {
					changeSummary.SetChanged(true)
					assert.True(changeSummary.IsChanged())
					assert.EqualValues(nLoops, changeSummary.GetNumberOfChanges())
				}

				for i := 0; i < 2; i++ {
					changeSummary.SetChanged(false)
					assert.False(changeSummary.IsChanged())
					assert.EqualValues(0, changeSummary.GetNumberOfChanges())
				}
			},
		)
	}
}
