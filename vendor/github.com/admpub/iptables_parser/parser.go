package iptables_parser

import (
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	vd "github.com/admpub/iptables_parser/validate_dns"
)

// Line represents a line in a iptables dump, e.g. generated with iptables-save.
// It is either Comment, Header, Default or Rule.
type Line interface {
	String() string
}

// Comment represents a comment in an iptables dump. Comments start with #.
type Comment struct {
	Content string
}

func (c Comment) String() string {
	return "#" + c.Content
}

// Header represents a header in an iptables dump and introduce a new table. They start with *.
type Header struct {
	Content string
}

func (h Header) String() string {
	return "*" + h.Content
}

// Policy represents a build-in policy. They can be parsed from
// iprables-save looking like ":FORWARD DROP [0:100]" They start with :.
// They can also be parsed from "iptables -S" looking like "-N|-P chain [target]".
// In the latter case, UserDefined will be set. For user defined policies, Action
// should be an empty string "" or "-".
type Policy struct {
	Chain       string
	Action      string
	UserDefined *bool // nil if unknown
	Counter     *Counter
}

func (d Policy) String() string {
	prefix := ":"
	if d.UserDefined != nil {
		if *d.UserDefined {
			prefix = "-N "
		} else {
			prefix = "-P "
		}
	}
	if d.Counter != nil {
		return fmt.Sprintf("%s%s %s %s", prefix, d.Chain, d.Action, d.Counter.String())
	}
	return fmt.Sprintf("%s%s %s", prefix, d.Chain, d.Action)
}

// Rule represents a rule in an iptables dump. Normally the start with -A.
// The parser treats the -A flag like any other flag, thus does not require
// the -A flag as the leading flag.
type Rule struct {
	Chain       string       // Name of the chain
	Source      *DNSOrIPPair // Will be nil, if -s flag was not set.
	Destination *DNSOrIPPair // Will be nil, if -s flag was not set.
	InInterf    *StringPair  // Will be nil, if -i flag was not set.
	OutInterf   *StringPair  // Will be nil, if -o flag was not set.
	Protocol    *StringPair  // Be aware that the protocol names can be different depending on your system.
	Fragment    *bool        // Will be nil, if flag was not set.
	IPv4        bool         // False, if flag was not set.
	IPv6        bool         // False, if flag was not set.
	Jump        *Target      // Will be nil, if -j flag was not set.
	Goto        *Target      // Will be nil, if -g flag was not set.
	Counter     *Counter     // Will be nil, if no counter was parsed.
	Matches     []Match      // Matches need to be a slice because order can matter. See man iptables-extension.
}

// NewRuleFromSpec returns a rule from a given rulespec and chain name.
// It will return nil and an error, if the rulespec does not resemble
// a valid rule, or contains unknown, or not implemented extensions.
func NewRuleFromSpec(chain string, rulespec ...string) (*Rule, error) {
	return NewRuleFromString(fmt.Sprintf("-A %s %s", chain, strings.Join(rulespec, " ")))
}

// NewRuleFromString returns a rule for the given string. It can only handle
// appended rules with the "-A <chain name>" flag.
// It will return nil and an error, if the given string does not resemble
// a valid rule, or contains unknown, or not implemented extensions.
func NewRuleFromString(s string) (*Rule, error) {
	return NewParser(strings.NewReader(s)).ParseRule()
}

// NewFromString takes a string a parses it until the EOF or NEWLINE
// to return a Header, Policy or Rule. It will return an error otherwise.
func NewFromString(s string) (Line, error) {
	return NewParser(strings.NewReader(s)).Parse()
}

// String returns the rule as a String, similar to how iptables-save prints the rules
// Note: Don't use this functions to compare rules because the order of flags can be
// different and different flags can have equal meanings.
func (r Rule) String() (s string) {
	s = fmt.Sprintf("-A %s %s", r.Chain, strings.Join(enquoteIfWS(r.Spec()), " "))
	if r.Counter != nil {
		s = fmt.Sprintf("%s %s", r.Counter.String(), s)
	}
	return
}

