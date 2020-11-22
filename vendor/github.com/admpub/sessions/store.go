// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sessions

import (
	"encoding/base32"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/admpub/securecookie"
	"github.com/webx-top/echo"
)

// Store is an interface for custom session stores.
//
// See CookieStore and FilesystemStore for examples.
type Store interface {
	// Get should return a cached session.
	Get(ctx echo.Context, name string) (*Session, error)

	// New should create and return a new session.
	//
	// Note that New should never return a nil session, even in the case of
	// an error if using the Registry infrastructure to cache the session.
	New(ctx echo.Context, name string) (*Session, error)
	Reload(ctx echo.Context, s *Session) error

	// Save should persist session to the underlying store implementation.
	Save(ctx echo.Context, s *Session) error
}

// CookieStore ----------------------------------------------------------------

// NewCookieStore returns a new CookieStore.
//
// Keys are defined in pairs to allow key rotation, but the common case is
// to set a single authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for
// encryption. The encryption key can be set to nil or omitted in the last
// pair, but the authentication key is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes.
// The encryption key, if set, must be either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256 modes.
//
// Use the convenience function securecookie.GenerateRandomKey() to create
// strong keys.
func NewCookieStore(keyPairs ...[]byte) *CookieStore {
	cs := &CookieStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
	}
	return cs
}

// CookieStore stores sessions using secure cookies.
type CookieStore struct {
	Codecs []securecookie.Codec
}

// Get returns a session for the given name after adding it to the registry.
//
// It returns a new session if the sessions doesn't exist. Access IsNew on
// the session to check if it is an existing session or a new one.
//
// It returns a new session and an error if the session exists but could
// not be decoded.
func (s *CookieStore) Get(ctx echo.Context, name string) (*Session, error) {
	return GetRegistry(ctx).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// The difference between New() and Get() is that calling New() twice will
// decode the session data twice, while Get() registers and reuses the same
// decoded session after the first call.
func (s *CookieStore) New(ctx echo.Context, name string) (*Session, error) {
	session := NewSession(s, name)
	session.IsNew = true
	var err error
	if v := ctx.GetCookie(name); len(v) > 0 {
		err = securecookie.DecodeMultiWithMaxAge(
			name, v, &session.Values,
			ctx.CookieOptions().MaxAge,
			s.Codecs...)
		if err == nil {
			session.IsNew = false
		}
	}
	return session, err
}

func (s *CookieStore) Reload(_ echo.Context, _ *Session) error {
	return nil
}

// Save adds a single session to the response.
func (s *CookieStore) Save(ctx echo.Context, session *Session) error {
	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values,
		s.Codecs...)
	if err != nil {
		return err
	}
	SetCookie(ctx, session.Name(), encoded)
	return nil
}

// MaxAge sets the maximum age for the store and the underlying cookie
// implementation. Individual sessions can be deleted by setting Options.MaxAge
// = -1 for that session.
func (s *CookieStore) MaxAge(age int) {
	// Set the maxAge for each securecookie instance.
	for _, codec := range s.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

// FilesystemStore ------------------------------------------------------------

var fileMutex sync.RWMutex

// NewFilesystemStore returns a new FilesystemStore.
//
// The path argument is the directory where sessions will be saved. If empty
// it will use os.TempDir().
//
// See NewCookieStore() for a description of the other parameters.
func NewFilesystemStore(path string, keyPairs ...[]byte) *FilesystemStore {
	if path == "" {
		path = os.TempDir()
	}
	fs := &FilesystemStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		path:   path,
	}
	return fs
}

// FilesystemStore stores sessions in the filesystem.
//
// It also serves as a reference for custom stores.
//
// This store is still experimental and not well tested. Feedback is welcome.
type FilesystemStore struct {
	Codecs []securecookie.Codec
	path   string
}

// MaxLength restricts the maximum length of new sessions to l.
// If l is 0 there is no limit to the size of a session, use with caution.
// The default for a new FilesystemStore is 4096.
func (s *FilesystemStore) MaxLength(l int) {
	for _, c := range s.Codecs {
		if codec, ok := c.(*securecookie.SecureCookie); ok {
			codec.MaxLength(l)
		}
	}
}

// Get returns a session for the given name after adding it to the registry.
//
// See CookieStore.Get().
func (s *FilesystemStore) Get(ctx echo.Context, name string) (*Session, error) {
	return GetRegistry(ctx).Get(s, name)
}

// New returns a session for the given name without adding it to the registry.
//
// See CookieStore.New().
func (s *FilesystemStore) New(ctx echo.Context, name string) (*Session, error) {
	session := NewSession(s, name)
	session.IsNew = true
	var err error
	if v := ctx.GetCookie(name); len(v) > 0 {
		err = securecookie.DecodeMultiWithMaxAge(
			name, v, &session.ID,
			ctx.CookieOptions().MaxAge,
			s.Codecs...)
		if err == nil {
			err = s.load(ctx, session)
			if err == nil {
				session.IsNew = false
			}
		}
	}
	return session, err
}

func (s *FilesystemStore) Reload(ctx echo.Context, session *Session) error {
	err := s.load(ctx, session)
	if err == nil {
		session.IsNew = false
	}
	return err
}

// Save adds a single session to the response.
func (s *FilesystemStore) Save(ctx echo.Context, session *Session) error {
	// Delete if max-age is < 0
	if ctx.CookieOptions().MaxAge < 0 {
		if err := s.erase(session); err != nil {
			return err
		}
		SetCookie(ctx, session.Name(), "", -1)
		return nil
	}
	if len(session.ID) == 0 {
		// Because the ID is used in the filename, encode it to
		// use alphanumeric characters only.
		session.ID = strings.TrimRight(
			base32.StdEncoding.EncodeToString(
				securecookie.GenerateRandomKey(32)), "=")
	}
	if err := s.save(session); err != nil {
		return err
	}
	encoded, err := securecookie.EncodeMulti(session.Name(), session.ID,
		s.Codecs...)
	if err != nil {
		return err
	}
	SetCookie(ctx, session.Name(), encoded)
	return nil
}

// delete session file
func (s *FilesystemStore) erase(session *Session) error {
	if len(session.ID) == 0 {
		return nil
	}
	filename := filepath.Join(s.path, "session_"+session.ID)
	fileMutex.RLock()
	defer fileMutex.RUnlock()

	err := os.Remove(filename)
	return err
}

// MaxAge sets the maximum age for the store and the underlying cookie
// implementation. Individual sessions can be deleted by setting Options.MaxAge
// = -1 for that session.
func (s *FilesystemStore) MaxAge(age int) {
	// Set the maxAge for each securecookie instance.
	for _, codec := range s.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

// save writes encoded session.Values to a file.
func (s *FilesystemStore) save(session *Session) error {
	encoded, err := securecookie.EncodeMulti(session.Name(), session.Values,
		s.Codecs...)
	if err != nil {
		return err
	}
	filename := filepath.Join(s.path, "session_"+session.ID)
	fileMutex.Lock()
	defer fileMutex.Unlock()
	return ioutil.WriteFile(filename, []byte(encoded), 0600)
}

// load reads a file and decodes its content into session.Values.
func (s *FilesystemStore) load(ctx echo.Context, session *Session) error {
	filename := filepath.Join(s.path, "session_"+session.ID)
	fileMutex.RLock()
	defer fileMutex.RUnlock()
	fdata, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err = securecookie.DecodeMultiWithMaxAge(
		session.Name(), string(fdata),
		&session.Values,
		ctx.CookieOptions().MaxAge,
		s.Codecs...); err != nil {
		return err
	}
	return nil
}
