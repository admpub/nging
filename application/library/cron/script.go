package cron

import (
	"os"
	"strings"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

func SaveScriptFile(m *dbschema.NgingTask) error {
	if !com.IsWindows {
		return nil
	}
	name := param.AsString(m.Id) + `.bat`
	if !strings.Contains(m.Command, "\n") {
		_ = common.RemoveCache(`taskscripts`, name)
		return nil
	}
	err := common.WriteCache(`taskscripts`, name, com.Str2bytes(m.Command))
	return err
}

func DeleteScriptFile(id uint) error {
	if !com.IsWindows {
		return nil
	}
	name := param.AsString(id) + `.bat`
	err := common.RemoveCache(`taskscripts`, name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}

func ScriptFile(id uint) string {
	if !com.IsWindows {
		return ``
	}
	name := param.AsString(id) + `.bat`
	return common.CacheFile(`taskscripts`, name)
}

func ScriptCommand(id uint, command string) string {
	if !com.IsWindows {
		return command
	}
	if !strings.Contains(command, "\n") {
		return command
	}
	scriptFile := ScriptFile(id)
	if len(scriptFile) == 0 {
		return command
	}
	return scriptFile
}
