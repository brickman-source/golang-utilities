/*
 * Copyright (c) 2018. Brickman.
 */

package baidu

import (
	"errors"
	"github.com/brickman-source/golang-utilities/http"
	"net/url"
	"time"
)

func (bd *Baidu) GetAccessTokenOpenApiByClient(apiKey, secretKey string) (token *BaiduToken, err error) {
	token = bd.loadTokenOpenApiFromCache(apiKey)
	if token == nil {
		bd.logf("GetAccessTokenBceByClient %v appId=%s appSecret=%s", bd, apiKey, secretKey)
		token, err = bd.getAccessTokenOpenApi(apiKey, secretKey)
		if err != nil {
			return
		}
	}
	return
}

func (bd *Baidu) getAccessTokenOpenApi(apiKey, secretKey string) (*BaiduToken, error) {
	ret := &BaiduToken{}

	getTokenURL, _ := url.Parse("https://openapi.baidu.com/oauth/2.0/token?grant_type=client_credentials")
	parameters := getTokenURL.Query()

	parameters.Set("client_id", apiKey)
	parameters.Set("client_secret", secretKey)

	getTokenURL.RawQuery = parameters.Encode()

	err := http.GetJSON(getTokenURL.String(), ret)
	if err != nil {
		return nil, err
	}
	if ret.Error != "" {
		return nil, errors.New(ret.ErrorDescription)
	}

	ret.ExpiresAt = time.Now().Unix() + ret.ExpiresIn

	bd.storeTokenOpenApiToCache(apiKey, ret, time.Second*time.Duration(ret.ExpiresIn))

	return ret, nil
}

func (bd *Baidu) storeTokenOpenApiToCache(apiKey string, cacheVal *BaiduToken, expiresIn time.Duration) {
	if bd.cache != nil {
		err := bd.cache.Set(
			"bd:access_token_open_api:"+bd.config.GetString("application.name")+":"+apiKey,
			cacheVal,
			expiresIn,
		)
		if err == nil {
			return
		}
	}
	bd.memory.Store("bd:access_token_open_api:"+apiKey, cacheVal)
}

func (bd *Baidu) loadTokenOpenApiFromCache(apiKey string) *BaiduToken {
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
		err := bd.cache.Get("bd:access_token_open_api:"+bd.config.GetString("application.name")+":"+apiKey, ret)
		if err == nil && isValidFunc(ret) {
			bd.logf("access token from cache: %v", ret)
			return ret
		}
	} else if val, ok := bd.memory.Load("bd:access_token_open_api:" + apiKey); ok && val != nil {
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
