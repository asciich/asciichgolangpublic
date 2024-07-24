package asciichgolangpublic

import (
	"fmt"
	"runtime/debug"
	"strings"

)

const SOFTWARE_NAME_UNDEFINED = "[software name not defined]"
const FALLBACK_SOFTWARE_NAME_UNDEFINED = "[default software name not defined]"

var globalSoftwareName = SOFTWARE_NAME_UNDEFINED
var globalFallbackSoftwareName = FALLBACK_SOFTWARE_NAME_UNDEFINED

var softwareVersion = SOFTWARE_NAME_UNDEFINED       // constant values can no be overwritten by ldflags
var softwareName = FALLBACK_SOFTWARE_NAME_UNDEFINED // constant values can no be overwritten by ldflags

type BinaryInfo struct {
}

func GetBinaryInfo() (binaryInfo *BinaryInfo) {
	return new(BinaryInfo)
}

func GetSoftwareNameString() (name string) {
	return GetBinaryInfo().GetSoftwareNameString()
}

func GetSoftwareVersionString() (version string) {
	return GetBinaryInfo().GetSoftwareVersionString()
}

func LogVersion() {
	GetBinaryInfo().LogInfo()
}

func NewBinaryInfo() (b *BinaryInfo) {
	return new(BinaryInfo)
}

func (b *BinaryInfo) GetGitHash() (gitHash string, err error) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", TracedError("ReadBuildInfo failed")
	}
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			return setting.Value, nil
		}
	}

	return "", TracedError("Revision not found")
}

func (b *BinaryInfo) GetGitHashOrErrorMessageOnError() (gitHash string) {
	gitHash, err := b.GetGitHash()
	if err != nil {
		errorMessage := fmt.Sprintf("BinaryInfo.LogInfo: '%v'", err)
		LogError(errorMessage)
		gitHash = errorMessage
	}

	return gitHash
}

func (b *BinaryInfo) GetInfoString() (infoString string) {
	return fmt.Sprintf(
		"Software '%v' version: %v ; git hash: '%v'",
		b.GetSoftwareName(),
		b.GetSoftwareVersionString(),
		b.GetGitHashOrErrorMessageOnError(),
	)
}

func (b *BinaryInfo) GetSoftwareName() (softwareName string) {
	if !b.IsSoftwareNameSet() {
		if b.IsFallbackSoftwareNameSet() {
			return globalFallbackSoftwareName
		}
	}

	return globalSoftwareName
}

func (b *BinaryInfo) GetSoftwareNameString() (version string) {
	return softwareName
}

func (b *BinaryInfo) GetSoftwareVersionString() (version string) {
	return softwareVersion
}

func (b *BinaryInfo) IsFallbackSoftwareNameSet() (isSet bool) {
	return globalFallbackSoftwareName != FALLBACK_SOFTWARE_NAME_UNDEFINED
}

func (b *BinaryInfo) IsSoftwareNameSet() (isSet bool) {
	return globalSoftwareName != SOFTWARE_NAME_UNDEFINED
}

func (b *BinaryInfo) LogInfo() {
	logMessage := b.GetInfoString()
	LogInfo(logMessage)
}

func (b *BinaryInfo) MustGetGitHash() (gitHash string) {
	gitHash, err := b.GetGitHash()
	if err != nil {
		LogGoErrorFatal(err)
	}

	return gitHash
}

func (b *BinaryInfo) MustSetFallbackSoftwareName(defaultName string) {
	err := b.SetFallbackSoftwareName(defaultName)
	if err != nil {
		LogGoErrorFatal(err)
	}
}

func (b *BinaryInfo) SetFallbackSoftwareName(defaultName string) (err error) {
	defaultName = strings.TrimSpace(defaultName)
	if len(defaultName) <= 0 {
		return TracedError("defaultName is empty string")
	}

	globalFallbackSoftwareName = defaultName

	return nil
}
