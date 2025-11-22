package gitlabutils

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/asciich/asciichgolangpublic"
	"github.com/asciich/asciichgolangpublic/pkg/files"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/filesoptions"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type DownloadMainReadmesOptions struct {
	// Where to store the README.md files:
	OuputPath string

	// Gitlab group of projects to download the REAMDE.md files:
	GitlabGroupUrl string

	// Ignore repos with no README.md file
	IgnoreNoReadmeMd bool
}

// Download the main README.md of all projects in a gitlab group
func DownloadMainReadmes(ctx context.Context, options *DownloadMainReadmesOptions) error {
	logging.LogInfoByCtxf(ctx, "Download main README.md files started.")

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	gitlab, err := NewAuthenticatedGitlab(ctx, options.GitlabGroupUrl)
	if err != nil {
		return err
	}

	group, err := gitlab.GetGroupByPath(ctx, options.GitlabGroupUrl)
	if err != nil {
		return err
	}

	projects, err := group.ListProjects(ctx, &asciichgolangpublic.GitlabListProjectsOptions{
		Recursive: true,
	})
	if err != nil {
		return err
	}

	outDir, err := files.GetLocalDirectoryByPath(options.OuputPath)
	if err != nil {
		return err
	}

	err = outDir.Create(ctx, &filesoptions.CreateOptions{})
	if err != nil {
		return err
	}

	type infoFileContent struct {
		BaseName          string `json:"basename"`
		URL               string `json:"url"`
		GitlabInstance    string `json:"gitlab_instance"`
		GitlabProjectPath string `json:"gitalb_project_path"`
		PathInGitlabRepo  string `json:"path_in_gitlab_repo"`
	}

	gitlabFqdn, err := gitlab.GetFqdn()
	if err != nil {
		return err
	}

	for _, p := range projects {
		url, err := p.GetProjectUrl(ctx)
		if err != nil {
			return err
		}

		logging.LogInfoByCtxf(ctx, "Collect main README.md of %s started.", url)

		projectPath, err := p.GetPath(ctx)
		if err != nil {
			return err
		}

		const pathInGitlab = "README.md"

		readme, err := p.GetFileInDefaultBranch(ctx, pathInGitlab)
		if err != nil {
			return err
		}

		readmeContent, err := readme.GetContentAsString(ctx)
		if err != nil {
			if options.IgnoreNoReadmeMd {
				if errors.Is(err, asciichgolangpublic.ErrGitlabRepositoryFileDoesNotExist) {
					logging.LogWarnByCtxf(ctx, "The repository %s has no README.md and is therefore ignored.", url)
					continue
				}
			}
			return err
		}

		o, err := outDir.CreateSubDirectory(ctx, projectPath, &filesoptions.CreateOptions{})
		if err != nil {
			return err
		}

		_, err = o.WriteStringToFile(ctx, "README.md", readmeContent, &filesoptions.WriteOptions{})
		if err != nil {
			return err
		}

		info := &infoFileContent{
			BaseName:          pathInGitlab,
			URL:               url,
			GitlabInstance:    gitlabFqdn,
			GitlabProjectPath: projectPath,
			PathInGitlabRepo:  pathInGitlab,
		}

		infoContent, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			return tracederrors.TracedErrorf("Failed to marshal info.json: %w", err)
		}

		_, err = o.WriteStringToFile(ctx, "info.json", string(infoContent), &filesoptions.WriteOptions{})
		if err != nil {
			return err
		}

		logging.LogInfoByCtxf(ctx, "Collect main README.md of %s finished.", url)
	}

	logging.LogInfoByCtxf(ctx, "Download main README.md files finished.")

	return nil
}
