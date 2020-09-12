package http

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetRequestBody(r *http.Request) []byte {
	if r.Method == "POST" || r.Method == "PUT" {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			println("ioutil.ReadAll(r.Body) err:", err.Error())
			return nil
		}
		return data
	} else {
		println("method:", r.Method)
	}
	return nil
}

func AddURLQuery(reqURL string, key, value string) string {
	x, err := url.Parse(reqURL)
	if err != nil {
		return reqURL
	}
	q := x.Query()
	q.Set(key, value)
	x.RawQuery = q.Encode()
	return x.String()
}
