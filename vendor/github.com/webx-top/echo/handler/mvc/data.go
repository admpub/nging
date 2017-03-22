package mvc

type Data interface {
	Assign(key string, val interface{})
	Assignx(values *map[string]interface{})
	SetTmplFuncs()
	Render(tmpl string, code ...int) error
	String() string
	Set(code int, args ...interface{})
	Gets() (code int, info interface{}, zone interface{}, data interface{})
	GetData() interface{}
}
