package iptables_parser

import (
	"errors"
	"fmt"
	"strings"
)

func (p *Parser) parseMatch(ms *[]Match) (state, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return sError, fmt.Errorf("unexpected token %q encountered, expected identifier", lit)
	}

	var m Match
	m.Flags = make(map[string]Flag)
	m.Type = lit
	var err error
	s := sStart
	switch lit {
	case "comment":
		s, err = p.parseComment(&m.Flags)
	case "tcp":
		s, err = p.parseTcp(&m.Flags)
	case "addrtype":
		s, err = p.parseAddrtype(&m.Flags)
	case "udp":
		s, err = p.parseUdp(&m.Flags)
	case "statistic":
		s, err = p.parseStatistic(&m.Flags)
	case "multiport":
		s, err = p.parseMultiport(&m.Flags)
	default:
		if _, ok := matchModules[lit]; ok {
			return sError, fmt.Errorf("match modules %q is not implemented", lit)
		}
		return sError, fmt.Errorf("unknown match modules: %q", lit)
	}
	if err != nil {
		return sError, err
	}
	*ms = append(*ms, m)
	return s, nil
}

func (p *Parser) parseStatistic(f *map[string]Flag) (state, error) {
	s := sStart
	for tok, lit := p.scanIgnoreWhitespace(); tok != EOF; tok, lit = p.scanIgnoreWhitespace() {
		for nextValue := false; !nextValue; {
			nextValue = true
			switch s {
			case sStart:
				switch tok {
				case NOT:
					s = sINotF
				case FLAG:
					s = sIF
					nextValue = false
				default:
					return sError, fmt.Errorf("unexpected token %q, expected flag, or \"!\"", lit)
				}
			case sINotF:
				switch {
				case lit == "--probability":
					_, lit := p.scanIgnoreWhitespace()
					(*f)["probability"] = Flag{
						Not:    true,
						Values: []string{lit},
					}
					s = sStart
				case lit == "--every":
					_, lit := p.scanIgnoreWhitespace()
					(*f)["every"] = Flag{
						Not:    true,
						Values: []string{lit},
					}
					s = sStart
				default:
					// The end of the match statement is reached.
					// Since we already scanned the ! charackter,
					// we have to return the sNot state, or
					// unscanIgnoreWhitespace twice (this can fail
					// because of a fixed sized buffer, that is full
					// of Whitespaces).
					p.unscan(1) // IgnoreWhitespace(2) // unscan 2
					return sNot, nil
				}
			case sIF:
				switch {
				case lit == "--mode":
					_, lit := p.scanIgnoreWhitespace()
					(*f)["mode"] = Flag{
						Values: []string{lit},
					}
					s = sStart
				case lit == "--probability":
					_, lit := p.scanIgnoreWhitespace()
					(*f)["probability"] = Flag{
						Values: []string{lit},
					}
					s = sStart
				case lit == "--packet":
					_, lit := p.scanIgnoreWhitespace()
					(*f)["packet"] = Flag{
						Values: []string{lit},
					}
					s = sStart
				case lit == "--every":
					_, lit := p.scanIgnoreWhitespace()
					(*f)["every"] = Flag{
						Values: []string{lit},
					}
					s = sStart
				default:
					// The end of the match statement is reached.
					p.unscan(1)
					return sStart, nil
				}

			default:
				return sStart, errors.New("unexpected error parsing match extension")
			}
		}
	}
	return sStart, nil
}

func (p *Parser) parseUdp(f *map[string]Flag) (state, error) {
	s := sStart
	for tok, lit := p.scanIgnoreWhitespace(); tok != EOF; tok, lit = p.scanIgnoreWhitespace() {
		for nextValue := false; !nextValue; {
			nextValue = true
			switch s {
			case sStart:
				switch tok {
				case NOT:
					s = sINotF
				case FLAG:
					s = sIF
					nextValue = false
				default:
					return sError, fmt.Errorf("unexpected token %q, expected flag, or \"!\"", lit)
				}
			case sINotF:
				switch {
				case lit == "--dport" || lit == "--destination-port":
					(*f)["destination-port"] = Flag{
						Not:    true,
						Values: []string{p.parsePort()},
					}
					s = sStart
				case lit == "--sport" || lit == "--source-port":
					(*f)["source-port"] = Flag{
						Not:    true,
						Values: []string{p.parsePort()},
					}
					s = sStart
				default:
					// The end of the match statement is reached.
					// Since we already scanned the ! charackter,
					// we have to return the sNot state, or
					// unscanIgnoreWhitespace twice (this can fail
					// because of a fixed sized buffer, that is full
					// of Whitespaces).
					p.unscan(1) // IgnoreWhitespace(2) // unscan 2
					return sNot, nil
				}
			case sIF:
				switch {
				case lit == "--dport" || lit == "--destination-port":
					(*f)["destination-port"] = Flag{
						Values: []string{p.parsePort()},
					}
					s = sStart
				case lit == "--sport" || lit == "--source-port":
					(*f)["source-port"] = Flag{
						Values: []string{p.parsePort()},
					}
					s = sStart
				default:
					// The end of the match statement is reached.
					p.unscan(1)
					return sStart, nil
				}

			default:
				return sStart, errors.New("unexpected error parsing match extension")
			}
		}
	}
	return sStart, nil
}

