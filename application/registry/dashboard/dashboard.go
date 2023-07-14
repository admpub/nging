package dashboard

import (
	"strings"

	"github.com/webx-top/echo"
)

var Default = &Dashboards{
	Backend: &Dashboard{
		Cards:         cards,
		Blocks:        &blocks,
		GlobalHeads:   &globalHeads,
		GlobalFooters: &globalFooters,
		TopButtons:    &topButtons,
		Extend:        map[string]*Dashboard{},
	},
	Frontend: &Dashboard{
		Extend: map[string]*Dashboard{},
	},
}

func New() *Dashboard {
	return &Dashboard{
		Cards:          &Cards{},
		Blocks:         &Blocks{},
		GlobalHeads:    &GlobalHeads{},
		GlobalFooters:  &GlobalFooters{},
		TopButtons:     &Buttons{},
		GroupedButtons: map[string]*Buttons{},
		Extend:         map[string]*Dashboard{},
	}
}

type Dashboards struct {
	Backend  *Dashboard
	Frontend *Dashboard
}

const (
	TypeCards          = `Cards`
	TypeBlocks         = `Blocks`
	TypeGlobalHeads    = `GlobalHeads`
	TypeGlobalFooters  = `GlobalFooters`
	TypeTopButtons     = `TopButtons`
	TypeGroupedButtons = `GroupedButtons`
)

type Dashboard struct {
	Cards          *Cards
	Blocks         *Blocks
	GlobalHeads    *GlobalHeads
	GlobalFooters  *GlobalFooters
	TopButtons     *Buttons
	GroupedButtons map[string]*Buttons
	Extend         map[string]*Dashboard
}

type IReady interface {
	Ready(ctx echo.Context) error
}

func (d *Dashboard) SetExtend(name string, dashboard *Dashboard) {
	d.Extend[name] = dashboard
}

func (d *Dashboard) GetExtend(name string) (dashboard *Dashboard) {
	dashboard = d.Extend[name]
	return
}

func (d *Dashboard) GetOrNewExtend(name string) (dashboard *Dashboard) {
	var ok bool
	dashboard, ok = d.Extend[name]
	if !ok {
		dashboard = New()
		d.Extend[name] = dashboard
	}
	return
}

func (d *Dashboard) SetGroupedButtons(group string, btn *Buttons) {
	d.GroupedButtons[group] = btn
}

func (d *Dashboard) GetButtonsByGroup(group string) (btn *Buttons) {
	btn = d.GroupedButtons[group]
	return
}

func (d *Dashboard) GetOrNewGroupedButtons(name string) (btn *Buttons) {
	var ok bool
	btn, ok = d.GroupedButtons[name]
	if !ok {
		btn = &Buttons{}
		d.GroupedButtons[name] = btn
	}
	return
}

func (d *Dashboard) Get(c echo.Context, dtype string) interface{} {
	parts := strings.SplitN(dtype, `#`, 2)
	dtype = parts[0]
	switch dtype {
	case TypeCards:
		if d.Cards == nil {
			return nil
		}
		cards := d.Cards.Build(c)
		return cards
	case TypeBlocks:
		if d.Blocks == nil {
			return nil
		}
		d.Blocks.Ready(c)
		return d.Blocks
	case TypeGlobalHeads:
		if d.GlobalHeads == nil {
			return nil
		}
		d.GlobalHeads.Ready(c)
		return d.GlobalHeads
	case TypeGlobalFooters:
		if d.GlobalFooters == nil {
			return nil
		}
		d.GlobalFooters.Ready(c)
		return d.GlobalFooters
	case TypeTopButtons:
		if d.TopButtons == nil {
			return nil
		}
		d.TopButtons.Ready(c)
		return d.TopButtons
	case TypeGroupedButtons:
		if d.GroupedButtons == nil {
			return nil
		}
		if len(parts) == 2 {
			return d.GetButtonsByGroup(parts[1])
		}
		return d.GroupedButtons
	default:
		return nil
	}
}
