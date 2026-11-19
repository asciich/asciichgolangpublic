# repoinfoutils

Handle repository infos.
The infos are mainly defined in the `.repoinfo.yml` at the toplevel of the repository.

## Full `.repoinfo.yml` documentation

All options are optional. Only fill in when required.

```yaml
---
repository_info:
    # Hostname of the GitLab instance hosting the repository.
    gitlab_instance: "gitlab.com"

    # SSH port of the GitLab instance. If unset the default port 2222 is used.
    gitlab_port: 2222

    # Path of the repository on the GitLab instance.
    # The '.git' suffix is not required. If it is added by mistake it is removed.
    gitlab_path: "<project>/<subproject..>/<repo_name>"
```

## Functions

- `GetGitlabInstance`: Returns the value of `gitlab_instance`.
- `GetGitlabPort`: Returns the value of `gitlab_port`.
- `GetGitlabPortOrDefaultIfUnset`: Returns the value of `gitlab_port` or the default port `2222` if it is not set.
- `GetGitlabPath`: Returns the value of `gitlab_path` without the `.git` suffix.
- `GetSshRemoteUrl`: Returns the SSH remote URL in the form of `ssh://git@<gitlab_instance>/<gitlab_path>.git`.

All functions accept either the path to the `.repoinfo.yml` file itself or the path to the directory containing it.
If a directory is given `.repoinfo.yml` is used. If it does not exist `.repoinfo.yaml` is used.

## Examples

- [Example_GetGitlabInstance_test.go](Example_GetGitlabInstance_test.go)
- [Example_GetGitlabPort_test.go](Example_GetGitlabPort_test.go)
- [Example_GetGitlabPortOrDefaultIfUnset_test.go](Example_GetGitlabPortOrDefaultIfUnset_test.go)
- [Example_GetGitlabPath_test.go](Example_GetGitlabPath_test.go)
- [Example_GetSshRemoteUrl_test.go](Example_GetSshRemoteUrl_test.go)

## Specification

See [repoinfoutils.spec.md](repoinfoutils.spec.md).
