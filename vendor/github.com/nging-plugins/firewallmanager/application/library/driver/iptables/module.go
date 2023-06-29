/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package iptables

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/webx-top/echo/param"
)

type Moduler interface {
	Args() []string
	Strings() []string
	ModuleStrings() []string
	String() string
}

var ModuleList = []string{`comment`, `string`, `time`, `connlimit`, `limit`}

type ModuleComment struct {
	Comment string // 注释
}

func (m *ModuleComment) Args() []string {
	var rs []string
	if len(m.Comment) == 0 {
		return rs
	}
	rs = append(rs, m.ModuleStrings()...)
	rs = append(rs, m.Strings()...)
	return rs
}

func (m *ModuleComment) Strings() []string {
	var rs []string
	if len(m.Comment) > 0 {
		rs = append(rs, `--comment`, m.Comment)
	}
	return rs
}

func (m *ModuleComment) ModuleStrings() []string {
	return []string{`-m`, `comment`}
}

func (m *ModuleComment) String() string {
	return strings.Join(m.ModuleStrings(), ` `) + ` ` + strings.Join(m.Strings(), ` `)
}

type ModuleString struct {
	Find string // 指定需要匹配的字符串。
	Algo string // 指定对应的匹配算法，可用算法为bm、kmp，此选项为必选项。
}

func (m *ModuleString) Args() []string {
	var rs []string
	rs = append(rs, m.ModuleStrings()...)
	rs = append(rs, m.Strings()...)
	return rs
}

func (m *ModuleString) Strings() []string {
	var rs []string
	if len(m.Find) > 0 {
		rs = append(rs, `--string`, m.Find)
	}
	if len(m.Algo) == 0 {
		m.Algo = `bm`
	}
	rs = append(rs, `--algo`, m.Algo)
	return rs
}

func (m *ModuleString) ModuleStrings() []string {
	return []string{`-m`, `string`}
}

func (m *ModuleString) String() string {
	return strings.Join(m.ModuleStrings(), ` `) + ` ` + strings.Join(m.Strings(), ` `)
}

type ModuleTime struct {
	Date      [2]string // 2006-01-02
	Time      [2]string // 15:04:05
	Weekdays  []uint    // 1-7
	Monthdays []uint    // 1-28/30/31
	KernelTZ  bool      // KernelTZ 为 false 的情况下，以上参数时间的时区为 UTC。否则为本地机器时区。
}

func joinUint(vals []uint, sep string) string {
	r := make([]string, len(vals))
	for i, v := range vals {
		r[i] = param.AsString(v)
	}
	return strings.Join(r, sep)
}

func (m *ModuleTime) Args() []string {
	var rs []string
	args := m.Strings()
	if len(args) == 0 {
		return rs
	}
	rs = append(rs, m.ModuleStrings()...)
	rs = append(rs, args...)
	return rs
}

func (m *ModuleTime) Strings() []string {
	var rs []string
	if len(m.Date[0]) > 0 {
		rs = append(rs, `--datestart`, m.Date[0])
	}
	if len(m.Date[1]) > 0 {
		rs = append(rs, `--datestop`, m.Date[1])
	}
	if len(m.Time[0]) > 0 {
		rs = append(rs, `--timestart`, m.Time[0])
	}
	if len(m.Time[1]) > 0 {
		rs = append(rs, `--timestop`, m.Time[1])
	}
	if len(m.Monthdays) > 0 {
		rs = append(rs, `--monthdays`, joinUint(m.Monthdays, `,`))
	}
	if len(m.Weekdays) > 0 {
		rs = append(rs, `--weekdays`, joinUint(m.Weekdays, `,`))
	}
	if m.KernelTZ {
		rs = append(rs, `--kerneltz`)
	}
	return rs
}

func (m *ModuleTime) ModuleStrings() []string {
	return []string{`-m`, `time`}
}

func (m *ModuleTime) String() string {
	return strings.Join(m.ModuleStrings(), ` `) + ` ` + strings.Join(m.Strings(), ` `)
}

// ModuleConnLimit 限制每个IP的最大连接数
type ModuleConnLimit struct {
	Upto  uint64 // 如果连接数低于或等于此值，则匹配
	Above uint64 // 如果连接数高于此值，则匹配
	Mask  uint16 // 此选项不能单独使用，在使用–connlimit-above选项时，配合此选项，则可以针对”某类IP段内的一定数量的IP”进行连接数量的限制。例如 24 或 27。
}

func (m *ModuleConnLimit) Args() []string {
	var rs []string
	args := m.Strings()
	if len(args) == 0 {
		return rs
	}
	rs = append(rs, m.ModuleStrings()...)
	rs = append(rs, args...)
	return rs
}

