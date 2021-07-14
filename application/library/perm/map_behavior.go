package perm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/abh/errorutil"
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

type CheckedBehavior struct {
	Value   interface{}
	Checked bool
}

func (b BehaviorPerms) CheckBehavior(perm string) *CheckedBehavior {
	v, y := b[perm]
	if !y {
		return &CheckedBehavior{}
	}
	return &CheckedBehavior{Value: v.Value, Checked: true}
}

func JSONBytesParseError(err error, jsonBytes []byte) error {
	var offset int64 = -1
	switch rErr := err.(type) {
	case *json.UnmarshalTypeError:
		offset = rErr.Offset
	case *json.SyntaxError:
		offset = rErr.Offset
	}
	if offset > -1 {
		byteReader := bytes.NewReader(jsonBytes)
		line, col, highlight := errorutil.HighlightBytePosition(byteReader, offset)
		extra := fmt.Sprintf(":\nError at line %d, column %d (offset %d):\n%s",
			line, col, offset, highlight)
		return fmt.Errorf("error parsing json object%s\n%w",
			extra, err)
	}
	return fmt.Errorf("%w: %s", err, string(jsonBytes))
}

func ParseBehavior(permBehaviors string, behaviors *Behaviors) (BehaviorPerms, error) {
	b := BehaviorPerms{}
	data := map[string]json.RawMessage{}
	dataBytes := []byte(permBehaviors)
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		err = JSONBytesParseError(err, dataBytes)
		return b, err
	}
	for name, jsonBytes := range data {
		item := behaviors.GetItem(name)
		if item == nil {
			continue
		}
		behavior, ok := item.X.(*Behavior)
		if !ok {
			continue
		}
		var value interface{}
		if behavior.Value != nil {
			rv := reflect.New(reflect.Indirect(reflect.ValueOf(behavior.Value)).Type())
			if rv.CanInterface() {
				value = rv.Interface()
			}
		}
		if value == nil && behavior.valueInitor != nil {
			value = behavior.valueInitor()
		}
		if err := json.Unmarshal(jsonBytes, &value); err != nil {
			err = JSONBytesParseError(err, jsonBytes)
			return b, err
		}
		b.Add(behavior, value)
	}
	return b, nil
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
				dataBytes := []byte(values[0])
				if err := json.Unmarshal(dataBytes, recv); err != nil {
					return ``, JSONBytesParseError(err, dataBytes)
				}
				data[name] = recv
			}
		case `slice`:
			data[name] = values
		default:
			if len(values) > 0 {
				if behavior.valueInitor != nil {
					recv := behavior.valueInitor()
					v := reflect.Indirect(reflect.ValueOf(recv))
					switch v.Kind() {
					case reflect.Slice, reflect.Map, reflect.Struct, reflect.Array:
						dataBytes := []byte(values[0])
						if err := json.Unmarshal(dataBytes, recv); err != nil {
							return ``, JSONBytesParseError(err, dataBytes)
						}
						data[name] = recv
					default:
						data[name] = values[0]
					}
				} else {
					data[name] = values[0]
				}
			}
		}
	}
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
