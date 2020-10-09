package wechat

import (
	"errors"
	"github.com/brickman-source/golang-utilities/http"
	"net/url"
)

func (wx *Wechat) SendTemplateMsg(param interface{}) error {
	accessToken, err := wx.GetAccessTokenByClient(
		wx.config.GetString("wechat.offiaccount.appId"),
		wx.config.GetString("wechat.offiaccount.appSecret"))
	if err != nil {
		return err
	}

	requestURL, _ := url.Parse("https://api.weixin.qq.com/cgi-bin/message/template/send")
	parameters := requestURL.Query()
	parameters.Set("access_token", accessToken.AccessToken)
	requestURL.RawQuery = parameters.Encode()

	ret := &struct {
		Errcode int    `json:"errcode" xml:"errcode"`
		Errmsg  string `json:"errmsg" xml:"errmsg"`
		Msgid   int    `json:"msgid" xml:"msgid"`
	}{}
	err = http.PostJSON(
		requestURL.String(),
		param,
		ret)
	if err != nil {
		return err
	}

	if ret.Errcode != 0 {
		return errors.New(ret.Errmsg)
	}
	return nil
}
