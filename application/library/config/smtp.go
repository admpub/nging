package config

import (
	"net/smtp"
	"strconv"
)

type SMTPConfig struct {
	Identity string `json:"identity"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *SMTPConfig) Address() string {
	if s.Port == 0 {
		s.Port = 25
	}
	return s.Host + `:` + strconv.Itoa(s.Port)
}

func (s *SMTPConfig) Auth() smtp.Auth {
	return smtp.PlainAuth(s.Identity, s.Username, s.Password, s.Host)
}
