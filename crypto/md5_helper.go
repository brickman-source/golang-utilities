package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

type Data []byte

func (d Data) GetMd5() []byte {
	h := md5.New()
	h.Write([]byte(d))
	src := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}
func (d Data) GetMd5String() string {
	h := md5.New()
	h.Write([]byte(d))
	return hex.EncodeToString(h.Sum(nil))
}

type String string

func (d String) GetMd5() []byte {
	h := md5.New()
	h.Write([]byte(d))
	src := h.Sum(nil)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}
func (d String) GetMd5String() string {
	h := md5.New()
	h.Write([]byte(d))
	return hex.EncodeToString(h.Sum(nil))
}
