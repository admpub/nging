package biz

import (
	"fmt"
	"net"
	"runtime"
	"strings"

	utils "github.com/admpub/nftablesutils"
	"github.com/google/nftables"
	"github.com/google/nftables/binaryutil"
	"github.com/google/nftables/expr"
	"github.com/vishvananda/netns"
)

const loIface = "lo"

var _ INFTables = &NFTables{}

func New(tableFamily nftables.TableFamily, c Config, managerPorts []uint16) *NFTables {
	return &NFTables{
		tableFamily:  tableFamily,
		cfg:          &c,
		managerPorts: managerPorts,
	}
}

// NFTables struct.
type NFTables struct {
	cfg         *Config
	tableFamily nftables.TableFamily

	originNetNS netns.NsHandle
	targetNetNS netns.NsHandle

	wanIface string
	wanIP    net.IP
	myIface  string
	myPort   uint16

	tFilter  *nftables.Table
	cInput   *nftables.Chain
	cForward *nftables.Chain
	cOutput  *nftables.Chain

	tNAT         *nftables.Table
	cPrerouting  *nftables.Chain
	cPostrouting *nftables.Chain

	filterSetTrustIP     *nftables.Set
	filterSetMyManagerIP *nftables.Set
	filterSetMyForwardIP *nftables.Set

	managerPorts []uint16

	applied bool
}

// Init nftables firewall.
func (nft *NFTables) Init() error {
	if nft.tableFamily == 0 {
		nft.tableFamily = nftables.TableFamilyIPv4
	}
	cfg := nft.cfg
	// obtain default interface name, ip address and gateway ip address
	var wanIface string
	var wanIP net.IP
	var err error
	if nft.tableFamily == nftables.TableFamilyIPv6 {
		wanIface, _, wanIP, err = utils.IPv6Addr()
	} else {
		wanIface, _, wanIP, err = utils.IPAddr()
	}
	if err != nil {
		err = fmt.Errorf(`failed to obtain default interface name: %w`, err)
	}

	defaultPolicy := nftables.ChainPolicyDrop
	if strings.ToLower(cfg.DefaultPolicy) == "accept" {
		defaultPolicy = nftables.ChainPolicyAccept
	}

	tFilter := &nftables.Table{
		Family: nft.tableFamily,
		Name:   cfg.TablePrefix + TableFilter,
	}
	cInput := &nftables.Chain{
		Name:     ChainInput,
		Table:    tFilter,
		Type:     nftables.ChainTypeFilter,
		Priority: nftables.ChainPriorityFilter,
		Hooknum:  nftables.ChainHookInput,
		Policy:   &defaultPolicy,
	}
	cForward := &nftables.Chain{
		Name:     ChainForward,
		Table:    tFilter,
		Type:     nftables.ChainTypeFilter,
		Priority: nftables.ChainPriorityFilter,
		Hooknum:  nftables.ChainHookForward,
		Policy:   &defaultPolicy,
	}
	cOutput := &nftables.Chain{
		Name:     ChainOutput,
		Table:    tFilter,
		Type:     nftables.ChainTypeFilter,
		Priority: nftables.ChainPriorityFilter,
		Hooknum:  nftables.ChainHookOutput,
		Policy:   &defaultPolicy,
	}

	tNAT := &nftables.Table{
		Family: nft.tableFamily,
		//Family: nftables.TableFamilyIPv4,
		Name: cfg.TablePrefix + TableNAT,
	}
	cPrerouting := &nftables.Chain{
		Name:     ChainPreRouting,
		Table:    tNAT,
		Type:     nftables.ChainTypeNAT,
		Priority: nftables.ChainPriorityNATDest,
		Hooknum:  nftables.ChainHookPrerouting,
	}
	cPostrouting := &nftables.Chain{
		Name:     ChainPostRouting,
		Table:    tNAT,
		Type:     nftables.ChainTypeNAT,
		Priority: nftables.ChainPriorityNATSource,
		Hooknum:  nftables.ChainHookPostrouting,
	}

	filterSetTrustIP := &nftables.Set{ // input / output IP whitelist
		Name:    "trust_ipset",
		Table:   tFilter,
		KeyType: nftables.TypeIPAddr,
	}
	filterSetMyManagerIP := &nftables.Set{ // input / output IP whitelist
		Name:    "my_manager_ipset",
		Table:   tFilter,
		KeyType: nftables.TypeIPAddr,
	}
	filterSetMyForwardIP := &nftables.Set{ // forward IP whitelist
		Name:    "my_forward_ipset",
		Table:   tFilter,
		KeyType: nftables.TypeIPAddr,
	}
	if tFilter.Family == nftables.TableFamilyIPv6 {
		filterSetTrustIP.KeyType = nftables.TypeIP6Addr
		filterSetMyManagerIP.KeyType = nftables.TypeIP6Addr
		filterSetMyForwardIP.KeyType = nftables.TypeIP6Addr
	}

	nft.wanIface = wanIface
	nft.wanIP = wanIP
	nft.myIface = cfg.MyIface
	nft.myPort = cfg.MyPort

	nft.tFilter = tFilter
	nft.cInput = cInput
	nft.cForward = cForward
	nft.cOutput = cOutput

	nft.tNAT = tNAT
	nft.cPrerouting = cPrerouting
	nft.cPostrouting = cPostrouting

	nft.filterSetTrustIP = filterSetTrustIP
	nft.filterSetMyManagerIP = filterSetMyManagerIP
	nft.filterSetMyForwardIP = filterSetMyForwardIP
	return err
}

