package gerberos

import (
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

type Match struct {
	Time   time.Time
	Line   string
	IP     string
	IPv6   bool
	Regexp *regexp.Regexp
}

func (r *Rule) MatchSimple(line string) (*Match, error) {
	for _, re := range r.regexp {
		m := re.FindStringSubmatch(line)
		if len(m) == 0 {
			continue
		}

		sm := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				sm[name] = m[i]
			}
		}
		h := sm["ip"]
		h = strings.Trim(h, "[]")
		ph := net.ParseIP(h)
		if ph == nil {
			return nil, fmt.Errorf(`failed to parse matched IP "%s"`, h)
		}

		return &Match{
			Line:   line,
			Time:   time.Now(),
			IP:     h,
			IPv6:   ph.To4() == nil,
			Regexp: re,
		}, nil
	}

	return nil, fmt.Errorf(`line "%s" does not match any regexp`, line)
}

func (r *Rule) MatchAggregate(line string) (*Match, error) {
	a := r.aggregate

	for _, re := range a.regexp {
		m := re.FindStringSubmatch(line)
		if len(m) == 0 {
			continue
		}

		sm := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				sm[name] = m[i]
			}
		}
		id := sm["id"]

		a.registryMutex.Lock()
		if ip, e := a.registry[id]; e {
			delete(a.registry, id)
			a.registryMutex.Unlock()

			return &Match{
				Line:   line,
				Time:   time.Now(),
				IP:     ip.String(),
				IPv6:   ip.To4() == nil,
				Regexp: re,
			}, nil
		}
		a.registryMutex.Unlock()
	}

	for _, re := range r.regexp {
		m := re.FindStringSubmatch(line)
		if len(m) == 0 {
			continue
		}

		sm := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				sm[name] = m[i]
			}
		}
		h := sm["ip"]
		h = strings.Trim(h, "[]")
		pip := net.ParseIP(h)
		if pip == nil {
			return nil, fmt.Errorf(`failed to parse matched IP "%s"`, h)
		}

		id := sm["id"]
		if id == "" {
			return nil, fmt.Errorf(`failed to match ID`)
		}

		a.registryMutex.Lock()
		a.registry[id] = pip
		if r.runner.Configuration.Verbose {
			log.Printf(`%s: added ID "%s" with IP %s to registry`, r.name, id, pip)
		}
		a.registryMutex.Unlock()

		go func(id string) {
			time.Sleep(a.interval)
			a.registryMutex.Lock()
			if ip, e := a.registry[id]; e {
				delete(a.registry, id)
				if r.runner.Configuration.Verbose {
					log.Printf(`%s: removed ID "%s" with IP %s from registry`, r.name, id, ip)
				}
			}
			a.registryMutex.Unlock()
		}(id)

		return nil, errors.New("incomplete aggregate")
	}

	return nil, fmt.Errorf(`line "%s" does not match any regexp`, line)
}

func (r *Rule) Match(line string) (*Match, error) {
	if r.aggregate != nil {
		return r.MatchAggregate(line)
	}

	return r.MatchSimple(line)
}

func (m Match) stringSimple() string {
	ipv := "IPv4"
	if m.IPv6 {
		ipv = "IPv6"
	}

	return fmt.Sprintf(`time = %s, IP = "%s", %s`, m.Time.Format(time.RFC3339), m.IP, ipv)
}

func (m Match) StringExtended() string {
	return fmt.Sprintf(`%s, line = "%s", regexp = "%s"`, m, m.Line, m.Regexp)
}

func (m Match) String() string {
	return m.stringSimple()
}
