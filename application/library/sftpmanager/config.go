package sftpmanager

import (
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/web-terminal/library/ssh"
	"github.com/pkg/sftp"
)

// Config 配置
type Config struct {
	Host       string `db:"host" bson:"host" comment:"主机名" json:"host" xml:"host"`
	Port       int    `db:"port" bson:"port" comment:"端口" json:"port" xml:"port"`
	Username   string `db:"username" bson:"username" comment:"用户名" json:"username" xml:"username"`
	Password   string `db:"password" bson:"password" comment:"密码" json:"password" xml:"password"`
	PrivateKey string `db:"private_key" bson:"private_key" comment:"私钥内容" json:"private_key" xml:"private_key"`
	Passphrase string `db:"passphrase" bson:"passphrase" comment:"私钥口令" json:"passphrase" xml:"passphrase"`
}

func (c *Config) Connect() (*sftp.Client, error) {
	account := &ssh.AccountConfig{
		User:     c.Username,
		Password: config.DefaultConfig.Decode(c.Password),
	}
	if len(c.PrivateKey) > 0 {
		account.PrivateKey = []byte(c.PrivateKey)
	}
	if len(c.Passphrase) > 0 {
		account.Passphrase = []byte(config.DefaultConfig.Decode(c.Passphrase))
	}
	config, err := ssh.NewSSHConfig(nil, nil, account)
	if err != nil {
		return nil, err
	}
	sshClient := ssh.New(config)
	err = sshClient.Connect(c.Host, c.Port)
	if err != nil {
		return nil, err
	}
	return sftp.NewClient(sshClient.Client)
}

func DefaultConnector(c *Config) (*sftp.Client, error) {
	return c.Connect()
}
