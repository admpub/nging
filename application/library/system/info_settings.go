package system

import (
	"strings"

	"github.com/webx-top/echo"
)

func NewSettings() *Settings {
	return &Settings{}
}

type AlarmThreshold struct {
	Memory float64
	CPU    float64
}

type Settings struct {
	MonitorOn      bool           // 是否开启监控
	AlarmOn        bool           // 是否开启告警
	AlarmThreshold AlarmThreshold // 告警阀值
	ReportEmail    string         // 如有多个邮箱，则一行一个
}

func (s *Settings) FromStore(h echo.H) *Settings {
	s.MonitorOn = h.Bool(`MonitorOn`)
	s.AlarmOn = h.Bool(`AlarmOn`)
	v := h.Store(`AlarmThreshold`)
	s.AlarmThreshold.CPU = v.Float64(`CPU`)
	s.AlarmThreshold.Memory = v.Float64(`Memory`)
	s.ReportEmail = strings.TrimSpace(h.String(`ReportEmail`))
	return s
}
