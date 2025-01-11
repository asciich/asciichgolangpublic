package kvm
/* TODO enable again
import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLibvirtXmlsGetMacAddress(t *testing.T) {
	gitRepo := MustGetLocalGitRepositoryByPath(".")
	testDir := gitRepo.MustGetSubDirectory("testdata", "LibvirtXmls", "GetMacAddress")

	type TestCase struct {
		testDir *Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range MustGetLegacyDir(testDir).MustGetSubDirectoryList(&ListDirectoryOptions{NonRecursive: true, ReturnRelativePaths: true}) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				inputString := tt.testDir.ReadFileInDirectoryAsString("input.xml")
				expectedMac := tt.testDir.MustGetFirstLineOfFileInDirectoryAsString("expected_mac.txt")

				macAddress := LibvirtXmls().MustGetMacAddressFromXmlString(inputString)

				assert.EqualValues(expectedMac, macAddress)
			},
		)
	}
}
*/