package gerberos

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var (
	DefaultChainName  = "gerberos"
	DefaultTable4Name = "gerberos4"
	DefaultTable6Name = "gerberos6"
)

type Backend interface {
	Initialize() error
	Ban(ip string, ipv6 bool, d time.Duration) error
	Finalize() error
}

var backends = map[string]func(*Runner) Backend{
	`ipset`: NewIpsetBackend,
	`nft`:   NewNftBackend,
	`test`:  NewTestBackend,
}

func RegisterBackend(name string, bfn func(*Runner) Backend) {
	backends[name] = bfn
}

func NewIpsetBackend(rn *Runner) Backend {
	return &ipsetBackend{runner: rn}
}

func NewNftBackend(rn *Runner) Backend {
	return &nftBackend{runner: rn}
}

func NewTestBackend(rn *Runner) Backend {
	return &testBackend{runner: rn}
}

type ipsetBackend struct {
	runner     *Runner
	chainName  string
	ipset4Name string
	ipset6Name string
}

func (b *ipsetBackend) deleteIpsetsAndIptablesEntries() error {
	if s, ec, _ := b.runner.Executor.Execute("iptables", "-D", b.chainName, "-j", "DROP", "-m", "set", "--match-set", b.ipset4Name, "src"); ec > 2 {
		return fmt.Errorf(`failed to delete iptables entry for set "%s": %s`, b.ipset4Name, s)
	}
	if s, ec, _ := b.runner.Executor.Execute("iptables", "-D", "INPUT", "-j", b.chainName); ec > 2 {
		return fmt.Errorf(`failed to delete iptables entry for chain "%s": %s`, b.chainName, s)
	}
	if s, ec, _ := b.runner.Executor.Execute("iptables", "-X", b.chainName); ec > 2 {
		return fmt.Errorf(`failed to delete iptables chain "%s": %s`, b.chainName, s)
	}
	if s, ec, _ := b.runner.Executor.Execute("ip6tables", "-D", b.chainName, "-j", "DROP", "-m", "set", "--match-set", b.ipset6Name, "src"); ec > 2 {
		return fmt.Errorf(`failed to delete ip6tables entry for set "%s": %s`, b.ipset6Name, s)
	}
	if s, ec, _ := b.runner.Executor.Execute("ip6tables", "-D", "INPUT", "-j", b.chainName); ec > 2 {
		return fmt.Errorf(`failed to delete ip6tables entry for chain "%s": %s`, b.chainName, s)
	}
	if s, ec, _ := b.runner.Executor.Execute("ip6tables", "-X", b.chainName); ec > 2 {
		return fmt.Errorf(`failed to delete ip6tables chain "%s": %s`, b.chainName, s)
	}
	time.Sleep(250 * time.Millisecond) // Workaround for potential kernel lock problems
	if s, ec, _ := b.runner.Executor.Execute("ipset", "destroy", b.ipset4Name); ec > 1 {
		return fmt.Errorf(`failed to destroy ipset "%s": %s`, b.ipset4Name, s)
	}
	time.Sleep(250 * time.Millisecond) // Workaround for potential kernel lock problems
	if s, ec, _ := b.runner.Executor.Execute("ipset", "destroy", b.ipset6Name); ec > 1 {
		return fmt.Errorf(`failed to destroy ipset "%s": %s`, b.ipset6Name, s)
	}

	return nil
}

func (b *ipsetBackend) createIpsets() error {
	time.Sleep(250 * time.Millisecond) // Workaround for potential kernel lock problems
	if s, ec, _ := b.runner.Executor.Execute("ipset", "create", b.ipset4Name, "hash:ip", "timeout", "0"); ec != 0 {
		return fmt.Errorf(`failed to create ipset "%s": %s`, b.ipset4Name, s)
	}
	time.Sleep(250 * time.Millisecond) // Workaround for potential kernel lock problems
	if s, ec, _ := b.runner.Executor.Execute("ipset", "create", b.ipset6Name, "hash:ip", "family", "inet6", "timeout", "0"); ec != 0 {
		return fmt.Errorf(`failed to create ipset "%s": %s`, b.ipset6Name, s)
	}

	return nil
}

