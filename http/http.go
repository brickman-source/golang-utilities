package http

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/brickman-source/golang-utilities/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	goHttp "net/http"
	"net/http/cookiejar"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func GetData(reqURL string) ([]byte, error) {
	tr := &goHttp.Transport{
		TLSClientConfig: &tls.Config{},
	}
	cookieJar, _ := cookiejar.New(nil)

	client := &goHttp.Client{Transport: tr, Jar: cookieJar, Timeout: time.Duration(time.Second * 30)}
	request, _ := goHttp.NewRequest("GET", reqURL, nil)

	response, err := client.Do(request)

	if err != nil {
		// log.Errorf( "Get data error:%s", err.Error())
		return nil, err
	}

	defer response.Body.Close()

	var body []byte

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		body, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
	default:
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
	}

	if response.StatusCode == 200 {
		return body, nil
	} else {
		return body, errors.WithCode(nil, fmt.Sprintf("HTTP_%d", response.StatusCode), response.Status)
	}
}
func GetDataWithHeaders(reqURL string, reqHeader goHttp.Header) (goHttp.Header, []byte, error) {
	tr := &goHttp.Transport{
		TLSClientConfig: &tls.Config{},
	}
	cookieJar, _ := cookiejar.New(nil)

	client := &goHttp.Client{Transport: tr, Jar: cookieJar}
	request, _ := goHttp.NewRequest("GET", reqURL, nil)

	if reqHeader != nil {
		for key, value := range reqHeader {
			if len(value) > 0 {
				request.Header.Set(key, value[0])
			}
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()
	var body []byte
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, nil, err
		}
		defer reader.Close()
		body, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, nil, err
		}
	default:
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, nil, err
		}
	}

	if response.StatusCode == 200 {
		return response.Header, body, nil
	} else {
		return response.Header, body, errors.WithCode(nil, fmt.Sprintf("HTTP_%d", response.StatusCode), response.Status)
	}
}
func GetJSON(url string, out interface{}) error {
	//resp, err := goHttp.Get(url)

	client := &goHttp.Client{}
	request, _ := goHttp.NewRequest("GET", url, nil)

	//request.Header.Set("Accept-Encoding", "gzip")

	resp, err := client.Do(request)

	if err != nil {
		return err
	}
	return readJSON(resp, out)
}

func PutDataWithHeaders(reqURL string, reqHeader goHttp.Header, bodyType string, data []byte) (goHttp.Header, []byte, error) {
	tr := &goHttp.Transport{
		TLSClientConfig: &tls.Config{},
	}
	client := &goHttp.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
	}

	request, _ := goHttp.NewRequest("PUT", reqURL, bytes.NewReader(data))

	if reqHeader != nil {
		for key, value := range reqHeader {
			if len(value) > 0 {
				request.Header.Set(key, value[0])
			}
		}
	}
	request.Header.Set("Content-Type", bodyType)
	response, err := client.Do(request)

	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()

	var body []byte

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, nil, err
		}
		defer reader.Close()

		body, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, nil, err
		}
	default:
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, nil, err
		}
	}
	if response.StatusCode == 200 {
		return response.Header, body, nil
	} else {
		return response.Header, body, errors.WithCode(nil, fmt.Sprintf("HTTP_%d", response.StatusCode), response.Status)
	}
}

func PostData(reqURL string, bodyType string, data []byte) ([]byte, error) {
	timeout := 15 * time.Second
	client := &goHttp.Client{
		Timeout: timeout,
	}

	response, err := client.Post(reqURL, bodyType, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var body []byte

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		body, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, err
		}
	default:
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
	}
	if response.StatusCode == 200 {
		return body, nil
	} else {
		return body, errors.WithCode(nil, fmt.Sprintf("HTTP_%d", response.StatusCode), response.Status)
	}
}
func PostDataWithHeaders(reqURL string, reqHeader goHttp.Header, bodyType string, data []byte, timeout time.Duration) (goHttp.Header, []byte, error) {
	tr := &goHttp.Transport{
		TLSClientConfig: &tls.Config{},
	}
	client := &goHttp.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	request, _ := goHttp.NewRequest("POST", reqURL, bytes.NewReader(data))

	if reqHeader != nil {
		for key, value := range reqHeader {
			if len(value) > 0 {
				request.Header.Set(key, value[0])
			}
		}
	}
	request.Header.Set("Content-Type", bodyType)
	response, err := client.Do(request)

	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()

	var body []byte

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, nil, err
		}
		defer reader.Close()

		body, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, nil, err
		}
	default:
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, nil, err
		}
	}
	if response.StatusCode == 200 {
		return response.Header, body, nil
	} else {
		return response.Header, body, errors.WithCode(nil, fmt.Sprintf("HTTP_%d", response.StatusCode), response.Status)
	}
}
func PostJSON(reqURL string, objToSend interface{}, out interface{}) error {
	jsonData, err := json.Marshal(objToSend)
	if err != nil {
		// log.Errorf( "marshal json err:%s", err.Error())
		return err
	}

	resp, err := goHttp.Post(reqURL, "application/json;charset=utf-8", bytes.NewReader(jsonData))
	if err != nil {
		// log.Errorf( "post json data to(%s)  err:%s", reqURL, err.Error())
		return err
	}
	return readJSON(resp, out)
}
func PostDataSsl(reqURL string, dataToSend, certPEMBlock, keyPEMBlock []byte) (respData []byte, err error) {
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return nil, err
	}

	tr := &goHttp.Transport{
		TLSClientConfig: &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true},
	}
	client := &goHttp.Client{Transport: tr}

	ret, err := client.Post(reqURL, "application/x-www-form-urlencoded", bytes.NewReader(dataToSend))
	if err != nil {
		return nil, err
	}
	defer func() {
		err = ret.Body.Close()
		if err != nil {
			// log.Errorf( err.Error())
		}
	}()

	data, err := ioutil.ReadAll(ret.Body)

	if err != nil {
		return nil, err
	}

	return data, err
}

