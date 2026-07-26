package kindutils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubernetesinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// SharedClusterName is the name of the shared KinD cluster used by all tests that can run concurrently.
const SharedClusterName = "asciichgolangpublic"

// clusterCreateLock is a mutex to coordinate cluster creation within a single process.
var clusterCreateLock sync.Mutex

// getClusterLockFilePath returns the path to the file-based lock for cluster creation.
func getClusterLockFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", tracederrors.TracedErrorf("Unable to get home directory: %v", err)
	}
	return filepath.Join(homeDir, ".kind-cluster-create.lock"), nil
}

// acquireClusterCreateLock acquires a file-based lock to ensure only one process creates the shared cluster.
// This function will block until the lock is acquired.
func acquireClusterCreateLock(ctx context.Context) (releaseFunc func(), err error) {
	lockFilePath, err := getClusterLockFilePath()
	if err != nil {
		return nil, err
	}

	// Try to acquire the lock with timeout
	maxWaitTime := 5 * time.Minute
	startTime := time.Now()

	for {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return nil, tracederrors.TracedErrorf("Context cancelled while waiting for cluster create lock: %v", ctx.Err())
		default:
		}

		// Try to create the lock file exclusively
		file, err := os.OpenFile(lockFilePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
		if err == nil {
			// Successfully created lock file
			pid := os.Getpid()
			_, writeErr := fmt.Fprintf(file, "%d\n%s", pid, time.Now().Format(time.RFC3339))
			if writeErr != nil {
				file.Close()
				os.Remove(lockFilePath)
				return nil, tracederrors.TracedErrorf("Unable to write to lock file: %v", writeErr)
			}
			file.Close()

			logging.LogInfoByCtxf(ctx, "Acquired cluster create lock (PID: %d).", pid)

			// Return release function
			return func() {
				os.Remove(lockFilePath)
				logging.LogInfoByCtxf(ctx, "Released cluster create lock.")
			}, nil
		}

		// Check if we've waited too long
		if time.Since(startTime) > maxWaitTime {
			return nil, tracederrors.TracedErrorf("Timeout waiting for cluster create lock after %v", maxWaitTime)
		}

		// Check if lock is stale (older than 1 hour)
		info, statErr := os.Stat(lockFilePath)
		if statErr == nil {
			if time.Since(info.ModTime()) > time.Hour {
				logging.LogInfoByCtxf(ctx, "Lock file appears stale (older than 1 hour), removing it.")
				os.Remove(lockFilePath)
				continue
			}
		}

		// Wait a bit before retrying
		select {
		case <-ctx.Done():
			return nil, tracederrors.TracedErrorf("Context cancelled while waiting for cluster create lock: %v", ctx.Err())
		case <-time.After(500 * time.Millisecond):
		}
	}
}

// GetOrCreateSharedCluster gets the shared cluster if it exists, or creates it if it doesn't.
// This function uses file-based locking to ensure only one process creates the cluster.
func GetOrCreateSharedCluster(ctx context.Context) (cluster kubernetesinterfaces.KubernetesCluster, err error) {
	clusterCreateLock.Lock()
	defer clusterCreateLock.Unlock()

	kind, err := GetLocalCommandExecutorKind()
	if err != nil {
		return nil, err
	}

	// First check if cluster already exists
	exists, err := kind.ClusterByNameExists(ctx, SharedClusterName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Shared cluster '%s' already exists. Reusing it.", SharedClusterName)
		return kind.GetClusterByName(SharedClusterName)
	}

	// Cluster doesn't exist, acquire file-based lock to create it
	releaseLock, err := acquireClusterCreateLock(ctx)
	if err != nil {
		return nil, err
	}
	defer releaseLock()

	// Double-check if another process created the cluster while we were waiting for the lock
	exists, err = kind.ClusterByNameExists(ctx, SharedClusterName)
	if err != nil {
		return nil, err
	}

	if exists {
		logging.LogInfoByCtxf(ctx, "Shared cluster '%s' was created by another process. Reusing it.", SharedClusterName)
		return kind.GetClusterByName(SharedClusterName)
	}

	// Create the cluster
	logging.LogInfoByCtxf(ctx, "Creating shared cluster '%s'...", SharedClusterName)
	cluster, err = CreateCluster(ctx, SharedClusterName)
	if err != nil {
		return nil, tracederrors.TracedErrorf("Failed to create shared cluster '%s': %v", SharedClusterName, err)
	}

	logging.LogInfoByCtxf(ctx, "Shared cluster '%s' created successfully.", SharedClusterName)
	return cluster, nil
}

// DeleteSharedCluster deletes the shared cluster.
// This should only be called in cleanup scenarios or when explicitly needed.
func DeleteSharedCluster(ctx context.Context) error {
	clusterCreateLock.Lock()
	defer clusterCreateLock.Unlock()

	kind, err := GetLocalCommandExecutorKind()
	if err != nil {
		return err
	}

	exists, err := kind.ClusterByNameExists(ctx, SharedClusterName)
	if err != nil {
		return err
	}

	if !exists {
		logging.LogInfoByCtxf(ctx, "Shared cluster '%s' does not exist. Skip deletion.", SharedClusterName)
		return nil
	}

	logging.LogInfoByCtxf(ctx, "Deleting shared cluster '%s'...", SharedClusterName)
	return DeleteClusterByName(ctx, SharedClusterName)
}
