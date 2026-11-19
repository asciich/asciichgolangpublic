# repoinfoutils specifications

## Rules

- Either the complete path to the file `.repoinfo.yml`, or the path to the directory containing it (which is usually the repo root) is accepted as `repoInfoPath`.
- If `repoInfoPath` points to a directory check for `.repoinfo.yml`. If this does not exist check for `.repoinfo.yaml`.

## Implementation

- `GetGitlabInstance(ctx context.Context, repoInfoPath string) (string, error)` returns the value of `gitlab_instance`.
- `GetGitlabPort(ctx context.Context, repoInfoPath string) (int, error)` returns the value of `gitlab_port`.
- `GetGitlabPortOrDefaultIfUnset(ctx context.Context, repoInfoPath string) (int, error)`:
    - Return the value of `GetGitlabPort` if `gitlab_port` is set.
    - Default port `2222` instead.
- `GetGitlabPath(ctx context.Context, repoInfoPath string) (string, error)` return the value of `gitlab_path`:
    - It must not include the `.git` suffix, even if it is provided by mistake in the file.
- `GetSshRemoteUrl(ctx context.Context, repoInfoPath string) (string, error)` return the remote URL in the form of `ssh://git@<GitlabInstance>/<GitlabPath>.git`.

## Testing

All functions must be fully unit tested.

## Documentation

The example `.repoinfo.yml` in the [README.md](README.md) must reflect all possible options.
