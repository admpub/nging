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

func FindInSet(key string, value string, useFulltextIndex ...bool) db.Compound {
	if len(useFulltextIndex) > 0 && useFulltextIndex[0] {
		v := CleanFulltextOperator(value)
		v = strings.ReplaceAll(v, `,`, ``)
		return match(`"`+v+`," ",`+v+`," ",`+v+`"`, true, key)
	}
	key = strings.Replace(key, "`", "``", -1)
	return db.Raw("FIND_IN_SET(?,`"+key+"`)", value)
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

func MatchAnyFields(fields []string, keywords string, idFields ...string) *db.Compounds {
	return SearchFields(fields, keywords, idFields...)
}

func MatchAnyField(field string, keywords string, idFields ...string) *db.Compounds {
	if len(field) == 0 {
		return db.NewCompounds()
	}
	field = strings.Trim(field, `,`)
	fields := strings.Split(field, `,`)
	return SearchFields(fields, keywords, idFields...)
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
	return cd.Add(searchAllFields(fields, keywords, idFields...).Or())
}

func MatchAllFields(fields []string, keywords string, idFields ...string) *db.Compounds {
	return searchAllFields(fields, keywords, idFields...)
}

func MatchAllField(field string, keywords string, idFields ...string) *db.Compounds {
	return SearchField(field, keywords, idFields...)
}

// SearchField 搜索某个字段(多个字段同时匹配)
// @param field 字段名。支持搜索多个字段，各个字段之间用半角逗号“,”隔开
// @param keywords 关键词
// @param idFields 如要搜索id字段需要提供id字段名
// @author swh <swh@admpub.com>
func SearchField(field string, keywords string, idFields ...string) *db.Compounds {
	if strings.Contains(field, ",") {
		fields := strings.Split(field, ",")
		return searchAllFields(fields, keywords, idFields...)
	}
	return searchAllField(field, keywords, idFields...)
}

func searchAllField(field string, keywords string, idFields ...string) *db.Compounds {
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
	fieldMode := parseFieldOp([]string{field})[0]
	var safelyMatchValues []string
	for _, v := range kws {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		if strings.Contains(v, "||") {
			vals := strings.Split(v, "||")
			if fieldMode.operator == OperatorMatch {
				for _, val := range vals {
					safelyMatchValues = append(safelyMatchValues, buildSafelyMatchValues(val, false)...)
				}
				continue
			}
			if fieldMode.isLikeQuery() {
				for key, val := range vals {
					val = com.AddSlashes(val, '_', '%')
					vals[key] = val
				}
			}
			cond := fieldMode.buildCondOther(vals, vals)
			cd.Add(cond.Or())
			continue
		}
		if fieldMode.operator == OperatorMatch {
			safelyMatchValues = append(safelyMatchValues, buildSafelyMatchValues(v, true)...)
			continue
		}
		if fieldMode.isLikeQuery() {
			v = com.AddSlashes(v, '_', '%')
		}
		vals := []string{v}
		fieldMode.buildCondOther(vals, vals, cd)
	}
	if len(safelyMatchValues) > 0 {
		keys := strings.Split(fieldMode.field, `,`)
		cd.Add(match(strings.Join(safelyMatchValues, ` `), true, keys...))
	}
	return cd
}

func searchAllFields(fields []string, keywords string, idFields ...string) *db.Compounds {
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
	fieldModes := parseFieldOp(fields)
	var hasLike, hasMatch bool
	for _, fieldMode := range fieldModes {
		if hasLike && hasMatch {
			break
		}
		if fieldMode.isMatchQuery() {
			hasMatch = true
			continue
		}
		if fieldMode.isLikeQuery() {
			hasLike = true
			continue
		}
	}
	safelyMatchValues := map[string][]string{}
	for _, v := range kws {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		var originalValues []string
		var likeValues []string
		var matchValues []string
		orQuery := strings.Contains(v, "||")
		if orQuery {
			originalValues = strings.Split(v, "||")
			if hasLike {
				for _, val := range originalValues {
					val = com.AddSlashes(val, '_', '%')
					likeValues = append(likeValues, val)
				}
			}
			if hasMatch {
				for _, val := range originalValues {
					matchValues = append(matchValues, buildSafelyMatchValues(val, false)...)
				}
			}
		} else {
			originalValues = append(originalValues, v)
			if hasLike {
				v = com.AddSlashes(v, '_', '%')
				likeValues = append(likeValues, v)
			}
			if hasMatch {
				matchValues = append(matchValues, buildSafelyMatchValues(v, true)...)
			}
		}
		_cond := db.NewCompounds()
		for _, f := range fieldModes {
			if f.buildCondMatch(matchValues, &safelyMatchValues) {
				continue
			}
			c := f.buildCondOther(likeValues, originalValues)
			if orQuery {
				_cond.Add(c.Or())
			} else {
				_cond.Add(c.And())
			}
		}
		cd.From(_cond)
	}
	if len(safelyMatchValues) > 0 {
		for _, f := range fieldModes {
			values, ok := safelyMatchValues[f.field]
			if ok {
				keys := strings.Split(f.field, `,`)
				cd.Add(match(strings.Join(values, ` `), true, keys...))
			}
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
	var dateStart, dateEnd, seperator string
	if len(seperators) > 0 {
		seperator = seperators[0]
	}
	if len(seperator) == 0 {
		seperator = ` - `
	}
	dataRange := strings.Split(keywords, seperator)
	dateStart = dataRange[0]
	if len(dataRange) > 1 {
		dateEnd = dataRange[1]
	}
	startDateAndTime := com.FixDateTimeString(dateStart)
	switch len(startDateAndTime) {
	case 2:
		dateStart = strings.Join(startDateAndTime, ` `)
	case 1:
		dateStart = startDateAndTime[0] + ` 00:00:00`
	default:
		return cond
	}
	//日期范围
	dateStartTs := com.StrToTime(dateStart)
	if dateStartTs <= 0 {
		return cond
	}
	cond.AddKV(field, db.Gte(dateStartTs))
	if len(dateEnd) > 0 {
		endDateAndTime := com.FixDateTimeString(dateEnd)
		switch len(endDateAndTime) {
		case 2:
			dateEnd = strings.Join(endDateAndTime, ` `)
		case 1:
			dateEnd = endDateAndTime[0] + ` 23:59:59`
		default:
			return cond
		}
		dateEndTs := com.StrToTime(dateEnd)
		if dateEndTs <= 0 {
			return cond
		}
		cond.AddKV(field, db.Lte(dateEndTs))
	}
	return cond
}
