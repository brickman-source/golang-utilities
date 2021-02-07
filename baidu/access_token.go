/*
 * Copyright (c) 2018. Brickman Source.
 */

package baidu

import (
	"errors"
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/json"
	"github.com/brickman-source/golang-utilities/log"
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

func (bd *Baidu) GetAccessTokenByClient(apiKey, secretKey string) (ret *BaiduToken, err error) {
	token := bd.loadTokenFromCache(apiKey)
	if token == nil {
		log.Infof("GetAccessTokenByClient %v appId=%s appSecret=%s", bd, apiKey, secretKey)
		token, err = bd.getAccessTokenByClient(apiKey, secretKey)
		if err != nil {
			return
		}
	}
	return
}

func (bd *Baidu) getAccessTokenByClient(apiKey, secretKey string) (*BaiduToken, error) {
	ret := &BaiduToken{}

	getTokenURL, _ := url.Parse("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials")
	parameters := getTokenURL.Query()

	parameters.Set("client_id", apiKey)
	parameters.Set("client_secret", secretKey)

	getTokenURL.RawQuery = parameters.Encode()

	log.Infof("%s getAccessTokenByClient:%s", bd.config.GetString("application.name"), getTokenURL.String())

	err := http.GetJSON(getTokenURL.String(), ret)
	if err != nil {
		return nil, err
	}
	if ret.Error != "" {
		return nil, errors.New(ret.ErrorDescription)
	}

	ret.ExpiresAt = time.Now().Unix() + ret.ExpiresIn

	log.Infof("%s getAccessTokenByClient new token: %v %v", apiKey, bd.config.GetString("application.name"), ret)

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
	if bd.cache != nil {
		ret := &BaiduToken{}
		err := bd.cache.Get("bd:access_token:"+bd.config.GetString("application.name")+":"+apiKey, ret)
		if err == nil {
			log.Infof("access token from cache: %v", ret)
			return ret
		}
	}
	if val, ok := bd.memory.Load("bd:access_token:" + apiKey); ok {
		ret := &BaiduToken{}
		err := json.Unmarshal([]byte(val.(string)), ret)
		if err == nil {
			log.Infof("access token from cache: %v", ret)
			if ret.ExpiresAt <= time.Now().Unix()-1000 {
				return nil
			}
			return ret
		}
	}
	return nil
}
