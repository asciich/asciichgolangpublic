package commandlineinterface

import "strings"

func IsLinePromptOnly(line string) (isPromptOnly bool) {
	stripped := strings.TrimSpace(line)

	if stripped == "" {
		return false
	}

	if strings.Contains(stripped, "\n") {
		return false
	}

	if !strings.HasSuffix(stripped, "$") {
		return false
	}

	return true
}
