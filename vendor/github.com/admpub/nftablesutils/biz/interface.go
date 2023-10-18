package biz

import (
	"net"
	"time"

	"github.com/google/nftables"
)

type INFTables interface {
	// UpdateTrustIPs updates filterSetTrustIP.
	UpdateTrustIPs(del, add []net.IP) error

	// UpdateManagerIPs updates filterSetManagerIP.
	UpdateManagerIPs(del, add []net.IP) error

	// UpdateMyForwardWanIPs updates filterSetForwardIP.
	UpdateForwardWanIPs(del, add []net.IP) error

	// Ban adding ip to backlist.
	Ban(add []net.IP, timeout time.Duration) error

	// Cleanup rules to default policy filtering.
	Cleanup() error

	// WanIP returns ip address of wan interface.
	WanIP() net.IP

	// IfacesIPs returns ip addresses list of additional ifaces.
	IfacesIPs() ([]net.IP, error)

	// -- table & chain & set --

	TableFilter() *nftables.Table
	ChainInput() *nftables.Chain
	ChainForward() *nftables.Chain
	ChainOutput() *nftables.Chain

	TableNAT() *nftables.Table
	ChainPrerouting() *nftables.Chain
	ChainPostrouting() *nftables.Chain

	FilterSetTrustIP() *nftables.Set
	FilterSetManagerIP() *nftables.Set
	FilterSetForwardIP() *nftables.Set
	FilterSetBlacklistIP() *nftables.Set

	Do(f func(conn *nftables.Conn) error) error
}
