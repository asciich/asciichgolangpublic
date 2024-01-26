package asciichgolangpublic

import (
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
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				assert.EqualValues(
					tt.expectedIsPointer,
					Pointers().IsPointer(tt.pointerToCheck),
				)
			},
		)
	}
}
