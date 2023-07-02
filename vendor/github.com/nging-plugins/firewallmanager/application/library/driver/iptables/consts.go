package iptables

import "github.com/nging-plugins/firewallmanager/application/library/enums"

const CommentPrefix = `NgingStatic`

const (
	FilterChainInput    = `NgingFilterInput`
	FilterChainOutput   = `NgingFilterOutput`
	FilterChainForward  = `NgingFilterForward`
	NATChainPreRouting  = `NgingNATPreRouting`
	NATChainPostRouting = `NgingNATPostRouting`
)

var (
	RefFilterChains = map[string]string{
		FilterChainInput:   enums.ChainInput,
		FilterChainOutput:  enums.ChainOutput,
		FilterChainForward: enums.ChainForward,
	}
	RefNATChains = map[string]string{
		NATChainPreRouting:  enums.ChainPreRouting,
		NATChainPostRouting: enums.ChainPostRouting,
	}
	FilterChains = []string{FilterChainInput, FilterChainOutput, FilterChainForward}
	NATChains    = []string{NATChainPreRouting, NATChainPostRouting}
)
