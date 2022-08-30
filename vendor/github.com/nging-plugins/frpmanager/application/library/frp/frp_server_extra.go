package frp

import (
	"encoding/json"

	"github.com/webx-top/echo/param"

	plugin "github.com/admpub/frp/pkg/plugin/server"
	frpLog "github.com/admpub/frp/pkg/util/log"
	"github.com/admpub/nging/v4/application/library/common"
)

func NewServerConfigExtra() *ServerConfigExtra {
	return &ServerConfigExtra{
		PluginOptions:    map[string]plugin.HTTPPluginOptions{},
		unmarshaledExtra: param.Store{},
	}
}

type ServerConfigExtra struct {
	PluginOptions    map[string]plugin.HTTPPluginOptions `json:"pluginOptions"`
	Extra            json.RawMessage                     `json:"extra,omitempty"`
	unmarshaledExtra param.Store
}

func (s *ServerConfigExtra) Parse(extra string) error {
	if len(extra) > 0 {
		jsonBytes := []byte(extra)
		err := json.Unmarshal(jsonBytes, s)
		if err != nil {
			err = common.JSONBytesParseError(err, jsonBytes)
			frpLog.Error(`failed to parse ServerConfigExtra: %v`, err)
			return err
		}
		if len(s.Extra) > 0 {
			err := json.Unmarshal(s.Extra, &s.unmarshaledExtra)
			if err != nil {
				err = common.JSONBytesParseError(err, jsonBytes)
				frpLog.Error(`failed to parse ServerConfigExtra.Extra: %v`, err)
				return err
			}
		}
	}
	return nil
}

func (s *ServerConfigExtra) String() string {
	jsonBytes, _ := json.Marshal(s)
	return string(jsonBytes)
}

func (s *ServerConfigExtra) UnmarshaledExtra() param.Store {
	return s.unmarshaledExtra
}
