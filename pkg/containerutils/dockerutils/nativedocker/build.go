package nativedocker

import (
	"archive/tar"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/moby/moby/client"

	"github.com/asciich/asciichgolangpublic/pkg/containerutils/containerinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/containerutils/dockerutils/dockeroptions"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/nativefiles"
	"github.com/asciich/asciichgolangpublic/pkg/filesutils/tempfiles"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// BuildImage builds a Docker image using the native Docker API.
func (d *Docker) BuildImage(ctx context.Context, options *dockeroptions.BuildContainerOptions) (containerinterfaces.Image, error) {
	if options == nil {
		return nil, tracederrors.TracedErrorNil("options")
	}

	if err := options.Validate(); err != nil {
		return nil, err
	}

	logging.LogInfoByCtxf(ctx, "Building Docker image '%s' started.", options.ImageNameAndTag)

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, tracederrors.TracedErrorf("unable to create docker client: %w", err)
	}
	defer cli.Close()

	// Prepare build options
	buildOptions := client.ImageBuildOptions{
		Tags:        []string{options.ImageNameAndTag},
		NoCache:     options.NoCache,
		PullParent:  options.PullParentImages,
		Remove:      true,
		ForceRemove: true,
	}

	// Convert build args to the required format (map[string]*string)
	if options.AdditionalBuildArgs != nil {
		buildOptions.BuildArgs = make(map[string]*string)
		for k, v := range options.AdditionalBuildArgs {
			value := v
			buildOptions.BuildArgs[k] = &value
		}
	}

	// Create a tar archive containing the build context
	var buildContextPath string
	var cleanupFunc func() error

	if options.DockerfilePath != "" {
		dockerfilePath := options.DockerfilePath
		buildContextPath = options.BuildContextPath

		if buildContextPath == "" {
			buildContextPath = filepath.Dir(dockerfilePath)
		}

		dockerfileName := filepath.Base(dockerfilePath)
		if dockerfileName != "Dockerfile" || filepath.Dir(dockerfilePath) != buildContextPath {
			tempDir, err := tempfiles.CreateTempDir(ctx)
			if err != nil {
				return nil, tracederrors.TracedErrorf("Failed to create temp directory: %w", err)
			}

			content, err := os.ReadFile(dockerfilePath)
			if err != nil {
				return nil, tracederrors.TracedErrorf("Failed to read Dockerfile: %w", err)
			}

			tempDockerfilePath := filepath.Join(tempDir, "Dockerfile")
			if err := nativefiles.WriteString(ctx, tempDockerfilePath, string(content)); err != nil {
				return nil, tracederrors.TracedErrorf("Failed to write Dockerfile to temp: %w", err)
			}

			buildContextPath = tempDir
			buildOptions.Dockerfile = "Dockerfile"

			cleanupFunc = func() error {
				return nativefiles.Delete(ctx, tempDir, nil)
			}
		} else {
			buildOptions.Dockerfile = "Dockerfile"
		}
	} else if options.DockerfileContent != "" {
		tempDir, err := tempfiles.CreateTempDir(ctx)
		if err != nil {
			return nil, tracederrors.TracedErrorf("Failed to create temp directory: %w", err)
		}

		dockerfilePath := filepath.Join(tempDir, "Dockerfile")
		if err := nativefiles.WriteString(ctx, dockerfilePath, options.DockerfileContent); err != nil {
			return nil, tracederrors.TracedErrorf("Failed to write Dockerfile: %w", err)
		}

		buildContextPath = tempDir
		buildOptions.Dockerfile = "Dockerfile"

		cleanupFunc = func() error {
			return nativefiles.Delete(ctx, tempDir, nil)
		}
	}

	if buildContextPath == "" {
		buildContextPath = "."
	}

	tarReader, err := createTarArchive(buildContextPath)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create tar archive: %w", err)
	}
	defer tarReader.Close()

	response, err := cli.ImageBuild(ctx, tarReader, buildOptions)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to build image '%s': %w", options.ImageNameAndTag, err)
	}
	defer response.Body.Close()

	buildOutput, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to read build output: %w", err)
	}

	logging.LogInfoByCtxf(ctx, "Build output: %s", string(buildOutput))

	if cleanupFunc != nil {
		if cleanupErr := cleanupFunc(); cleanupErr != nil {
			logging.LogErrorByCtxf(ctx, "Failed to cleanup temp directory: %v", cleanupErr)
		}
	}

	logging.LogChangedByCtxf(ctx, "Built Docker image '%s'.", options.ImageNameAndTag)

	return d.GetImageByName(options.ImageNameAndTag)
}

func createTarArchive(dir string) (io.ReadCloser, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if _, err := tw.Write(content); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return io.NopCloser(buf), nil
}
