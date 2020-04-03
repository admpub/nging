package mysql

import (
	"regexp"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/db"
)

var (
	searchMultiKwRule   = regexp.MustCompile(`[\s]+`)                        //多个关键词
	splitMultiIDRule    = regexp.MustCompile(`[^\d-]+`)                      //多个Id
	searchCompareRule   = regexp.MustCompile(`^[><!][=]?[\d]+(?:\.[\d]+)?$`) //多个Id
	searchIDRule        = regexp.MustCompile(`^[\s\d-,]+$`)                  //多个Id
	searchParagraphRule = regexp.MustCompile(`"[^"]+"`)                      //段落
)

func CleanFulltextOperator(v string) string {
	if com.StrIsAlphaNumeric(v) {
		return v
	}
	v = strings.ReplaceAll(v, `'`, ``)
	v = strings.ReplaceAll(v, `+`, ``)
	v = strings.ReplaceAll(v, `-`, ``)
	v = strings.ReplaceAll(v, `*`, ``)
	v = strings.ReplaceAll(v, `"`, ``)
	v = strings.ReplaceAll(v, `\`, ``)
	return v
}

func FindInSet(key string, value string, useFulltextIndex ...bool) db.Compound {
	key = strings.Replace(key, "`", "``", -1)
	if len(useFulltextIndex) > 0 && useFulltextIndex[0] {
		v := CleanFulltextOperator(value)
		v = strings.ReplaceAll(v, `,`, ``)
		return db.Raw("MATCH(`" + key + "`) AGAINST ('\"" + v + ",\" \"," + v + ",\" \"," + v + "\"' IN BOOLEAN MODE)")
	}
	return db.Raw("FIND_IN_SET(?,`"+key+"`)", value)
}

func Match(value string, keys ...string) db.Compound {
	for idx, key := range keys {
		key = strings.ReplaceAll(key, "`", "``")
		keys[idx] = "`" + key + "`"
	}
	v := CleanFulltextOperator(value)
	return db.Raw("MATCH(" + strings.Join(keys, ",") + ") AGAINST ('" + v + "')")
}

func CompareField(idField string, keywords string) db.Compound {
	if len(keywords) == 0 || len(idField) == 0 {
		return db.EmptyCond
	}
	var op string
	if len(keywords) > 2 && keywords[1] == '=' {
		op = keywords[0:2]
	} else {
		op = keywords[0:1]
	}
	switch op {
	case `>=`:
		return db.Cond{idField: db.Gte(keywords[2:])}
	case `==`:
		return db.Cond{idField: keywords[2:]}
	case `<=`:
		return db.Cond{idField: db.Lte(keywords[2:])}
	case `!=`:
		return db.Cond{idField: db.NotEq(keywords[2:])}
	case `>`:
		return db.Cond{idField: db.Gt(keywords[1:])}
	case `<`:
		return db.Cond{idField: db.Lt(keywords[1:])}
	case `=`:
		return db.Cond{idField: keywords[1:]}
	}
	return db.EmptyCond
}

func IsCompareField(keywords string) bool {
	return len(searchCompareRule.FindString(keywords)) > 0
}

func IsRangeField(keywords string) bool {
	return len(searchIDRule.FindString(keywords)) > 0
}

// SearchFields 搜索某个字段(多个字段任一匹配)
// @param fields 字段名
// @param keywords 关键词
// @param idFields 如要搜索id字段需要提供id字段名
// @author swh <swh@admpub.com>
func SearchFields(fields []string, keywords string, idFields ...string) *db.Compounds {
	cd := db.NewCompounds()
	if len(keywords) == 0 || len(fields) == 0 {
		return cd
	}
	var idField string
	if len(idFields) > 0 {
		idField = idFields[0]
	}
	keywords = strings.TrimSpace(keywords)
	if len(idField) > 0 {
		switch {
		case IsCompareField(keywords):
			return cd.Add(CompareField(idField, keywords))
		case IsRangeField(keywords):
			return RangeField(idField, keywords)
		}
	}
	var paragraphs []string
	keywords = searchParagraphRule.ReplaceAllStringFunc(keywords, func(paragraph string) string {
		paragraph = strings.Trim(paragraph, `"`)
		paragraphs = append(paragraphs, paragraph)
		return ""
	})
	kws := searchMultiKwRule.Split(keywords, -1)
	kws = append(kws, paragraphs...)
	cond := db.NewCompounds()
	for _, v := range kws {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		var values []string
		if strings.Contains(v, "||") {
			for _, val := range strings.Split(v, "||") {
				val = com.AddSlashes(val, '_', '%')
				values = append(values, val)
			}
		} else {
			v = com.AddSlashes(v, '_', '%')
			values = append(values, v)
		}
		_cond := db.NewCompounds()
		for _, f := range fields {
			var isEq bool
			if len(f) > 1 && f[0] == '=' {
				isEq = true
				f = f[1:]
			}
			c := db.NewCompounds()
			for _, val := range values {
				if isEq {
					c.AddKV(f, val)
				} else {
					c.AddKV(f, db.Like(`%`+val+`%`))
				}
			}
			_cond.Add(c.Or())
		}
		cond.From(_cond)
	}
	return cd.Add(cond.Or())
}

// SearchField 搜索某个字段(多个字段同时匹配)
// @param field 字段名。支持搜索多个字段，各个字段之间用半角逗号“,”隔开
// @param keywords 关键词
// @param idFields 如要搜索id字段需要提供id字段名
// @author swh <swh@admpub.com>
func SearchField(field string, keywords string, idFields ...string) *db.Compounds {
	cd := db.NewCompounds()
	if len(keywords) == 0 || len(field) == 0 {
		return cd
	}
	var idField string
	if len(idFields) > 0 {
		idField = idFields[0]
	}
	keywords = strings.TrimSpace(keywords)
	if len(idField) > 0 {
		switch {
		case IsCompareField(keywords):
			return cd.Add(CompareField(idField, keywords))
		case IsRangeField(keywords):
			return RangeField(idField, keywords)
		}
	}
	var paragraphs []string
	keywords = searchParagraphRule.ReplaceAllStringFunc(keywords, func(paragraph string) string {
		paragraph = strings.Trim(paragraph, `"`)
		paragraphs = append(paragraphs, paragraph)
		return ""
	})
	kws := searchMultiKwRule.Split(keywords, -1)
	kws = append(kws, paragraphs...)
	if strings.Contains(field, ",") {
		fs := strings.Split(field, ",")
		for _, v := range kws {
			v = strings.TrimSpace(v)
			if len(v) == 0 {
				continue
			}
			var values []string
			if strings.Contains(v, "||") {
				for _, val := range strings.Split(v, "||") {
					val = com.AddSlashes(val, '_', '%')
					values = append(values, val)
				}
			} else {
				v = com.AddSlashes(v, '_', '%')
				values = append(values, v)
			}
			_cond := db.NewCompounds()
			for _, f := range fs {
				var isEq bool
				if len(f) > 1 && f[0] == '=' {
					isEq = true
					f = f[1:]
				}
				c := db.NewCompounds()
				for _, val := range values {
					if isEq {
						c.AddKV(f, val)
					} else {
						c.AddKV(f, db.Like(`%`+val+`%`))
					}
				}
				_cond.Add(c.Or())
			}
			cd.From(_cond)
		}
		return cd
	}
	var isEq bool
	if len(field) > 1 && field[0] == '=' {
		isEq = true
		field = field[1:]
	}
	for _, v := range kws {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		if strings.Contains(v, "||") {
			vals := strings.Split(v, "||")
			cond := db.NewCompounds()
			for _, val := range vals {
				if isEq {
					cd.AddKV(field, v)
				} else {
					val = com.AddSlashes(val, '_', '%')
					cond.AddKV(field, db.Like(`%`+val+`%`))
				}
			}
			cd.Add(cond.Or())
			continue
		}
		if isEq {
			cd.AddKV(field, v)
		} else {
			v = com.AddSlashes(v, '_', '%')
			cd.AddKV(field, db.Like(`%`+v+`%`))
		}
	}
	return cd
}

// RangeField 字段范围查询
func RangeField(idField string, keywords string) *db.Compounds {
	cd := db.NewCompounds()
	if len(keywords) == 0 || len(idField) == 0 {
		return cd
	}
	keywords = strings.TrimSpace(keywords)
	kws := splitMultiIDRule.Split(keywords, -1)
	cond := db.NewCompounds()
	for _, v := range kws {
		length := len(v)
		if length < 1 {
			continue
		}
		if p := strings.Index(v, "-"); p > 0 {
			if length < 2 {
				continue
			}
			if v[length-1] == '-' {
				v = strings.Trim(v, "-")
				if len(v) == 0 {
					continue
				}
				cond.AddKV(idField, db.Gte(v))
				continue
			}

			v = strings.Trim(v, "-")
			if len(v) == 0 {
				continue
			}
			vs := strings.SplitN(v, "-", 2)
			cond.AddKV(idField, db.Between(vs[0], vs[1]))
		} else {
			cond.AddKV(idField, v)
		}
	}
	return cd.Add(cond.Or())
}

// EqField 单字段相等查询
func EqField(field string, keywords string) db.Compound {
	if len(keywords) == 0 || len(field) == 0 {
		return db.EmptyCond
	}
	keywords = strings.TrimSpace(keywords)
	return db.Cond{field: keywords}
}

// GenDateRange 生成日期范围条件
// 生成日期范围条件
// @param field 字段名。支持搜索多个字段，各个字段之间用半角逗号“,”隔开
// @param keywords 关键词
func GenDateRange(field string, keywords string, seperators ...string) *db.Compounds {
	cond := db.NewCompounds()
	if len(keywords) == 0 || len(field) == 0 {
		return cond
	}
	var skwd, skwdExt, seperator string
	if len(seperators) > 0 {
		seperator = seperators[0]
	}
	if len(seperator) == 0 {
		seperator = ` - `
	}
	dataRange := strings.Split(keywords, seperator)
	skwd = dataRange[0]
	if len(dataRange) > 1 {
		skwdExt = dataRange[1]
	}
	//日期范围
	dateBegin := com.StrToTime(skwd + ` 00:00:00`)
	cond.AddKV(field, db.Gte(dateBegin))
	if len(skwdExt) > 0 {
		dateEnd := com.StrToTime(skwd + ` 23:59:59`)
		cond.AddKV(field, db.Lte(dateEnd))
	}
	return cond
}
