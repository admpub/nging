/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package handler

import (
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/errors"
	"github.com/webx-top/echo"
)

func Paging(ctx echo.Context) (page int, size int) {
	return common.Paging(ctx)
}

func Ok(v string) errors.Successor {
	return common.Ok(v)
}

func Err(ctx echo.Context, err error) (ret interface{}) {
	return common.Err(ctx, err)
}

func ok(ctx echo.Context, msg string) {
	ctx.Session().AddFlash(Ok(msg))
}

func fail(ctx echo.Context, msg string) {
	ctx.Session().AddFlash(msg)
}
