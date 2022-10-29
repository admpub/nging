package request

type Login struct {
	User string `validate:"required,username"`
	Pass string `validate:"required,min=8,max=64"`
	Code string `validate:"required"`
}
