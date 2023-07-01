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
	SET_TRUST   = 1 // add filterSetTrustIP
	SET_MANAGER = 2 // add filterSetMyManagerIP
	SET_FORWARD = 4 // add filterSetMyForwardIP
	SET_ALL     = 8 // add filterSetTrustIP filterSetMyManagerIP filterSetMyForwardIP
)
