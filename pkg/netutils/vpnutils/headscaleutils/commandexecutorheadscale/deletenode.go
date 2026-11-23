package commandexecutorheadscale

import (
	"context"
	"strconv"

	"github.com/asciich/asciichgolangpublic/pkg/commandexecutor/commandexecutorinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/parameteroptions"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// DeleteNodeByName ensures the node with the given name is deleted.
//
// Idempotent: If no node with this name exists, nothing is changed.
// Returns an error if more than one node has this name.
func DeleteNodeByName(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, nodeName string) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if nodeName == "" {
		return tracederrors.TracedErrorEmptyString("nodeName")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete headscale node '%s' on '%s' started.", nodeName, hostDescription)

	nodes, err := ListNodes(ctx, commandExecutor)
	if err != nil {
		return err
	}

	ids := []uint64{}
	for _, node := range nodes {
		name, err := node.GetName()
		if err != nil {
			return err
		}

		if name != nodeName {
			continue
		}

		id, err := node.GetID()
		if err != nil {
			return err
		}

		ids = append(ids, id)
	}

	if len(ids) > 1 {
		return tracederrors.TracedErrorf(
			"Unable to delete headscale node '%s' on '%s': '%d' nodes have this name (ids: %v). Use the node id instead.",
			nodeName, hostDescription, len(ids), ids,
		)
	}

	if len(ids) == 0 {
		logging.LogInfoByCtxf(ctx, "Headscale node '%s' on '%s' already absent.", nodeName, hostDescription)
		return nil
	}

	err = deleteNodeByIdUnchecked(ctx, commandExecutor, ids[0])
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Deleted headscale node '%s' (id '%d') on '%s'.", nodeName, ids[0], hostDescription)

	return nil
}

// DeleteNodeById ensures the node with the given id is deleted.
//
// Idempotent: If no node with this id exists, nothing is changed.
func DeleteNodeById(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, id uint64) error {
	if commandExecutor == nil {
		return tracederrors.TracedErrorNil("commandExecutor")
	}

	if id == 0 {
		return tracederrors.TracedError("id must not be 0")
	}

	hostDescription, err := commandExecutor.GetHostDescription()
	if err != nil {
		return err
	}

	logging.LogInfoByCtxf(ctx, "Delete headscale node with id '%d' on '%s' started.", id, hostDescription)

	nodes, err := ListNodes(ctx, commandExecutor)
	if err != nil {
		return err
	}

	exists := false
	for _, node := range nodes {
		nodeId, err := node.GetID()
		if err != nil {
			return err
		}

		if nodeId == id {
			exists = true
			break
		}
	}

	if !exists {
		logging.LogInfoByCtxf(ctx, "Headscale node with id '%d' on '%s' already absent.", id, hostDescription)
		return nil
	}

	err = deleteNodeByIdUnchecked(ctx, commandExecutor, id)
	if err != nil {
		return err
	}

	logging.LogChangedByCtxf(ctx, "Deleted headscale node with id '%d' on '%s'.", id, hostDescription)

	return nil
}

// deleteNodeByIdUnchecked runs the delete command without checking if the node exists.
func deleteNodeByIdUnchecked(ctx context.Context, commandExecutor commandexecutorinterfaces.CommandExecutor, id uint64) error {
	_, err := commandExecutor.RunCommandAndGetStdoutAsString(ctx, &parameteroptions.RunCommandOptions{
		Command: []string{"headscale", "nodes", "delete", "--identifier", strconv.FormatUint(id, 10), "--force"},
	})
	if err != nil {
		return tracederrors.TracedErrorf("Failed to delete headscale node with id '%d': %w", id, err)
	}

	return nil
}
