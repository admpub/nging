package webauthn

import (
	"encoding/gob"
	"sync"

	"github.com/admpub/log"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/webx-top/echo"
)

func init() {
	gob.Register(&webauthn.SessionData{})
}

func New(user UserHandler, lazyInit ...func(echo.Context) *webauthn.Config) *Server {
	a := &Server{
		user: user,
	}
	if len(lazyInit) > 0 {
		a.lazyInit = lazyInit[0]
	}
	return a
}

type Server struct {
	webAuthn *webauthn.WebAuthn
	lazyInit func(echo.Context) *webauthn.Config
	onceInit sync.Once
	user     UserHandler
	lock     sync.RWMutex
}

func (s *Server) Init(cfg *webauthn.Config) error {
	var err error
	s.lock.Lock()
	s.webAuthn, err = webauthn.New(cfg)
	s.lock.Unlock()
	return err
}

func (s *Server) Object(ctx echo.Context) *webauthn.WebAuthn {
	s.lock.RLock()
	a := s.webAuthn
	s.lock.RUnlock()
	if a == nil && s.lazyInit != nil {
		s.onceInit.Do(func() {
			cfg := s.lazyInit(ctx)
			err := s.Init(cfg)
			if err != nil {
				log.Error(err)
				return
			}
			a = s.webAuthn
		})
	}
	return a
}

func (s *Server) RegisterRoute(r echo.RouteRegister) {

	g := r.Group(`/webauthn`)
	g.Post("/register/begin/:username", s.handleBeginRegistration).SetName(`webauthn.beginRegister`)
	g.Post("/register/finish/:username", s.handleFinishRegistration).SetName(`webauthn.finishRegister`)
	g.Post("/login/begin/:username", s.handleBeginLogin).SetName(`webauthn.beginLogin`)
	g.Post("/login/finish/:username", s.handleFinishLogin).SetName(`webauthn.finishLogin`)
	g.Post("/unbind/begin/:username", s.handleBeginUnbind).SetName(`webauthn.beginUnbind`)
	g.Post("/unbind/finish/:username", s.handleFinishUnbind).SetName(`webauthn.finishUnbind`)

}

func (s *Server) RegisterRouteForLogin(r echo.RouteRegister) {
	g := r.Group(`/webauthn`)
	g.Post("/login/begin/:username", s.handleBeginLogin).SetName(`webauthn.beginLogin`)
	g.Post("/login/finish/:username", s.handleFinishLogin).SetName(`webauthn.finishLogin`)
}

// - register -

func (s *Server) handleBeginRegistration(ctx echo.Context) error {

	// get username/friendly name
	username := ctx.Param("username")

	// get user
	user, err := s.user.GetUser(ctx, username, TypeRegister, StageBegin)
	if err != nil {
		return err
	}

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = credentialExcludeList(ctx, user)
	}

	// generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := s.Object(ctx).BeginRegistration(
		user,
		registerOptions,
	)

	if err != nil {
		return err
	}

	ctx.Session().Set(sessionKeyRegister, sessionData)

	return ctx.JSON(options)
}

func (s *Server) handleFinishRegistration(ctx echo.Context) error {

	// get username
	username := ctx.Param("username")

	// get user
	user, err := s.user.GetUser(ctx, username, TypeRegister, StageFinish)
	if err != nil {
		return err
	}

	// load the session data
	sessionData, ok := ctx.Session().Get(sessionKeyRegister).(*webauthn.SessionData)
	if !ok {
		return echo.ErrBadRequest
	}

	credential, err := s.Object(ctx).FinishRegistration(user, *sessionData, ctx.Request().StdRequest())
	if err != nil {
		return err
	}

	err = s.user.Register(ctx, user, credential)
	if err != nil {
		return err
	}

	ctx.Session().Delete(sessionKeyRegister)

	return ctx.JSON("Registration Success")
}

// - login -

func (s *Server) handleBeginLogin(ctx echo.Context) error {

	// get username
	username := ctx.Param("username")

	// get user
	user, err := s.user.GetUser(ctx, username, TypeLogin, StageBegin)
	if err != nil {
		return err
	}

	// generate PublicKeyCredentialRequestOptions, session data
	options, sessionData, err := s.Object(ctx).BeginLogin(user)
	if err != nil {
		return err
	}

	// store session data as marshaled JSON
	ctx.Session().Set(sessionKeyLogin, sessionData)

	return ctx.JSON(options)
}

func (s *Server) handleFinishLogin(ctx echo.Context) error {

	// get username
	username := ctx.Param("username")

	// get user
	user, err := s.user.GetUser(ctx, username, TypeLogin, StageFinish)
	if err != nil {
		return err
	}

	// load the session data
	sessionData, ok := ctx.Session().Get(sessionKeyLogin).(*webauthn.SessionData)
	if !ok {
		return echo.ErrBadRequest
	}

	// in an actual implementation, we should perform additional checks on
	// the returned 'credential', i.e. check 'credential.Authenticator.CloneWarning'
	// and then increment the credentials counter
	credential, err := s.Object(ctx).FinishLogin(user, *sessionData, ctx.Request().StdRequest())
	if err != nil {
		return err
	}

	// handle successful login
	err = s.user.Login(ctx, user, credential)
	if err != nil {
		return err
	}

	ctx.Session().Delete(sessionKeyLogin)

	return ctx.JSON("Login Success")
}

// - unbind -

func (s *Server) handleBeginUnbind(ctx echo.Context) error {

	// get username
	username := ctx.Param("username")

	// get user
	user, err := s.user.GetUser(ctx, username, TypeUnbind, StageBegin)
	if err != nil {
		return err
	}

	// generate PublicKeyCredentialRequestOptions, session data
	options, sessionData, err := s.Object(ctx).BeginLogin(user)
	if err != nil {
		return err
	}

	// store session data as marshaled JSON
	ctx.Session().Set(sessionKeyUnbind, sessionData)

	return ctx.JSON(options)
}

func (s *Server) handleFinishUnbind(ctx echo.Context) error {

	// get username
	username := ctx.Param("username")

	// get user
	user, err := s.user.GetUser(ctx, username, TypeUnbind, StageFinish)
	if err != nil {
		return err
	}

	// load the session data
	sessionData, ok := ctx.Session().Get(sessionKeyUnbind).(*webauthn.SessionData)
	if !ok {
		return echo.ErrBadRequest
	}

	// in an actual implementation, we should perform additional checks on
	// the returned 'credential', i.e. check 'credential.Authenticator.CloneWarning'
	// and then increment the credentials counter
	credential, err := s.Object(ctx).FinishLogin(user, *sessionData, ctx.Request().StdRequest())
	if err != nil {
		return err
	}

	// handle successful unbind
	err = s.user.Unbind(ctx, user, credential)
	if err != nil {
		return err
	}

	ctx.Session().Delete(sessionKeyUnbind)

	return ctx.JSON("Unbind Success")
}
