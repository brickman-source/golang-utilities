/*
 * Copyright (c) 2020. Brickman Source.
 */

package baidu

import (
	"errors"
	"github.com/brickman-source/golang-utilities/converter"
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/log"
	"net/url"
	"strings"
	"time"
)

type BaiduToken struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	RefreshToken     string `json:"refresh_token" xml:"refresh_token"`
	ExpiresIn        int    `json:"expires_in" xml:"expires_in"`
	SessionKey       string `json:"session_key" xml:"session_key"`
	AccessToken      string `json:"access_token" xml:"access_token"`
	Scope            string `json:"scope" xml:"scope"`
	SessionSecret    string `json:"session_secret" xml:"session_secret"`
}

func (bd *Baidu) GetAccessTokenByClient(apiKey, secretKey string) (*BaiduToken, error) {
	ret := &BaiduToken{}
	cacheKey := "bd:access_token:" + bd.config.GetString("application.name") + ":" + apiKey
	err := bd.cache.Get(cacheKey, ret)
	if err == nil {
		log.Infof(nil, "access token from cache: %v", ret)
		return ret, nil
	} else {
		log.Errorf(nil, "cannot get token %v", cacheKey)
	}
	return ret, err
}

func (bd *Baidu) getAccessTokenByClient(apiKey, secretKey string) (*BaiduToken, error) {
	ret := &BaiduToken{}
	cacheKey := "bd:access_token:" + bd.config.GetString("application.name") + ":" + apiKey
	err := bd.cache.Get(cacheKey, ret)
	if err != nil {
		log.Errorf(nil, "%s GetAccessTokenByClient appid err:%s", bd.config.GetString("application.name"), err.Error())
	} else {
		log.Infof(nil, "%s GetAccessTokenByClient old token: %v", bd.config.GetString("application.name"), ret.AccessToken)
	}

	getTokenURL, _ := url.Parse("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials")
	parameters := getTokenURL.Query()

	parameters.Set("client_id", apiKey)
	parameters.Set("client_secret", secretKey)

	getTokenURL.RawQuery = parameters.Encode()

	log.Infof(nil, "%s getAccessTokenByClient:%s", bd.config.GetString("application.name"), getTokenURL.String())

	err = http.GetJSON(getTokenURL.String(), ret)
	if err != nil {
		return nil, err
	}
	if ret.Error != "" {
		return nil, errors.New(ret.ErrorDescription)
	}
	log.Infof(nil, "%s getAccessTokenByClient new token: %v %v", cacheKey, bd.config.GetString("application.name"), ret)
	err = bd.cache.Set(cacheKey, ret, time.Second*time.Duration(ret.ExpiresIn))
	if err != nil {
		log.Errorf(nil, "set cache err: %v", err)
	}
	return ret, nil
}

func (bd *Baidu) fetchAccessTokensLoop() {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf(nil, "fetchAccessTokensLoop err: %v", r)
		}
	}()
	t := time.NewTimer(time.Second)
	for {
		select {
		case <-t.C:
			miniPrograms := bd.config.GetMapSlice("bd.baiduConfigurations")
			if miniPrograms == nil {
				return
			}
			minExpiresIn := 999999999
			for _, miniP := range miniPrograms {
				var apiKey, secretKey string
				for s, i := range miniP {
					if strings.ToLower(s) == "apikey" {
						apiKey = converter.AsString(i, "")
					}
					if strings.ToLower(s) == "secretkey" {
						secretKey = converter.AsString(i, "")
					}
				}
				//apiKey := helper.GetValueAsString(miniP, "apiKey", "")
				//secretKey := helper.GetValueAsString(miniP, "secretKey", "")
				token, err := bd.getAccessTokenByClient(
					apiKey,
					secretKey,
				)
				if err != nil {
					log.Errorf(nil, "fetch access token err: %v", err)
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
		case <-bd.quit:
			return
		}
	}
}