// Spec returns the rule specifications of the rule.
// The rulespec does not contain the chain name.
// Different rule specs can descibe the same rule, so
// don't use the rulespec to compare rules.
// The rule spec can be used to append, insert or delete
// rules with coreos' go-iptables module.
func (r Rule) Spec() (ret []string) {
	if r.Source != nil {
		ret = append(ret, r.Source.Spec("-s")...)
	}
	if r.Destination != nil {
		ret = append(ret, r.Destination.Spec("-d")...)
	}
	if r.InInterf != nil {
		ret = append(ret, r.InInterf.Spec("-i")...)
	}
	if r.OutInterf != nil {
		ret = append(ret, r.OutInterf.Spec("-o")...)
	}
	if r.Protocol != nil {
		ret = append(ret, r.Protocol.Spec("-p")...)
	}
	if r.Fragment != nil {
		if *r.Fragment {
			ret = append(ret, "-f")
		} else {
			ret = append(ret, "!", "-f")
		}
	}
	if r.IPv4 {
		ret = append(ret, "-4")
	}
	if r.IPv6 {
		ret = append(ret, "-6")
	}
	if len(r.Matches) > 0 {
		for _, m := range r.Matches {
			ret = append(ret, m.Spec()...)
		}
	}
	if r.Jump != nil {
		ret = append(ret, r.Jump.Spec("-j")...)
	}
	if r.Goto != nil {
		ret = append(ret, r.Goto.Spec("-g")...)
	}
	return
}

// EqualTo returns true, if the rules are
// equal to each other.
func (r Rule) EqualTo(r2 Rule) bool {
	return reflect.DeepEqual(r, r2)
}

// DNSOrIPPair either holds an IP or DNS and a flag.
// The boolean not-flag is used when an address or
// DNS name is reverted with a "!" character.
type DNSOrIPPair struct {
	Value DNSOrIP
	Not   bool
}

// String returns the part of the iptables rule. It requires its flag as string
// to generate the correct string, e.g. "! -s 10.0.0.1/32".
func (d DNSOrIPPair) String(f string) string {
	return strings.Join(d.Spec(f), " ")
}

// Spec returns a DNSOrIPPair how coreos' iptables package would expect it.
func (d DNSOrIPPair) Spec(f string) []string {
	s := []string{"!", f, d.Value.String()}
	if !d.Not {
		return s[1:]
	}
	return s
}

// DNSOrIP represents either a DNS name or an IP address.
// IPs, as they are more specific, are preferred.
type DNSOrIP struct {
	// DNS must be a valid RFC 1123 subdomain.
	// +optional
	dNS string
	// IP must be a valid IP address.
	// +optional
	iP net.IPNet
}

// Set IP if string is a valid IP address, or DNS if string is a valid DNS name,
// else return error.
func (d *DNSOrIP) Set(s string) error {
	sn := s
	// TODO: this can probably be done in a nicer way.
	if !strings.Contains(sn, "/") {
		sn = sn + "/32"
	}
	if _, ipnet, err := net.ParseCIDR(sn); err == nil {
		d.iP = *ipnet
		d.dNS = ""
		return nil
	}
	if !vd.IsDNS(s) {
		return fmt.Errorf("%q is not a valid DNS name", s)
	}
	d.dNS = s
	return nil
}

func (d *DNSOrIP) String() string {
	if d.dNS != "" {
		return d.dNS
	}
	return d.iP.String()
}

// NewDNSOrIP takes a string and return a DNSOrIP, or an error.
// It tries to parse it as an IP, if this fails it will check,
// whether the input is a valid DNS name.
func NewDNSOrIP(s string) (*DNSOrIP, error) {
	ret := &DNSOrIP{}
	if err := ret.Set(s); err != nil {
		return nil, err
	}
	return ret, nil
}

