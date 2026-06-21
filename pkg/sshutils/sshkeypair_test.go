package sshutils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/sshutils"
)

func TestValidateEmpty(t *testing.T) {
	ctx := contextutils.ContextVerbose()

	keyPair := &sshutils.SSHKeyPair{}
	err := keyPair.Validate(ctx)
	require.Error(t, err)
}
