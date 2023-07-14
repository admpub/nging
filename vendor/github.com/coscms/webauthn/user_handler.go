package webauthn

import (
	"encoding/binary"

	"github.com/webx-top/echo"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

type UserHandler interface {
	GetUser(ctx echo.Context, username string, opType Type, stage Stage) (webauthn.User, error)
	Register(ctx echo.Context, user webauthn.User, cred *webauthn.Credential) error
	Login(ctx echo.Context, user webauthn.User, cred *webauthn.Credential) error
	Unbind(ctx echo.Context, user webauthn.User, cred *webauthn.Credential) error
}

func credentialExcludeList(ctx echo.Context, user webauthn.User) []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range user.WebAuthnCredentials() {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}

func WebAuthnID(uid uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(uid))
	return buf
}
