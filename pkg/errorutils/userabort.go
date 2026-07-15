package errorutils

import "strings"

// isUserAbort returns true if the error represents a user-initiated abort.
func IsUserAbort(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "abort") ||
		strings.Contains(msg, "interrupt") ||
		strings.Contains(msg, "no item")
}
