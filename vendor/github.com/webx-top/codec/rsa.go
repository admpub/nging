package codec

func NewRSA() *RSA {
	return &RSA{}
}

type RSA struct {
	publicKey  *RSAPublicKey
	privateKey *RSAPrivateKey
}

func (r *RSA) SetPublicKey(publicKey []byte) (err error) {
	r.publicKey, err = NewRSAPublicKey(publicKey)
	return err
}

func (r *RSA) PublicKey() *RSAPublicKey {
	return r.publicKey
}

func (r *RSA) SetPrivateKey(privateKey []byte) (err error) {
	r.privateKey, err = NewRSAPrivateKey(privateKey)
	return err
}

func (r *RSA) PrivateKey() *RSAPrivateKey {
	return r.privateKey
}