func (p *Parser) parseAddrtype(f *map[string]Flag) (state, error) {
	s := sStart
	for tok, lit := p.scanIgnoreWhitespace(); tok != EOF; tok, lit = p.scanIgnoreWhitespace() {
		for nextValue := false; !nextValue; {
			nextValue = true
			switch s {
			case sStart:
				switch tok {
				case NOT:
					s = sINotF
				case FLAG:
					s = sIF
					nextValue = false
				default:
					return sError, fmt.Errorf("unexpected token %q, expected flag, or \"!\"", lit)
				}
			case sINotF:
				switch lit {
				case "--src-type":
					_, lit := p.scanIgnoreWhitespace()
					(*f)["src-type"] = Flag{
						Not:    true,
						Values: []string{lit},
					}
					s = sStart
				case "--dst-type":
					_, lit := p.scanIgnoreWhitespace()
					(*f)["dst-type"] = Flag{
						Not:    true,
						Values: []string{lit},
					}
					s = sStart
				default:
					// The end of the match statement is reached.
					// Since we already scanned the ! charackter,
					// we have to return the sNot state, or
					// unscanIgnoreWhitespace twice (this can fail
					// because of a fixed sized buffer, that is full
					// of Whitespaces).
					p.unscan(1) // IgnoreWhitespace(2) // unscan 2
					return sNot, nil
				}
			case sIF:
				switch lit {
				case "--src-type":
					_, lit := p.scanIgnoreWhitespace()
					(*f)["src-type"] = Flag{
						Values: []string{lit},
					}
					s = sStart
				case "--dst-type":
					_, lit := p.scanIgnoreWhitespace()
					(*f)["dst-type"] = Flag{
						Values: []string{lit},
					}
					s = sStart
				case "--limit-iface-in":
					_, lit := p.scanIgnoreWhitespace()
					(*f)["limit-iface-in"] = Flag{
						Values: []string{lit},
					}
					s = sStart
				case "--limit-iface-out":
					_, lit := p.scanIgnoreWhitespace()
					(*f)["limit-iface-out"] = Flag{
						Values: []string{lit},
					}
					s = sStart
				default:
					// The end of the match statement is reached.
					p.unscan(1)
					return sStart, nil
				}

			default:
				return sStart, errors.New("unexpected error parsing match extension")
			}
		}
	}
	return sStart, nil
}

func (p *Parser) parseTcp(f *map[string]Flag) (state, error) {
	s := sStart
	for tok, lit := p.scanIgnoreWhitespace(); tok != EOF; tok, lit = p.scanIgnoreWhitespace() {
		for nextValue := false; !nextValue; {
			nextValue = true
			switch s {
			case sStart:
				switch tok {
				case NOT:
					s = sINotF
				case FLAG:
					s = sIF
					nextValue = false
				default:
					return sError, fmt.Errorf("unexpected token %q, expected flag, or \"!\"", lit)
				}
			case sINotF:
				switch {
				case lit == "--dport" || lit == "--destination-port":
					(*f)["destination-port"] = Flag{
						Not:    true,
						Values: []string{p.parsePort()},
					}
					s = sStart
				case lit == "--sport" || lit == "--source-port":
					(*f)["source-port"] = Flag{
						Not:    true,
						Values: []string{p.parsePort()},
					}
					s = sStart
				case lit == "--syn":
					(*f)["syn"] = Flag{
						Not: true,
					}
					s = sStart
				case lit == "--tcp-option":
					_, l := p.scanIgnoreWhitespace()
					(*f)["tcp-option"] = Flag{
						Not:    true,
						Values: []string{l},
					}
					s = sStart
				case lit == "--tcp-flags":
					t, _ := p.scanIgnoreWhitespace()
					p.unscan(1)
					list := p.parseList()
					if t == EOF {
						return sError, errors.New("unexpected EOF")
					}

					_, l2 := p.scanIgnoreWhitespace()
					(*f)["tcp-flags"] = Flag{
						Not:    true,
						Values: []string{strings.Join(list, ","), l2},
					}
					s = sStart
				default:
					// The end of the match statement is reached.
					// Since we already scanned the ! charackter,
					// we have to return the sNot state, or
					// unscanIgnoreWhitespace twice (this can fail
					// because of a fixed sized buffer, that is full
					// of Whitespaces).
					p.unscan(1) // IgnoreWhitespace(2) // unscan 2
					return sNot, nil
				}
			case sIF:
				switch {
				case lit == "--dport" || lit == "--destination-port":
					(*f)["destination-port"] = Flag{
						Values: []string{p.parsePort()},
					}
					s = sStart
				case lit == "--sport" || lit == "--source-port":
					(*f)["source-port"] = Flag{
						Values: []string{p.parsePort()},
					}
					s = sStart
				case lit == "--syn":
					(*f)["syn"] = Flag{}
					s = sStart
				case lit == "--tcp-option":
					_, l := p.scanIgnoreWhitespace()
					(*f)["tcp-option"] = Flag{
						Values: []string{l},
					}
					s = sStart
				case lit == "--tcp-flags":
					t, _ := p.scanIgnoreWhitespace()
					p.unscan(1)
					list := p.parseList()
					if t == EOF {
						return sError, errors.New("unexpected EOF")
					}

					_, l2 := p.scanIgnoreWhitespace()
					(*f)["tcp-flags"] = Flag{
						Values: []string{strings.Join(list, ","), l2},
					}
					s = sStart
				default:
					// The end of the match statement is reached.
					p.unscan(1)
					return sStart, nil
				}

			default:
				return sStart, errors.New("unexpected error parsing match extension")
			}
		}
	}
	return sStart, nil
}

