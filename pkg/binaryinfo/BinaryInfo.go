package binaryinfo

import (
	"fmt"
	"runtime/debug"
	"strings"
)

const SOFTWARE_NAME_UNDEFINED = "[software name not defined]"
const SOFTWARE_VERSION_UNDEFINED = "[software version not defined]"

var SoftwareVersion = "" // constant values can no be overwritten by ldflags
var SoftwareName = ""    // constant values can no be overwritten by ldflags

var fallbackSoftwareName = ""

var ErrReadBuildInfoFailed = fmt.Errorf("read build info failed")
var ErrRevisionNotFound = fmt.Errorf("revision not found")

func GetGitHash() (gitHash string, err error) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", ErrReadBuildInfoFailed
	}
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			return setting.Value, nil
		}
	}

	return "", ErrRevisionNotFound
}

func GetGitHashOrErrorMessageOnError() (gitHash string) {
	gitHash, err := GetGitHash()
	if err != nil {
		errorMessage := fmt.Sprintf("BinaryInfo.LogInfo: '%s'", err.Error())
		gitHash = errorMessage
	}

	return gitHash
}

func GetInfoString() (infoString string) {
	return fmt.Sprintf(
		"Software '%s' version: '%s' ; git hash: '%s'",
		GetSoftwareName(),
		GetSoftwareVersionString(),
		GetGitHashOrErrorMessageOnError(),
	)
}

func GetSoftwareName() string {
	if SoftwareName != "" {
		return SoftwareName
	}

	if fallbackSoftwareName != "" {
		return fallbackSoftwareName
	}

	return SOFTWARE_NAME_UNDEFINED
}

func GetSoftwareVersionString() (version string) {
	if SoftwareVersion != "" {
		return SoftwareVersion
	}

	return SOFTWARE_VERSION_UNDEFINED
}

// Print the software version on stdout
func PrintInfo() {
	fmt.Println(GetInfoString())
}

func SetFallbackSoftwareName(fallbackName string) (err error) {
	fallbackName = strings.TrimSpace(fallbackName)
	if len(fallbackName) <= 0 {
		return fmt.Errorf("fallbackName is empty string")
	}

	fallbackSoftwareName = fallbackName

	return nil
}
