package biz

import (
	"net"

	"github.com/google/nftables"
)

type INFTables interface {
	// UpdateTrustIPs updates filterSetTrustIP.
	UpdateTrustIPs(del, add []net.IP) error

	// UpdateMyManagerIPs updates filterSetMyManagerIP.
	UpdateMyManagerIPs(del, add []net.IP) error

	// UpdateMyForwardWanIPs updates filterSetMyForwardIP.
	UpdateMyForwardWanIPs(del, add []net.IP) error

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
	FilterSetMyManagerIP() *nftables.Set
	FilterSetMyForwardIP() *nftables.Set

	Do(f func(conn *nftables.Conn) error) error
}
