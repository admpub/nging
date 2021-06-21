package codec

func NewRSA() *RSA {
	return &RSA{}
}

type RSA struct {
	publicKey  *RSAPublicKey
	privateKey *RSAPrivateKey
}

func (r *RSA) SetPublicKey(pubKey string) (err error) {
	r.publicKey, err = NewRSAPublicKey(pubKey)
	return err
}

func (r *RSA) PublicKey() *RSAPublicKey {
	return r.publicKey
}

func (r *RSA) SetPrivateKey(privKey string) (err error) {
	r.privateKey, err = NewRSAPrivateKey(privKey)
	return err
}

func (r *RSA) PrivateKey() *RSAPrivateKey {
	return r.privateKey
}
