package nativegitlab

import (
	"context"
	"net/http"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// ProjectExists checks if the project with the given path (e.g. "my-group/my-subgroup/my-project") exists.
//
// A project which exists but is not visible to the user of the nativeClient is reported as not existing,
// since GitLab answers with "404 Not Found" in this case as well.
func ProjectExists(ctx context.Context, nativeClient *gitlab.Client, projectPath string) (bool, error) {
	if nativeClient == nil {
		return false, tracederrors.TracedErrorNil("nativeClient")
	}

	if projectPath == "" {
		return false, tracederrors.TracedErrorEmptyString("projectPath")
	}

	logging.LogInfoByCtxf(ctx, "Check if gitlab project '%s' exists started.", projectPath)

	ret := true

	// The client escapes the path (e.g. the '/' characters) itself when a string is used as project ID.
	_, response, err := nativeClient.Projects.GetProject(projectPath, nil, gitlab.WithContext(ctx))
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			ret = false
		} else {
			return false, tracederrors.TracedErrorf("Failed to get gitlab project '%s': %w", projectPath, err)
		}
	}

	if ret {
		logging.LogInfoByCtxf(ctx, "Check if gitlab project '%s' exists finished. Project exists.", projectPath)
	} else {
		logging.LogInfoByCtxf(ctx, "Check if gitlab project '%s' exists finished. Project does not exist.", projectPath)
	}

	return ret, nil
}
