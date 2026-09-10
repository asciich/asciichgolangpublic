package iputils

import (
    "context"
    "net"

    "github.com/asciich/asciichgolangpublic/pkg/logging"
    "github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// IsValidIP checks if the given string is a valid IP address (IPv4 or IPv6).
// Returns true if the string is a valid IP address, false otherwise.
// Returns an error only for empty strings.
func IsValidIP(ctx context.Context, ip string) (bool, error) {
    if ip == "" {
        return false, tracederrors.TracedErrorEmptyString("ip")
    }

    isValid := net.ParseIP(ip) != nil

    if isValid {
        logging.LogInfoByCtxf(ctx, "IP address '%s' is valid.", ip)
    } else {
        logging.LogInfoByCtxf(ctx, "IP address '%s' is NOT valid.", ip)
    }

    return isValid, nil
}

// CheckValidIP checks if the given string is a valid IP address (IPv4 or IPv6).
// Returns nil if valid, ErrEmptyString for empty strings, or ErrInvalidIP for invalid IPs.
func CheckValidIP(ctx context.Context, ip string) error {
    if ip == "" {
        logging.LogInfoByCtxf(ctx, "IP address check failed: empty string.")
        return tracederrors.TracedErrorEmptyString("ip")
    }

    isValid := net.ParseIP(ip) != nil

    if isValid {
        logging.LogInfoByCtxf(ctx, "IP address '%s' is valid.", ip)
        return nil
    }

    logging.LogInfoByCtxf(ctx, "IP address '%s' is NOT valid.", ip)
    return ErrInvalidIP
}

// IsValidIPv4 checks if the given string is a valid IPv4 address.
// Returns true if the string is a valid IPv4 address, false otherwise.
// Returns an error only for empty strings.
func IsValidIPv4(ctx context.Context, ip string) (bool, error) {
    if ip == "" {
        return false, tracederrors.TracedErrorEmptyString("ip")
    }

    parsedIP := net.ParseIP(ip)
    isValid := parsedIP != nil && parsedIP.To4() != nil

    if isValid {
        logging.LogInfoByCtxf(ctx, "IPv4 address '%s' is valid.", ip)
    } else {
        logging.LogInfoByCtxf(ctx, "IPv4 address '%s' is NOT valid.", ip)
    }

    return isValid, nil
}

// CheckValidIPv4 checks if the given string is a valid IPv4 address.
// Returns nil if valid, ErrEmptyString for empty strings, or ErrInvalidIPv4 for invalid IPs.
func CheckValidIPv4(ctx context.Context, ip string) error {
    if ip == "" {
        logging.LogInfoByCtxf(ctx, "IPv4 address check failed: empty string.")
        return tracederrors.TracedErrorEmptyString("ip")
    }

    parsedIP := net.ParseIP(ip)

    if parsedIP != nil && parsedIP.To4() != nil {
        logging.LogInfoByCtxf(ctx, "IPv4 address '%s' is valid.", ip)
        return nil
    }

    logging.LogInfoByCtxf(ctx, "IPv4 address '%s' is NOT valid.", ip)
    return ErrInvalidIPv4
}

// IsValidIPv6 checks if the given string is a valid IPv6 address.
// Returns true if the string is a valid IPv6 address, false otherwise.
// Returns an error only for empty strings.
func IsValidIPv6(ctx context.Context, ip string) (bool, error) {
    if ip == "" {
        return false, tracederrors.TracedErrorEmptyString("ip")
    }

    parsedIP := net.ParseIP(ip)
    isValid := parsedIP != nil && parsedIP.To4() == nil

    if isValid {
        logging.LogInfoByCtxf(ctx, "IPv6 address '%s' is valid.", ip)
    } else {
        logging.LogInfoByCtxf(ctx, "IPv6 address '%s' is NOT valid.", ip)
    }

    return isValid, nil
}

// CheckValidIPv6 checks if the given string is a valid IPv6 address.
// Returns nil if valid, ErrEmptyString for empty strings, or ErrInvalidIPv6 for invalid IPs.
func CheckValidIPv6(ctx context.Context, ip string) error {
    if ip == "" {
        logging.LogInfoByCtxf(ctx, "IPv6 address check failed: empty string.")
        return tracederrors.TracedErrorEmptyString("ip")
    }

    parsedIP := net.ParseIP(ip)

    if parsedIP != nil && parsedIP.To4() == nil {
        logging.LogInfoByCtxf(ctx, "IPv6 address '%s' is valid.", ip)
        return nil
    }

    logging.LogInfoByCtxf(ctx, "IPv6 address '%s' is NOT valid.", ip)
    return ErrInvalidIPv6
}
