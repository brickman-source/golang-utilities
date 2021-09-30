package wechat

import (
	"errors"
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/brickman-source/golang-utilities/map/helper"
	"net/url"
	"time"
)

type AccessTokenResult struct {
	Errcode int    `json:"errcode,omitempty" bson:"-"`
	Errmsg  string `json:"errmsg,omitempty" bson:"-"`

	AccessToken string `json:"access_token,omitempty" bson:"accessToken,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty" bson:"expiresIn,omitempty"`
	RfreshToken string `json:"rfresh_token,omitempty" bson:"rfreshToken,omitempty"`

	OpenId string `json:"openid,omitempty" bson:"openid,omitempty"`
	Scope  string `json:"scope,omitempty" bson:"scope,omitempty"`
}

func (wx *Wechat) GetAccessTokenByCode(code string) (*AccessTokenResult, error) {
	requestURL, _ := url.Parse("https://api.weixin.qq.com/sns/oauth2/access_token")
	parameters := requestURL.Query()

	parameters.Set("appid", wx.config.GetString("wechat.offiaccount.appId"))
	parameters.Set("secret", wx.config.GetString("wechat.offiaccount.appSecret"))
	parameters.Set("code", code)
	parameters.Set("grant_type", "authorization_code")

	requestURL.RawQuery = parameters.Encode()

	ret := &AccessTokenResult{}

	err := http.GetJSON(requestURL.String(), ret)

	if err != nil {
		return nil, err
	}
	if ret.Errcode != 0 {
		return nil, errors.New(ret.Errmsg)
	}
	return ret, nil
}

func (wx *Wechat) GetAccessTokenByClient(appId, appSecret string) (*AccessTokenResult, error) {
	ret := &AccessTokenResult{}
	cacheKey := "wx:access_token:" + wx.config.GetString("application.name") + ":" + appId
	err := wx.cache.Get(cacheKey, ret)
	if err == nil {
		log.Infof("access token from cache: %v", ret)
		return ret, nil
	} else {
		log.Errorf("cannot get token %v", cacheKey)
	}
	return ret, err
}

func (wx *Wechat) getAccessTokenByClient(appId, appSecret string) (*AccessTokenResult, error) {
	ret := &AccessTokenResult{}
	cacheKey := "wx:access_token:" + wx.config.GetString("application.name") + ":" + appId
	err := wx.cache.Get(cacheKey, ret)
	if err != nil {
		log.Errorf("%s GetAccessTokenBceByClient appid err:%s", wx.config.GetString("application.name"), err.Error())
	} else {
		log.Infof("%s GetAccessTokenBceByClient old token: %v", wx.config.GetString("application.name"), ret.AccessToken)
	}

	getTokenURL, _ := url.Parse("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential")
	parameters := getTokenURL.Query()

	parameters.Set("appid", appId)
	parameters.Set("secret", appSecret)

	getTokenURL.RawQuery = parameters.Encode()

	log.Infof("%s getAccessTokenByClient:%s", wx.config.GetString("application.name"), getTokenURL.String())

	err = http.GetJSON(getTokenURL.String(), ret)
	if err != nil {
		return nil, err
	}
	if ret.Errcode != 0 {
		return nil, errors.New(ret.Errmsg)
	}
	log.Infof("%s getAccessTokenByClient new token: %v %v", cacheKey, wx.config.GetString("application.name"), ret)
	err = wx.cache.Set(cacheKey, ret, time.Second*time.Duration(ret.ExpiresIn))
	if err != nil {
		log.Errorf("set cache err: %v", err)
	}
	return ret, nil
}

func (wx *Wechat) fetchAccessTokensLoop() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("fetchAccessTokensLoop err: %v", r)
		}
	}()
	t := time.NewTimer(time.Second)
	for {
		select {
		case <-t.C:
			miniPrograms := wx.config.GetMapSlice("wx.miniPrograms")
			if miniPrograms == nil {
				return
			}
			minExpiresIn := 99999999
			for _, miniP := range miniPrograms {
				appId := helper.GetValueAsString(miniP, "appid", "")
				appSecret := helper.GetValueAsString(miniP, "appsecret", "")
				token, err := wx.getAccessTokenByClient(appId, appSecret)
				if err != nil {
					log.Errorf("fetch access token err: %v", err)
				}
				if token != nil {
					if token.ExpiresIn < minExpiresIn {
						minExpiresIn = token.ExpiresIn
					}
				}
			}
			if minExpiresIn < 120 {
				minExpiresIn = 120
			} else {
				minExpiresIn -= 60
			}
			t.Reset(time.Second * time.Duration(minExpiresIn))
		case <-wx.quit:
			return
		}
	}
}