func (b *ipsetBackend) createIptablesEntries() error {
	if s, ec, _ := b.runner.Executor.Execute("iptables", "-N", b.chainName); ec != 0 {
		return fmt.Errorf(`failed to create iptables chain "%s": %s`, b.chainName, s)
	}
	if s, ec, _ := b.runner.Executor.Execute("iptables", "-I", b.chainName, "-j", "DROP", "-m", "set", "--match-set", b.ipset4Name, "src"); ec != 0 {
		return fmt.Errorf(`failed to create iptables entry for set "%s": %s`, b.ipset4Name, s)
	}
	if s, ec, _ := b.runner.Executor.Execute("iptables", "-I", "INPUT", "-j", b.chainName); ec != 0 {
		return fmt.Errorf(`failed to create iptables entry for chain "%s": %s`, b.chainName, s)
	}
	if s, ec, _ := b.runner.Executor.Execute("ip6tables", "-N", b.chainName); ec != 0 {
		return fmt.Errorf(`failed to create ip6tables chain "%s": %s`, b.chainName, s)
	}
	if s, ec, _ := b.runner.Executor.Execute("ip6tables", "-I", b.chainName, "-j", "DROP", "-m", "set", "--match-set", b.ipset6Name, "src"); ec != 0 {
		return fmt.Errorf(`failed to create ip6tables entry for set "%s": %s`, b.ipset6Name, s)
	}
	if s, ec, _ := b.runner.Executor.Execute("ip6tables", "-I", "INPUT", "-j", b.chainName); ec != 0 {
		return fmt.Errorf(`failed to create ip6tables entry for chain "%s": %s`, b.chainName, s)
	}

	return nil
}

func (b *ipsetBackend) saveIpsets() error {
	f, err := os.Create(b.runner.Configuration.SaveFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, _, err := b.runner.Executor.ExecuteWithStd(nil, f, "ipset", "save"); err != nil {
		return err
	}

	// Always ensure file is saved to disk. This should prevent loss of banned IPs on shutdown.
	return f.Sync()
}

func (b *ipsetBackend) restoreIpsets() error {
	f, err := os.Open(b.runner.Configuration.SaveFilePath)
	if err != nil {
		return err
	}

	defer func() {
		f.Close()
		if err := os.Remove(b.runner.Configuration.SaveFilePath); err != nil {
			log.Printf("failed to delete save file: %s", err)
		}
	}()

	if _, _, err := b.runner.Executor.ExecuteWithStd(f, nil, "ipset", "restore"); err != nil {
		return err
	}

	return nil
}

func (b *ipsetBackend) Initialize() error {
	b.chainName = DefaultChainName
	b.ipset4Name = DefaultTable4Name
	b.ipset6Name = DefaultTable6Name

	// Check privileges
	if s, _, err := b.runner.Executor.Execute("ipset", "list"); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return errors.New("ipset: command not found")
		}
		return fmt.Errorf("ipset: insufficient privileges: %s", s)
	}
	if s, _, err := b.runner.Executor.Execute("iptables", "-L"); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return errors.New("iptables: command not found")
		}
		return fmt.Errorf("iptables: insufficient privileges: %s", s)
	}
	if s, _, err := b.runner.Executor.Execute("ip6tables", "-L"); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return errors.New("ip6tables: command not found")
		}
		return fmt.Errorf("ip6tables: insufficient privileges: %s", s)
	}

	// Initialize ipsets and ip(6)tables entries
	if err := b.deleteIpsetsAndIptablesEntries(); err != nil {
		return fmt.Errorf("failed to delete ipsets and iptables entries: %w", err)
	}
	if b.runner.Configuration.SaveFilePath != "" {
		if err := b.restoreIpsets(); err != nil {
			if err := b.createIpsets(); err != nil {
				return fmt.Errorf("failed to create ipsets: %w", err)
			}
		} else {
			log.Printf(`restored ipsets from "%s"`, b.runner.Configuration.SaveFilePath)
		}
	} else {
		log.Printf("warning: not persisting ipsets")
		if err := b.createIpsets(); err != nil {
			return fmt.Errorf("failed to create ipsets: %w", err)
		}
	}
	if err := b.createIptablesEntries(); err != nil {
		return fmt.Errorf("failed to create ip(6)tables entries: %w", err)
	}

	return nil
}

