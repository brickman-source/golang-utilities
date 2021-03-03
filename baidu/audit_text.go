/*
 * Copyright (c) 2018. Brickman Source.
 */

package baidu

import (
	"encoding/json"
	"github.com/brickman-source/golang-utilities/http"
	"net/url"
)

func (bd *Baidu) AuditText(str string, appId, appSecret string) (*BaiduAuditResult, error) {
	accessToken, err := bd.GetAccessTokenByClient(appId, appSecret)
	if err != nil {
		bd.logf( "cannot get access token(%v): %v", appId, err.Error())
		return nil, err
	}

	detectURL, _ := url.Parse(`https://aip.baidubce.com/rest/2.0/solution/v1/text_censor/v2/user_defined`)
	detectQuery := detectURL.Query()
	detectQuery.Set("access_token", accessToken.AccessToken)
	detectURL.RawQuery = detectQuery.Encode()

	params := url.Values{
		"text": {str},
	}
	bdRecognizeResultData, err := http.PostData(detectURL.String(), "application/x-www-form-urlencoded", []byte(params.Encode()))

	if err != nil {
		bd.logf( "recognize picture err: %v", err)
		return nil, err
	}
	baiduAuditTextResult := &BaiduAuditResult{}
	err = json.Unmarshal(bdRecognizeResultData, baiduAuditTextResult)
	if err != nil {
		bd.logf( "recognize picture err: %v", err)
		return nil, err
	}
	return baiduAuditTextResult, nil
}
