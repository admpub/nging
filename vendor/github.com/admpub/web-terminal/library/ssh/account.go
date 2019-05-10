package ssh

type AccountConfig struct {
	User       string
	Password   string
	PrivateKey []byte
	Passphrase []byte
}
