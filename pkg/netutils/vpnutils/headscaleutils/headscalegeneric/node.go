package headscalegeneric

import (
	"context"
	"time"

	"github.com/asciich/asciichgolangpublic/pkg/netutils/iputils"
	"github.com/asciich/asciichgolangpublic/pkg/netutils/vpnutils/headscaleutils/headscaleinterfaces"
	"github.com/asciich/asciichgolangpublic/pkg/tracederrors"
)

// Compile time check: Node must implement the headscaleinterfaces.Node interface.
var _ headscaleinterfaces.Node = (*Node)(nil)

// Node is a generic, purely in-memory representation of a headscale node.
//
// It does not perform any IO. It is meant to be filled by the native or
// commandexecutor implementations (e.g. after parsing `headscale nodes list -o json`).
type Node struct {
	id          uint64
	hostname    string
	name        string
	machineKey  string
	nodeKey     string
	user        string
	tags        []string
	ipAddresses []string
	ephemeral   bool
	lastSeen    *time.Time
	expiry      *time.Time
	online      bool
}

// NewNode returns a new, empty Node.
func NewNode() *Node {
	return &Node{
		tags:        []string{},
		ipAddresses: []string{},
	}
}

// normalizeTime returns a copy of the given time.
// nil and the zero time (0001-01-01 00:00:00, headscale's "not set" value) are mapped to nil.
func normalizeTime(t *time.Time) *time.Time {
	if t == nil || t.IsZero() {
		return nil
	}

	ret := *t

	return &ret
}

// GetDeepCopy returns a deep copy of the node.
// Slices and time pointers are copied, so the copy shares no memory with the original.
func (n *Node) GetDeepCopy() (*Node, error) {
	if n == nil {
		return nil, tracederrors.TracedErrorNil("node")
	}

	ret := *n

	ret.tags = append([]string{}, n.tags...)
	ret.ipAddresses = append([]string{}, n.ipAddresses...)
	ret.lastSeen = normalizeTime(n.lastSeen)
	ret.expiry = normalizeTime(n.expiry)

	return &ret, nil
}

// ---------- ID ----------

func (n *Node) GetID() (uint64, error) {
	if n == nil {
		return 0, tracederrors.TracedErrorNil("node")
	}

	if n.id == 0 {
		return 0, tracederrors.TracedError("id of the node is not set")
	}

	return n.id, nil
}

func (n *Node) SetID(id uint64) error {
	if n == nil {
		return tracederrors.TracedErrorNil("node")
	}

	if id == 0 {
		return tracederrors.TracedError("id must not be 0")
	}

	n.id = id

	return nil
}

// ---------- Hostname ----------

func (n *Node) GetHostname() (string, error) {
	if n == nil {
		return "", tracederrors.TracedErrorNil("node")
	}

	if n.hostname == "" {
		return "", tracederrors.TracedError("hostname of the node is not set")
	}

	return n.hostname, nil
}

func (n *Node) SetHostname(hostname string) error {
	if n == nil {
		return tracederrors.TracedErrorNil("node")
	}

	if hostname == "" {
		return tracederrors.TracedErrorEmptyString("hostname")
	}

	n.hostname = hostname

	return nil
}

// ---------- Name ----------

func (n *Node) GetName() (string, error) {
	if n == nil {
		return "", tracederrors.TracedErrorNil("node")
	}

	if n.name == "" {
		return "", tracederrors.TracedError("name of the node is not set")
	}

	return n.name, nil
}

func (n *Node) SetName(name string) error {
	if n == nil {
		return tracederrors.TracedErrorNil("node")
	}

	if name == "" {
		return tracederrors.TracedErrorEmptyString("name")
	}

	n.name = name

	return nil
}

// ---------- MachineKey ----------

func (n *Node) GetMachineKey() (string, error) {
	if n == nil {
		return "", tracederrors.TracedErrorNil("node")
	}

	if n.machineKey == "" {
		return "", tracederrors.TracedError("machine key of the node is not set")
	}

	return n.machineKey, nil
}

func (n *Node) SetMachineKey(machineKey string) error {
	if n == nil {
		return tracederrors.TracedErrorNil("node")
	}

	if machineKey == "" {
		return tracederrors.TracedErrorEmptyString("machineKey")
	}

	n.machineKey = machineKey

	return nil
}

// ---------- NodeKey ----------

func (n *Node) GetNodeKey() (string, error) {
	if n == nil {
		return "", tracederrors.TracedErrorNil("node")
	}

	if n.nodeKey == "" {
		return "", tracederrors.TracedError("node key of the node is not set")
	}

	return n.nodeKey, nil
}

