package x509utils

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/commandexecutor"
	"github.com/asciich/asciichgolangpublic/files"
	"github.com/asciich/asciichgolangpublic/logging"
	"github.com/asciich/asciichgolangpublic/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/testutils"
)

func mustRepoRoot(ctx context.Context) (repoRootDir files.Directory) {
	const verbose = true

	repoRootPath, err := commandexecutor.Bash().RunCommandAndGetStdoutAsString(
		commandexecutor.WithLiveOutputOnStdoutIfVerbose(ctx),
		&parameteroptions.RunCommandOptions{
			Command: []string{"git", "-C", ".", "rev-parse", "--show-toplevel"},
		},
	)
	if err != nil {
		logging.LogGoErrorFatalWithTrace(err)
	}

	repoRootPath = strings.TrimSpace(repoRootPath)

	repoRootDir, err = files.GetLocalDirectoryByPath(repoRootPath)
	if err != nil {
		logging.LogGoErrorFatalWithTrace(err)
	}

	return repoRootDir
}

func TestX509CertificateLoadFromFilePath(t *testing.T) {
	testDir := mustRepoRoot(getCtx()).MustGetSubDirectory("testdata", "X509Certificate", "LoadFromFilePath")

	type TestCase struct {
		testDir files.Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range mustutils.Must(testDir.ListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true})) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				inputFile := tt.testDir.MustGetFileInDirectory("input")
				expectedSubjectFile := tt.testDir.MustGetFileInDirectory("expected_subject")
				expectedIssuerStringFile := tt.testDir.MustGetFileInDirectory("expected_issuer_string")

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				subject := cert.MustGetSubjectString()
				expectedSubject := expectedSubjectFile.MustReadFirstLineAndTrimSpace()
				require.EqualValues(expectedSubject, subject)

				issuert := cert.MustGetIssuerString()
				expectedIssuer := expectedIssuerStringFile.MustReadFirstLineAndTrimSpace()
				require.EqualValues(expectedIssuer, issuert)
			},
		)
	}
}

func TestX509CertificateGetAsPemString(t *testing.T) {
	testDir := mustRepoRoot(getCtx()).MustGetSubDirectory("testdata", "X509Certificate", "GetAsPemString")

	type TestCase struct {
		testDir files.Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range mustutils.Must(testDir.ListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true})) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				inputFile := tt.testDir.MustGetFileInDirectory("input")
				expectedSubjectFile := tt.testDir.MustGetFileInDirectory("expected_pem")

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				pemString := cert.MustGetAsPemString()
				expectedPemString := expectedSubjectFile.MustReadAsString()

				require.EqualValues(expectedPemString, pemString)
			},
		)
	}
}

func TestX509CertificateIsRootCa(t *testing.T) {
	testDir := mustRepoRoot(getCtx()).MustGetSubDirectory("testdata", "X509Certificate", "IsRootCa")

	type TestCase struct {
		testDir files.Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range mustutils.Must(testDir.ListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true})) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				const verbose bool = true

				inputFile := tt.testDir.MustGetFileInDirectory("input")
				expectedSubjectFile := tt.testDir.MustGetFileInDirectory("expected_is_root_ca")

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				isRootCa := cert.MustIsRootCa(verbose)
				expectedIsRootCa := expectedSubjectFile.MustReadAsBool()

				require.EqualValues(expectedIsRootCa, isRootCa)
			},
		)
	}
}

func TestX509CertificateIsV1(t *testing.T) {
	testDir := mustRepoRoot(getCtx()).MustGetSubDirectory("testdata", "X509Certificate", "IsV1")

	type TestCase struct {
		testDir files.Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range mustutils.Must(testDir.ListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true})) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				inputFile := tt.testDir.MustGetFileInDirectory("input")
				expectedSubjectFile := tt.testDir.MustGetFileInDirectory("expected_is_v1")

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				isV1 := cert.MustIsV1()
				expectedIsV1 := expectedSubjectFile.MustReadAsBool()

				require.EqualValues(expectedIsV1, isV1)
			},
		)
	}
}

func TestX509CertificateIsV3(t *testing.T) {
	testDir := mustRepoRoot(getCtx()).MustGetSubDirectory("testdata", "X509Certificate", "IsV3")

	type TestCase struct {
		testDir files.Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range mustutils.Must(testDir.ListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true})) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				require := require.New(t)

				inputFile := tt.testDir.MustGetFileInDirectory("input")
				expectedSubjectFile := tt.testDir.MustGetFileInDirectory("expected_is_v3")

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				isV3 := cert.MustIsV3()
				expectedIsV3 := expectedSubjectFile.MustReadAsBool()

				require.EqualValues(expectedIsV3, isV3)
			},
		)
	}
}