// StringPair is a string with a flag.
// It is used to represent flags that specify a string value
// and can be negated with a "!".
type StringPair struct {
	Not   bool
	Value string
}

func (sp StringPair) String(f string) string {
	return strings.Join(sp.Spec(f), " ")
}

// Spec returns a StringPair how coreos' iptables package would expect it.
func (sp StringPair) Spec(f string) []string {
	ret := []string{"!", f, sp.Value}
	if !sp.Not {
		return ret[1:]
	}
	return ret
}

// Counter represents the package and byte counters.
type Counter struct {
	packets uint64
	bytes   uint64
}

func (c Counter) String() string {
	return fmt.Sprintf("[%d:%d]", c.packets, c.bytes)
}

// Match represents one match expression from the iptables-extension.
// See man iptables-extenstion for more info.
type Match struct {
	Type  string
	Flags map[string]Flag
}

func (m Match) String() string {
	return strings.Join(m.Spec(), " ")
}

// Spec returns a Match how coreos' iptables package would expect it.
func (m Match) Spec() []string {
	ret := make([]string, 2, 2+len(m.Flags)*2)
	ret[0], ret[1] = "-m", m.Type
	for k, val := range m.Flags {
		ret = append(ret, val.Spec("--"+k)...)
	}
	return ret
}

// Flag is flag, e.g. --dport 8080. It can be negated with a leading !.
// Sometimes a flag is followed by several arguments.
type Flag struct {
	Not    bool
	Values []string
}

func (fl Flag) String(f string) string {
	return strings.Join(fl.Spec(f), " ")
}

// Spec returns a Flag how coreos' iptables package would expect it.
func (fl Flag) Spec(f string) []string {
	ret := []string{"!", f}
	ret = append(ret, fl.Values...)
	if !fl.Not {
		return ret[1:]
	}
	return ret
}

// Target represents a Target Extension. See iptables-extensions(8).
type Target struct {
	Name  string
	Flags map[string]Flag
}

func (t Target) String(name string) string {
	return strings.Join(t.Spec(name), " ")
}

// Spec returns a Target how coreos' iptables package would expect it.
func (t Target) Spec(f string) []string {
	ret := make([]string, 2, 2+len(t.Flags)*2)
	ret[0], ret[1] = f, t.Name
	for k, val := range t.Flags {
		ret = append(ret, val.Spec("--"+k)...)
	}
	return ret
}

// BUFSIZE is the max buffer size of the ring buffer in the parser.
const BUFSIZE = 16

// Parser represents a parser.
type Parser struct {
	s   *scanner
	buf struct {
		toks [BUFSIZE]Token  // token buffer
		lits [BUFSIZE]string // literal buffer
		p    int             // current position in the buffer (max=BUF_SIZE)
		n    int             // offset (max=BUF_SIZE)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: newScanner(r)}
}

// Parse parses one line and returns a Rule, Comment, Header or DEFAULT.
func (p *Parser) Parse() (l Line, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recover from panic: %v", r)
			return
		}
	}()
	tok, lit := p.scanIgnoreWhitespace()
	switch tok {
	case COMMENTLINE:
		return Comment{Content: lit}, nil
	case HEADER:
		return Header{Content: lit}, nil
	case FLAG:
		p.unscan(1)
		return p.parseRule()
	case COLON:
		return p.parseDefault(p.s.scanLine())
	case EOF:
		return nil, io.EOF // ErrEOF
	case NEWLINE:
		return nil, errors.New("empty line")
	default:
		return nil, fmt.Errorf("unexpected format of first token: %s, skipping rest %q of the line", lit, p.s.scanLine())
	}
}

func (p *Parser) ParseRule() (*Rule, error) {
	l, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse rule: %w", err)
	}
	switch r := l.(type) {
	case Rule:
		return &r, nil
	default:
		return nil, errors.New("failed to cast Line to type Rule")
	}
}

