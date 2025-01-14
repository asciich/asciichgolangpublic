package asciichgolangpublic

import (
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/datetime"
	"github.com/asciich/asciichgolangpublic/errors"
	"github.com/asciich/asciichgolangpublic/logging"
)

type SignalMessengersService struct {
}

func NewSignalMessengers() (s *SignalMessengersService) {
	return new(SignalMessengersService)
}

func NewSignalMessengersService() (s *SignalMessengersService) {
	return new(SignalMessengersService)
}

func SignalMessengers() (s *SignalMessengersService) {
	return NewSignalMessengers()
}

func (s *SignalMessengersService) MustParseCreationDateFromSignalPictureBaseName(baseName string) (creationDate *time.Time) {
	creationDate, err := s.ParseCreationDateFromSignalPictureBaseName(baseName)
	if err != nil {
		logging.LogGoErrorFatal(err)
	}

	return creationDate
}

func (s SignalMessengersService) ParseCreationDateFromSignalPictureBaseName(baseName string) (creationDate *time.Time, err error) {
	if baseName == "" {
		return nil, errors.TracedError("baseName is empty string")
	}

	if !strings.HasPrefix(baseName, "signal-") {
		return nil, errors.TracedErrorf("baseName '%s' is not a singal picture base name", baseName)
	}

	dateString := strings.TrimPrefix(baseName, "signal-")
	layoutString := "2006-01-02-15-04-05"

	if len(dateString) < len(layoutString) {
		return nil, errors.TracedErrorf("To short dateString: '%s'", dateString)
	}

	dateString = dateString[:len(layoutString)]
	creationDate, err = datetime.Dates().ParseStringWithGivenLayout(dateString, layoutString)
	if err != nil {
		return nil, err
	}

	return creationDate, nil
}
