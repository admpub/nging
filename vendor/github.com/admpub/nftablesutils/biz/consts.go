package biz

const (
	TableFilter = `filter`
	TableNAT    = `nat`
	TableMangle = `mangle`
	TableRaw    = `raw`
)

const (
	ChainInput       = `INPUT`
	ChainOutput      = `OUTPUT`
	ChainForward     = `FORWARD`
	ChainPreRouting  = `PREROUTING`
	ChainPostRouting = `POSTROUTING`
)

const (
	ApplyTypeHTTP = `http`
	ApplyTypeSMTP = `smtp`
	ApplyTypeDNS  = `smtp`
)

var ApplyAll = []string{
	ApplyTypeHTTP,
	ApplyTypeSMTP,
	ApplyTypeDNS,
}

const (
	SET_TRUST     = 1   // add filterSetTrustIP
	SET_MANAGER   = 2   // add filterSetManagerIP
	SET_FORWARD   = 4   // add filterSetForwardIP
	SET_BLACKLIST = 8   // add filterSetBlacklistIP
	SET_ALL       = 512 // add filterSetTrustIP filterSetManagerIP filterSetForwardIP filterSetBlacklistIP
)

const (
	RULE_LOCAL_IFACE        = 1
	RULE_WAN_IFACE          = 2
	RULE_SDN                = 4
	RULE_SDN_FORWARD        = 8
	RULE_NAT                = 16
	RULE_BLACKLIST          = 32
	RULE_INPUT_LOCAL_IFACE  = 64
	RULE_OUTPUT_LOCAL_IFACE = 128
	RULE_ALL                = 512
)
