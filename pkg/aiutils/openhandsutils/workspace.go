package openhandsutils

import (
	"context"
	"encoding/json"
	"net/url"
	"sort"

	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

func (o *Openhands) DeleteAllWorkspaces(ctx context.Context) error {
	logging.LogInfoByCtxf(ctx, "Delete all openhands workspaces started.")

	type Workspace struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Path string `json:"path"`
	}

	type WorkspacesResponse struct {
		Workspaces       []Workspace `json:"workspaces"`
		WorkspaceParents []string    `json:"workspaceParents"`
	}

	body, err := o.DoApiRequest(ctx, "GET", "/api/workspaces", nil)
	if err != nil {
		return err
	}

	var response WorkspacesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return tracederrors.TracedErrorf("Unable to parse workspaces response: %w", err)
	}

	for _, workspace := range response.Workspaces {
		logging.LogInfoByCtxf(ctx, "Deleting openhands workspace '%s' with path '%s'.", workspace.Name, workspace.Path)

		_, err := o.DoApiRequest(
			ctx,
			"DELETE",
			"/api/workspaces?path="+url.QueryEscape(workspace.Path),
			nil,
		)
		if err != nil {
			return tracederrors.TracedErrorf("Failed to delete workspace '%s': %w", workspace.Name, err)
		}

		logging.LogChangedByCtxf(ctx, "Deleted workspace '%s' with path '%s'.", workspace.Name, workspace.Path)
	}

	logging.LogInfoByCtxf(ctx, "Delete all openhands workspaces finished.")

	return nil
}

func (o *Openhands) ListWorkspaceNames(ctx context.Context) ([]string, error) {
	logging.LogInfoByCtxf(ctx, "List openhands workspace names started.")

	type Workspace struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Path string `json:"path"`
	}

	type WorkspacesResponse struct {
		Workspaces       []Workspace `json:"workspaces"`
		WorkspaceParents []string    `json:"workspaceParents"`
	}

	body, err := o.DoApiRequest(
		ctx,
		"GET",
		"/api/workspaces",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var response WorkspacesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Unable to parse workspaces response: %w", err)
	}

	names := []string{}
	for _, workspace := range response.Workspaces {
		names = append(names, workspace.Name)
	}

	sort.Strings(names)

	logging.LogInfoByCtxf(ctx, "List openhands workspace names finished.")

	return names, nil
}
func (o *Openhands) WorkspaceExists(ctx context.Context, workspaceName string) (bool, error) {
	if workspaceName == "" {
		return false, tracederrors.TracedErrorEmptyString("workspaceName")
	}

	existingNames, err := o.ListWorkspaceNames(ctx)
	if err != nil {
		return false, err
	}

	for _, name := range existingNames {
		if name == workspaceName {
			return true, nil
		}
	}

	return false, nil
}

func (o *Openhands) CreateWorkspace(ctx context.Context, workspaceName string, path string) error {
	if workspaceName == "" {
		return tracederrors.TracedErrorEmptyString("workspaceName")
	}

	if path == "" {
		return tracederrors.TracedErrorEmptyString("path")
	}

	logging.LogInfoByCtxf(ctx, "Create openhands workspace '%s' started.", workspaceName)

	exists, err := o.WorkspaceExists(ctx, workspaceName)
	if err != nil {
		return err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Openhands workspace '%s' already exists.", workspaceName)
		return nil
	}

	type Workspace struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Path string `json:"path"`
	}

	type WorkspacesPayload struct {
		Workspaces       []Workspace `json:"workspaces"`
		WorkspaceParents []string    `json:"workspaceParents"`
	}

	payloadData := WorkspacesPayload{
		Workspaces: []Workspace{
			{
				ID:   "/" + workspaceName,
				Name: workspaceName,
				Path: path,
			},
		},
		WorkspaceParents: []string{},
	}

	payload, err := json.Marshal(payloadData)
	if err != nil {
		return tracederrors.TracedErrorf("Failed to prepare payload to create workspace: %w", err)
	}

	_, err = o.DoApiRequest(ctx, "POST", "/api/workspaces", payload)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Created openhands workspace '%s' with path '%s'.", workspaceName, path)

	logging.LogInfoByCtxf(ctx, "Create openhands workspace '%s' finished.", workspaceName)

	return nil
}
