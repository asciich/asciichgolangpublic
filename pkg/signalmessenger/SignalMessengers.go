package signalmessenger

import (
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/datetime"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/tracederrors"
)

func MustParseCreationDateFromSignalPictureBaseName(baseName string) (creationDate *time.Time) {
	creationDate, err := ParseCreationDateFromSignalPictureBaseName(baseName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return creationDate
}

func ParseCreationDateFromSignalPictureBaseName(baseName string) (creationDate *time.Time, err error) {
	if baseName == "" {
		return nil, tracederrors.TracedError("baseName is empty string")
	}

	if !strings.HasPrefix(baseName, "signal-") {
		return nil, tracederrors.TracedErrorf("baseName '%s' is not a singal picture base name", baseName)
	}

	dateString := strings.TrimPrefix(baseName, "signal-")
	layoutString := "2006-01-02-15-04-05"

	if len(dateString) < len(layoutString) {
		return nil, tracederrors.TracedErrorf("To short dateString: '%s'", dateString)
	}

	dateString = dateString[:len(layoutString)]
	creationDate, err = datetime.Dates().ParseStringWithGivenLayout(dateString, layoutString)
	if err != nil {
		return nil, err
	}

	return creationDate, nil
}
