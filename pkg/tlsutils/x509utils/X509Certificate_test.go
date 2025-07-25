package x509utils

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorbashoo"
	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/mustutils"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/testutils"
)

func mustRepoRoot(ctx context.Context) (repoRootDir filesinterfaces.Directory) {
	repoRootPath, err := commandexecutorbashoo.Bash().RunCommandAndGetStdoutAsString(
		commandexecutorgeneric.WithLiveOutputOnStdoutIfVerbose(ctx),
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
	testDir, err := mustRepoRoot(getCtx()).GetSubDirectory("testdata", "X509Certificate", "LoadFromFilePath")
	require.NoError(t, err)

	type TestCase struct {
		testDir filesinterfaces.Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range mustutils.Must(testDir.ListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true})) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				inputFile, err := tt.testDir.GetFileInDirectory("input")
				require.NoError(t, err)

				expectedSubjectFile, err := tt.testDir.GetFileInDirectory("expected_subject")
				require.NoError(t, err)

				expectedIssuerStringFile, err := tt.testDir.GetFileInDirectory("expected_issuer_string")
				require.NoError(t, err)

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				subject := cert.MustGetSubjectString()
				expectedSubject := expectedSubjectFile.MustReadFirstLineAndTrimSpace()
				require.EqualValues(t, expectedSubject, subject)

				issuert := cert.MustGetIssuerString()
				expectedIssuer := expectedIssuerStringFile.MustReadFirstLineAndTrimSpace()
				require.EqualValues(t, expectedIssuer, issuert)
			},
		)
	}
}

func TestX509CertificateGetAsPemString(t *testing.T) {
	testDir, err := mustRepoRoot(getCtx()).GetSubDirectory("testdata", "X509Certificate", "GetAsPemString")
	require.NoError(t, err)

	type TestCase struct {
		testDir filesinterfaces.Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range mustutils.Must(testDir.ListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true})) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				inputFile, err := tt.testDir.GetFileInDirectory("input")
				require.NoError(t, err)

				expectedSubjectFile, err := tt.testDir.GetFileInDirectory("expected_pem")
				require.NoError(t, err)

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				pemString := cert.MustGetAsPemString()
				expectedPemString := expectedSubjectFile.MustReadAsString()

				require.EqualValues(t, expectedPemString, pemString)
			},
		)
	}
}

func TestX509CertificateIsRootCa(t *testing.T) {
	testDir, err := mustRepoRoot(getCtx()).GetSubDirectory("testdata", "X509Certificate", "IsRootCa")
	require.NoError(t, err)

	type TestCase struct {
		testDir filesinterfaces.Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range mustutils.Must(testDir.ListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true})) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				const verbose bool = true

				inputFile, err := tt.testDir.GetFileInDirectory("input")
				require.NoError(t, err)

				expectedSubjectFile, err := tt.testDir.GetFileInDirectory("expected_is_root_ca")
				require.NoError(t, err)

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				isRootCa := cert.MustIsRootCa(verbose)
				expectedIsRootCa := expectedSubjectFile.MustReadAsBool()

				require.EqualValues(t, expectedIsRootCa, isRootCa)
			},
		)
	}
}

func TestX509CertificateIsV1(t *testing.T) {
	testDir, err := mustRepoRoot(getCtx()).GetSubDirectory("testdata", "X509Certificate", "IsV1")
	require.NoError(t, err)

	type TestCase struct {
		testDir filesinterfaces.Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range mustutils.Must(testDir.ListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true})) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				inputFile, err := tt.testDir.GetFileInDirectory("input")
				require.NoError(t, err)

				expectedSubjectFile, err := tt.testDir.GetFileInDirectory("expected_is_v1")
				require.NoError(t, err)

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				isV1 := cert.MustIsV1()
				expectedIsV1 := expectedSubjectFile.MustReadAsBool()

				require.EqualValues(t, expectedIsV1, isV1)
			},
		)
	}
}

func TestX509CertificateIsV3(t *testing.T) {
	testDir, err := mustRepoRoot(getCtx()).GetSubDirectory("testdata", "X509Certificate", "IsV3")
	require.NoError(t, err)

	type TestCase struct {
		testDir filesinterfaces.Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range mustutils.Must(testDir.ListSubDirectories(&parameteroptions.ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true})) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			testutils.MustFormatAsTestname(tt),
			func(t *testing.T) {
				inputFile, err := tt.testDir.GetFileInDirectory("input")
				require.NoError(t, err)

				expectedSubjectFile, err := tt.testDir.GetFileInDirectory("expected_is_v3")
				require.NoError(t, err)

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				isV3 := cert.MustIsV3()
				expectedIsV3 := expectedSubjectFile.MustReadAsBool()

				require.EqualValues(t, expectedIsV3, isV3)
			},
		)
	}
}