func (p *Parser) parseMultiport(f *map[string]Flag) (state, error) {
	s := sStart
	for tok, lit := p.scanIgnoreWhitespace(); tok != EOF; tok, lit = p.scanIgnoreWhitespace() {
		for nextValue := false; !nextValue; {
			nextValue = true
			switch s {
			case sStart:
				switch tok {
				case NOT:
					s = sINotF
				case FLAG:
					s = sIF
					nextValue = false
				default:
					return sError, fmt.Errorf("unexpected token %q, expected flag, or \"!\"", lit)
				}
			case sINotF:
				switch {
				case lit == "--dports" || lit == "--destination-ports":
					(*f)["destination-ports"] = Flag{
						Not:    true,
						Values: p.parsePorts(),
					}
					s = sStart
				case lit == "--sports" || lit == "--source-ports":
					(*f)["source-ports"] = Flag{
						Not:    true,
						Values: p.parsePorts(),
					}
					s = sStart
				case lit == "--ports":
					(*f)["ports"] = Flag{
						Not:    true,
						Values: p.parsePorts(),
					}
					s = sStart
				default:
					p.unscan(1)
					return sNot, nil
				}
			case sIF:
				switch {
				case lit == "--dports" || lit == "--destination-ports":
					(*f)["destination-ports"] = Flag{
						Values: p.parsePorts(),
					}
					s = sStart
				case lit == "--sports" || lit == "--source-ports":
					(*f)["source-ports"] = Flag{
						Values: p.parsePorts(),
					}
					s = sStart
				case lit == "--ports":
					(*f)["ports"] = Flag{
						Not:    true,
						Values: p.parsePorts(),
					}
					s = sStart
				default:
					// The end of the match statement is reached.
					p.unscan(1)
					return sStart, nil
				}

			default:
				return sStart, errors.New("unexpected error parsing match extension")
			}
		}
	}
	return sStart, nil
}

func (p *Parser) parsePort() string {
	_, l := p.scanIgnoreWhitespace()
	if t, _ := p.scanIgnoreWhitespace(); t == COLON {
		_, c := p.scan()
		l = l + ":" + c
	} else {
		p.unscan(1)
	}
	return l
}

func (p *Parser) parsePorts() []string {
	var ports []string
	_, l := p.scanIgnoreWhitespace()
	ports = append(ports, l)
	for {
		if t, _ := p.scanIgnoreWhitespace(); t == COMMA {
			_, c := p.scan()
			ports = append(ports, c)
		} else {
			p.unscan(1)
			break
		}
	}
	return ports
}

func (p *Parser) parseList() (strs []string) {
	const (
		sC state = iota*2 + 1
	)
	s := sStart
	for tok, lit := p.scan(); tok != EOF; tok, lit = p.scanIgnoreWhitespace() {
		switch s {
		case sStart:
			if tok != IDENT && tok != QUOTATION {
				p.unscan(1)
				return
			}
			strs = append(strs, lit)
			s = sC
		case sC:
			if tok != COMMA {
				p.unscan(1)
				return
			}
			s = sStart
		default:
			panic("this should not have happended")
		}
	}
	return
}

func (p *Parser) parseComment(f *map[string]Flag) (state, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != FLAG {
		return sError, errors.New("unexpected token, expected flag \"--comment\"")
	}
	switch lit {
	case "--comment":
		_, lit := p.scanIgnoreWhitespace()
		(*f)["comment"] = Flag{
			Values: []string{lit},
		}
	default:
		return sError, fmt.Errorf("unexpected flag %q, expected \"--comment\"", lit)
	}
	return sStart, nil
}
