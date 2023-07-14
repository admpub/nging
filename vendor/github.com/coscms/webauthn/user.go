package webauthn

import "github.com/go-webauthn/webauthn/webauthn"

var _ webauthn.User = &User{}

type User struct {
	ID          uint64
	Name        string
	DisplayName string
	Icon        string
	Credentials []webauthn.Credential
}

// User ID according to the Relying Party
func (u *User) WebAuthnID() []byte {
	return WebAuthnID(u.ID)
}

// User Name according to the Relying Party
func (u *User) WebAuthnName() string {
	return u.Name
}

// Display Name of the user
func (u *User) WebAuthnDisplayName() string {
	return u.DisplayName
}

// User's icon url
func (u *User) WebAuthnIcon() string {
	return u.Icon
}

// Credentials owned by the user
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}