func (nft *NFTables) ApplyDefault() error {
	return nft.apply()
}

// networkNamespaceBind target by name.
func (nft *NFTables) networkNamespaceBind() (*nftables.Conn, error) {
	if nft.cfg.NetworkNamespace == "" {
		return &nftables.Conn{NetNS: int(nft.originNetNS)}, nil
	}

	// Lock the OS Thread so we don't accidentally switch namespaces
	runtime.LockOSThread()

	origin, err := netns.Get()
	if err != nil {
		nft.networkNamespaceRelease()
		return nil, fmt.Errorf(`failed to netns.Get: %w`, err)
	}

	nft.originNetNS = origin

	target, err := netns.GetFromName(nft.cfg.NetworkNamespace)
	if err != nil {
		nft.networkNamespaceRelease()
		return nil, fmt.Errorf(`failed to netns.GetFromName(%q): %w`, nft.cfg.NetworkNamespace, err)
	}

	// switch to target network namespace
	err = netns.Set(target)
	if err != nil {
		nft.networkNamespaceRelease()
		return nil, fmt.Errorf(`failed to netns.Set(%q): %w`, nft.cfg.NetworkNamespace, err)
	}
	nft.targetNetNS = target

	return &nftables.Conn{NetNS: int(nft.targetNetNS)}, nil
}

// networkNamespaceRelease to origin.
func (nft *NFTables) networkNamespaceRelease() error {
	if nft.cfg.NetworkNamespace == "" {
		return nil
	}

	// finally unlock os thread
	defer runtime.UnlockOSThread()

	// switch back to the original namespace
	err := netns.Set(nft.originNetNS)
	if err != nil {
		return err
	}

	// close fd to origin and dev ns
	nft.originNetNS.Close()
	nft.targetNetNS.Close()

	nft.targetNetNS = 0

	return nil
}

func (nft *NFTables) ApplyBase(c *nftables.Conn) error {
	// add filter table
	// cmd: nft add table ip filter
	c.AddTable(nft.tFilter)
	// add input chain of filter table
	// cmd: nft add chain ip filter input \
	// { type filter hook input priority 0 \; policy drop\; }
	c.AddChain(nft.cInput)
	// add forward chain
	// cmd: nft add chain ip filter forward \
	// { type filter hook forward priority 0 \; policy drop\; }
	c.AddChain(nft.cForward)
	// add output chain
	// cmd: nft add chain ip filter output \
	// { type filter hook output priority 0 \; policy drop\; }
	c.AddChain(nft.cOutput)

	// add nat table
	// cmd: nft add table ip nat
	c.AddTable(nft.tNAT)

	// add prerouting chain
	// cmd: nft add chain ip nat prerouting \
	// { type nat hook prerouting priority -100 \; }
	c.AddChain(nft.cPrerouting)

	// add postrouting chain
	// cmd: nft add chain ip nat postrouting \
	// { type nat hook postrouting priority 100 \; }
	c.AddChain(nft.cPostrouting)

	if nft.cfg.DisableInitSet {
		return nil
	}

	// Init sets.
	return nft.InitSet(c)
}

