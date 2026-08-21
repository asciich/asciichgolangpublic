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
- For filenames:
    - Do not use `_methods<.extension>` suffix. Instead of `add_methods.go` use `add.go`.

## Package Organization: Native and CommandExecutor Implementations

- For packages where functions interact with files, network, or external systems, provide **two implementation subpackages**:
    1. **`native<packagename>`**: Implements the logic using native Go libraries. Works only on the local machine.
    2. **`commandexecutor<packagename>`**: Implements the logic using the `commandexecutor` in combination with shell commands. This allows execution on remote machines (e.g. over SSH).
- **Both subpackages must implement the same set of functions** so the functionality stays the same regardless of whether it is executed locally or on a remote machine.
- The `commandexecutor...` functions take an additional `commandexecutor` parameter compared to their `native...` counterparts.

### Naming Conventions

- The native subpackage is named `native<packagename>` (e.g. `nativex509utils`, `nativetruststoreutils`).
- The commandexecutor subpackage is named `commandexecutor<packagename>` (e.g. `commandexecutorx509utils`, `commandexecutortruststoreutils`).
- The parent package (e.g. `x509utils`) provides **convenience functions** that delegate to the `native...` implementation. This allows callers who only need local execution to use the simpler API without explicitly choosing an implementation.

### Adding a New Function

1. Implement the function in **both** `native<packagename>` and `commandexecutor<packagename>`.
2. Add a convenience wrapper in the parent package that calls the `native...` implementation.
3. Add an `Example_<FunctionName>_test.go` in the parent package containing a well-commented example test for documentation purposes.

### Generic Functions

- Functions that operate purely on in-memory data (e.g. `GetCommonName(cert *x509.Certificate) (string, error)`) and do **not** require file or network access are placed in a dedicated `generic<packagename>` subpackage (e.g. `genericx509utils`).
- These functions do not need dual implementations since they have no I/O dependencies.

### Object-Oriented Packages (`oo` suffix)

- For packages that provide procedural/functional APIs (functions taking explicit parameters), an **object-oriented wrapper package** can be provided with the suffix `oo`.
- Naming: `<packagename>oo` (e.g. `nativefilesoo`, `commandexecutorexecoo`, `commandexecutorbashoo`).
- The `oo` package defines **structs** that hold state (e.g. a file path, a connection) and provides **methods** on those structs that delegate to the underlying procedural package.
- The methods retrieve the necessary parameters from the struct's fields instead of requiring them as explicit function arguments.
- Example:
    - Procedural style in `nativefiles`:
        ```golang
        func Delete(ctx context.Context, pathToDelete string, options *filesoptions.DeleteOptions) error {
            // implementation
        }
        ```
    - Object-oriented style in `nativefilesoo`:
        ```golang
        func (f *File) Delete(ctx context.Context, options *filesoptions.DeleteOptions) error {
            path, err := f.GetPath()
            if err != nil {
                return err
            }
            return nativefiles.Delete(ctx, path, options)
        }
        ```
- The `oo` package **must not** reimplement logic. It only wraps calls to the underlying procedural package.
- This pattern applies to both `native...oo` and `commandexecutor...oo` packages.

### Package Structure Example

```
pkg/<domain>/<packagename>/
├── commandexecutor<packagename>/   # Shell-based implementation (works remotely)
├── native<packagename>/            # Native Go implementation (local only)
├── generic<packagename>/           # Pure in-memory helpers (no I/O)
├── <packagename>options/           # Options structs used by the functions
├── Example_<Function>_test.go      # Documentation examples (one per convenience function)
├── <packagename>_test.go           # Cross-implementation tests (loop over all impls)
├── <packagename>.go                # Convenience functions delegating to native impl
├── README.md
└──<packagename>.spec.md
```

### Cross-Implementation Testing

- Tests in the parent package directory validate that **all implementations behave identically**.
- This is achieved by defining an implementation struct (e.g. `x509Implementation`) that holds function references for each implementation variant.
- A `getX509Implementations()` function returns a slice of all implementations to test, including:
    - The convenience functions from the parent package itself.
    - The `native...` implementation.
    - The `commandexecutor...` implementations (one per supported executor, e.g. `exec` and `bash`).
- Every test iterates over all implementations and runs the same assertions for each:
    ```golang
    func Test_SomeFunction(t *testing.T) {
        implementations := getX509Implementations()

        for _, impl := range implementations {
            impl := impl

            t.Run(impl.Name+"_description of test case", func(t *testing.T) {
                // test logic using impl.SomeFunction(...)
            })
        }
    }
    ```
- Additionally, dedicated tests validate that all implementations return the **exact same result** for identical inputs.

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
                    _, ok = host.(*commandexecutorhost.CommandExecutorHost)
                    require.True(t, ok, "GetHostByHostname('example.com') must return a CommandExecutorHost")
                })
            }
            ```

## Container handing

- Avoid installing things in containers. Use container images which already include what is needed. Relying on the package servers is not a good practice.
