package mustutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func functionReturningOneValueAndAnError() (string, error) {
	return "", nil
}

func functionReturningTwoValuesAndAnError() (string, int, error) {
	return "", 123, nil
}

func functionReturningThreeValuesAndAnError() (string, int, bool, error) {
	return "", 123, false, nil
}

func functionReturningFourValuesAndAnError() (string, int, bool, uint64, error) {
	return "", 123, false, 17, nil
}

func TestMustUtils_Must(t *testing.T) {
	assert.EqualValues(
		t,
		"",
		Must(functionReturningOneValueAndAnError()),
	)
}

func TestMustUtils_Must2(t *testing.T) {
	v1, v2 := Must2(functionReturningTwoValuesAndAnError())

	assert.EqualValues(t, "", v1)
	assert.EqualValues(t, 123, v2)
}

func TestMustUtils_Must3(t *testing.T) {
	v1, v2, v3 := Must3(functionReturningThreeValuesAndAnError())

	assert.EqualValues(t, "", v1)
	assert.EqualValues(t, 123, v2)
	assert.EqualValues(t, false, v3)
}

func TestMustUtils_Must4(t *testing.T) {
	v1, v2, v3, v4 := Must4(functionReturningFourValuesAndAnError())

	assert.EqualValues(t, "", v1)
	assert.EqualValues(t, 123, v2)
	assert.EqualValues(t, false, v3)
	assert.EqualValues(t, uint64(17), v4)
}
