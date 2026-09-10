package iputils

import "errors"

// Error types for IP validation
var (
    ErrInvalidIP      = errors.New("invalid IP address")
    ErrInvalidIPv4    = errors.New("invalid IPv4 address")
    ErrInvalidIPv6    = errors.New("invalid IPv6 address")
)

// IsInvalidIPError returns true if the error is an IP validation error (IPv4 or IPv6).
func IsInvalidIPError(err error) bool {
    return errors.Is(err, ErrInvalidIP) || errors.Is(err, ErrInvalidIPv4) || errors.Is(err, ErrInvalidIPv6)
}

// IsInvalidIPv4Error returns true if the error is an ErrInvalidIPv4 error.
func IsInvalidIPv4Error(err error) bool {
    return errors.Is(err, ErrInvalidIPv4)
}

// IsInvalidIPv6Error returns true if the error is an ErrInvalidIPv6 error.
func IsInvalidIPv6Error(err error) bool {
    return errors.Is(err, ErrInvalidIPv6)
}
