# constitution for asciichgolangpublic

## Implementation

- For execution on localhost the `native...` are used. Native in this case means natively implemented in golang without external `exec`
- For working with jumphosts the CommandExecutor approach is used. Instead of using go native libraries shell commands are used.
    - Do not relay on the the exit code 1 to check if a resource exists:
        - Instead of `ls abc.txt` where return code 0 means `abc.txt` exists and all others mean either it does not exists or there was an error...
        - Use a `sh` like `sh -c 'ls abc.txt &>/dev/null && echo yes || echo no'` and check for `yes`, `no` or other values (which indicate an error). This way we are sure the execution itself was performed.
- The default variable name in small functions is `ret`:
    - Use:
        ```golang
        return ret, nil
        ```
    - Insead of:
        ```golang
        return toReturn, nil
        ```
- Use `defer` instead of cleanup calls at the end of a function. 
  This rule counts as well for test cases.
  Since all functions are implemented in an idempotent way it is safe to use `defer`:
    - Use:
        ```golang
        defer kubernetes.DeleteClusterRoleByName(ctx, roleName)
        ```
    - Instead of:
        ```golang
        // Cleanup
        mustutils.Must0(kubernetes.DeleteClusterRoleByName(ctx, roleName))
        }// end of the function
        ```
- Avoid creating `tracederrors` without assigning them to a variable or return them.
    - Example:
        - This is highly likely a mistake:
            ```golang
            tracederrors.TracedErrorf("ReadCloser process for command '%s' finished with error: %w", fullCommandJoined, err) // missing return statement in front of the line
            ```
        - And must be returned:
            ```golang
            return tracederrors.TracedErrorf("ReadCloser process for command '%s' finished with error: %w", fullCommandJoined, err) // correct
            ```


## Documentation

- All `Example_*_test.go` files are meant for documentation. They must all be listed and linked in the corresponding `README.md`. This is done as a bullet point list. 
- All package names mentioned in a markdown file must be a link to the README.md of the mentioned package.


## Testing

- In tests use the `require` package for condition checks.
- Use `require` to check for no error:
    - Use:
        ```golang
        require.NoError(t, err)
        ```
    - Instead of:
        ```golang
        if err != nil {
	    	panic(err)
	    }
        ```
- Unit Test Structure:
    - Use t.Run subtests for distinct assertions
        - When a single test function validates multiple distinct behaviors or cases, each case must be wrapped in its own t.Run subtest. This provides clearer test output, allows individual subtests to be run in isolation, and makes failures easier to identify.
        - Instead of:
            ```golang
            func TestHost_GetHostByHostnameReturnsCorrectType(t *testing.T) {
                host, err := hostsutils.GetHostByHostname("localhost")
                require.NoError(t, err)
                _, ok := host.(*nativehost.NativeHost)
                require.True(t, ok, "GetHostByHostname('localhost') must return a NativeHost")

                host, err = hostsutils.GetHostByHostname("example.com")
                require.NoError(t, err)
                _, ok = host.(*commandexecutorhost.CommandExecutorHost)
                require.True(t, ok, "GetHostByHostname('example.com') must return a CommandExecutorHost")
            }
            ```
        - Use:
            ```
            func TestHost_GetHostByHostnameReturnsCorrectType(t *testing.T) {
                t.Run("localhost must return NativeHost", func(t *testing.T) {
                    host, err := hostsutils.GetHostByHostname("localhost")
                    require.NoError(t, err)
                    _, ok := host.(*nativehost.NativeHost)
                    require.True(t, ok, "GetHostByHostname('localhost') must return a NativeHost")
                })

                t.Run("any other hostname must return a CommandExecutorHost", func(t *testing.T) {
                    host, err := hostsutils.GetHostByHostname("example.com")
                    require.NoError(t, err)
                    _, ok := host.(*commandexecutorhost.CommandExecutorHost)
                    require.True(t, ok, "GetHostByHostname('example.com') must return a CommandExecutorHost")
                })
            }
            ```

## Container handing

- Avoid installing things in containers. Use container images which already include what is needed. Relying on the package servers is not a good practice.
