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
package common

import (
	stdErr "errors"

	"github.com/admpub/nging/application/library/errors"
	"github.com/webx-top/echo"
	"github.com/webx-top/pagination"
)

var PageMaxSize = 1000

func Paging(ctx echo.Context) (page int, size int) {
	page = ctx.Formx(`page`).Int()
	size = ctx.Formx(`size`).Int()
	if page < 1 {
		page = 1
	}
	if size < 1 || size > PageMaxSize {
		size = 50
	}
	return
}

func PagingWithPagination(ctx echo.Context, delKeys ...string) (page int, size int, totalRows int, p *pagination.Pagination) {
	page, size = Paging(ctx)
	totalRows = ctx.Formx(`rows`).Int()
	p = pagination.New(ctx).SetAll(``, totalRows, page, 10, size).SetURL(map[string]string{
		`rows`: `rows`,
		`page`: `page`,
		`size`: `size`,
	}, delKeys...)
	return
}

func Ok(v string) errors.Successor {
	return errors.NewOk(v)
}

func Err(ctx echo.Context, err error) (ret interface{}) {
	if err == nil {
		flash := ctx.Flash()
		if flash != nil {
			if errMsg, ok := flash.(string); ok {
				ret = stdErr.New(errMsg)
			} else {
				ret = flash
			}
		}
	} else {
		ret = err
	}
	return
}
