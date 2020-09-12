package main

import (
	"flag"
	"github.com/brickman-source/golang-utilities/config"
	"github.com/brickman-source/golang-utilities/crypto"
	"github.com/brickman-source/golang-utilities/log"
)

func main() {
	aesKey := flag.String("aeskey", "", "")
	aesKeySalt := flag.String("aeskeysalt", "", "")
	typ := flag.String("type", "", "enc or dec")
	str := flag.String("str", "", "")

	flag.Parse()

	println("aeskey:", *aesKey)
	println("aeskeysalt:", *aesKeySalt)
	println("type:", *typ)
	println("str:", *str)

	cfg := &config.Config{
		AesKey: crypto.String(*aesKey + *aesKeySalt).GetMd5(),
	}
	if *typ == "enc" {
		encrypted := cfg.EncryptString(*str)
		log.Infof(nil, "encrypted: %s", encrypted)
	} else {
		decrypted := cfg.DecryptString(*str)
		log.Infof(nil, "decrypted: %s", decrypted)
	}
}
