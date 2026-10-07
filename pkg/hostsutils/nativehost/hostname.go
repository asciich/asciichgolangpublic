//go:build linux

package nativehost

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorexecoo"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/commandexecutorhostsutils"
	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

const (
	etcHostnamePath   = "/etc/hostname"
	maxHostnameLength = 64 // HOST_NAME_MAX on Linux
	maxLabelLength    = 63
)

// SetHostName sets the OS hostname of the local machine to the given value
// using native Go (no external commands).
//
// Both the static hostname (persisted in /etc/hostname) and the transient
// hostname (kernel, via sethostname(2)) are set, matching the behavior of
// "hostnamectl set-hostname". The operation is skipped when both are already
// set to the requested value (idempotent).
func (n *NativeHost) SetHostName(ctx context.Context, hostname string, options *hostsutilsoptions.SetHostnameOptions) error {
	if options == nil {
		options = &hostsutilsoptions.SetHostnameOptions{}
	}

	if options.UseSudo {
		exec := commandexecutorexecoo.Exec()
		return commandexecutorhostsutils.SetHostName(ctx, exec, hostname, options)
	}

	return n.setHostnameNatively(ctx, hostname, options)
}

func (n *NativeHost) setHostnameNatively(ctx context.Context, hostname string, options *hostsutilsoptions.SetHostnameOptions) error {
	if hostname == "" {
		return tracederrors.TracedErrorEmptyString("hostname")
	}

	err := validateHostname(hostname)
	if err != nil {
		return err
	}

	if options == nil {
		options = &hostsutilsoptions.SetHostnameOptions{}
	}

	hostDescription, err := n.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Set hostname to '%s' on '%s' started.", hostname, hostDescription)

	currentStaticHostname, err := readStaticHostname()
	if err != nil {
		return err
	}

	currentTransientHostname, err := os.Hostname()
	if err != nil {
		return tracederrors.TracedErrorf("Unable to get transient hostname: %w", err)
	}
	currentTransientHostname = strings.TrimSpace(currentTransientHostname)

	staticUpToDate := currentStaticHostname == hostname
	transientUpToDate := currentTransientHostname == hostname

	if staticUpToDate && transientUpToDate {
		logging.LogInfoByCtxf(ctx, "Hostname on '%s' is already set to '%s'. Skip setting hostname.", hostDescription, hostname)
		logging.LogInfoByCtxf(ctx, "Set hostname to '%s' on '%s' finished.", hostname, hostDescription)
		return nil
	}

	if os.Geteuid() != 0 {
		return tracederrors.TracedErrorf(
			"Setting hostname on '%s' requires root privileges (current euid=%d).",
			hostDescription,
			os.Geteuid(),
		)
	}

	if !staticUpToDate {
		if err = writeStaticHostname(hostname); err != nil {
			return err
		}
		logging.LogChangedByCtxf(ctx, "Set static hostname from '%s' to '%s' on '%s'.", currentStaticHostname, hostname, hostDescription)
	}

	if !transientUpToDate {
		if err = syscall.Sethostname([]byte(hostname)); err != nil {
			return tracederrors.TracedErrorf("Unable to set transient hostname to '%s': %w", hostname, err)
		}
		logging.LogChangedByCtxf(ctx, "Set transient hostname from '%s' to '%s' on '%s'.", currentTransientHostname, hostname, hostDescription)
	}

	logging.LogInfoByCtxf(ctx, "Set hostname to '%s' on '%s' finished.", hostname, hostDescription)

	return nil
}

// validateHostname checks the hostname against RFC 1123 rules and the Linux
// HOST_NAME_MAX limit, similar to the validation done by hostnamectl.
func validateHostname(hostname string) error {
	if len(hostname) > maxHostnameLength {
		return tracederrors.TracedErrorf(
			"Invalid hostname '%s': length %d exceeds maximum of %d characters.",
			hostname, len(hostname), maxHostnameLength,
		)
	}

	for _, label := range strings.Split(hostname, ".") {
		if label == "" {
			return tracederrors.TracedErrorf("Invalid hostname '%s': empty label (leading, trailing or double dot).", hostname)
		}

		if len(label) > maxLabelLength {
			return tracederrors.TracedErrorf("Invalid hostname '%s': label '%s' exceeds %d characters.", hostname, label, maxLabelLength)
		}

		if label[0] == '-' || label[len(label)-1] == '-' {
			return tracederrors.TracedErrorf("Invalid hostname '%s': label '%s' must not start or end with '-'.", hostname, label)
		}

		for _, r := range label {
			isValid := (r >= 'a' && r <= 'z') ||
				(r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') ||
				r == '-'
			if !isValid {
				return tracederrors.TracedErrorf("Invalid hostname '%s': invalid character '%c'.", hostname, r)
			}
		}
	}

	return nil
}

// readStaticHostname reads the static hostname from /etc/hostname.
// Returns an empty string if the file does not exist or contains no hostname.
// Empty lines and comment lines (starting with '#') are ignored.
func readStaticHostname() (string, error) {
	content, err := os.ReadFile(etcHostnamePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", tracederrors.TracedErrorf("Unable to read '%s': %w", etcHostnamePath, err)
	}

	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		return line, nil
	}

	return "", nil
}

// writeStaticHostname atomically writes the hostname to /etc/hostname
// (temp file + fsync + rename + directory fsync).
func writeStaticHostname(hostname string) (err error) {
	dir := filepath.Dir(etcHostnamePath)

	tmpFile, err := os.CreateTemp(dir, ".hostname-*")
	if err != nil {
		return tracederrors.TracedErrorf("Unable to create temporary file in '%s': %w", dir, err)
	}
	tmpPath := tmpFile.Name()

	// Clean up the temp file on any failure.
	defer func() {
		if err != nil {
			_ = tmpFile.Close()
			_ = os.Remove(tmpPath)
		}
	}()

	if _, err = tmpFile.WriteString(hostname + "\n"); err != nil {
		return tracederrors.TracedErrorf("Unable to write temporary hostname file '%s': %w", tmpPath, err)
	}

	if err = tmpFile.Chmod(0o644); err != nil {
		return tracederrors.TracedErrorf("Unable to chmod temporary hostname file '%s': %w", tmpPath, err)
	}

	if err = tmpFile.Sync(); err != nil {
		return tracederrors.TracedErrorf("Unable to sync temporary hostname file '%s': %w", tmpPath, err)
	}

	if err = tmpFile.Close(); err != nil {
		return tracederrors.TracedErrorf("Unable to close temporary hostname file '%s': %w", tmpPath, err)
	}

	if err = os.Rename(tmpPath, etcHostnamePath); err != nil {
		return tracederrors.TracedErrorf("Unable to rename '%s' to '%s': %w", tmpPath, etcHostnamePath, err)
	}

	// Persist the rename itself.
	dirHandle, err := os.Open(dir)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to open directory '%s' for sync: %w", dir, err)
	}
	defer dirHandle.Close()

	if err = dirHandle.Sync(); err != nil {
		return tracederrors.TracedErrorf("Unable to sync directory '%s': %w", dir, err)
	}

	return nil
}
