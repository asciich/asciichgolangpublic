package asciichgolangpublic

// Color codes for terminal.
//
// Source: https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
// Source: https://en.wikipedia.org/wiki/ANSI_escape_code

type TerminalColorsService struct{}

func NewTerminalColorsService() (t *TerminalColorsService) {
	return new(TerminalColorsService)
}

func TerminalColors() (terminalColors *TerminalColorsService) {
	return new(TerminalColorsService)
}

// This color code resets all color settings back default.
func (t *TerminalColorsService) GetCodeNoColor() (code string) {
	return "\033[0m"
}

func (t *TerminalColorsService) GetCodeBlack() (code string) {
	return "\033[0;30m"
}

func (t *TerminalColorsService) GetCodeBlue() (code string) {
	return "\033[0;34m"
}

func (t *TerminalColorsService) GetCodeBrightBlack() (code string) {
	return "\033[0;90m"
}

func (t *TerminalColorsService) GetCodeBrightBlue() (code string) {
	return "\033[0;94m"
}

func (t *TerminalColorsService) GetCodeBrightCyan() (code string) {
	return "\033[0;96m"
}

func (t *TerminalColorsService) GetCodeBrightGreen() (code string) {
	return "\033[0;92m"
}

func (t *TerminalColorsService) GetCodeBrightMagenta() (code string) {
	return "\033[0;95m"
}

func (t *TerminalColorsService) GetCodeBrightRed() (code string) {
	return "\033[0;91m"
}

func (t *TerminalColorsService) GetCodeBrightWhite() (code string) {
	return "\033[0;97m"
}

func (t *TerminalColorsService) GetCodeBrightYellow() (code string) {
	return "\033[0;93m"
}

func (t *TerminalColorsService) GetCodeCyan() (code string) {
	return "\033[0;36m"
}

func (t *TerminalColorsService) GetCodeGray() (code string) {
	return t.GetCodeBrightBlack()
}

func (t *TerminalColorsService) GetCodeGreen() (code string) {
	return "\033[0;32m"
}

func (t *TerminalColorsService) GetCodeMangenta() (code string) {
	return "\033[0;35m"
}

func (t *TerminalColorsService) GetCodeRed() (code string) {
	return "\033[0;31m"
}

func (t *TerminalColorsService) GetCodeWhite() (code string) {
	return "\033[0;37m"
}

func (t *TerminalColorsService) GetCodeYellow() (code string) {
	return "\033[0;33m"
}
