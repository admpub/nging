package factory

import (
	"context"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/sqlbuilder"
)

type ParamOption func(*Param)

func OContext(ctx context.Context) ParamOption {
	return func(p *Param) {
		p.ctx = ctx
	}
}

func OFactory(factory *Factory) ParamOption {
	return func(p *Param) {
		p.factory = factory
	}
}

func OIndex(index int) ParamOption {
	return func(p *Param) {
		p.index = index
	}
}

func OReadOnly(readOnly bool) ParamOption {
	return func(p *Param) {
		p.readOnly = readOnly
	}
}

func OCollection(collection string, alias ...string) ParamOption {
	return func(p *Param) {
		p.collection = collection
		if len(alias) > 0 {
			p.alias = alias[0]
		}
	}
}

func OAlias(alias string) ParamOption {
	return func(p *Param) {
		p.alias = alias
	}
}

func OMiddleware(middleware func(db.Result) db.Result, name ...string) ParamOption {
	return func(p *Param) {
		p.middleware = middleware
		if len(name) > 0 {
			p.middlewareName = name[0]
		}
	}
}

func OMiddlewareName(middlewareName string) ParamOption {
	return func(p *Param) {
		p.middlewareName = middlewareName
	}
}

func OMiddlewareSelector(middlewareSelector func(sqlbuilder.Selector) sqlbuilder.Selector, name ...string) ParamOption {
	return func(p *Param) {
		p.middlewareSelector = middlewareSelector
		if len(name) > 0 {
			p.middlewareName = name[0]
		}
	}
}

func OMiddlewareTx(middlewareTx func(*Transaction) error) ParamOption {
	return func(p *Param) {
		p.middlewareTx = middlewareTx
	}
}

func OResult(result interface{}) ParamOption {
	return func(p *Param) {
		p.result = result
	}
}

func OArgs(args ...interface{}) ParamOption {
	return func(p *Param) {
		p.args = args
	}
}

func OCols(cols ...interface{}) ParamOption {
	return func(p *Param) {
		p.cols = cols
	}
}

func OJoins(joins ...*Join) ParamOption {
	return func(p *Param) {
		p.joins = joins
	}
}

func OAddJoin(joinType string, collection string, alias string, condition string, args ...interface{}) ParamOption {
	return func(p *Param) {
		p.AddJoin(joinType, collection, alias, condition, args...)
	}
}

func OSend(save interface{}) ParamOption {
	return func(p *Param) {
		p.save = save
	}
}

func OOffset(offset int) ParamOption {
	return func(p *Param) {
		p.offset = offset
	}
}

func OPage(page int) ParamOption {
	return func(p *Param) {
		p.page = page
	}
}

func OSize(size int) ParamOption {
	return func(p *Param) {
		p.size = size
	}
}

func OMaxAge(maxAge int64) ParamOption {
	return func(p *Param) {
		p.maxAge = maxAge
	}
}

func OTrans(trans Transactioner) ParamOption {
	return func(p *Param) {
		p.SetTrans(trans)
	}
}

func OCachedKey(cachedKey string) ParamOption {
	return func(p *Param) {
		p.SetCachedKey(cachedKey)
	}
}

func OModel(model Model) ParamOption {
	return func(p *Param) {
		p.SetModel(model)
	}
}

func OTotal(total int64) ParamOption {
	return func(p *Param) {
		p.SetTotal(total)
	}
}
