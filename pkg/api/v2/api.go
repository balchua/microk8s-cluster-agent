package v2

import (
	"net"

	"github.com/canonical/microk8s-cluster-agent/pkg/snap"
)

// API implements the v2 API.
type API struct {
	// Snap interacts with the MicroK8s snap.
	Snap snap.Snap

	// ListControlPlaneNodeIPs is used in v2/join to list the IP addresses of the
	// known control plane nodes.
	ListControlPlaneNodeIPs ListControlPlaneNodeIPsFunc

	// LookupIP is net.LookupIP.
	LookupIP func(string) ([]net.IP, error)
}