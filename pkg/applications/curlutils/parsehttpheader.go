package curlutils

import (
	"strconv"
	"strings"

	"github.com/asciich/asciichgolangpublic/pkg/datatypes/slicesutils"
	"github.com/asciich/asciichgolangpublic/pkg/datatypes/stringsutils"
	"github.com/asciich/asciichgolangpublic/pkg/httputils/httpgeneric"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func ParseHttpHeader(headerData []byte) (*httpgeneric.HttpHeader, error) {
	if headerData == nil {
		return nil, tracederrors.TracedErrorNil("headerData")
	}

	headerDataString := string(headerData)
	if headerDataString == "" {
		return nil, tracederrors.TracedErrorEmptyString(headerDataString)
	}

	lines := stringsutils.SplitLines(headerDataString, true)
	lines = slicesutils.RemoveEmptyStrings(lines)

	header := &httpgeneric.HttpHeader{}

	if len(lines) > 0 {
		line := lines[0]
		splitted := strings.Split(line, " ")
		if len(splitted) < 3 {
			return nil, tracederrors.TracedErrorf("Failed to parse http status code from line='%s' of header='%s'.", line, headerDataString)
		}

		statusCodeString := splitted[1]
		statusCode, err := strconv.Atoi(statusCodeString)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to parse HTTP status code string '%s' from line='%s' of header='%s': %w", statusCodeString, line, headerDataString, err)
		}

		header.StatusCode = statusCode
	}

	if header.StatusCode == 0 {
		return nil, tracederrors.TracedErrorf("Failed to parse HTTP header. StatusCode is 0. header='%s'.", headerDataString)
	}

	return header, nil
}