func (b *ipsetBackend) Ban(ip string, ipv6 bool, d time.Duration) error {
	s := b.ipset4Name
	if ipv6 {
		s = b.ipset6Name
	}
	ds := int64(d.Seconds())
	if _, _, err := b.runner.Executor.Execute("ipset", "test", s, ip); err != nil {
		if _, _, err := b.runner.Executor.Execute("ipset", "add", s, ip, "timeout", fmt.Sprint(ds)); err != nil {
			return err
		}
	}
	return nil
}

func (b *ipsetBackend) Finalize() error {
	if b.runner.Configuration.SaveFilePath != "" {
		if err := b.saveIpsets(); err != nil {
			return fmt.Errorf(`failed to save ipsets to "%s": %w`, b.runner.Configuration.SaveFilePath, err)
		}
	}
	if err := b.deleteIpsetsAndIptablesEntries(); err != nil {
		return fmt.Errorf("failed to delete ipsets and ip(6)tables entries: %w", err)
	}
	return nil
}

type nftBackend struct {
	runner     *Runner
	table4Name string
	table6Name string
	set4Name   string
	set6Name   string
}

func (b *nftBackend) createTables() error {
	if s, _, err := b.runner.Executor.Execute("nft", "add", "table", "ip", b.table4Name); err != nil {
		return fmt.Errorf(`failed to add table "%s": %s`, b.table4Name, s)
	}
	if s, _, err := b.runner.Executor.Execute("nft", "add", "set", "ip", b.table4Name, b.set4Name, "{ type ipv4_addr; flags timeout; }"); err != nil {
		return fmt.Errorf(`failed to add ip set "%s": %s`, b.table4Name, s)
	}
	if s, _, err := b.runner.Executor.Execute("nft", "add", "chain", "ip", b.table4Name, "INPUT", "{ type filter hook input priority 0; policy accept; }"); err != nil {
		return fmt.Errorf(`failed to add input chain: %s`, s)
	}
	if s, _, err := b.runner.Executor.Execute("nft", "flush", "chain", "ip", b.table4Name, "INPUT"); err != nil {
		return fmt.Errorf(`failed to flush input chain: %s`, s)
	}
	if s, _, err := b.runner.Executor.Execute("nft", "add", "rule", "ip", b.table4Name, "INPUT", "ip", "saddr", "@"+b.set4Name, "reject"); err != nil {
		return fmt.Errorf(`failed to add rule: %s`, s)
	}
	if s, _, err := b.runner.Executor.Execute("nft", "add", "table", "ip6", b.table6Name); err != nil {
		return fmt.Errorf(`failed to create ip6 table "%s": %s`, b.table6Name, s)
	}
	if s, _, err := b.runner.Executor.Execute("nft", "add", "set", "ip6", b.table6Name, b.set6Name, "{ type ipv6_addr; flags timeout; }"); err != nil {
		return fmt.Errorf(`failed to add ip set "%s": %s`, b.table6Name, s)
	}
	if s, _, err := b.runner.Executor.Execute("nft", "add", "chain", "ip6", b.table6Name, "INPUT", "{ type filter hook input priority 0; policy accept; }"); err != nil {
		return fmt.Errorf(`failed to add input chain: %s`, s)
	}
	if s, _, err := b.runner.Executor.Execute("nft", "flush", "chain", "ip6", b.table6Name, "INPUT"); err != nil {
		return fmt.Errorf(`failed to flush input chain: %s`, s)
	}
	if s, _, err := b.runner.Executor.Execute("nft", "add", "rule", "ip6", b.table6Name, "INPUT", "ip6", "saddr", "@"+b.set6Name, "reject"); err != nil {
		return fmt.Errorf(`failed to add rule: %s`, s)
	}

	return nil
}

