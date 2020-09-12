package crypto

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/brickman-source/golang-utilities/errors"
)

var (
	ErrKeyMustBePEMEncoded = errors.New("Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key")
	ErrNotRSAPrivateKey    = errors.New("Key is not a valid RSA private key")
)

func RSAEncrypt(pubKey *rsa.PublicKey, data []byte) ([]byte, error) {
	var err error
	//key, err := ParseRSAPublicKeyFromPEM(keyPem)
	//if err != nil {
	//	return nil, err
	//}
	klen := pubKey.N.BitLen()/8 - 11
	if len(data) <= klen {
		var bb []byte
		if bb, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, data); err != nil {
			return nil, err
		}
		return bb, nil
	}
	var buf bytes.Buffer
	var bb []byte
	for i, w, r := 0, 0, len(data); r > 0; i, r = i+w, r-w {
		if r <= klen {
			if bb, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, data[i:]); err != nil {
				return nil, err
			}
			buf.Write(bb)
			w = r
		} else {
			if bb, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, data[i:i+klen]); err != nil {
				return nil, err
			}
			buf.Write(bb)
			w = klen
		}
	}
	return buf.Bytes(), nil
}
func RSADecrypt(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	var err error

	//key, err := ParseRSAPrivateKeyFromPEM(keyPem)
	//if err != nil {
	//	return nil, err
	//}

	klen := privateKey.N.BitLen() / 8
	if len(data) <= klen {
		var bb []byte
		if bb, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, data); err != nil {
			return nil, err
		}
		return bb, nil
	}
	var buf bytes.Buffer
	var bb []byte
	for i, w, r := 0, 0, len(data); r > 0; i, r = i+w, r-w {
		if r <= klen {
			if bb, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, data[i:]); err != nil {
				return nil, err
			}
			buf.Write(bb)
			w = r
		} else {
			if bb, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, data[i:i+klen]); err != nil {
				return nil, err
			}
			buf.Write(bb)
			w = klen
		}
	}
	return buf.Bytes(), nil
}
func RSASign(privateKey *rsa.PrivateKey, data []byte) (string, error) {
	var bb []byte
	var err error

	h := sha1.New()
	h.Write(data)
	digest := h.Sum(nil)

	bb, err = rsa.SignPKCS1v15(nil, privateKey, crypto.SHA1, digest)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bb), nil
}
func RSA256Sign(privateKey *rsa.PrivateKey, data []byte) (string, error) {
	var bb []byte
	var err error

	h := sha256.New()
	h.Write(data)
	digest := h.Sum(nil)

	bb, err = rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, digest)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bb), nil
}

func HmacSha1(source, secret string) string {
	key := []byte(secret)
	hasher := hmac.New(sha1.New, key)
	hasher.Write([]byte(source))
	signedBytes := hasher.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}

func RSAVerify(pubKey *rsa.PublicKey, data []byte, sign string) error {
	bs, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	hash := crypto.SHA1
	h := hash.New()
	h.Write(data)
	hashed := h.Sum(nil)

	return rsa.VerifyPKCS1v15(pubKey, hash, hashed, bs)
}

func RSA256Verify(pubKey *rsa.PublicKey, data []byte, sign string) error {
	bs, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	hash := crypto.SHA256
	h := hash.New()
	h.Write(data)
	hashed := h.Sum(nil)

	return rsa.VerifyPKCS1v15(pubKey, hash, hashed, bs)
}

func ParseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}
func ParseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, err
		}
	}

	var pkey *rsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}

func ParseRSAPublicKey(key []byte) (*rsa.PublicKey, error) {
	var err error

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(key); err != nil {
		if cert, err := x509.ParseCertificate(key); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, err
		}
	}

	var pkey *rsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}
