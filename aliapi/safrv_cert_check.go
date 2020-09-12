package aliapi

import (
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/log"
	gohttp "net/http"
	"net/url"
)

const (
	IDCardCert_API_URL = "https://idcardcert.market.alicloudapi.com/idCardCert"
)

func (api *AliApiHelper) CheckIdCardNameMatch(idCardNo, name string) (ret []byte, err error) {
	reqURL, _ := url.Parse(IDCardCert_API_URL)
	q := reqURL.Query()
	q.Set("idCard", idCardNo)
	q.Set("name", name)
	reqURL.RawQuery = q.Encode()

	log.Debugf(nil, "check id name match:%s", reqURL.String())
	_, respData, err := http.GetDataWithHeaders(
		reqURL.String(),
		gohttp.Header{
			"Authorization": []string{"APPCODE " + api.cfg.GetString("aliapi.ocr.idcard.appcode")},
		},
	)
	return respData, err
}

type CheckIdCardNameMatchResult struct {
	AddrCode   string `json:"addrCode" xml:"addrCode"`
	Area       string `json:"area" xml:"area"`
	Birthday   string `json:"birthday" xml:"birthday"`
	City       string `json:"city" xml:"city"`
	IDCard     string `json:"idCard" xml:"idCard"`
	LastCode   string `json:"lastCode" xml:"lastCode"`
	Msg        string `json:"msg" xml:"msg"`
	Name       string `json:"name" xml:"name"`
	Prefecture string `json:"prefecture" xml:"prefecture"`
	Province   string `json:"province" xml:"province"`
	Sex        string `json:"sex" xml:"sex"`
	Status     string `json:"status" xml:"status"`
}
