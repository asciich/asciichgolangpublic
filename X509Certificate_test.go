package asciichgolangpublic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestX509CertificateLoadFromFilePath(t *testing.T) {
	gitRepo := MustGetLocalGitRepositoryByPath(".")
	testDir := gitRepo.MustGetSubDirectory("testdata", "X509Certificate", "LoadFromFilePath")

	type TestCase struct {
		testDir Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range testDir.MustGetSubDirectories(&ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true}) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				inputFile := tt.testDir.MustGetFileInDirectory("input")
				expectedSubjectFile := tt.testDir.MustGetFileInDirectory("expected_subject")
				expectedIssuerStringFile := tt.testDir.MustGetFileInDirectory("expected_issuer_string")

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				subject := cert.MustGetSubjectString()
				expectedSubject := expectedSubjectFile.MustReadFirstLineAndTrimSpace()
				assert.EqualValues(expectedSubject, subject)

				issuert := cert.MustGetIssuerString()
				expectedIssuer := expectedIssuerStringFile.MustReadFirstLineAndTrimSpace()
				assert.EqualValues(expectedIssuer, issuert)
			},
		)
	}
}

func TestX509CertificateGetAsPemString(t *testing.T) {
	gitRepo := MustGetLocalGitRepositoryByPath(".")
	testDir := gitRepo.MustGetSubDirectory("testdata", "X509Certificate", "GetAsPemString")

	type TestCase struct {
		testDir Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range testDir.MustGetSubDirectories(&ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true}) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				inputFile := tt.testDir.MustGetFileInDirectory("input")
				expectedSubjectFile := tt.testDir.MustGetFileInDirectory("expected_pem")

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				pemString := cert.MustGetAsPemString()
				expectedPemString := expectedSubjectFile.MustReadAsString()

				assert.EqualValues(expectedPemString, pemString)
			},
		)
	}
}

func TestX509CertificateIsRootCa(t *testing.T) {
	gitRepo := MustGetLocalGitRepositoryByPath(".")
	testDir := gitRepo.MustGetSubDirectory("testdata", "X509Certificate", "IsRootCa")

	type TestCase struct {
		testDir Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range testDir.MustGetSubDirectories(&ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true}) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				const verbose bool = true

				inputFile := tt.testDir.MustGetFileInDirectory("input")
				expectedSubjectFile := tt.testDir.MustGetFileInDirectory("expected_is_root_ca")

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				isRootCa := cert.MustIsRootCa(verbose)
				expectedIsRootCa := expectedSubjectFile.MustReadAsBool()

				assert.EqualValues(expectedIsRootCa, isRootCa)
			},
		)
	}
}

func TestX509CertificateIsV1(t *testing.T) {
	gitRepo := MustGetLocalGitRepositoryByPath(".")
	testDir := gitRepo.MustGetSubDirectory("testdata", "X509Certificate", "IsV1")

	type TestCase struct {
		testDir Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range testDir.MustGetSubDirectories(&ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true}) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				inputFile := tt.testDir.MustGetFileInDirectory("input")
				expectedSubjectFile := tt.testDir.MustGetFileInDirectory("expected_is_v1")

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				isV1 := cert.MustIsV1()
				expectedIsV1 := expectedSubjectFile.MustReadAsBool()

				assert.EqualValues(expectedIsV1, isV1)
			},
		)
	}
}

func TestX509CertificateIsV3(t *testing.T) {
	gitRepo := MustGetLocalGitRepositoryByPath(".")
	testDir := gitRepo.MustGetSubDirectory("testdata", "X509Certificate", "IsV3")

	type TestCase struct {
		testDir Directory
	}

	tests := []TestCase{}
	for _, testCaseDir := range testDir.MustGetSubDirectories(&ListDirectoryOptions{Recursive: false, ReturnRelativePaths: true}) {
		tests = append(tests, TestCase{testCaseDir})
	}

	for _, tt := range tests {
		t.Run(
			MustFormatAsTestname(tt),
			func(t *testing.T) {
				assert := assert.New(t)

				inputFile := tt.testDir.MustGetFileInDirectory("input")
				expectedSubjectFile := tt.testDir.MustGetFileInDirectory("expected_is_v3")

				cert := MustGetX509CertificateFromFilePath(inputFile.MustGetLocalPath())

				isV3 := cert.MustIsV3()
				expectedIsV3 := expectedSubjectFile.MustReadAsBool()

				assert.EqualValues(expectedIsV3, isV3)
			},
		)
	}
}
