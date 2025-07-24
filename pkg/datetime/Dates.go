package datetime

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type DatesService struct{}

func Dates() (d *DatesService) {
	return NewDatesService()
}

func NewDatesService() (d *DatesService) {
	return new(DatesService)
}

func (d *DatesService) FormatAsYYYYmmdd_HHMMSSString(input *time.Time) (formated string, err error) {
	if input == nil {
		return "", tracederrors.TracedError("input is nil")
	}

	formated = input.Format("20060102_150405")
	return formated, nil
}

func (d *DatesService) GetCurrentYearAsString() (year string) {
	yearInt, _, _ := time.Now().Date()
	year = strconv.Itoa(yearInt)

	return year
}

func (d *DatesService) LayoutStringYYYYmmdd_HHMMSS() (layout string) {
	return "20060102_150405"
}

func (d *DatesService) MustFormatAsYYYYmmdd_HHMMSSString(input *time.Time) (formated string) {
	formated, err := d.FormatAsYYYYmmdd_HHMMSSString(input)
	if err != nil {
		log.Panic(err)
	}

	return formated
}

func (d *DatesService) MustParseString(input string) (date *time.Time) {
	date, err := d.ParseString(input)
	if err != nil {
		log.Panic(err)
	}

	return date
}

func (d *DatesService) MustParseStringPrefixAsDate(input string) (parsed *time.Time) {
	parsed, err := d.ParseStringPrefixAsDate(input)
	if err != nil {
		log.Panic(err)
	}

	return parsed
}

func (d *DatesService) MustParseStringWithGivenLayout(input string, layout string) (date *time.Time) {
	date, err := d.ParseStringWithGivenLayout(input, layout)
	if err != nil {
		log.Panic(err)
	}

	return date
}

func (d *DatesService) ParseString(input string) (date *time.Time, err error) {
	input = strings.TrimSpace(input)

	if input == "" {
		return nil, tracederrors.TracedError("input is empty string")
	}

	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		time.UnixDate,
		"Mon Jan _2 15:04:05 PM MST 2006", // Unix date with AM/PM
		"20060102_150405",
		"20060102-150405",
		time.RubyDate,
	}

	for _, layout := range layouts {
		date, err = d.ParseStringWithGivenLayout(input, layout)
		if err != nil {
			if strings.Contains(err.Error(), "Unable to parse as date ") {
				continue
			}

			return nil, err
		}

		return date, nil
	}

	return nil, tracederrors.TracedErrorf("Unable to parse date '%s'", input)
}

func (d *DatesService) ParseStringPrefixAsDate(input string) (parsed *time.Time, err error) {
	if input == "" {
		return nil, tracederrors.TracedError("input is empty string")
	}

	layoutString := d.LayoutStringYYYYmmdd_HHMMSS()
	if len(input) >= len(layoutString) {
		stringToCheck := input[0:len(layoutString)]

		parsed, err = d.ParseString(stringToCheck)
		if err != nil {
			return nil, err
		}
	}

	if parsed == nil {
		return nil, tracederrors.TracedErrorf(
			"Unable to parse prefix of '%s' as date.",
			input,
		)
	}

	return parsed, nil
}

func (d *DatesService) ParseStringWithGivenLayout(input string, layout string) (date *time.Time, err error) {
	input = strings.TrimSpace(input)
	layout = strings.TrimSpace(layout)

	if input == "" {
		return nil, tracederrors.TracedError("input is empty string")
	}

	if layout == "" {
		return nil, tracederrors.TracedError("layout is empty string")
	}

	parsed, err := time.Parse(layout, input)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to parse as date '%s' with given layout '%s'", input, layout)
	}

	return &parsed, nil
}
