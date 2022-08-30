package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/confl"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/nging-plugins/frpmanager/application/dbschema"
	"github.com/nging-plugins/frpmanager/application/library/frp"
)

const FRPConfigExtension = `.json` //`.yaml`

func ConfigFile(id uint, isServer bool) string {
	configFile := `server`
	if !isServer {
		configFile = `client`
	}
	configFile = filepath.Join(echo.Wd(), `config`, `frp`, configFile)
	err := os.MkdirAll(configFile, os.ModePerm)
	if err != nil {
		log.Error(err)
	}
	return filepath.Join(configFile, fmt.Sprintf(`%d`, id)+FRPConfigExtension)
}

func PidFile(id string, isServer bool) string {
	pidFile := `server`
	if !isServer {
		pidFile = `client`
	}
	pidFile = filepath.Join(echo.Wd(), `data/pid/frp`, pidFile)
	err := os.MkdirAll(pidFile, os.ModePerm)
	if err != nil {
		log.Error(err)
	}
	return filepath.Join(pidFile, id+`.pid`)
}

func SaveConfigFile(data interface{}) (err error) {
	var configFile string
	switch v := data.(type) {
	case *dbschema.NgingFrpServer:
		configFile = ConfigFile(v.Id, true)
		if v.Disabled == `Y` {
			if !com.FileExists(configFile) {
				return nil
			}
			return os.Remove(configFile)
		}
		if len(v.Plugins) > 0 {
			copied := *v
			serverConfigExtra := frp.NewServerConfigExtra()
			serverConfigExtra.PluginOptions = frp.ServerPluginOptions(strings.Split(copied.Plugins, `,`)...)
			if len(copied.Extra) > 0 {
				serverConfigExtra.Extra = []byte(copied.Extra)
			}
			copied.Extra = serverConfigExtra.String()
			data = copied
		}
	case *dbschema.NgingFrpClient:
		configFile = ConfigFile(v.Id, false)
		if v.Disabled == `Y` {
			if !com.FileExists(configFile) {
				return nil
			}
			return os.Remove(configFile)
		}
		cfg := frp.NewClientConfig()
		cfg.NgingFrpClient = v
		cfg.Extra, err = frp.Table2Config(v)
		if err != nil {
			return
		}
		cfg.NgingFrpClient.Extra = ``
		data = cfg
	default:
		return fmt.Errorf(`unsupport save config: %T`, v)
	}
	var b []byte
	if strings.HasSuffix(configFile, `.json`) {
		b, err = json.MarshalIndent(data, ``, `  `)
	} else {
		b, err = confl.Marshal(data)
	}
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, b, os.ModePerm)
}
