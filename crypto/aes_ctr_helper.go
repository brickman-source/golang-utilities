package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func AESEncryptCTR(src []byte, key []byte, iv []byte) (dst []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	stream := cipher.NewCTR(block, iv)
	encrypted := make([]byte, len(src))
	stream.XORKeyStream(encrypted, src)
	return encrypted, nil
}

func AESDecryptCTR(src []byte, key []byte, iv []byte) (dst []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	stream := cipher.NewCTR(block, iv)
	decrypted := make([]byte, len(src))
	stream.XORKeyStream(decrypted, src)
	return decrypted, nil
}

