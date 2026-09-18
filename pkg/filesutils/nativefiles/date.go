package nativefiles

import (
	"context"
	"time"

	"golang.org/x/sys/unix"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// GetBirthDate returns the birth (creation) time of the file or directory at
// the given path. Birth time is not exposed by the portable os.FileInfo API,
// so it is retrieved via the statx(2) syscall (STATX_BTIME). Note that not all
// filesystems report a birth time; in that case an error is returned.
func GetBirthDate(ctx context.Context, path string) (time.Time, error) {
	if path == "" {
		return time.Time{}, tracederrors.TracedErrorEmptyString("path")
	}

	var stat unix.Statx_t
	err := unix.Statx(unix.AT_FDCWD, path, 0, unix.STATX_BTIME, &stat)
	if err != nil {
		return time.Time{}, tracederrors.TracedErrorf("Failed to statx '%s' to get birth date: %w", path, err)
	}

	// Even when statx succeeds the filesystem may not provide a birth time.
	// The STATX_BTIME bit in the returned mask tells us whether Btime is valid.
	if stat.Mask&unix.STATX_BTIME == 0 {
		return time.Time{}, tracederrors.TracedErrorf(
			"Birth date is not available for '%s' on this filesystem.", path,
		)
	}

	birthDate := time.Unix(int64(stat.Btime.Sec), int64(stat.Btime.Nsec)).UTC()

	logging.LogInfoByCtxf(ctx, "Birth date of '%s' is '%v'.", path, birthDate)

	return birthDate, nil
}

// GetBirthDateRFC3339 returns the birth (creation) time of the file or
// directory at the given path formatted as an RFC3339 string.
func GetBirthDateRFC3339(ctx context.Context, path string) (string, error) {
	if path == "" {
		return "", tracederrors.TracedErrorEmptyString("path")
	}

	birthDate, err := GetBirthDate(ctx, path)
	if err != nil {
		return "", err
	}

	ret := birthDate.Format(time.RFC3339)

	logging.LogInfoByCtxf(ctx, "Birth date of '%s' in RFC3339 is '%s'.", path, ret)

	return ret, nil
}
