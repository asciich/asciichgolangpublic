# iputils specifications

This are the specifications for the [`iputils` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- Input strings are NOT trimmed - whitespace makes an IP address invalid.
- Empty string input returns `tracederrors.TracedErrorEmptyString("ip")`.
- Every IP validation attempt must be logged:
    - For valid IPs: Log which IP was validated as valid.
    - For invalid IPs: Log which IP was validated as invalid.

## Function Behavior

### IsValid* Functions

- Return `(bool, error)` tuple.
- Return `true, nil` for valid IPs.
- Return `false, tracederrors.TracedErrorEmptyString("ip")` for empty strings.
- Return `false, nil` for invalid IPs (non-empty).

### CheckValid* Functions

- Return `error` only.
- Return `nil` for valid IPs.
- Return `tracederrors.TracedErrorEmptyString("ip")` for empty strings.
- Return specific error (`ErrInvalidIP`, `ErrInvalidIPv4`, or `ErrInvalidIPv6`) for invalid IPs.

## Testing

- Validate correct IP addresses (IPv4 and IPv6) return `true`/`nil`.
- Validate invalid IP addresses return `false`/appropriate error.
- Validate empty strings return `tracederrors.TracedErrorEmptyString("ip")`.
- Validate strings with leading/trailing whitespace are rejected as invalid.
- Verify logging output matches validation attempts.
- Verify error helper functions (`IsInvalidIPError`, `IsInvalidIPv4Error`, `IsInvalidIPv6Error`) work correctly.
- Verify error wrapping compatibility with `errors.Is()`.
