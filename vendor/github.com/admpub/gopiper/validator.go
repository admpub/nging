package gopiper

import (
	"reflect"
	"regexp"
	"strconv"
	"unicode/utf8"

	"github.com/admpub/regexp2"
	"github.com/webx-top/com"
)

func init() {
	// 验证器类型的过滤器统一用下划线开头（验证不通过时，返回ErrInvalidContent错误便于后续处理）
	RegisterFilter("_required", required, "非空", `_required`, ``)
	RegisterFilter("_email", email, "E-mail地址", `_email`, ``)
	RegisterFilter("_username", username, "用户名(字母/数字/汉字)", `_username`, ``)
	RegisterFilter("_singleline", singleline, "单行文本", `_singleline`, ``)
	RegisterFilter("_mutiline", mutiline, "多行文本", `_mutiline`, ``)
	RegisterFilter("_url", url, "URL", `_url`, ``)
	RegisterFilter("_chinese", chinese, "全是汉字", `_chinese`, ``)
	RegisterFilter("_haschinese", haschinese, "包含汉字", `_haschinese`, ``)
	RegisterFilter("_minsize", minsize, "最小长度", `_minsize(5)`, ``)
	RegisterFilter("_maxsize", maxsize, "最大长度", `_maxsize(5)`, ``)
	RegisterFilter("_size", size, "匹配长度", `_size(5)`, ``)
	RegisterFilter("_alpha", alpha, "字母", `_alpha`, ``)
	RegisterFilter("_alphanum", alphanum, "字母或数字", `_alphanum`, ``)
	RegisterFilter("_numeric", numeric, "纯数字", `_numeric`, ``)
	RegisterFilter("_match", match, "正则匹配", `_match([a-z]+)`, ``)
	RegisterFilter("_unmatch", unmatch, "正则不匹配", `_unmatch([a-z]+)`, ``)
	RegisterFilter("_match2", match2, "正则匹配(兼容Perl5和.NET)", `_match2([a-z]+)`, ``)
	RegisterFilter("_unmatch2", unmatch2, "正则不匹配(兼容Perl5和.NET)", `_unmatch2([a-z]+)`, ``)
}

func required(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if len(v) == 0 {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}

func email(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if !com.IsEmailRFC(v) {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}

func username(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if !com.IsUsername(v) {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}

func singleline(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if !com.IsSingleLineText(v) {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}

func mutiline(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if !com.IsMultiLineText(v) {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}

func url(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {

	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if !com.IsURL(v) {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}

func chinese(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if !com.IsChinese(v) {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}

func haschinese(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if !com.HasChinese(v) {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}

func minsize(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	minSize, _ := strconv.Atoi(params.String())
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if utf8.RuneCountInString(v) < minSize {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}

func maxsize(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	maxSize, _ := strconv.Atoi(params.String())
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if utf8.RuneCountInString(v) > maxSize {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}

func size(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	size, _ := strconv.Atoi(params.String())
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if utf8.RuneCountInString(v) != size {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}

func alpha(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		for _, v := range v {
			if !com.IsAlpha(v) {
				return v, ErrInvalidContent
			}
		}
		return v, nil
	})
}

func alphanum(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		for _, v := range v {
			if !com.IsAlphaNumeric(v) {
				return v, ErrInvalidContent
			}
		}
		return v, nil
	})
}

func numeric(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		for _, v := range v {
			if !com.IsNumeric(v) {
				return v, ErrInvalidContent
			}
		}
		return v, nil
	})
}

func match(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	re, err := regexp.Compile(params.String())
	if err != nil {
		return src.Interface(), err
	}
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if !re.MatchString(v) {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}

func unmatch(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	re, err := regexp.Compile(params.String())
	if err != nil {
		return src.Interface(), err
	}
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if re.MatchString(v) {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}

func match2(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	re, err := regexp2.Compile(params.String(), 0)
	if err != nil {
		return src.Interface(), err
	}
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if ok, _ := re.MatchString(v); !ok {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}

func unmatch2(pipe *PipeItem, src *reflect.Value, params *reflect.Value) (interface{}, error) {
	re, err := regexp2.Compile(params.String(), 0)
	if err != nil {
		return src.Interface(), err
	}
	return _filterValue(src.Interface(), func(v string) (interface{}, error) {
		if ok, _ := re.MatchString(v); ok {
			return v, ErrInvalidContent
		}
		return v, nil
	})
}
