package frp

import (
	"encoding/json"

	plugin "github.com/admpub/frp/pkg/plugin/server"
	frpLog "github.com/admpub/frp/pkg/util/log"
	"github.com/admpub/nging/v3/application/library/common"
)

func NewServerConfigExtra() *ServerConfigExtra {
	return &ServerConfigExtra{
		PluginOptions: map[string]plugin.HTTPPluginOptions{},
	}
}

type ServerConfigExtra struct {
	PluginOptions map[string]plugin.HTTPPluginOptions `json:"pluginOptions"`
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
	}
	return nil
}

func (s *ServerConfigExtra) String() string {
	jsonBytes, _ := json.Marshal(s)
	return string(jsonBytes)
}
