package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func TestRespData(t *testing.T) {
	var respData interface{}
	respData = map[string]interface{}{}
	recv := respData.(map[string]interface{})
	recv[`add`] = `addv`
	assert.Equal(t, `addv`, respData.(map[string]interface{})[`add`])

	respData = param.Store{}
	mapr := respData.(param.Store)
	recv = mapr
	recv[`add2`] = `addv2`
	assert.Equal(t, `addv2`, respData.(param.Store)[`add2`])

	respData = echo.NewData(nil)
	data := respData.(echo.Data)
	data.SetData(echo.H{`a`: `b`})
	mapr = data.GetData().(param.Store)
	recv = mapr
	recv[`add3`] = `addv3`
	assert.Equal(t, `addv3`, respData.(echo.Data).GetData().(param.Store)[`add3`])
}
