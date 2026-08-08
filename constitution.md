# constitution for asciichgolangpublic

## Implementation

- For execution on localhost the `native...` are used. Native in this case means natively implemented in golang without external `exec`
- For working with jumphosts the CommandExecutor approach is used. Instead of using go native libraries shell commands are used.
    - Do not relay on the the exit code 1 to check if a resource exists:
        - Instead of `ls abc.txt` where return code 0 means `abc.txt` exists and all others mean either it does not exists or there was an error...
        - Use a `sh` like `sh -c 'ls abc.txt &>/dev/null && echo yes || echo no'` and check for `yes`, `no` or other values (which indicate an error). This way we are sure the execution itself was performed.

## Documentation

- All `Example_*_test.go` files are meant for documentation. They must all be listed and linked in the corresponding `README.md`. This is done as a bullet point list. 
