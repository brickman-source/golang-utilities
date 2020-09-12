package sms

import (
	"io/ioutil"
	"path"
	"runtime"
)

func getBitsAppKey() string {
	_, filename, _, _ := runtime.Caller(0)
	data, _ := ioutil.ReadFile(path.Join(path.Dir(filename), "qiniu_key.txt"))
	return string(data)
}

func getBitsSecretKey() string {
	_, filename, _, _ := runtime.Caller(0)
	data, _ := ioutil.ReadFile(path.Join(path.Dir(filename), "qiniu_secret.txt"))
	return string(data)
}
