package binaryinfo

import (
	"fmt"
	"runtime/debug"
	"strings"

	"gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/logging"
	"gitlab.asciich.ch/tools/asciichgolangpublic.git/tracederrors"
)

const SOFTWARE_NAME_UNDEFINED = "[software name not defined]"
const FALLBACK_SOFTWARE_NAME_UNDEFINED = "[default software name not defined]"

var globalSoftwareName = SOFTWARE_NAME_UNDEFINED
var globalFallbackSoftwareName = FALLBACK_SOFTWARE_NAME_UNDEFINED

var softwareVersion = SOFTWARE_NAME_UNDEFINED       // constant values can no be overwritten by ldflags
var softwareName = FALLBACK_SOFTWARE_NAME_UNDEFINED // constant values can no be overwritten by ldflags

func LogVersion() {
	LogInfo()
}

func GetGitHash() (gitHash string, err error) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", tracederrors.TracedError("ReadBuildInfo failed")
	}
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			return setting.Value, nil
		}
	}

	return "", tracederrors.TracedError("Revision not found")
}

func GetGitHashOrErrorMessageOnError() (gitHash string) {
	gitHash, err := GetGitHash()
	if err != nil {
		errorMessage := fmt.Sprintf("BinaryInfo.LogInfo: '%v'", err)
		gitHash = errorMessage
	}

	return gitHash
}

func GetInfoString() (infoString string) {
	return fmt.Sprintf(
		"Software '%v' version: %v ; git hash: '%v'",
		GetSoftwareName(),
		GetSoftwareVersionString(),
		GetGitHashOrErrorMessageOnError(),
	)
}

func GetSoftwareName() (softwareName string) {
	if !IsSoftwareNameSet() {
		if IsFallbackSoftwareNameSet() {
			return globalFallbackSoftwareName
		}
	}

	return globalSoftwareName
}

func GetSoftwareNameString() (version string) {
	return softwareName
}

func GetSoftwareVersionString() (version string) {
	return softwareVersion
}

func IsFallbackSoftwareNameSet() (isSet bool) {
	return globalFallbackSoftwareName != FALLBACK_SOFTWARE_NAME_UNDEFINED
}

func IsSoftwareNameSet() (isSet bool) {
	return globalSoftwareName != SOFTWARE_NAME_UNDEFINED
}

func LogInfo() {
	logMessage := GetInfoString()
	logging.LogInfo(logMessage)
}

func MustGetGitHash() (gitHash string) {
	gitHash, err := GetGitHash()
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return gitHash
}

func MustSetFallbackSoftwareName(defaultName string) {
	err := SetFallbackSoftwareName(defaultName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}
}

func SetFallbackSoftwareName(defaultName string) (err error) {
	defaultName = strings.TrimSpace(defaultName)
	if len(defaultName) <= 0 {
		return tracederrors.TracedError("defaultName is empty string")
	}

	globalFallbackSoftwareName = defaultName

	return nil
}
