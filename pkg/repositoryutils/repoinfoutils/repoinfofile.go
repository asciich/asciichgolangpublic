package repoinfoutils

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
	"gopkg.in/yaml.v3"
)

// Names of the repo info file in the order they are looked up in a directory.
var repoInfoFileNames = []string{
	".repoinfo.yml",
	".repoinfo.yaml",
}

// Content of the `.repoinfo.yml` file.
type repoInfoFile struct {
	RepositoryInfo repositoryInfo `yaml:"repository_info"`
}

// The `repository_info` section of the `.repoinfo.yml` file.
//
// Optional numeric values are pointers to be able to distinguish between "unset" and the zero value.
type repositoryInfo struct {
	GitlabInstance string `yaml:"gitlab_instance"`
	GitlabPort     *int   `yaml:"gitlab_port"`
	GitlabPath     string `yaml:"gitlab_path"`
}

// Resolves the path to the repo info file.
//
// repoInfoPath can either be the path to the file itself or the path to the directory containing the file.
func resolveRepoInfoFilePath(repoInfoPath string) (string, error) {
	if repoInfoPath == "" {
		return "", tracederrors.TracedErrorEmptyString("repoInfoPath")
	}

	info, err := os.Stat(repoInfoPath)
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to access repoInfoPath '%s': %w", repoInfoPath, err)
	}

	if !info.IsDir() {
		return repoInfoPath, nil
	}

	for _, name := range repoInfoFileNames {
		candidate := filepath.Join(repoInfoPath, name)

		candidateInfo, err := os.Stat(candidate)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}

			return "", tracederrors.TracedErrorf("Unable to access '%s': %w", candidate, err)
		}

		if candidateInfo.IsDir() {
			continue
		}

		return candidate, nil
	}

	return "", tracederrors.TracedErrorf(
		"No repo info file (one of %v) found in directory '%s'.",
		repoInfoFileNames,
		repoInfoPath,
	)
}

// Loads and parses the repo info file.
func loadRepositoryInfo(ctx context.Context, repoInfoPath string) (*repositoryInfo, error) {
	filePath, err := resolveRepoInfoFilePath(repoInfoPath)
	if err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Load repo info file '%s' started.", filePath)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to read repo info file '%s': %w", filePath, err)
	}

	parsed := repoInfoFile{}
	err = yaml.Unmarshal(content, &parsed)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to parse repo info file '%s': %w", filePath, err)
	}

	logging.LogInfoByCtxf(ctx, "Load repo info file '%s' finished.", filePath)

	return &parsed.RepositoryInfo, nil
}
