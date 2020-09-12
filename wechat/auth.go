package wechat

import (
	"net/url"
)

func (wx *Wechat) GetAuthURL(redirectUri, scope string) (string, error) {
	var cbScope string
	if scope == "base" {
		cbScope = "base"
		scope = "snsapi_base"
	} else if scope == "userinfo" {
		cbScope = "userinfo"
		scope = "snsapi_userinfo"
	} else {
		cbScope = "base"
		scope = "snsapi_base"
	}
	callbackURL, _ := url.Parse(wx.config.GetString("wechat.auth.callbackURL")) // url.Parse(fmt.Sprintf("%s/auth/callback", wx.config.GetString("host.binding")))
	query := callbackURL.Query()
	query.Set("redirect_uri", redirectUri)
	query.Set("scope", cbScope)
	callbackURL.RawQuery = query.Encode()

	requestURL, _ := url.Parse("https://open.weixin.qq.com/connect/oauth2/authorize")

	parameters := requestURL.Query()
	parameters.Set("appid", wx.config.GetString("wechat.appId"))
	parameters.Set("redirect_uri", callbackURL.String())
	parameters.Set("response_type", "code")
	parameters.Set("scope", scope)
	parameters.Set("state", "0")

	requestURL.Fragment = "wechat_redirect"

	requestURL.RawQuery = parameters.Encode()

	return requestURL.String(), nil
}
