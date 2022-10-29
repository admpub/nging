package request

type Login struct {
	User string `validate:"required,username"`
	Pass string `validate:"required,min=8,max=64"`
	Code string `validate:"required"`
}

type Register struct {
	InvitationCode       string `validate:"required,min=16,max=32"`
	Username             string `validate:"required,username"`
	Email                string `validate:"required,email"`
	Password             string `validate:"required,min=8,max=64"`
	ConfirmationPassword string `validate:"required,eqfield=Password"`
}