func (nft *NFTables) InitSet(c *nftables.Conn) error {
	// add trust_ipset
	// cmd: nft add set ip filter trust_ipset { type ipv4_addr\; }
	// --
	// set trust_ipset {
	//         type ipv4_addr
	// }
	err := c.AddSet(nft.filterSetTrustIP, nil)
	if err != nil {
		return fmt.Errorf(`nft.AddSet(%q): %w`, nft.filterSetTrustIP.Name, err)
	}

	// add mymanager_ipset
	// cmd: nft add set ip filter mymanager_ipset { type ipv4_addr\; }
	// --
	// set mymanager_ipset {
	//         type ipv4_addr
	// }
	err = c.AddSet(nft.filterSetMyManagerIP, nil)
	if err != nil {
		return fmt.Errorf(`nft.AddSet(%q): %w`, nft.filterSetMyManagerIP.Name, err)
	}

	// add myforward_ipset
	// cmd: nft add set ip filter myforward_ipset { type ipv4_addr\; }
	// --
	// set myforward_ipset {
	//         type ipv4_addr
	// }
	err = c.AddSet(nft.filterSetMyForwardIP, nil)
	if err != nil {
		return fmt.Errorf(`nft.AddSet(%q): %w`, nft.filterSetMyForwardIP.Name, err)
	}
	return err
}

// apply rules
func (nft *NFTables) apply() error {
	if !nft.cfg.Enabled {
		return nil
	}

	// bind network namespace if it was set in config
	c, err := nft.networkNamespaceBind()
	if err != nil {
		return fmt.Errorf(`nft.networkNamespaceBind: %w`, err)
	}

	// release network namespace finally
	defer nft.networkNamespaceRelease()
	if nft.cfg.ClearRuleset {
		c.FlushRuleset()
	} else {
		c.FlushTable(nft.tFilter)
		c.FlushTable(nft.tNAT)
		_ = c.Flush()
	}
	//
	// Init Tables and Chains.
	//
	err = nft.ApplyBase(c)
	if err != nil {
		return err
	}

	//
	// Init filter rules.
	//

	nft.inputLocalIfaceRules(c)
	nft.outputLocalIfaceRules(c)
	if err = nft.applyCommonRules(c, nft.wanIface); err != nil {
		return err
	}
	err = nft.sdnRules(c)
	if err != nil {
		return fmt.Errorf(`nft.sdnRules: %w`, err)
	}
	err = nft.sdnForwardRules(c)
	if err != nil {
		return fmt.Errorf(`nft.sdnForwardRules: %w`, err)
	}
	nft.natRules(c)

	for _, iface := range nft.cfg.Ifaces {
		if iface == nft.wanIface {
			continue
		}

		if err = nft.applyCommonRules(c, iface); err != nil {
			return err
		}
	}

	// apply configuration
	err = c.Flush()
	if err != nil {
		return err
	}
	nft.applied = true

	return nil
}