var (
	matchModules     map[string]struct{}
	targetExtensions map[string]struct{}
)

func init() {
	matchModules = make(map[string]struct{})
	for _, m := range []string{"addrtype", "ah", "bpf", "cgroup", "cluster", "comment", "connbytes", "connlabel", "connlimit", "connmark", "conntrack", "cpu", "dccp", "devgroup", "dscp", "dst", "ecn", "esp", "eui64", "frag", "hashlimit", "hbh", "helper", "hl", "icmp", "icmp6", "iprange", "ipv6header", "ipvs", "length", "limit", "mac", "mark", "mh", "multiport", "nfacct", "osf", "owner", "physdev", "pkttype", "policy", "quota", "rateest", "realm", "recent", "rpfilter", "rt", "sctp", "set", "socket", "state", "statistics", "string", "tcp", "tcpmss", "time", "tos", "ttl", "u32", "udp"} {
		matchModules[m] = struct{}{}
	}
	targetExtensions = make(map[string]struct{})
	for _, e := range []string{"AUDIT", "CHECKSUM", "CLASSIFY", "CLUSTERIP", "CONNMARK", "CONNSECMARK", "CT", "DNAT", "DNPT", "DSCP", "ECN", "HL", "HMARK", "IDLETIMER", "LED", "LOG", "MARK", "MASQUERADE", "NETMAP", "NFLOG", "NFQUEUE", "NOTRACK", "RATEEST", "REDIRECT", "REJECT", "SECMARK", "SET", "SNAT", "SNPT", "SYNPROXY", "TCPMSS", "TCPOPTSTRIP", "TEE", "TOS", "TPROXY", "TRACE", "TTL", "ULOG"} {
		targetExtensions[e] = struct{}{}
	}
}

var (
	regDefault *regexp.Regexp = regexp.MustCompile(`^\s*(\S+)\s+(\S+)\s+(\[\d*\:\d*\])\s*$`)
	regCounter *regexp.Regexp = regexp.MustCompile(`^\[(\d*)\:(\d*)\]$`)
)

func (p *Parser) parseDefault(lit string) (Line, error) {
	var r Policy
	r.Chain = string(regDefault.ReplaceAll([]byte(lit), []byte("$1")))
	a := regDefault.ReplaceAll([]byte(lit), []byte("$2"))
	r.Action = string(a)
	cs := regDefault.ReplaceAll([]byte(lit), []byte("$3"))
	c, err := parseCounter(cs)
	if err != nil {
		return nil, err
	}

	r.Counter = &c
	return r, nil
}

// parseCounter parses something like "[0:100]"
func parseCounter(bytes []byte) (Counter, error) {
	var c Counter
	pc := regCounter.ReplaceAll(bytes, []byte("$1"))
	i, err := strconv.ParseUint(string(pc), 10, 0)
	if err != nil {
		return c, fmt.Errorf("Could not parse counter: %w", err)
	}
	c.packets = i
	pc = regCounter.ReplaceAll(bytes, []byte("$2"))
	i, err = strconv.ParseUint(string(pc), 10, 0)
	if err != nil {
		return c, fmt.Errorf("Could not parse counter: %w", err)
	}
	c.bytes = i
	return c, nil
}

// State for the state machine
type state int

const (
	// Only use even numbers to have some local states, that can be odd numbers.
	sStart state = iota * 2 // Start state
	sIF                     // Interpret a flag
	sINotF                  // Interprete flag with NOT
	sNot                    // NOT state
	sA                      // append rule
	sN                      // user defined chain
	sP                      // policy for build-in chain
	sError
)

