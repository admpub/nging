package dashboard

var Default = &Dashboards{
	Backend: &Dashboard{
		Cards:         cards,
		Blocks:        &blocks,
		GlobalHeads:   &globalHeads,
		GlobalFooters: &globalFooters,
		TopButtons:    &topButtons,
	},
	Frontend: &Dashboard{},
}

type Dashboards struct {
	Backend  *Dashboard
	Frontend *Dashboard
}

type Dashboard struct {
	Cards         *Cards
	Blocks        *Blocks
	GlobalHeads   *GlobalHeads
	GlobalFooters *GlobalFooters
	TopButtons    *TopButtons
}
