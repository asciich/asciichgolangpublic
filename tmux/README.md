# tmux

Helper functions to orchestrate [tmux](https://github.com/tmux/tmux/wiki).

An example usage to orchestrate the input to another binary can be found in the `TestTemuxWindow_WaitOutputMatchesRegex` defined in [TmuxWindow_test.go](TmuxWindow_test.go).

## For developers

Run tests:

```bash
bash -c "cd '$(git rev-parse --show-toplevel)' && go test -v ./tmux/..."
```