// sdnRules to apply.
func (nft *NFTables) sdnRules(c *nftables.Conn) error {
	if len(nft.myIface) == 0 {
		return nil
	}
	// cmd: nft add rule ip filter input meta iifname "wg0" ip protocol icmp \
	// icmp type echo-request ct state new accept
	// --
	// iifname "wg0" icmp type echo-request ct state new accept
	exprs := make([]expr.Any, 0, 12)
	exprs = append(exprs, utils.SetIIF(nft.myIface)...)
	switch nft.tableFamily {
	case nftables.TableFamilyIPv4:
		exprs = append(exprs, utils.SetProtoICMP()...)
		exprs = append(exprs, utils.SetICMPTypeEchoRequest()...)
	case nftables.TableFamilyIPv6:
		exprs = append(exprs, utils.SetProtoICMPv6()...)
		exprs = append(exprs, utils.SetICMPv6TypeEchoRequest()...)
	}
	exprs = append(exprs, utils.SetConntrackStateNew()...)
	exprs = append(exprs, utils.ExprAccept())
	rule := &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	// cmd: nft add rule ip filter input meta iifname "wg0" ip protocol icmp \
	// ct state { established, related } accept
	// --
	// iifname "wg0" ip protocol icmp ct state { established, related } accept
	ctStateSet := utils.GetConntrackStateSet(nft.tFilter)
	elems := utils.GetConntrackStateSetElems(defaultStateWithOld)
	err := c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	exprs = make([]expr.Any, 0, 7)
	exprs = append(exprs, utils.SetIIF(nft.myIface)...)
	switch nft.tableFamily {
	case nftables.TableFamilyIPv4:
		exprs = append(exprs, utils.SetProtoICMP()...)
	case nftables.TableFamilyIPv6:
		exprs = append(exprs, utils.SetProtoICMPv6()...)
	}
	exprs = append(exprs, utils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, utils.ExprAccept())

	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	// cmd: nft add rule ip filter input meta iifname "wg0" \
	// ip protocol tcp tcp dport { 80, 8080 } ip saddr @mymanager_ipset \
	// ct state { new, established } accept
	// --
	// iifname "wg0" tcp dport { https, 8443 } ip saddr @mymanager_ipset ct state { established, new } accept
	ctStateSet = utils.GetConntrackStateSet(nft.tFilter)
	elems = utils.GetConntrackStateSetElems(defaultStateWithNew)
	err = c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	portSet := utils.GetPortSet(nft.tFilter)
	portSetElems := make([]nftables.SetElement, len(nft.managerPorts))
	for i, p := range nft.managerPorts {
		portSetElems[i] = nftables.SetElement{Key: binaryutil.BigEndian.PutUint16(p)}
	}
	err = c.AddSet(portSet, portSetElems)
	if err != nil {
		return err
	}

	exprs = make([]expr.Any, 0, 9)
	exprs = append(exprs, utils.SetIIF(nft.myIface)...)
	exprs = append(exprs, utils.SetProtoTCP()...)
	exprs = append(exprs, utils.SetDPortSet(portSet)...)
	switch nft.tableFamily {
	case nftables.TableFamilyIPv4:
		exprs = append(exprs, utils.SetSAddrSet(nft.filterSetMyManagerIP)...)
	case nftables.TableFamilyIPv6:
		exprs = append(exprs, utils.SetSAddrIPv6Set(nft.filterSetMyManagerIP)...)
	}
	exprs = append(exprs, utils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, utils.ExprAccept())
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cInput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	// cmd: nft add rule ip filter output meta oifname "wg0" ip protocol icmp \
	// ct state { new, established } accept
	// --
	// oifname "wg0" ip protocol icmp ct state { established, new } accept
	ctStateSet = utils.GetConntrackStateSet(nft.tFilter)
	elems = utils.GetConntrackStateSetElems(defaultStateWithNew)
	err = c.AddSet(ctStateSet, elems)
	if err != nil {
		return err
	}

	exprs = make([]expr.Any, 0, 7)
	exprs = append(exprs, utils.SetOIF(nft.myIface)...)
	switch nft.tableFamily {
	case nftables.TableFamilyIPv4:
		exprs = append(exprs, utils.SetProtoICMP()...)
	case nftables.TableFamilyIPv6:
		exprs = append(exprs, utils.SetProtoICMPv6()...)
	}
	exprs = append(exprs, utils.SetConntrackStateSet(ctStateSet)...)
	exprs = append(exprs, utils.ExprAccept())

	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	// cmd: nft add rule ip filter output meta oifname "wg0" \
	// ip protocol tcp tcp sport { 80, 8080 } ip daddr @mymanager_ipset \
	// ct state established accept
	// --
	// oifname "wg0" tcp sport { https, 8443 } ct state established accept
	portSet = utils.GetPortSet(nft.tFilter)
	portSetElems = make([]nftables.SetElement, len(nft.managerPorts))
	for i, p := range nft.managerPorts {
		portSetElems[i] = nftables.SetElement{Key: binaryutil.BigEndian.PutUint16(p)}
	}
	err = c.AddSet(portSet, portSetElems)
	if err != nil {
		return err
	}

	exprs = make([]expr.Any, 0, 10)
	exprs = append(exprs, utils.SetOIF(nft.myIface)...)
	exprs = append(exprs, utils.SetProtoTCP()...)
	exprs = append(exprs, utils.SetSPortSet(portSet)...)
	switch nft.tableFamily {
	case nftables.TableFamilyIPv4:
		exprs = append(exprs, utils.SetDAddrSet(nft.filterSetMyManagerIP)...)
	case nftables.TableFamilyIPv6:
		exprs = append(exprs, utils.SetDAddrIPv6Set(nft.filterSetMyManagerIP)...)
	}
	exprs = append(exprs, utils.SetConntrackStateEstablished()...)
	exprs = append(exprs, utils.ExprAccept())
	rule = &nftables.Rule{
		Table: nft.tFilter,
		Chain: nft.cOutput,
		Exprs: exprs,
	}
	c.AddRule(rule)

	return nil
}

// sdnForwardRules to apply.
func (nft *NFTables) sdnForwardRules(c *nftables.Conn) error {
	err := nft.forwardSMTPRules(c)
	if err != nil {
		return err
	}

	return nft.forwardInterfaceRules(c)
}

