package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	//"log"
)

var (
	EmptyAESIV = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

func AESEncryptCBC(src []byte, key []byte, iv []byte) (dst []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	src, _ = Pkcs7Padding(src, aes.BlockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	dst = make([]byte, len(src))
	blockMode.CryptBlocks(dst, src)
	return
}

func AESDecryptCBC(src []byte, key []byte, iv []byte) (dst []byte, err error) {
	if len(src)%aes.BlockSize != 0 {
		return nil, errors.New("src is not correct")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	dst = make([]byte, len(src))

	blockMode.CryptBlocks(dst, src)
	//log.Println("dst:", dst)

	dst, err = Pkcs7UnPadding(dst, block.BlockSize())

	return
}
