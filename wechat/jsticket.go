/*
 * Copyright (c) 2019. Brickman Source.
 */

package wechat

import (
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/log"
	"time"
)

type TicketInfo struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"` // seconds
}

func (wx *Wechat) GetJSTicket(appId, appSecret string) (string, error) {
	var ret string

	cacheKey := wx.config.GetString("application.name") + ".wx_js_ticket." + appId
	err := wx.cache.Get(cacheKey, &ret)
	if err == nil && ret != "" {
		return ret, nil
	}
	if err != nil {
		log.Errorf( "%s GetJSTicket appid err:%s", wx.config.GetString("application.name"), err.Error())
	}

	accessToken, err := wx.GetAccessTokenByClient(appId, appSecret)
	if err != nil {
		return "", err
	}
	log.Infof( "[GetJSTicket]accessToken:%s", accessToken.AccessToken)
	result := &TicketInfo{}

	err = http.GetJSON(`https://api.weixin.qq.com/cgi-bin/ticket/getticket?type=jsapi&access_token=`+accessToken.AccessToken, result)
	if err != nil {
		return "", err
	}

	if result.ExpiresIn > 30 {
		result.ExpiresIn = result.ExpiresIn - 30
	}
	_ = wx.cache.Set(cacheKey, &result.Ticket, time.Second*time.Duration(result.ExpiresIn))

	return result.Ticket, nil
}
