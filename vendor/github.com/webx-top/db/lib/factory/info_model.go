package factory

func NewMI(short string, make func(connID int) Model, comment string) *ModelInstancer {
	return &ModelInstancer{
		Short:   short,
		Make:    make,
		Comment: comment,
	}
}

// ModelInstancer 模型实例化
type ModelInstancer struct {
	Short   string //表名称(不含前缀)
	Make    func(connID int) Model
	Comment string //注释
}

// ModelInstancers 模型实例化
type ModelInstancers map[string]*ModelInstancer // structName => *ModelInstancer

// Register 模型构造函数登记
func (m ModelInstancers) Register(instancers map[string]*ModelInstancer) {
	for structName, instancer := range instancers {
		m[structName] = instancer
	}
}

type TableNamers map[string]func(obj interface{}) string

// Register 注册表名称生成器(表名不带前缀)
func (t TableNamers) Register(namers map[string]func(obj interface{}) string) {
	for table, namer := range namers {
		t[table] = namer
	}
}

// Get 获取表名称生成器(表名不带前缀)
func (t TableNamers) Get(table string) func(obj interface{}) string {
	if namer, ok := t[table]; ok {
		return namer
	}
	return DefaultTableNamer(table)
}
