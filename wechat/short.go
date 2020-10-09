package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brickman-source/golang-utilities/http"
	"net/url"
)

func (wx *Wechat) GetShortURL(longURL string) (string, error) {
	accessToken, err := wx.GetAccessTokenByClient(
		wx.config.GetString("wechat.offiaccount.appId"),
		wx.config.GetString("wechat.offiaccount.appSecret"))
	if err != nil {
		return longURL, nil
	}

	requestURL, _ := url.Parse("https://api.weixin.qq.com/cgi-bin/shorturl")
	parameters := requestURL.Query()
	parameters.Set("access_token", accessToken.AccessToken)
	requestURL.RawQuery = parameters.Encode()

	ret := &struct {
		Errcode  int    `json:"errcode" xml:"errcode"`
		Errmsg   string `json:"errmsg" xml:"errmsg"`
		ShortURL string `json:"short_url" xml:"short_url"`
	}{}
	respData, err := http.PostData(
		requestURL.String(),
		"",
		[]byte(fmt.Sprintf(`{"action":"long2short","long_url":"%s"}`, longURL)),
	)
	if err != nil {
		return longURL, err
	}
	err = json.Unmarshal(respData, ret)
	if err != nil {
		return longURL, err
	}
	if ret.Errcode != 0 {
		return longURL, errors.New(ret.Errmsg)
	}
	return ret.ShortURL, nil
}
