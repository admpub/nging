package perm

import (
	"encoding/json"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

type BehaviorPerms map[string]Behavior

func (b BehaviorPerms) Add(be *Behavior, value ...interface{}) {
	if be == nil {
		return
	}
	bCopy := *be
	if len(value) > 0 {
		bCopy.Value = value[0]
	}
	if bCopy.Value == nil && be.valueInitor != nil {
		bCopy.Value = be.valueInitor()
	}
	b[be.Name] = bCopy
}

func (b BehaviorPerms) Get(name string) Behavior {
	r, _ := b[name]
	return r
}

func (b BehaviorPerms) CheckBehavior(perm string) (interface{}, bool) {
	v, y := b[perm]
	if !y {
		return nil, false
	}
	return v.Value, true
}

func ParseBehavior(permBehaviors string, behaviors *Behaviors) BehaviorPerms {
	data := echo.H{}
	json.Unmarshal([]byte(permBehaviors), &data)
	b := BehaviorPerms{}
	for name, value := range data {
		item := behaviors.GetItem(name)
		if item == nil {
			continue
		}
		behavior, ok := item.X.(*Behavior)
		if !ok {
			continue
		}
		b.Add(behavior, value)
	}
	return b
}

func SerializeBehaviorValues(permBehaviors map[string][]string, behaviors *Behaviors) (string, error) {
	data := echo.H{}
	for name, values := range permBehaviors {
		item := behaviors.GetItem(name)
		if item == nil {
			continue
		}
		behavior, ok := item.X.(*Behavior)
		if !ok {
			continue
		}
		if behavior.formValueEncoder != nil {
			if val, err := behavior.formValueEncoder(values); err == nil {
				data[name] = val
			}
			continue
		}
		switch behavior.ValueType {
		case `list`:
			data[name] = strings.Join(values, ",")
		case `number`, `float64`, `float`:
			if len(values) > 0 {
				data[name] = param.AsFloat64(values[0])
			}
		case `float32`:
			if len(values) > 0 {
				data[name] = param.AsFloat32(values[0])
			}
		case `int32`:
			if len(values) > 0 {
				data[name] = param.AsInt32(values[0])
			}
		case `int`:
			if len(values) > 0 {
				data[name] = param.AsInt(values[0])
			}
		case `int64`:
			if len(values) > 0 {
				data[name] = param.AsInt64(values[0])
			}
		case `uint32`:
			if len(values) > 0 {
				data[name] = param.AsUint32(values[0])
			}
		case `uint`:
			if len(values) > 0 {
				data[name] = param.AsUint(values[0])
			}
		case `uint64`:
			if len(values) > 0 {
				data[name] = param.AsUint64(values[0])
			}
		case `json`:
			if len(values) > 0 {
				var recv interface{}
				if behavior.valueInitor != nil {
					recv = behavior.valueInitor()
				} else {
					recv = &echo.H{}
				}
				json.Unmarshal([]byte(values[0]), recv)
				data[name] = recv
			}
		case `slice`:
			data[name] = values
		default:
			if len(values) > 0 {
				data[name] = values[0]
			}
		}
	}
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
