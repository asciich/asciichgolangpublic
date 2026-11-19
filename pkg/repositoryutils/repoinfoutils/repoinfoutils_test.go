package repoinfoutils_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/repositoryutils/repoinfoutils"
)

const fullRepoInfo = `---
repository_info:
  gitlab_instance: "gitlab.example.com"
  gitlab_port: 2200
  gitlab_path: "tools/my_repo"
`

// Writes the given content as fileName into a new temporary directory and returns the directory.
func createRepoInfoDir(t *testing.T, fileName string, content string) string {
	t.Helper()

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, fileName), []byte(content), 0o644))

	return dir
}

func Test_Rules_RepoInfoPath(t *testing.T) {
	ctx := contextutils.ContextSilent()

	t.Run("directory with .repoinfo.yml", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", fullRepoInfo)

		instance, err := repoinfoutils.GetGitlabInstance(ctx, dir)
		require.NoError(t, err)
		require.EqualValues(t, "gitlab.example.com", instance)
	})

	t.Run("complete path to the file", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", fullRepoInfo)

		instance, err := repoinfoutils.GetGitlabInstance(ctx, filepath.Join(dir, ".repoinfo.yml"))
		require.NoError(t, err)
		require.EqualValues(t, "gitlab.example.com", instance)
	})

	t.Run("directory with .repoinfo.yaml only", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yaml", fullRepoInfo)

		instance, err := repoinfoutils.GetGitlabInstance(ctx, dir)
		require.NoError(t, err)
		require.EqualValues(t, "gitlab.example.com", instance)
	})

	t.Run(".repoinfo.yml is preferred over .repoinfo.yaml", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", fullRepoInfo)
		require.NoError(t, os.WriteFile(
			filepath.Join(dir, ".repoinfo.yaml"),
			[]byte("repository_info:\n  gitlab_instance: \"other.example.com\"\n"),
			0o644,
		))

		instance, err := repoinfoutils.GetGitlabInstance(ctx, dir)
		require.NoError(t, err)
		require.EqualValues(t, "gitlab.example.com", instance)
	})

	t.Run("additional unknown options are ignored", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", fullRepoInfo+"  gitlab_public_repository: true\nrelease_type: \"golang-package\"\n")

		instance, err := repoinfoutils.GetGitlabInstance(ctx, dir)
		require.NoError(t, err)
		require.EqualValues(t, "gitlab.example.com", instance)
	})

	t.Run("empty repoInfoPath", func(t *testing.T) {
		_, err := repoinfoutils.GetGitlabInstance(ctx, "")
		require.Error(t, err)
	})

	t.Run("not existing repoInfoPath", func(t *testing.T) {
		_, err := repoinfoutils.GetGitlabInstance(ctx, filepath.Join(t.TempDir(), "does-not-exist"))
		require.Error(t, err)
	})

	t.Run("directory without repo info file", func(t *testing.T) {
		_, err := repoinfoutils.GetGitlabInstance(ctx, t.TempDir())
		require.Error(t, err)
	})

	t.Run("invalid yaml", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", "repository_info: [this is: not valid")

		_, err := repoinfoutils.GetGitlabInstance(ctx, dir)
		require.Error(t, err)
	})
}

func Test_GetGitlabInstance(t *testing.T) {
	ctx := contextutils.ContextSilent()

	t.Run("set", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", fullRepoInfo)

		instance, err := repoinfoutils.GetGitlabInstance(ctx, dir)
		require.NoError(t, err)
		require.EqualValues(t, "gitlab.example.com", instance)
	})

	t.Run("unset", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", "repository_info:\n  gitlab_path: \"a/b\"\n")

		_, err := repoinfoutils.GetGitlabInstance(ctx, dir)
		require.Error(t, err)
	})
}

func Test_GetGitlabPort(t *testing.T) {
	ctx := contextutils.ContextSilent()

	t.Run("set", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", fullRepoInfo)

		port, err := repoinfoutils.GetGitlabPort(ctx, dir)
		require.NoError(t, err)
		require.EqualValues(t, 2200, port)
	})

	t.Run("unset", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", "repository_info:\n  gitlab_instance: \"gitlab.com\"\n")

		_, err := repoinfoutils.GetGitlabPort(ctx, dir)
		require.Error(t, err)
	})

	t.Run("port zero is invalid", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", "repository_info:\n  gitlab_port: 0\n")

		_, err := repoinfoutils.GetGitlabPort(ctx, dir)
		require.Error(t, err)
	})

	t.Run("port too high is invalid", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", "repository_info:\n  gitlab_port: 65536\n")

		_, err := repoinfoutils.GetGitlabPort(ctx, dir)
		require.Error(t, err)
	})

	t.Run("port is not a number", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", "repository_info:\n  gitlab_port: \"abc\"\n")

		_, err := repoinfoutils.GetGitlabPort(ctx, dir)
		require.Error(t, err)
	})
}

