package sqlbuilder

import (
	"reflect"
	"strings"

	"github.com/webx-top/echo/param"
)

// parse neq(field,value)
func parsePipe(pipeName string) Pipe {
	pos := strings.Index(pipeName, "(")
	if pos > 0 {
		param := pipeName[pos+1:]
		param = strings.TrimSuffix(param, ")")
		funcName := pipeName[0:pos]
		if gen, ok := PipeGeneratorList[funcName]; ok {
			return gen(param)
		}
		return nil
	}
	pipe, ok := PipeList[pipeName]
	if !ok {
		return nil
	}
	return pipe
}

type Pipe func(row reflect.Value, val interface{}) interface{}
type Pipes map[string]Pipe
type PipeGenerators map[string]func(params string) Pipe

func (pipes *Pipes) Add(name string, pipe Pipe) {
	(*pipes)[name] = pipe
}

func (gens *PipeGenerators) Add(name string, generator func(params string) Pipe) {
	(*gens)[name] = generator
}

var (
	PipeGeneratorList = PipeGenerators{
		`neq`: func(params string) Pipe { // name:value
			args := strings.SplitN(params, `:`, 2)
			var (
				fieldName     string
				expectedValue string
			)
			switch len(args) {
			case 2:
				fieldName = strings.TrimSpace(args[0])
				expectedValue = strings.TrimSpace(args[1])
			default:
				return nil
			}
			return func(row reflect.Value, v interface{}) interface{} {
				fieldValue := mapper.FieldByName(row, fieldName).Interface()
				if expectedValue != param.AsString(fieldValue) {
					return v
				}
				return nil
			}
		},
		`eq`: func(params string) Pipe { // name,value
			args := strings.SplitN(params, `:`, 2)
			var (
				fieldName     string
				expectedValue string
			)
			switch len(args) {
			case 2:
				fieldName = strings.TrimSpace(args[0])
				expectedValue = strings.TrimSpace(args[1])
			default:
				return nil
			}
			return func(row reflect.Value, v interface{}) interface{} {
				fieldValue := mapper.FieldByName(row, fieldName).Interface()
				if expectedValue == param.AsString(fieldValue) {
					return v
				}
				return nil
			}
		},
	}
	PipeList = Pipes{
		`split`: func(_ reflect.Value, v interface{}) interface{} {
			items := strings.Split(v.(string), `,`)
			result := []interface{}{}
			for _, item := range items {
				item = strings.TrimSpace(item)
				if len(item) == 0 {
					continue
				}
				result = append(result, item)
			}
			return result
		},
		`gtZero`: func(_ reflect.Value, v interface{}) interface{} {
			i := param.AsUint64(v)
			if i > 0 {
				return i
			}
			return nil
		},
		`notEmpty`: func(_ reflect.Value, v interface{}) interface{} {
			s := v.(string)
			if len(s) > 0 {
				return s
			}
			return nil
		},
	}
)
