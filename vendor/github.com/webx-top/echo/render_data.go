package echo

import (
	"html/template"
	"strings"
	"time"

	"github.com/admpub/humanize"
	"github.com/admpub/timeago"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"
)

type FormatRender func(ctx Context, data interface{}) error
type DataWrapper func(Context, interface{}) interface{}

func DefaultRenderDataWrapper(ctx Context, data interface{}) interface{} {
	return NewRenderData(ctx, data)
}

func NewRenderData(ctx Context, data interface{}) *RenderData {
	return &RenderData{
		ctx:    ctx,
		now:    com.NewTime(time.Now()),
		Func:   ctx.Funcs(),
		Data:   data,
		Stored: param.MapReadonly(ctx.Stored()),
	}
}

type RenderData struct {
	ctx    Context
	now    *com.Time
	Func   template.FuncMap
	Data   interface{}
	Stored param.MapReadonly
}

func (r *RenderData) Now() *com.Time {
	return r.now
}

func (r *RenderData) UnixTime() int64 {
	return r.now.Unix()
}

func (r *RenderData) T(format string, args ...interface{}) string {
	return r.ctx.T(format, args...)
}

func (r *RenderData) Lang() LangCode {
	return r.ctx.Lang()
}

func (r *RenderData) Get(key string, defaults ...interface{}) interface{} {
	return r.ctx.Get(key, defaults...)
}

func (r *RenderData) Set(key string, value interface{}) string {
	r.ctx.Set(key, value)
	return ``
}

func (r *RenderData) Cookie() Cookier {
	return r.ctx.Cookie()
}

func (r *RenderData) Session() Sessioner {
	return r.ctx.Session()
}

func (r *RenderData) Form(key string, defaults ...string) string {
	return r.ctx.Form(key, defaults...)
}

func (r *RenderData) Formx(key string, defaults ...string) param.String {
	return r.ctx.Formx(key, defaults...)
}

func (r *RenderData) Query(key string, defaults ...string) string {
	return r.ctx.Query(key, defaults...)
}

func (r *RenderData) Queryx(key string, defaults ...string) param.String {
	return r.ctx.Queryx(key, defaults...)
}

func (r *RenderData) FormValues(key string) []string {
	return r.ctx.FormValues(key)
}

func (r *RenderData) FormxValues(key string) param.StringSlice {
	return r.ctx.FormxValues(key)
}

func (r *RenderData) QueryValues(key string) []string {
	return r.ctx.QueryValues(key)
}

func (r *RenderData) QueryxValues(key string) param.StringSlice {
	return r.ctx.QueryxValues(key)
}

func (r *RenderData) Param(key string, defaults ...string) string {
	return r.ctx.Param(key, defaults...)
}

func (r *RenderData) Paramx(key string, defaults ...string) param.String {
	return r.ctx.Paramx(key, defaults...)
}

func (r *RenderData) URL() engine.URL {
	return r.ctx.Request().URL()
}

func (r *RenderData) URI() string {
	return r.ctx.Request().URI()
}

func (r *RenderData) Site() string {
	return r.ctx.Site()
}

func (r *RenderData) SiteURI() string {
	return r.Site() + strings.TrimPrefix(r.URI(), `/`)
}

func (r *RenderData) Referer() string {
	return r.ctx.Referer()
}

func (r *RenderData) Header() engine.Header {
	return r.ctx.Request().Header()
}

func (r *RenderData) Flash(keys ...string) interface{} {
	return r.ctx.Flash(keys...)
}

func (r *RenderData) HasAnyRequest() bool {
	return r.ctx.HasAnyRequest()
}

func (r *RenderData) DurationFormat(t interface{}, args ...string) *com.Durafmt {
	return tplfunc.DurationFormat(r.Lang().String(), t, args...)
}

func (r *RenderData) TsHumanize(startTime interface{}, endTime ...interface{}) string {
	humanizer, err := humanize.New(r.Lang().String())
	if err != nil {
		return err.Error()
	}
	var (
		startDate = tplfunc.ToTime(startTime)
		endDate   time.Time
	)
	if len(endTime) > 0 {
		endDate = tplfunc.ToTime(endTime[0])
	}
	if endDate.IsZero() {
		endDate = time.Now().Local()
	}
	return humanizer.TimeDiff(endDate, startDate, 0)
}

func (r *RenderData) CaptchaForm(args ...interface{}) template.HTML {
	return tplfunc.CaptchaFormWithURLPrefix(r.ctx.Echo().Prefix(), args...)
}

func (r *RenderData) MakeURL(h interface{}, args ...interface{}) string {
	return r.ctx.Echo().URL(h, args...)
}

func (r *RenderData) Ext() string {
	return r.ctx.DefaultExtension()
}

func (r *RenderData) Fetch(tmpl string, data interface{}) template.HTML {
	b, e := r.ctx.Fetch(tmpl, data)
	if e != nil {
		return template.HTML(e.Error())
	}
	return template.HTML(string(b))
}

func (r *RenderData) TimeAgo(v interface{}, options ...string) string {
	if datetime, ok := v.(string); ok {
		return timeago.Take(datetime, r.Lang().Format(false, `-`))
	}
	var option string
	if len(options) > 0 {
		option = options[0]
	}
	return timeago.Timestamp(param.AsInt64(v), r.Lang().Format(false, `-`), option)
}

func (r *RenderData) Prefix() string {
	return r.ctx.Route().Prefix
}

func (r *RenderData) Path() string {
	return r.ctx.Path()
}

func (r *RenderData) Queries() map[string][]string {
	return r.ctx.Queries()
}

func (r *RenderData) Domain() string {
	return r.ctx.Domain()
}

func (r *RenderData) Port() int {
	return r.ctx.Port()
}

func (r *RenderData) Scheme() string {
	return r.ctx.Scheme()
}

func (r *RenderData) RequestURI() string {
	return r.ctx.RequestURI()
}

func (r *RenderData) GetNextURL(varNames ...string) string {
	return GetNextURL(r.ctx, varNames...)
}

func (r *RenderData) ReturnToCurrentURL(varNames ...string) string {
	return ReturnToCurrentURL(r.ctx, varNames...)
}

func (r *RenderData) WithNextURL(urlStr string, varNames ...string) string {
	return WithNextURL(r.ctx, urlStr, varNames...)
}

func (r *RenderData) GetOtherURL(urlStr string, next string) string {
	return GetOtherURL(r.ctx, next)
}
