package dockeroptions

import (
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

type BuildContainerOptions struct {
	// ImageNameAndTag specifies the name and tag for the built image (e.g., "myimage:latest")
	ImageNameAndTag string

	// DockerfilePath specifies the path to an existing Dockerfile
	// Either DockerfilePath or DockerfileContent must be set, but not both
	DockerfilePath string

	// DockerfileContent specifies the Dockerfile content directly
	// Either DockerfilePath or DockerfileContent must be set, but not both
	DockerfileContent string

	// BuildContextPath specifies the build context directory (defaults to current directory if not set)
	BuildContextPath string

	// AdditionalBuildArgs specifies build-time variables
	AdditionalBuildArgs map[string]string

	// NoCache forces rebuilding of all intermediate containers
	NoCache bool

	// PullParentImages attempts to pull the parent images during build
	PullParentImages bool
}

func NewBuildContainerOptions() *BuildContainerOptions {
	return new(BuildContainerOptions)
}

func (o *BuildContainerOptions) SetImageNameAndTag(imageNameAndTag string) error {
	if imageNameAndTag == "" {
		return tracederrors.TracedErrorf("imageNameAndTag is empty string")
	}
	o.ImageNameAndTag = imageNameAndTag
	return nil
}

func (o *BuildContainerOptions) SetDockerfilePath(dockerfilePath string) error {
	if dockerfilePath == "" {
		return tracederrors.TracedErrorf("dockerfilePath is empty string")
	}
	o.DockerfilePath = dockerfilePath
	return nil
}

func (o *BuildContainerOptions) SetDockerfileContent(dockerfileContent string) error {
	if dockerfileContent == "" {
		return tracederrors.TracedErrorf("dockerfileContent is empty string")
	}
	o.DockerfileContent = dockerfileContent
	return nil
}

func (o *BuildContainerOptions) SetBuildContextPath(buildContextPath string) error {
	o.BuildContextPath = buildContextPath
	return nil
}

func (o *BuildContainerOptions) SetAdditionalBuildArgs(args map[string]string) error {
	o.AdditionalBuildArgs = args
	return nil
}

func (o *BuildContainerOptions) SetNoCache(noCache bool) error {
	o.NoCache = noCache
	return nil
}

func (o *BuildContainerOptions) SetPullParentImages(pull bool) error {
	o.PullParentImages = pull
	return nil
}

func (o *BuildContainerOptions) Validate() error {
	if o.ImageNameAndTag == "" {
		return tracederrors.TracedErrorf("ImageNameAndTag is required")
	}

	// Either DockerfilePath or DockerfileContent must be set, but not both
	hasPath := o.DockerfilePath != ""
	hasContent := o.DockerfileContent != ""

	if !hasPath && !hasContent {
		return tracederrors.TracedErrorf("Either DockerfilePath or DockerfileContent must be set")
	}

	if hasPath && hasContent {
		return tracederrors.TracedErrorf("Only one of DockerfilePath or DockerfileContent can be set, not both")
	}

	return nil
}