func (m *ModuleConnLimit) Strings() []string {
	var rs []string
	if m.Above > 0 {
		rs = append(rs, `--connlimit-above`, param.AsString(m.Above))
	} else if m.Upto > 0 {
		rs = append(rs, `--connlimit-upto`, param.AsString(m.Upto))
	}

	if m.Mask > 0 && (m.Above > 0 || m.Upto > 0) {
		rs = append(rs, `--connlimit-mask`, param.AsString(m.Mask))
	}
	return rs
}

func (m *ModuleConnLimit) ModuleStrings() []string {
	return []string{`-m`, `connlimit`}
}

func (m *ModuleConnLimit) String() string {
	return strings.Join(m.ModuleStrings(), ` `) + ` ` + strings.Join(m.Strings(), ` `)
}

func ParseConnLimit(limitStr string) (*ModuleConnLimit, error) {
	e := &ModuleConnLimit{}
	limitStr = strings.TrimSpace(limitStr)
	var err error
	if strings.HasSuffix(limitStr, `+`) {
		limitStr = strings.TrimSuffix(limitStr, `+`)
		e.Above, err = strconv.ParseUint(limitStr, 10, 64)
	} else {
		e.Upto, err = strconv.ParseUint(limitStr, 10, 64)
	}
	return e, err
}

// ParseLimits parse ModuleLimit
// rateStr := `1+/bytes/second`
func ParseLimits(rateStr string, burst uint) (*ModuleLimit, error) {
	e := &ModuleLimit{
		Limit: 0,
		Unit:  `second`,
		Burst: burst,
	}
	var err error
	var isLimitBytes bool
	parts := strings.SplitN(rateStr, `/`, 3)
	switch len(parts) {
	case 3:
		parts[2] = strings.TrimSpace(parts[2])
		if len(parts[2]) > 0 {
			switch parts[2][0] {
			case 's': // second
				e.Unit = `second`
			case 'm': // minute
				e.Unit = `minute`
			case 'h': // hour
				e.Unit = `hour`
			case 'd': // day
				e.Unit = `day`
			case 'w': // week
				e.Unit = `week`
			}
		}
		fallthrough
	case 2:
		parts[1] = strings.TrimSpace(parts[1])
		if len(parts[1]) > 0 {
			switch parts[1][0] {
			case 'p': // pkts
				// ok
			case 'b': // bytes
				isLimitBytes = true
			}
		}
		fallthrough
	case 1:
		parts[0] = strings.TrimSpace(parts[0])
		parts[0] = strings.TrimSuffix(parts[0], `+`)
		e.Limit, err = strconv.ParseUint(parts[0], 10, 64)
		if err != nil {
			err = fmt.Errorf(`failed to ParseUint(%q) from %q: %w`, parts[0], rateStr, err)
		} else {
			if isLimitBytes { // 限制字节时，转换为大致的包数量，假设每个包1500bytes
				e.Limit = e.Limit / 1500
			}
		}
	}
	return e, err
}

func ParseHashLimits(rateStr string, burst uint) (*ModuleHashLimit, error) {
	e := &ModuleHashLimit{
		Upto:         0,
		Above:        0,
		Unit:         `second`,
		Burst:        burst,
		Mode:         ``,
		Mask:         0,
		Name:         ``,
		Buckets:      0,
		MaxEntries:   0,
		ExpireMs:     0,
		GcIntervalMs: 0,
	}
	var err error
	var isLimitBytes bool
	parts := strings.SplitN(rateStr, `/`, 3)
	switch len(parts) {
	case 3:
		parts[2] = strings.TrimSpace(parts[2])
		if len(parts[2]) > 0 {
			switch parts[2][0] {
			case 's': // second
				e.Unit = `second`
			case 'm': // minute
				e.Unit = `minute`
			case 'h': // hour
				e.Unit = `hour`
			case 'd': // day
				e.Unit = `day`
			case 'w': // week
				e.Unit = `week`
			}
		}
		fallthrough
	case 2:
		parts[1] = strings.TrimSpace(parts[1])
		if len(parts[1]) > 0 {
			switch parts[1][0] {
			case 'p': // pkts
				// ok
			case 'b': // bytes
				isLimitBytes = true
			}
		}
		fallthrough
	case 1:
		parts[0] = strings.TrimSpace(parts[0])
		if strings.HasSuffix(parts[0], `+`) {
			parts[0] = strings.TrimSuffix(parts[0], `+`)
			e.Above, err = strconv.ParseUint(parts[0], 10, 64)
		} else {
			e.Upto, err = strconv.ParseUint(parts[0], 10, 64)
		}
		if err != nil {
			err = fmt.Errorf(`failed to ParseUint(%q) from %q: %w`, parts[0], rateStr, err)
		} else {
			if isLimitBytes { // 限制字节时，转换为大致的包数量，假设每个包1500bytes
				if e.Above > 0 {
					e.Above = e.Above / 1500
				} else {
					e.Upto = e.Upto / 1500
				}
			}
		}
	}
	return e, err
}

