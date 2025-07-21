package logging

import "runtime"

type LogSettings struct {
	ColorDisabled bool
}

func NewLogSettings() (l *LogSettings) {
	return new(LogSettings)
}

func (l *LogSettings) GetColorDisabled() (colorDisabled bool) {

	return l.ColorDisabled
}

func (l *LogSettings) IsColorDisabled() (colorDisabled bool) {
	if runtime.GOOS == "windows" {
		// Color logging currently not implemented for Windows.
		return false
	}

	return l.ColorDisabled
}

func (l *LogSettings) IsColorEnabled() (colorEnabled bool) {
	return !l.ColorDisabled
}

func (l *LogSettings) SetColorDisabled(colorDisabled bool) {
	l.ColorDisabled = colorDisabled
}

func (l *LogSettings) SetColorEnabled(colorEnabled bool) {
	l.SetColorDisabled(!colorEnabled)
}
