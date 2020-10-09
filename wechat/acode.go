package wechat

import (
	"encoding/json"
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/log"
	"net/url"
)

func (wx *Wechat) GetWxAcodeUnlimit(scene, page string, appId, appSecret string) ([]byte, error) {
	accessToken, err := wx.GetAccessTokenByClient(appId, appSecret)
	if err != nil {
		log.Errorf( "cannot get access token: %v", err.Error())
		return nil, err
	}

	requestURL, _ := url.Parse("https://api.weixin.qq.com/wxa/getwxacodeunlimit")
	parameters := requestURL.Query()
	parameters.Set("access_token", accessToken.AccessToken)
	requestURL.RawQuery = parameters.Encode()

	jData, err := json.Marshal(&struct {
		Scene string `json:"scene"`
		Page  string `json:"page"`
	}{
		Scene: scene,
		Page:  page,
	})
	if err != nil {
		jData = []byte("{}")
	}

	log.Infof( "get wx acode unlimit: %v %v", requestURL.String(), string(jData))
	respData, err := http.PostData(
		requestURL.String(),
		"application/json;charset=UTF-8",
		jData,
	)

	if err != nil {
		log.Errorf( "get wx acode unlimit err: %v %v %v", err, requestURL.String(), string(jData))
	}
	return respData, err
}