func (b *nftBackend) deleteTables() error {
	if s, _, err := b.runner.Executor.Execute("nft", "delete", "table", "ip", b.table4Name); err != nil {
		return fmt.Errorf(`failed to delete table "%s": %s`, b.table4Name, s)
	}
	if s, _, err := b.runner.Executor.Execute("nft", "delete", "table", "ip6", b.table6Name); err != nil {
		return fmt.Errorf(`failed to delete table "%s": %s`, b.table6Name, s)
	}

	return nil
}

func (b *nftBackend) saveSets() error {
	f, err := os.Create(b.runner.Configuration.SaveFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, _, err := b.runner.Executor.ExecuteWithStd(nil, f, "nft", "list", "set", "ip", b.table4Name, b.set4Name); err != nil {
		return err
	}

	if _, _, err := b.runner.Executor.ExecuteWithStd(nil, f, "nft", "list", "set", "ip6", b.table6Name, b.set6Name); err != nil {
		return err
	}

	// Always ensure file is saved to disk. This should prevent loss of banned IPs on shutdown.
	return f.Sync()
}

func (b *nftBackend) restoreSets() error {
	_, _, err := b.runner.Executor.Execute("nft", "-f", b.runner.Configuration.SaveFilePath)

	return err
}

func (b *nftBackend) Initialize() error {
	b.table4Name = DefaultTable4Name
	b.table6Name = DefaultTable6Name
	b.set4Name = "set4"
	b.set6Name = "set6"

	// Check privileges
	if s, _, err := b.runner.Executor.Execute("nft", "list", "ruleset"); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return errors.New("nft: command not found")
		}
		return fmt.Errorf("nft: insufficient privileges: %s", s)
	}

	if err := b.createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	if b.runner.Configuration.SaveFilePath != "" {
		if err := b.restoreSets(); err != nil {
			log.Printf(`failed to restore sets from "%s": %s`, b.runner.Configuration.SaveFilePath, err)
		} else {
			log.Printf(`restored sets from "%s"`, b.runner.Configuration.SaveFilePath)
		}
	} else {
		log.Printf("warning: not persisting sets")
	}

	return nil
}

func (b *nftBackend) Ban(ip string, ipv6 bool, d time.Duration) error {
	ds := int64(d.Seconds())

	t, tn, sn := "ip", b.table4Name, b.set4Name
	if ipv6 {
		t, tn, sn = "ip6", b.table6Name, b.set6Name
	}
	if s, ec, err := b.runner.Executor.Execute("nft", "add", "element", t, tn, sn, fmt.Sprintf("{ %s timeout %ds }", ip, ds)); err != nil {
		if ec == 1 {
			// This IP is probably already in set. Ignore the error. This is to be reworked
			// when support for nft < v1.0.0 is dropped. However, since Ubuntu 20.04 only has
			// v0.9.3, this is needed.
			return nil
		}
		return fmt.Errorf(`failed to add element to set "%s": %s`, b.set6Name, s)
	}

	return nil
}

func (b *nftBackend) Finalize() error {
	if b.runner.Configuration.SaveFilePath != "" {
		if err := b.saveSets(); err != nil {
			return fmt.Errorf(`failed to save sets to "%s": %w`, b.runner.Configuration.SaveFilePath, err)
		}
	}

	if err := b.deleteTables(); err != nil {
		return fmt.Errorf("failed to delete tables: %w", err)
	}

	return nil
}

type testBackend struct {
	runner        *Runner
	initializeErr error
	banErr        error
	finalizeErr   error
}

func (b *testBackend) Initialize() error {
	return b.initializeErr
}

func (b *testBackend) Ban(ip string, ipv6 bool, d time.Duration) error {
	return b.banErr
}

func (b *testBackend) Finalize() error {
	return b.finalizeErr
}
