package ssh

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func New(config *ssh.ClientConfig) *SSH {
	return &SSH{
		Config: config,
	}
}

type SSH struct {
	Config  *ssh.ClientConfig
	Client  *ssh.Client
	Session *ssh.Session
}

func (s *SSH) Connect(ip string, ports ...int) error {
	port := 22
	if len(ports) > 0 {
		port = ports[0]
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf(`%s:%d`, ip, port), s.Config)
	if err != nil {
		return errors.New("Failed to dial: " + err.Error())
	}
	s.Client = client

	session, err := client.NewSession()
	if err != nil {
		return errors.New("Failed to create session: " + err.Error())
	}
	s.Session = session

	return nil
}

func (s *SSH) Close() error {
	if s.Session == nil {
		return nil
	}
	return s.Session.Close()
}