func (p *Parser) parseRule() (Line, error) {
	var r Rule
	s := sStart
	var err error
	for tok, lit := p.scanIgnoreWhitespace(); tok != EOF && tok != NEWLINE; tok, lit = p.scanIgnoreWhitespace() {
		for nextValue := false; !nextValue; {
			nextValue = true
			switch s {
			case sStart:
				switch tok {
				case FLAG:
					s = sIF
					nextValue = false
				case NOT:
					s = sNot
				default:
					s = sError
					break
				}
			case sIF:
				switch {
				case isSrc(lit):
					r.Source = new(DNSOrIPPair)
					s, err = p.parseAddr(r.Source, false)
				case isDest(lit):
					r.Destination = new(DNSOrIPPair)
					s, err = p.parseAddr(r.Destination, false)
				case lit == "-p" || lit == "--protocol":
					r.Protocol = new(StringPair)
					s, err = p.parseProtocol(r.Protocol, false)
				case isMatch(lit):
					s, err = p.parseMatch(&r.Matches)
				case lit == "-A" || lit == "--append":
					s = sA
				case lit == "-P" || lit == "--policy":
					s = sP
				case lit == "-N" || lit == "--new-chain":
					s = sN
				case lit == "-j" || lit == "--jump":
					r.Jump = new(Target)
					s, err = p.parseTarget(r.Jump)
				case lit == "-g" || lit == "--goto":
					r.Goto = new(Target)
					s, err = p.parseTarget(r.Goto)
				case lit == "-i" || lit == "--in-interface":
					r.InInterf = new(StringPair)
					s, err = p.parseStringPair(r.InInterf, false)
				case lit == "-o" || lit == "--out-interface":
					r.OutInterf = new(StringPair)
					s, err = p.parseStringPair(r.OutInterf, false)
				case lit == "-f" || lit == "--fragment":
					_true := true
					r.Fragment = &(_true)
					s = sStart
				case lit == "-4" || lit == "--ipv4":
					r.IPv4 = true
					s = sStart
				case lit == "-6" || lit == "--ipv6":
					r.IPv6 = true
					s = sStart
				default:
					err = fmt.Errorf("unknown flag %q found", lit)
					s = sError
				}
			case sINotF:
				switch {
				case isSrc(lit):
					r.Source = new(DNSOrIPPair)
					s, err = p.parseAddr(r.Source, true)
				case isDest(lit):
					r.Destination = new(DNSOrIPPair)
					s, err = p.parseAddr(r.Destination, true)
				case lit == "-p" || lit == "--protocol":
					r.Protocol = new(StringPair)
					s, err = p.parseProtocol(r.Protocol, true)
				case lit == "-i" || lit == "--in-interface":
					r.InInterf = new(StringPair)
					s, err = p.parseStringPair(r.InInterf, true)
				case lit == "-o" || lit == "--out-interface":
					r.OutInterf = new(StringPair)
					s, err = p.parseStringPair(r.OutInterf, true)
				case lit == "-f" || lit == "--fragment":
					_false := false
					r.Fragment = &(_false)
					s = sStart
				default:
					err = fmt.Errorf("encountered unknown flag %q, or flag can not be negated with \"!\"", lit)
					s = sError
				}
			case sA:
				r.Chain = lit
				s = sStart
			case sN:
				r.Chain = lit
				s = sStart
				p.unscan(1)
				return p.parseUserDefinedPolicy()
			case sP:
				r.Chain = lit
				s = sStart
				p.unscan(1)
				return p.parseDefaultPolicy()
			case sNot:
				switch tok {
				case FLAG:
					nextValue = false
					s = sINotF
				default:
					err = fmt.Errorf("unexpected token %q, expected identifier", lit)
					s = sError
				}
			case sError:
				return nil, fmt.Errorf("failed to parse line, skipping rest %q of the line: %w", p.s.scanLine(), err)
			default:
				nextValue = true

			}
			// Avoid scanning the next token, if an error occured.
			nextValue = nextValue && err == nil
		}
	}
	return r, nil
}

func (p *Parser) parseUserDefinedPolicy() (Line, error) {
	return p.parsePolicy(true)
}

func (p *Parser) parseDefaultPolicy() (Line, error) {
	return p.parsePolicy(false)
}

