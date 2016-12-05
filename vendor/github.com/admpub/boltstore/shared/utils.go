package shared

import (
	"time"

	"github.com/admpub/boltstore/shared/protobuf"
	"github.com/gogo/protobuf/proto"
)

// Session converts the byte slice to the session struct value.
func Session(data []byte) (protobuf.Session, error) {
	session := protobuf.Session{}
	err := proto.Unmarshal(data, &session)
	return session, err
}

// Expired checks if the session is expired.
func Expired(session protobuf.Session) bool {
	return *session.ExpiresAt > 0 && *session.ExpiresAt <= time.Now().Unix()
}

// NewSession creates and returns a session data.
func NewSession(values []byte, maxAge int) *protobuf.Session {
	expiresAt := time.Now().Unix() + int64(maxAge)
	return &protobuf.Session{Values: values, ExpiresAt: &expiresAt}
}