// natRules to apply.
func (nft *NFTables) natRules(c *nftables.Conn) {
	nft.natInterfaceRules(c)
}

// UpdateTrustIPs updates filterSetTrustIP.
func (nft *NFTables) UpdateTrustIPs(del, add []net.IP) error {
	if !nft.applied {
		return nil
	}

	return nft.updateIPSet(nft.filterSetTrustIP, del, add)
}

// UpdateMyManagerIPs updates filterSetMyManagerIP.
func (nft *NFTables) UpdateMyManagerIPs(del, add []net.IP) error {
	if !nft.applied {
		return nil
	}

	return nft.updateIPSet(nft.filterSetMyManagerIP, del, add)
}

// UpdateMyForwardWanIPs updates filterSetMyForwardIP.
func (nft *NFTables) UpdateMyForwardWanIPs(del, add []net.IP) error {
	if !nft.applied {
		return nil
	}

	return nft.updateIPSet(nft.filterSetMyForwardIP, del, add)
}

func (nft *NFTables) updateIPSet(set *nftables.Set, del, add []net.IP) error {
	// bind network namespace if it was set in config
	c, err := nft.networkNamespaceBind()
	if err != nil {
		return err
	}
	// release network namespace finally
	defer nft.networkNamespaceRelease()

	if len(del) > 0 {
		elements := make([]nftables.SetElement, len(del))
		for i, v := range del {
			elements[i] = nftables.SetElement{Key: v}
		}
		err = c.SetDeleteElements(set, elements)
		if err != nil {
			return err
		}
	}

	if len(add) > 0 {
		elements := make([]nftables.SetElement, len(add))
		for i, v := range add {
			elements[i] = nftables.SetElement{Key: v}
		}
		err = c.SetAddElements(set, elements)
		if err != nil {
			return err
		}
	}

	return c.Flush()
}

// Cleanup rules to default policy filtering.
func (nft *NFTables) Cleanup() error {
	if !nft.cfg.Enabled {
		return nil
	}
	// bind network namespace if it was set in config
	c, err := nft.networkNamespaceBind()
	if err != nil {
		return err
	}
	// release network namespace finally
	defer nft.networkNamespaceRelease()

	if nft.cfg.ClearRuleset {
		c.FlushRuleset()
	} else {
		c.FlushTable(nft.tFilter)
		c.FlushTable(nft.tNAT)
		_ = c.Flush()
	}
	_ = c.Flush()
	nft.applied = false

	return nil
}

// WanIP returns ip address of wan interface.
func (nft *NFTables) WanIP() net.IP {
	return nft.wanIP
}

// IfacesIPs returns ip addresses list of additional ifaces.
func (nft *NFTables) IfacesIPs() ([]net.IP, error) {
	ips := make([]net.IP, 0, len(nft.cfg.Ifaces))

	for _, v := range nft.cfg.Ifaces {
		if v == nft.wanIface || v == nft.myIface {
			continue
		}

		iface, err := net.InterfaceByName(v)
		if err != nil {
			return nil, err
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipnet.IP.To4()
			if ip != nil {
				ips = append(ips, ip)
			}
		}
	}

	return ips, nil
}

func (nft *NFTables) TableFilter() *nftables.Table {
	return nft.tFilter
}

func (nft *NFTables) ChainInput() *nftables.Chain {
	return nft.cInput
}

func (nft *NFTables) ChainForward() *nftables.Chain {
	return nft.cForward
}

func (nft *NFTables) ChainOutput() *nftables.Chain {
	return nft.cOutput
}

func (nft *NFTables) TableNAT() *nftables.Table {
	return nft.tNAT
}

func (nft *NFTables) ChainPostrouting() *nftables.Chain {
	return nft.cPostrouting
}

func (nft *NFTables) ChainPrerouting() *nftables.Chain {
	return nft.cPrerouting
}

func (nft *NFTables) FilterSetTrustIP() *nftables.Set {
	return nft.filterSetTrustIP
}

func (nft *NFTables) FilterSetMyManagerIP() *nftables.Set {
	return nft.filterSetMyManagerIP
}

func (nft *NFTables) FilterSetMyForwardIP() *nftables.Set {
	return nft.filterSetMyForwardIP
}

func (nft *NFTables) Do(f func(conn *nftables.Conn) error) error {
	// bind network namespace if it was set in config
	c, err := nft.networkNamespaceBind()
	if err != nil {
		return err
	}
	// release network namespace finally
	defer nft.networkNamespaceRelease()
	return f(c)
}
