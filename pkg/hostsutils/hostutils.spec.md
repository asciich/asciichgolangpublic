# hostutils specifications

## Implementation

- `GetHostByHostname("localhost")` must return a `NativeHost` while all other hostnames return a `CommandExecutorHost`.
    - This must be validated in a unittest.
 
## Testing

- Tests in the `hostutils` package are meant to test both implementations behave the same way.
    - Each test therefore uses a for loop to run the same test for both implementations available.
