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
