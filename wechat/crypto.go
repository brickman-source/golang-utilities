package wechat

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"log"
	"reflect"
	"strings"
)

func SignXmlMd5(v interface{}, apikey string) (string, bool) {
	var signStr bytes.Buffer

	vt := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)

	for i := 0; i < vt.NumField(); i++ {
		field := vt.Field(i)
		name := field.Name

		log.Println(name)
		keytemp := field.Tag.Get("xml")
		keymap := strings.Split(keytemp, ",")
		key := keymap[0]

		cField := reflect.Indirect(vv).FieldByName(name)
		if key != "xml" && key != "sign" {
			if cField.Type().Name() == "CData" {
				value := cField.Field(0).String()
				if value != "" {
					signStr.WriteString(key + "=" + value + "&")
				}
			} else {
				value := cField.String()
				if value != "" {
					signStr.WriteString(key + "=" + value + "&")
				}
			}

		}
	}
	signStr.WriteString("key=" + apikey)

	log.Println("signstr", signStr.String())

	hasher := md5.New()
	hasher.Write([]byte(signStr.String()))
	sign := hex.EncodeToString(hasher.Sum(nil))
	sign = strings.ToUpper(sign)
	log.Println("signret", sign)

	return sign, true
}
