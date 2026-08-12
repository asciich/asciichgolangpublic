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

## Documentation

- All `Example_*_test.go` files are meant for documentation. They must all be listed and linked in the corresponding `README.md`. This is done as a bullet point list. 

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