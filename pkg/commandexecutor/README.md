# commandexector

CommandExecutors have an important role in this library since they allow us to:
- Execute commands on remote systems like (jump-)hosts reachable over the SSH.
- Execute the same command on the local machine, inside a container (no matter if docker, kubernetes, ...)
- Can be used to short cut the amount of development time since we already know how to do things with redirects in bash or a specific tool...

But it has downsides:
- It's not real programming, it's abusing golang for scripting automation.
- It's a security risk. Calling exec (especially with unchecked user input as parameter) leads to security issues.

## Avoid exec calls.

To avoid exec calls on the local machine set the env var accordingly:
```bash
export ASCIICHGOLANGPUBLIC_AVOID_EXEC=1
```

## For developers

To run all tests use:
```bash
bash -c "cd commandexecutor && go test -v ./..."
```