func (n *Node) SetNodeKey(nodeKey string) error {
	if n == nil {
		return tracederrors.TracedErrorNil("node")
	}

	if nodeKey == "" {
		return tracederrors.TracedErrorEmptyString("nodeKey")
	}

	n.nodeKey = nodeKey

	return nil
}

// ---------- User ----------

func (n *Node) GetUser() (string, error) {
	if n == nil {
		return "", tracederrors.TracedErrorNil("node")
	}

	if n.user == "" {
		return "", tracederrors.TracedError("user of the node is not set")
	}

	return n.user, nil
}

func (n *Node) SetUser(user string) error {
	if n == nil {
		return tracederrors.TracedErrorNil("node")
	}

	if user == "" {
		return tracederrors.TracedErrorEmptyString("user")
	}

	n.user = user

	return nil
}

// ---------- Tags ----------

// GetTags returns a copy of the tags. Returns an empty slice if there are no tags.
func (n *Node) GetTags() ([]string, error) {
	if n == nil {
		return nil, tracederrors.TracedErrorNil("node")
	}

	return append([]string{}, n.tags...), nil
}

// SetTags replaces the tags. The given slice is copied.
func (n *Node) SetTags(tags []string) error {
	if n == nil {
		return tracederrors.TracedErrorNil("node")
	}

	for _, tag := range tags {
		if tag == "" {
			return tracederrors.TracedErrorEmptyString("tag")
		}
	}

	n.tags = append([]string{}, tags...)

	return nil
}

// ---------- IP addresses ----------

// GetIPAddresses returns a copy of the IP addresses.
// The addresses were validated by SetIPAddresses.
func (n *Node) GetIPAddresses() ([]string, error) {
	if n == nil {
		return nil, tracederrors.TracedErrorNil("node")
	}

	return append([]string{}, n.ipAddresses...), nil
}

// SetIPAddresses replaces the IP addresses. Every address is validated using the iputils package.
// The given slice is copied.
func (n *Node) SetIPAddresses(ctx context.Context, ipAddresses []string) error {
	if n == nil {
		return tracederrors.TracedErrorNil("node")
	}

	for _, ip := range ipAddresses {
		err := iputils.CheckValidIP(ctx, ip)
		if err != nil {
			return err
		}
	}

	n.ipAddresses = append([]string{}, ipAddresses...)

	return nil
}

// ---------- Ephemeral ----------

func (n *Node) IsEphemeral() (bool, error) {
	if n == nil {
		return false, tracederrors.TracedErrorNil("node")
	}

	return n.ephemeral, nil
}

func (n *Node) SetEphemeral(ephemeral bool) error {
	if n == nil {
		return tracederrors.TracedErrorNil("node")
	}

	n.ephemeral = ephemeral

	return nil
}

// ---------- LastSeen ----------

// GetLastSeen returns a copy of the last seen time.
// Returns nil if the node was never seen.
func (n *Node) GetLastSeen() (*time.Time, error) {
	if n == nil {
		return nil, tracederrors.TracedErrorNil("node")
	}

	return normalizeTime(n.lastSeen), nil
}

// SetLastSeen sets the last seen time.
// nil and the zero time (0001-01-01 00:00:00) both mean "never seen".
func (n *Node) SetLastSeen(lastSeen *time.Time) error {
	if n == nil {
		return tracederrors.TracedErrorNil("node")
	}

	n.lastSeen = normalizeTime(lastSeen)

	return nil
}

// ---------- Expiry ----------

// GetExpiry returns a copy of the key expiration.
// Returns nil if the key never expires.
func (n *Node) GetExpiry() (*time.Time, error) {
	if n == nil {
		return nil, tracederrors.TracedErrorNil("node")
	}

	return normalizeTime(n.expiry), nil
}

// SetExpiry sets the key expiration.
// nil and the zero time (0001-01-01 00:00:00) both mean "never expires".
func (n *Node) SetExpiry(expiry *time.Time) error {
	if n == nil {
		return tracederrors.TracedErrorNil("node")
	}

	n.expiry = normalizeTime(expiry)

	return nil
}

// ---------- Online ----------

func (n *Node) IsOnline() (bool, error) {
	if n == nil {
		return false, tracederrors.TracedErrorNil("node")
	}

	return n.online, nil
}

func (n *Node) SetOnline(online bool) error {
	if n == nil {
		return tracederrors.TracedErrorNil("node")
	}

	n.online = online

	return nil
}

// ---------- Expired ----------

// IsExpired returns true if the node key is expired.
// It is derived from the expiry: a node without expiry never expires.
func (n *Node) IsExpired() (bool, error) {
	if n == nil {
		return false, tracederrors.TracedErrorNil("node")
	}

	if n.expiry == nil {
		return false, nil
	}

	return n.expiry.Before(time.Now()), nil
}
