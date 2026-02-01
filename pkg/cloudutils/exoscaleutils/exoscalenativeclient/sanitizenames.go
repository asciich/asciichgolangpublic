package exoscalenativeclient

import (
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func SanitizeNames(domainName, recordName string) (string, string, error) {
	if domainName == "" {
		return "", "", tracederrors.TracedErrorEmptyString("domainName")
	}

	if recordName == "" {
		return "", "", tracederrors.TracedErrorEmptyString("recordName")
	}

	recordName = strings.TrimSuffix(recordName, "."+domainName)

	return domainName, recordName, nil
}
