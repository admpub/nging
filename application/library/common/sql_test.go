package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/db"
	"github.com/webx-top/echo/defaults"
)

func TestSelectPageCond(t *testing.T) {
	ctx := defaults.NewMockContext()
	ctx.Request().Form().Set(`searchValue`, `5,4,3,6,7`)
	cond := db.NewCompounds()
	sv := SelectPageCond(ctx, cond)
	assert.Equal(t, []string{`5`, `4`, `3`, `6`, `7`}, sv.PKValues)
	assert.Equal(t, "FIELD(`id`,'5','4','3','6','7')", sv.OrderByString())

	sv = SelectPageCond(ctx, cond, `user.id`)
	assert.Equal(t, []string{`5`, `4`, `3`, `6`, `7`}, sv.PKValues)
	assert.Equal(t, "FIELD(`user`.`id`,'5','4','3','6','7')", sv.OrderByString())

	sv = SelectPageCond(ctx, cond, "us`er.i`d")
	assert.Equal(t, []string{`5`, `4`, `3`, `6`, `7`}, sv.PKValues)
	assert.Equal(t, "FIELD(`us``er`.`i``d`,'5','4','3','6','7')", sv.OrderByString())
}
