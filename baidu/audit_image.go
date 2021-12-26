/*
 * Copyright (c) 2018. Brickman.
 */

package baidu

import (
	"encoding/base64"
	"encoding/json"
	"github.com/brickman-source/golang-utilities/http"
	"net/url"
)

func (bd *Baidu) AuditImage(data []byte, appId, appSecret string) (*BaiduAuditResult, error) {
	accessToken, err := bd.GetAccessTokenBceByClient(appId, appSecret)
	if err != nil {
		bd.logf("cannot get access token(%v): %v", appId, err.Error())
		return nil, err
	}

	detectURL, _ := url.Parse(`https://aip.baidubce.com/rest/2.0/solution/v1/img_censor/v2/user_defined`)
	detectQuery := detectURL.Query()
	detectQuery.Set("access_token", accessToken.AccessToken)
	detectURL.RawQuery = detectQuery.Encode()

	params := url.Values{
		"image":   {base64.StdEncoding.EncodeToString(data)},
		"imgType": {"0"},
	}
	bdRecognizeResultData, err := http.PostData(detectURL.String(), "application/x-www-form-urlencoded", []byte(params.Encode()))

	if err != nil {
		bd.logf("recognize picture err: %v", err)
		return nil, err
	}
	baiduAuditImageResult := &BaiduAuditResult{}
	err = json.Unmarshal(bdRecognizeResultData, baiduAuditImageResult)
	if err != nil {
		bd.logf("recognize picture err: %v", err)
		return nil, err
	}
	return baiduAuditImageResult, nil
	//if baiduAuditImageResult.ErrorCode != 0 || baiduAuditImageResult.ErrorMsg!= "" {
	//	bd.logf( "recognize picture err: %v %v", baiduAuditImageResult.ErrorCode, baiduAuditImageResult.ErrorMsg)
	//	return baiduAuditImageResult, errors.New(baiduAuditImageResult.ErrorMsg)
	//}
	//if baiduAuditImageResult.ConclusionType != 1 {
	//	bd.logf( "recognize picture result: %v %v", baiduAuditImageResult.ConclusionType, baiduAuditImageResult.Conclusion)
	//	conclusionMessage := make([]string, 0)
	//	for _, datum := range baiduPornRecognizeResult.Data {
	//		conclusionMessage =append(conclusionMessage, datum.Msg)
	//	}
	//	return ctx.JSON(200, &UploadImageResponse{
	//		ErrCode:    "image_error",
	//		ErrMessage: strings.Join(conclusionMessage, ","),
	//		Success:    false,
	//	})
	//}
}

type BaiduAuditResult struct {
	ErrorCode      int    `json:"error_code" xml:"error_code"`
	ErrorMsg       string `json:"error_msg" xml:"error_msg"`
	LogID          int64  `json:"log_id" xml:"log_id"`
	Conclusion     string `json:"conclusion" xml:"conclusion"`
	ConclusionType int    `json:"conclusionType" xml:"conclusionType"`
	Data           []struct {
		ConclusionType int      `json:"conclusionType,omitempty" xml:"conclusionType,omitempty"`
		Msg            string   `json:"msg" xml:"msg"`
		DatasetName    string   `json:"datasetName,omitempty" xml:"datasetName,omitempty"`
		Type           int      `json:"type" xml:"type"`
		SubType        int      `json:"subType" xml:"subType"`
		Conclusion     string   `json:"conclusion,omitempty" xml:"conclusion,omitempty"`
		Probability    float64  `json:"probability,omitempty" xml:"probability,omitempty"`
		Codes          []string `json:"codes,omitempty" xml:"codes,omitempty"`
		Stars          []struct {
			Probability float64 `json:"probability" xml:"probability"`
			Name        string  `json:"name" xml:"name"`
		} `json:"stars,omitempty" xml:"stars,omitempty"`
		Completeness float64 `json:"completeness,omitempty" xml:"completeness,omitempty"`
		Hits         []struct {
			Probability int      `json:"probability" xml:"probability"`
			DatasetName string   `json:"datasetName" xml:"datasetName"`
			Words       []string `json:"words" xml:"words"`
		} `json:"hits,omitempty" xml:"hits,omitempty"`
		Conclution     string `json:"conclution,omitempty" xml:"conclution,omitempty"`
		ConclutionType int    `json:"conclutionType,omitempty" xml:"conclutionType,omitempty"`
	} `json:"data" xml:"data"`
	RawData []struct {
		Type    int `json:"type" xml:"type"`
		Results []struct {
			Result []struct {
				Location struct {
					Top    int `json:"top" xml:"top"`
					Left   int `json:"left" xml:"left"`
					Width  int `json:"width" xml:"width"`
					Height int `json:"height" xml:"height"`
				} `json:"location" xml:"location"`
				Stars []struct {
					Probability float64 `json:"probability" xml:"probability"`
					Name        string  `json:"name" xml:"name"`
				} `json:"stars" xml:"stars"`
			} `json:"result" xml:"result"`
			LogID             int64  `json:"log_id" xml:"log_id"`
			IncludePolitician string `json:"include_politician" xml:"include_politician"`
			ResultNum         int    `json:"result_num" xml:"result_num"`
			ResultConfidence  string `json:"result_confidence" xml:"result_confidence"`
		} `json:"results" xml:"results"`
	} `json:"rawData" xml:"rawData"`
}
