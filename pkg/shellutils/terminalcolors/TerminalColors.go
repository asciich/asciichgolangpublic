package terminalcolors

// Color codes for terminal.
//
// Source: https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
// Source: https://en.wikipedia.org/wiki/ANSI_escape_code

// This color code resets all color settings back to default.
const CODE_NO_COLOR = "\033[0m"

const CODE_BLACK = "\033[0;30m"
const CODE_BLUE = "\033[0;34m"
const CODE_BRIGHT_BLACK = "\033[0;90m"
const CODE_BRIGHT_BLUE = "\033[0;94m"
const CODE_BRIGHT_CYAN = "\033[0;96m"
const CODE_BRIGHT_GREEN = "\033[0;92m"
const CODE_BRIGHT_MAGENTA = "\033[0;95m"
const CODE_BRIGHT_RED = "\033[0;91m"
const CODE_BRIGHT_WHITE = "\033[0;97m"
const CODE_BRIGHT_YELLOW = "\033[0;93m"
const CODE_CYAN = "\033[0;36m"
const CODE_GREEN = "\033[0;32m"
const CODE_MANGENTA = "\033[0;35m"
const CODE_RED= "\033[0;31m"
const CODE_WHITE="\033[0;37m"
const CODE_YELLOW = "\033[0;33m"

func GetCodeGray() string {
	return CODE_BRIGHT_BLACK
}
