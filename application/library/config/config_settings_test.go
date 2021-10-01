package config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
)

func _newTestConfig() echo.H {
	return echo.H{
		`test`: echo.H{
			`item1_1`: 1,
			`item1_2`: 1,
		},
		`test2`: echo.H{
			`item2_1`: 1,
			`item2_2`: 1,
		},
	}
}

func TestSettings(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	OnGroupSetSettings(`test`, func(diffs Diffs) error {
		echo.Dump(diffs)
		buf.WriteString(`1`)
		return nil
	})
	OnGroupSetSettings(`test2`, func(diffs Diffs) error {
		echo.Dump(diffs)
		buf.WriteString(`9`)
		return nil
	})
	OnKeySetSettings(`test.item1_1`, func(diff Diff) error {
		buf.WriteString(`2`)
		return nil
	})
	OnKeySetSettings(`test.item1_2`, func(diff Diff) error {
		buf.WriteString(`3`)
		return nil
	})
	OnKeySetSettings(`test2.item2_1`, func(diff Diff) error {
		buf.WriteString(`4`)
		return nil
	})
	OnKeySetSettings(`test2.item2_2`, func(diff Diff) error {
		buf.WriteString(`5`)
		return nil
	})
	oldConfigs := _newTestConfig()
	configs := _newTestConfig()
	err := FireInitSettings(configs)
	if err != nil {
		panic(err)
	}
	//assert.Equal(t, `192345`, buf.String())
	assert.Equal(t, 6, len(buf.String()))
	buf.Reset()

	st := NewSettings(NewConfig())
	configs = _newTestConfig()
	configs[`test`] = echo.H{
		`item1_1`: 2,
	}
	st.setConfigs(configs, oldConfigs)
	assert.Equal(t, `12`, buf.String())
	buf.Reset()

	delete(onGroupSetSettings, `test`)

	//configs = _newTestConfig()
	st.setConfigs(configs, oldConfigs)
	assert.Equal(t, ``, buf.String())
	buf.Reset()

	OnGroupSetSettings(`test`, func(diffs Diffs) error {
		echo.Dump(diffs)
		buf.WriteString(`1`)
		return nil
	})
	oldConfigs = _newTestConfig()
	configs = _newTestConfig()
	configs.SetMKey(`test.item1_2`, 2)
	assert.Equal(t, 2, configs.GetStore(`test`).Int(`item1_2`))
	println(`oldConfigs -------------`)
	echo.Dump(oldConfigs)
	println(`newConfigs -------------`)
	echo.Dump(configs)
	st.setConfigs(configs, oldConfigs)
	assert.Equal(t, `13`, buf.String())
	buf.Reset()
}
