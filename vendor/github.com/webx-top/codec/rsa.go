package codec

func NewRSA() *RSA {
	return &RSA{}
}

type RSA struct {
	publicKey  *RSAPublicKey
	privateKey *RSAPrivateKey
}

func (r *RSA) SetPublicKeyBytes(publicKey []byte) (err error) {
	r.publicKey, err = NewRSAPublicKey(publicKey)
	return err
}

func (r *RSA) SetPublicKey(publicKey *RSAPublicKey) *RSA {
	r.publicKey = publicKey
	return r
}

func (r *RSA) PublicKey() *RSAPublicKey {
	return r.publicKey
}

func (r *RSA) SetPrivateKeyBytes(privateKey []byte) (err error) {
	r.privateKey, err = NewRSAPrivateKey(privateKey)
	return err
}

func (r *RSA) SetPrivateKey(privateKey *RSAPrivateKey) *RSA {
	r.privateKey = privateKey
	return r
}

func (r *RSA) PrivateKey() *RSAPrivateKey {
	return r.privateKey
}
