package exitcodes

func ExitCodeOK() (exitCode int) {
	return 0
}

func ExitCodeTimeout() (exitCode int) {
	return 124
}
