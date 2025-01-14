package terminalcolors

// Color codes for terminal.
//
// Source: https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
// Source: https://en.wikipedia.org/wiki/ANSI_escape_code

// This color code resets all color settings back default.
func GetCodeNoColor() (code string) {
	return "\033[0m"
}

func GetCodeBlack() (code string) {
	return "\033[0;30m"
}

func GetCodeBlue() (code string) {
	return "\033[0;34m"
}

func GetCodeBrightBlack() (code string) {
	return "\033[0;90m"
}

func GetCodeBrightBlue() (code string) {
	return "\033[0;94m"
}

func GetCodeBrightCyan() (code string) {
	return "\033[0;96m"
}

func GetCodeBrightGreen() (code string) {
	return "\033[0;92m"
}

func GetCodeBrightMagenta() (code string) {
	return "\033[0;95m"
}

func GetCodeBrightRed() (code string) {
	return "\033[0;91m"
}

func GetCodeBrightWhite() (code string) {
	return "\033[0;97m"
}

func GetCodeBrightYellow() (code string) {
	return "\033[0;93m"
}

func GetCodeCyan() (code string) {
	return "\033[0;36m"
}

func GetCodeGray() (code string) {
	return GetCodeBrightBlack()
}

func GetCodeGreen() (code string) {
	return "\033[0;32m"
}

func GetCodeMangenta() (code string) {
	return "\033[0;35m"
}

func GetCodeRed() (code string) {
	return "\033[0;31m"
}

func GetCodeWhite() (code string) {
	return "\033[0;37m"
}

func GetCodeYellow() (code string) {
	return "\033[0;33m"
}
