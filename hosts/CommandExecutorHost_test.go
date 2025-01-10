package hosts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandExecutorHost_HostnameOfLocalhost(t *testing.T) {
	assert := assert.New(t)

	host := MustGetLocalCommandExecutorHost()

	assert.EqualValues(
		"localhost",
		host.MustGetHostName(),
	)
}