func (p *Parser) parsePolicy(d bool) (Line, error) {
	ret := Policy{
		UserDefined: new(bool), // create a new pointer in case the caller keeps using the input variable. This ain't rust.
	}
	*ret.UserDefined = d
	if tok, lit := p.scanIgnoreWhitespace(); tok != EOF && tok != NEWLINE {
		ret.Chain = lit
	} else {
		return nil, errors.New("unexpected end of line")
	}
	if tok, lit := p.scanIgnoreWhitespace(); tok != EOF && tok != NEWLINE {
		ret.Action = lit
	} else {
		return ret, nil
	}
	if tok, lit := p.scanIgnoreWhitespace(); tok != EOF && tok != NEWLINE {
		return nil, fmt.Errorf("found %q, expected EOF or newline", lit)
	}
	return ret, nil
}

// parseProtocol is not restricted on protocol types because the names
// can depend on the underlying system. E.g. ipv4 is called ipencap
// in Gentoo based systems.
func (p *Parser) parseProtocol(r *StringPair, not bool) (state, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok == NEWLINE || tok == EOF {
		return sError, errors.New("unexpected end of line while parsing protocol")
	}
	*r = StringPair{
		Not:   not,
		Value: lit,
	}
	return sStart, nil
}

func (p *Parser) parseAddr(r *DNSOrIPPair, not bool) (state, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok == NEWLINE || tok == EOF {
		return sError, errors.New("unexpected end of line while parsing address")
	}
	doi, err := NewDNSOrIP(lit)
	if err != nil {
		return sError, err
	}
	*r = DNSOrIPPair{Value: *doi, Not: not}
	return sStart, nil
}

// parseStringPair
func (p *Parser) parseStringPair(sp *StringPair, not bool) (state, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		*sp = StringPair{Value: "", Not: not}
		p.unscan(1)
		return sStart, errors.New("unexpected token, expected IDENT")
	}
	*sp = StringPair{Value: lit, Not: not}
	return sStart, nil
}

// mod is not the remainder %, but the modulo function with -a%b != -(a%b).
// Python has such implementation.
func mod(a, b int) int {
	return (a%b + b) % b
}

// scan returns the next token from the underlying scanner.
// If tokens bave been unscanned then read the previous one instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, return it.
	if p.buf.n != 0 {
		p.buf.n--
		return p.buf.toks[mod(p.buf.p-p.buf.n-1, BUFSIZE)], p.buf.lits[mod(p.buf.p-p.buf.n-1, BUFSIZE)]
	}
	// Otherwise read the next token from the scanner.
	tok, lit = p.s.scan()
	// Save it to the buffer in case we unscan later.
	p.buf.toks[p.buf.p], p.buf.lits[p.buf.p] = tok, lit
	p.buf.p++ // increase the pointer of the ring buffer.
	p.buf.p %= BUFSIZE
	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	for tok == WS {
		tok, lit = p.scan()
	}
	return
}

// unscan reverts the pointer on the buffer, callers should not unscan more then what was
// previously read, or values larger then BUF_SIZE.
func (p *Parser) unscan(n int) {
	if p.buf.n+n >= BUFSIZE {
		panic("size exceeds buffer")
	}
	p.buf.n += n
}

var hasWS *regexp.Regexp = regexp.MustCompile(`\s`)

func enquoteIfWS(s []string) []string {
	ret := make([]string, len(s))
	for i, e := range s {
		if hasWS.MatchString(e) {
			ret[i] = fmt.Sprintf("%q", e)
		} else {
			ret[i] = e
		}
	}
	return ret
}

func isSrc(s string) bool {
	return s == "-s" || s == "--src" || s == "--source"
}

func isDest(s string) bool {
	return s == "-d" || s == "--dest" || s == "--dst" || s == "--destination"
}

func isMatch(s string) bool {
	return s == "-m" || s == "--match"
}