// ModuleLimit 限制每个IP的最大发包数
type ModuleLimit struct {
	Limit uint64 // 指定令牌桶中生成新令牌的频率
	Unit  string // 时间单位 second、minute、hour、day
	Burst uint   // 指定令牌桶中令牌的最大数量
}

func (m *ModuleLimit) Args() []string {
	var rs []string
	args := m.Strings()
	if len(args) == 0 {
		return rs
	}
	rs = append(rs, m.ModuleStrings()...)
	rs = append(rs, args...)
	return rs
}

func (m *ModuleLimit) Strings() []string {
	var rs []string
	if m.Burst > 0 {
		rs = append(rs, `--limit-burst`, param.AsString(m.Burst))
	}
	if m.Limit > 0 && len(m.Unit) > 0 {
		rs = append(rs, `--limit`, param.AsString(m.Limit)+`/`+m.Unit)
	}
	return rs
}

func (m *ModuleLimit) ModuleStrings() []string {
	return []string{`-m`, `limit`}
}

func (m *ModuleLimit) String() string {
	return strings.Join(m.ModuleStrings(), ` `) + ` ` + strings.Join(m.Strings(), ` `)
}

type HashLimitMode string

const (
	HashLimitModeSrcIP   HashLimitMode = `srcip`
	HashLimitModeSrcPort HashLimitMode = `srcport`
	HashLimitModeDstIP   HashLimitMode = `dstip`
	HashLimitModeDstPort HashLimitMode = `dstport`
)

// ModuleHashLimit 限制每个IP的最大发包数
type ModuleHashLimit struct {
	Upto         uint64        // 如果速率低于或等于此值，则匹配
	Above        uint64        // 如果速率高于此值，则匹配。
	Unit         string        // 时间单位 second、minute、hour、day
	Burst        uint          // 指定令牌桶中令牌的最大数量
	Mode         HashLimitMode // 一个用逗号分隔的对象列表。如果没有给出–hashlimit-mode选项，’hashlimit’ 的行为就像 ‘limit’ 一样，但是在做哈希管理的代价很高。
	Mask         uint16        // 当mode设置为 srcip 或 dstip 时, 配置相应的掩码表示一个网段。例如8、16、24、32
	Name         string        // 定义这条hashlimit规则的名称, 所有的条目(entry)都存放在 /proc/net/ipt_hashlimit/{foo} 里。
	Buckets      uint          // 散列表的桶数（buckets）
	MaxEntries   uint          // 散列中的最大条目
	ExpireMs     uint          // hash规则失效时间, 单位毫秒(milliseconds)
	GcIntervalMs uint          // 垃圾回收器回收的间隔时间, 单位毫秒
}

func (m *ModuleHashLimit) Args() []string {
	var rs []string
	args := m.Strings()
	if len(args) == 0 {
		return rs
	}
	rs = append(rs, m.ModuleStrings()...)
	rs = append(rs, args...)
	return rs
}

func (m *ModuleHashLimit) Strings() []string {
	var rs []string
	if m.Burst > 0 {
		rs = append(rs, `--hashlimit-burst`, param.AsString(m.Burst))
	}
	unit := m.Unit
	if len(unit) == 0 {
		unit = `second`
	}
	if m.Upto > 0 {
		rs = append(rs, `--hashlimit-upto`, param.AsString(m.Upto)+`/`+unit)
	} else if m.Above > 0 {
		rs = append(rs, `--hashlimit-above`, param.AsString(m.Above)+`/`+unit)
	}
	if len(m.Mode) > 0 {
		rs = append(rs, `--hashlimit-mode`, string(m.Mode))
		if m.Mask > 0 {
			switch m.Mode {
			case HashLimitModeSrcIP:
				rs = append(rs, `--hashlimit-srcmask`, param.AsString(m.Mask))
			case HashLimitModeDstIP:
				rs = append(rs, `--hashlimit-dstmask`, param.AsString(m.Mask))
			}
		}
	}
	rs = append(rs, `--hashlimit-name`, m.Name)
	if m.Buckets > 0 {
		rs = append(rs, `--hashlimit-htable-size`, param.AsString(m.Buckets))
	}
	if m.MaxEntries > 0 {
		rs = append(rs, `--hashlimit-htable-max`, param.AsString(m.MaxEntries))
	}
	if m.ExpireMs > 0 {
		rs = append(rs, `--hashlimit-htable-expire`, param.AsString(m.ExpireMs))
	}
	if m.GcIntervalMs > 0 {
		rs = append(rs, `--hashlimit-htable-gcinterval`, param.AsString(m.GcIntervalMs))
	}
	return rs
}

func (m *ModuleHashLimit) ModuleStrings() []string {
	return []string{`-m`, `hashlimit`}
}

func (m *ModuleHashLimit) String() string {
	return strings.Join(m.ModuleStrings(), ` `) + ` ` + strings.Join(m.Strings(), ` `)
}
