package kuberneteshost

import (
	"context"

	"github.com/asciich/asciichgolangpublic/pkg/hostsutils/hostsutilsinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/kubernetesutils/kubeadmutils"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// IsControlPlaneInitialized returns true if the given host has already been
// initialized as a kubernetes control-plane (i.e. /etc/kubernetes/admin.conf
// exists).
func IsControlPlaneInitialized(ctx context.Context, host hostsutilsinterfaces.Host) (bool, error) {
	if host == nil {
		return false, tracederrors.TracedErrorNil("host")
	}

	return kubeadmutils.IsControlPlaneInitialized(ctx, host)
}

// GetControlPlaneJoinCommand returns a ready-to-run "kubeadm join" command that
// lets an additional node join the existing cluster as a control-plane node.
//
// Fresh control-plane certificates are (re-)uploaded so the returned command is
// valid even on an already running cluster (the certificate key expires ~2h
// after upload).
func GetControlPlaneJoinCommand(ctx context.Context, host hostsutilsinterfaces.Host) (string, error) {
	if host == nil {
		return "", tracederrors.TracedErrorNil("host")
	}

	return kubeadmutils.GetControlPlaneJoinCommandUsingCommandExecutor(ctx, host)
}

// JoinControlPlaneOptions contains the options used to join an additional
// control-plane node to an existing cluster.
type JoinControlPlaneOptions struct {
	// JoinCommand is the "kubeadm join ... --control-plane --certificate-key ..."
	// command obtained from GetControlPlaneJoinCommand.
	JoinCommand string

	// ControlPlaneEndpoint is the shared/virtual IP (VIP) of the cluster.
	ControlPlaneEndpoint string

	// EnableKubeVip deploys the KubeVip static-pod manifest on the joining node
	// so the VIP can float to it.
	EnableKubeVip bool

	// KubeVipInterface is the network interface KubeVip binds the VIP to.
	// Optional: auto-detected when empty.
	KubeVipInterface string

	// KubeVipVersion is the KubeVip image version.
	// Optional: a built-in default is used when empty.
	KubeVipVersion string
}

// JoinControlPlane joins the given host to the existing cluster as an
// additional control-plane node. The operation is skipped when the node is
// already initialized (idempotent).
func JoinControlPlane(ctx context.Context, host hostsutilsinterfaces.Host, options *JoinControlPlaneOptions) error {
	if host == nil {
		return tracederrors.TracedErrorNil("host")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	return kubeadmutils.JoinControlPlaneUsingCommandExecutor(ctx, host, &kubeadmutils.JoinControlPlaneOptions{
		JoinCommand:          options.JoinCommand,
		ControlPlaneEndpoint: options.ControlPlaneEndpoint,
		EnableKubeVip:        options.EnableKubeVip,
		KubeVipInterface:     options.KubeVipInterface,
		KubeVipVersion:       options.KubeVipVersion,
	})
}

// GetNodeJoinCommand returns a ready-to-run "kubeadm join" command that lets a
// worker node join the existing cluster.
//
// In contrast to GetControlPlaneJoinCommand this command joins the node as a
// plain worker: it contains neither the "--control-plane" flag nor a
// "--certificate-key", so no control-plane certificates are uploaded.
//
// A fresh bootstrap token is (re-)created so the returned command is valid even
// on an already running cluster (the default bootstrap token expires ~24h after
// creation).
func GetNodeJoinCommand(ctx context.Context, host hostsutilsinterfaces.Host) (string, error) {
	if host == nil {
		return "", tracederrors.TracedErrorNil("host")
	}

	return kubeadmutils.GetNodeJoinCommandUsingCommandExecutor(ctx, host)
}

// IsNodeJoined returns true if the given host has already joined a kubernetes
// cluster as a node (i.e. the kubelet config /etc/kubernetes/kubelet.conf
// exists).
//
// This works for both worker nodes and control-plane nodes, since every node
// that ran "kubeadm join" (or "kubeadm init") has a kubelet.conf. Use
// IsControlPlaneInitialized instead when you specifically need to know whether
// a host is a control-plane.
func IsNodeJoined(ctx context.Context, host hostsutilsinterfaces.Host) (bool, error) {
	if host == nil {
		return false, tracederrors.TracedErrorNil("host")
	}

	return kubeadmutils.IsNodeJoined(ctx, host)
}

// JoinNode joins the given host to the existing cluster as a worker node. The
// operation is skipped when the node has already joined (idempotent).
func JoinNode(ctx context.Context, host hostsutilsinterfaces.Host, options *JoinNodeOptions) error {
	if host == nil {
		return tracederrors.TracedErrorNil("host")
	}

	if options == nil {
		return tracederrors.TracedErrorNil("options")
	}

	return kubeadmutils.JoinNodeUsingCommandExecutor(ctx, host, &kubeadmutils.JoinNodeOptions{
		JoinCommand: options.JoinCommand,
	})
}