func PostXml(reqURL string, xmlToSend string, objReceived interface{}) (respData []byte, err error) {
	ret, err := goHttp.Post(reqURL, "application/x-www-form-urlencoded;charset=utf-8", strings.NewReader(xmlToSend))

	if err != nil {
		// log.Errorf( "post xml err:%s", err.Error())
		return nil, err
	}

	defer ret.Body.Close()

	data, err := ioutil.ReadAll(ret.Body)

	if err != nil {
		// log.Errorf( "post xml err:%s", err.Error())
		return nil, err
	}

	err = xml.Unmarshal(data, objReceived)

	return data, err
}

func PostFiles(reqURL string, values map[string][]string, progressReporter func(r int64)) (ret []byte, err error) {
	var b ProgressReader
	b.Reporter = progressReporter

	w := multipart.NewWriter(&b)
	for key, files := range values {
		for _, file := range files {
			var fw io.Writer
			fileName := filepath.Base(file)
			if fw, err = w.CreateFormFile(key, fileName); err != nil {
				return
			}
			f, err := os.Open(file)
			if err != nil {
				return nil, err
			}

			if _, err = io.Copy(fw, f); err != nil {
				return nil, err
			}
		}
	}

	w.Close()

	client := &goHttp.Client{}
	request, _ := goHttp.NewRequest("POST", reqURL, &b)

	request.Header.Set("Content-Type", w.FormDataContentType())

	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	var body []byte

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		var reader *gzip.Reader
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return
		}
		defer reader.Close()

		body, err = ioutil.ReadAll(reader)
		if err != nil {
			return
		}
	default:
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return
		}
	}

	if response.StatusCode == 200 {
		return body, nil
	} else {
		return body, errors.WithCode(nil, fmt.Sprintf("HTTP_%d", response.StatusCode), response.Status)
	}
}
func DownloadFile(reqURL string, filePath string) error {
	tr := &goHttp.Transport{
		//TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	cookieJar, _ := cookiejar.New(nil)

	client := &goHttp.Client{Transport: tr, Jar: cookieJar, Timeout: time.Duration(time.Second * 30)}
	request, _ := goHttp.NewRequest("GET", reqURL, nil)

	response, err := client.Do(request)

	if err != nil {
		// log.Errorf( "Get data error:%s", err.Error())
		return err
	}

	defer response.Body.Close()

	dir := path.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			return err
		}
		defer reader.Close()

		writer, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		defer writer.Close()

		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}
	default:
		writer, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		defer writer.Close()

		_, err = io.Copy(writer, response.Body)
		if err != nil {
			return err
		}
	}
	if response.StatusCode == 200 {
		return nil
	} else {
		return errors.WithCode(nil, fmt.Sprintf("HTTP_%d", response.StatusCode), response.Status)
	}
}

func DeleteDataWithHeaders(reqURL string, reqHeader goHttp.Header) (goHttp.Header, []byte, error) {
	tr := &goHttp.Transport{
		TLSClientConfig: &tls.Config{},
	}
	cookieJar, _ := cookiejar.New(nil)

	client := &goHttp.Client{Transport: tr, Jar: cookieJar}
	request, _ := goHttp.NewRequest("DELETE", reqURL, nil)

	if reqHeader != nil {
		for key, value := range reqHeader {
			if len(value) > 0 {
				request.Header.Set(key, value[0])
			}
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()
	var body []byte
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(response.Body)
		if err != nil {
			return nil, nil, err
		}
		defer reader.Close()
		body, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, nil, err
		}
	default:
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, nil, err
		}
	}

	if response.StatusCode == 200 {
		return response.Header, body, nil
	} else {
		return response.Header, body, errors.WithCode(nil, fmt.Sprintf("HTTP_%d", response.StatusCode), response.Status)
	}
}

func readJSON(resp *goHttp.Response, out interface{}) (err error) {
	defer resp.Body.Close()

	var reader io.ReadCloser
	//// log.Infof( "Content-Encoding:%s", resp.Header.Get("Content-Encoding"))
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		// log.Infof( "response:gzip")
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	if resp.StatusCode >= 400 {
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}

		return errors.WithCode(nil, fmt.Sprintf("HTTP_%d", resp.StatusCode), string(body))
	}

	if out == nil {
		io.Copy(ioutil.Discard, reader)
		return nil
	}

	decoder := json.NewDecoder(reader)
	return decoder.Decode(out)
}
