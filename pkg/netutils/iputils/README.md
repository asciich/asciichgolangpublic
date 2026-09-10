# iputils

The `iputils` package provides utilities for IP address validation (IPv4 and IPv6).

## Functions

### Validation Functions (bool, error)

* **IsValidIP()**: Check if a string is a valid IP address (IPv4 or IPv6)
* **IsValidIPv4()**: Check if a string is a valid IPv4 address
* **IsValidIPv6()**: Check if a string is a valid IPv6 address

### Check Functions (error only)

* **CheckValidIP()**: Validate IP address, returns specific error types
* **CheckValidIPv4()**: Validate IPv4 address, returns specific error types
* **CheckValidIPv6()**: Validate IPv6 address, returns specific error types

### Error Types

* **ErrInvalidIP**: Generic invalid IP address error
* **ErrInvalidIPv4**: Invalid IPv4 address error
* **ErrInvalidIPv6**: Invalid IPv6 address error
* Empty string errors use `tracederrors.TracedErrorEmptyString("ip")`

### Error Helpers

* **IsInvalidIPError()**: Check if error is any IP validation error
* **IsInvalidIPv4Error()**: Check if error is IPv4 validation error
* **IsInvalidIPv6Error()**: Check if error is IPv6 validation error

## Specifications

For specifications see [iputils.spec.md](iputils.spec.md)

## Examples

```go
ctx := context.Background()

// Using IsValidIP (returns bool, error)
isValid, err := iputils.IsValidIP(ctx, "192.168.1.1")
// isValid = true, err = nil

isValid, err = iputils.IsValidIP(ctx, " 192.168.1.1 ")
// isValid = false, err = nil (whitespace makes it invalid)

isValid, err = iputils.IsValidIP(ctx, "")
// isValid = false, err = tracederrors.TracedErrorEmptyString("ip")

// Using CheckValidIP (returns error only)
err := iputils.CheckValidIP(ctx, "192.168.1.1")
// err = nil (valid)

err = iputils.CheckValidIP(ctx, "invalid")
// err = iputils.ErrInvalidIP

err = iputils.CheckValidIP(ctx, "")
// err = tracederrors.TracedErrorEmptyString("ip")

// Using error helpers
if iputils.IsInvalidIPError(err) {
    // Handle any IP validation error
}

if iputils.IsInvalidIPv4Error(err) {
    // Handle specifically IPv4 error
}

// Check for empty string error
if errors.Is(err, tracederrors.ErrTracedErrorEmptyString) {
    // Handle empty string error
}

// CheckValidIPv4 and CheckValidIPv6 work similarly
err = iputils.CheckValidIPv4(ctx, "192.168.1.1")
// err = nil

err = iputils.CheckValidIPv4(ctx, "::1")
// err = iputils.ErrInvalidIPv4 (this is IPv6, not IPv4)

err = iputils.CheckValidIPv6(ctx, "2001:db8::1")
// err = nil

err = iputils.CheckValidIPv6(ctx, "192.168.1.1")
// err = iputils.ErrInvalidIPv6 (this is IPv4, not IPv6)
```
