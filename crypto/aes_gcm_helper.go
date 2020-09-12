package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"golang.org/x/crypto/scrypt"
	"log"
)

func AESEncryptGCM(data []byte, key [32]byte) ([]byte, error) {
	//	初始化 block cipher
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	//	设置 block cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	//	生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	//	封装、返回
	return gcm.Seal(nonce, nonce, data, nil), nil
}
func AESDecryptGCM(cipherText []byte, key [32]byte) (plaintext []byte, err error) {
	//	初始化 block cipher
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	//	设置 block cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	//	返回解开的包，注意这里的 nonce 是直接取的。
	return gcm.Open(nil,
		cipherText[:gcm.NonceSize()],
		cipherText[gcm.NonceSize():],
		nil,
	)
}
func NewKey(salt, password []byte) [32]byte {
	key, err := scrypt.Key(password, salt, 16384, 8, 1, 32)

	if err != nil {
		log.Fatal(err)
	}
	var ret [32]byte

	copy(ret[:], key[:32])

	return ret
}
