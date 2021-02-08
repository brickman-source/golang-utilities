/*
 * Copyright (c) 2018. Brickman Source.
 */

package baidu

import (
	"errors"
	"github.com/brickman-source/golang-utilities/http"
	"net/url"
	"time"
)

type BaiduToken struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	RefreshToken     string `json:"refresh_token" xml:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in" xml:"expires_in"`
	ExpiresAt        int64  `json:"expires_at" xml:"expires_at"`
	SessionKey       string `json:"session_key" xml:"session_key"`
	AccessToken      string `json:"access_token" xml:"access_token"`
	Scope            string `json:"scope" xml:"scope"`
	SessionSecret    string `json:"session_secret" xml:"session_secret"`
}

func (bd *Baidu) GetAccessTokenByClient(apiKey, secretKey string) (token *BaiduToken, err error) {
	token = bd.loadTokenFromCache(apiKey)
	if token == nil {
		bd.logf("GetAccessTokenByClient %v appId=%s appSecret=%s", bd, apiKey, secretKey)
		token, err = bd.getAccessToken(apiKey, secretKey)
		if err != nil {
			return
		}
	}
	return
}

func (bd *Baidu) getAccessToken(apiKey, secretKey string) (*BaiduToken, error) {
	ret := &BaiduToken{}

	getTokenURL, _ := url.Parse("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials")
	parameters := getTokenURL.Query()

	parameters.Set("client_id", apiKey)
	parameters.Set("client_secret", secretKey)

	getTokenURL.RawQuery = parameters.Encode()

	bd.logf("config %v %s getAccessToken:%s", bd.config, bd.config.GetString("application.name"), getTokenURL.String())

	err := http.GetJSON(getTokenURL.String(), ret)
	if err != nil {
		return nil, err
	}
	if ret.Error != "" {
		return nil, errors.New(ret.ErrorDescription)
	}

	ret.ExpiresAt = time.Now().Unix() + ret.ExpiresIn

	bd.logf("%s getAccessToken new token: %v %v", apiKey, bd.config.GetString("application.name"), ret)

	bd.storeTokenToCache(apiKey, ret, time.Second*time.Duration(ret.ExpiresIn))

	return ret, nil
}

func (bd *Baidu) storeTokenToCache(apiKey string, cacheVal *BaiduToken, expiresIn time.Duration) {
	if bd.cache != nil {
		err := bd.cache.Set(
			"bd:access_token:"+bd.config.GetString("application.name")+":"+apiKey,
			cacheVal,
			expiresIn,
		)
		if err == nil {
			return
		}
	}
	bd.memory.Store("bd:access_token:"+apiKey, cacheVal)
}

func (bd *Baidu) loadTokenFromCache(apiKey string) *BaiduToken {
	isValidFunc := func(t *BaiduToken) bool {
		if t.ExpiresAt <= time.Now().Unix()-1000 {
			bd.logf("token expired")
			return false
		}
		return true
	}
	if bd.cache != nil {
		bd.logf("cache is not null")
		ret := &BaiduToken{}
		err := bd.cache.Get("bd:access_token:"+bd.config.GetString("application.name")+":"+apiKey, ret)
		if err == nil && isValidFunc(ret) {
			bd.logf("access token from cache: %v", ret)
			return ret
		}
	} else if val, ok := bd.memory.Load("bd:access_token:" + apiKey); ok && val != nil{
		if ret, ok := val.(*BaiduToken); ok {
			bd.logf("access token from memory: %v", ret)
			if isValidFunc(ret) {
				return ret
			}
		}
	}
	bd.logf("didnt found token in cache or token is expired")
	return nil
}
