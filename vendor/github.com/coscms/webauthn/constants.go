package webauthn

type Type int
type Stage int

const (
	TypeRegister Type = iota + 1
	TypeLogin
	TypeUnbind
)

const (
	StageBegin Stage = iota + 1
	StageFinish
)

const (
	sessionKeyRegister = `webauthn.registration`
	sessionKeyLogin    = `webauthn.authentication`
	sessionKeyUnbind   = `webauthn.unbind`
)
