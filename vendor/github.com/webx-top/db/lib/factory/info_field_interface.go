package factory

type FieldInfor interface {
	GetName() string
	GetDataType() string
	IsUnsigned() bool
	IsPrimaryKey() bool
	IsAutoIncrement() bool
	GetMinValue() float64
	GetMaxValue() float64
	GetPrecision() int
	GetMaxSize() int
	GetOptions() []string
	GetDefaultValue() string
	GetComment() string
	GetGoType() string
	GetMyType() string
	GetGoName() string

	// - utils
	HTMLAttrBuilder(required bool) HTMLAttrs
	Validate(value interface{}) error
}