func Test_GetGitlabPortOrDefaultIfUnset(t *testing.T) {
	ctx := contextutils.ContextSilent()

	t.Run("set", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", fullRepoInfo)

		port, err := repoinfoutils.GetGitlabPortOrDefaultIfUnset(ctx, dir)
		require.NoError(t, err)
		require.EqualValues(t, 2200, port)
	})

	t.Run("unset returns default", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", "repository_info:\n  gitlab_instance: \"gitlab.com\"\n")

		port, err := repoinfoutils.GetGitlabPortOrDefaultIfUnset(ctx, dir)
		require.NoError(t, err)
		require.EqualValues(t, 2222, port)
	})

	t.Run("empty file returns default", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", "")

		port, err := repoinfoutils.GetGitlabPortOrDefaultIfUnset(ctx, dir)
		require.NoError(t, err)
		require.EqualValues(t, repoinfoutils.DefaultGitlabPort, port)
	})

	t.Run("invalid port is an error and not replaced by default", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", "repository_info:\n  gitlab_port: 0\n")

		_, err := repoinfoutils.GetGitlabPortOrDefaultIfUnset(ctx, dir)
		require.Error(t, err)
	})

	t.Run("not existing path", func(t *testing.T) {
		_, err := repoinfoutils.GetGitlabPortOrDefaultIfUnset(ctx, filepath.Join(t.TempDir(), "does-not-exist"))
		require.Error(t, err)
	})
}

func Test_GetGitlabPath(t *testing.T) {
	ctx := contextutils.ContextSilent()

	t.Run("set", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", fullRepoInfo)

		gitlabPath, err := repoinfoutils.GetGitlabPath(ctx, dir)
		require.NoError(t, err)
		require.EqualValues(t, "tools/my_repo", gitlabPath)
	})

	t.Run(".git suffix is removed", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", "repository_info:\n  gitlab_path: \"tools/my_repo.git\"\n")

		gitlabPath, err := repoinfoutils.GetGitlabPath(ctx, dir)
		require.NoError(t, err)
		require.EqualValues(t, "tools/my_repo", gitlabPath)
	})

	t.Run("unset", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", "repository_info:\n  gitlab_instance: \"gitlab.com\"\n")

		_, err := repoinfoutils.GetGitlabPath(ctx, dir)
		require.Error(t, err)
	})

	t.Run("only .git suffix", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", "repository_info:\n  gitlab_path: \".git\"\n")

		_, err := repoinfoutils.GetGitlabPath(ctx, dir)
		require.Error(t, err)
	})
}

func Test_GetSshRemoteUrl(t *testing.T) {
	ctx := contextutils.ContextSilent()

	t.Run("set", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", fullRepoInfo)

		url, err := repoinfoutils.GetSshRemoteUrl(ctx, dir)
		require.NoError(t, err)
		require.EqualValues(t, "ssh://git@gitlab.example.com/tools/my_repo.git", url)
	})

	t.Run(".git suffix in gitlab_path is not duplicated", func(t *testing.T) {
		dir := createRepoInfoDir(
			t,
			".repoinfo.yml",
			"repository_info:\n  gitlab_instance: \"gitlab.com\"\n  gitlab_path: \"tools/my_repo.git\"\n",
		)

		url, err := repoinfoutils.GetSshRemoteUrl(ctx, dir)
		require.NoError(t, err)
		require.EqualValues(t, "ssh://git@gitlab.com/tools/my_repo.git", url)
	})

	t.Run("instance unset", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", "repository_info:\n  gitlab_path: \"tools/my_repo\"\n")

		_, err := repoinfoutils.GetSshRemoteUrl(ctx, dir)
		require.Error(t, err)
	})

	t.Run("path unset", func(t *testing.T) {
		dir := createRepoInfoDir(t, ".repoinfo.yml", "repository_info:\n  gitlab_instance: \"gitlab.com\"\n")

		_, err := repoinfoutils.GetSshRemoteUrl(ctx, dir)
		require.Error(t, err)
	})

	t.Run("empty repoInfoPath", func(t *testing.T) {
		_, err := repoinfoutils.GetSshRemoteUrl(ctx, "")
		require.Error(t, err)
	})
}
