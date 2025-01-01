package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGnuPg_SignAndValidate(t *testing.T) {
	if ContinuousIntegration().IsRunningInGithub() {
		LogWarnf("Not available in Github CI.")
		return
	}

	tests := []struct {
		implementationName string
		contentString      string
		asciiArmor         bool
	}{
		{"localFile", "hello world", true},
		{"localCommandExecutorFile", "hello world", true},
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				toTest := getFileToTest(tt.implementationName)
				defer toTest.Delete(verbose)

				signatureFile := toTest.MustGetParentDirectory().MustGetFileInDirectory(
					toTest.MustGetBaseName() + ".asc",
				)
				defer signatureFile.Delete(verbose)

				assert.True(toTest.MustExists(verbose))
				assert.False(signatureFile.MustExists(verbose))

				GnuPG().MustSignFile(
					toTest,
					&GnuPGSignOptions{
						DetachedSign: true,
						AsciiArmor:   tt.asciiArmor,
						Verbose:      verbose,
					},
				)

				assert.True(toTest.MustExists(verbose))
				assert.True(signatureFile.MustExists(verbose))

				GnuPG().MustCheckSignatureValid(signatureFile, verbose)
			},
		)
	}
}
