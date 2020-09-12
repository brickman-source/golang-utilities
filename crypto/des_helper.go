package crypto

import (
	"crypto/cipher"
	"crypto/des"
)

//
//func main() {
//	originalText := "yolan"
//	fmt.Println(originalText)
//
//	key := []byte{0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC, 0xBC}
//
//	// encrypt value to base64
//	cryptoText := encrypt(key, originalText)
//	fmt.Println(cryptoText)
//
//}

var (
	EmptyDesIV = []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
)

// encrypt string to base64 crypto using des
func DesEncrypt(data []byte, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData := Pkcs5Padding(data, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, EmptyDesIV)
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)

	return cryted, nil
}

func DesDecrypt(cipherText, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)

	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, EmptyDesIV)
	origData := make([]byte, len(cipherText))
	blockMode.CryptBlocks(origData, cipherText)
	origData = Pkcs5UnPadding(origData)
	return origData, nil

}
