package filesgeneric

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrFileNotFound = errors.New("file not found")

// fileNotFoundMessageIndicators are lowercase message parts that indicate a missing file.
// "no such file or directory" covers GNU coreutils (e.g. "stat: cannot statx '...': No such file or directory")
// and BusyBox (e.g. "stat: can't stat '...': No such file or directory").
var fileNotFoundMessageIndicators = []string{
	"no such file or directory",
}

// errorMatchers are applied in order by GetAsError.
// Each matcher must return 'err' untouched if it does not match.
//
// To support a new error type:
//  1. Add a sentinel error (e.g. ErrPermissionDenied).
//  2. Add an IsErr... and a GetAs...ErrorIfMessageMatches function.
//  3. Register the GetAs...ErrorIfMessageMatches function here.
var errorMatchers = []func(err error) error{
	GetAsFileNotFoundErrorIfMessageMatches,
}

// GetAsError applies all known error matchers to 'err' and wraps it
// with every matching sentinel error (e.g. ErrFileNotFound).
//
// - nil is returned as nil.
// - Errors not matching any matcher are returned untouched.
// - The original error stays in the chain, so errors.Is(result, err) remains true.
func GetAsError(err error) error {
	if err == nil {
		return nil
	}

	for _, matcher := range errorMatchers {
		err = matcher(err)
	}

	return err
}

func IsErrFileNotFound(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrFileNotFound) {
		return true
	}

	return errors.Is(err, os.ErrNotExist)
}

// GetAsFileNotFoundErrorIfMessageMatches wraps 'err' so IsErrFileNotFound returns true
// if its message indicates a missing file (e.g. "No such file or directory" from a failed command).
//
// - nil is returned as nil.
// - Errors already detected by IsErrFileNotFound are returned untouched.
// - Errors not matching are returned untouched.
// - The original error stays in the chain, so errors.Is(result, err) remains true.
func GetAsFileNotFoundErrorIfMessageMatches(err error) error {
	if err == nil {
		return nil
	}

	if IsErrFileNotFound(err) {
		return err
	}

	return wrapIfMessageContainsAny(err, ErrFileNotFound, fileNotFoundMessageIndicators)
}

// wrapIfMessageContainsAny wraps 'err' with 'sentinel' if the error message contains
// any of the given lowercase 'indicators' (case insensitive).
// Otherwise 'err' is returned untouched.
func wrapIfMessageContainsAny(err error, sentinel error, indicators []string) error {
	if err == nil {
		return nil
	}

	msg := strings.ToLower(err.Error())
	for _, indicator := range indicators {
		if strings.Contains(msg, indicator) {
			return fmt.Errorf("%w: %w", sentinel, err)
		}
	}

	return err
}
