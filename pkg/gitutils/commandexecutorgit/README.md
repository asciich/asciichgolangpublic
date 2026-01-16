# commandexecutorgit


This is the GitRepository implementation based on a CommandExecutor (e.g. Bash, SSH...).
This means it's a wrapper around the "git" binary which needs to be available.
hile very inefficient this solution can manage git repository on remote hosts, inside containers...
which makes it very flexible.

When dealing with locally available repositories it's recommended to use the LocalGitRepository
implementation which uses go build in git functionality instead.