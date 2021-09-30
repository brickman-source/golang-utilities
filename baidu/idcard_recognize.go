package baidu

import (
	"encoding/base64"
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/json"
	"net/url"
)

type IdCardSide string

const (
	IdCardSide_Front IdCardSide = "front"
	IdCardSide_Back  IdCardSide = "back"
)

type IdCardRecognizeRequest struct {
	Image string `json:"image"`
	Side  string `json:"side"`
}

type IdCardRecognizeResponse struct {
	LogID         int64  `json:"log_id"`
	Direction     int    `json:"direction"`
	ImageStatus   string `json:"image_status"`
	Photo         string `json:"photo"`
	PhotoLocation struct {
		Width  int `json:"width"`
		Top    int `json:"top"`
		Left   int `json:"left"`
		Height int `json:"height"`
	} `json:"photo_location"`
	WordsResult struct {
		Address struct {
			Location struct {
				Left   int `json:"left"`
				Top    int `json:"top"`
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"location"`
			Words string `json:"words"`
		} `json:"住址"`
		IdCard struct {
			Location struct {
				Left   int `json:"left"`
				Top    int `json:"top"`
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"location"`
			Words string `json:"words"`
		} `json:"公民身份号码"`
		Birth struct {
			Location struct {
				Left   int `json:"left"`
				Top    int `json:"top"`
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"location"`
			Words string `json:"words"`
		} `json:"出生"`
		Name struct {
			Location struct {
				Left   int `json:"left"`
				Top    int `json:"top"`
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"location"`
			Words string `json:"words"`
		} `json:"姓名"`
		Sex struct {
			Location struct {
				Left   int `json:"left"`
				Top    int `json:"top"`
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"location"`
			Words string `json:"words"`
		} `json:"性别"`
		Nationality struct {
			Location struct {
				Left   int `json:"left"`
				Top    int `json:"top"`
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"location"`
			Words string `json:"words"`
		} `json:"民族"`
	} `json:"words_result"`
	WordsResultNum int `json:"words_result_num"`
}

func (bd *Baidu) IdCardRecognize(
	imageData []byte,
	side IdCardSide,
	detectRisk bool,
	appId, appSecret string) (*IdCardRecognizeResponse, error) {
	accessToken, err := bd.GetAccessTokenBceByClient(appId, appSecret)
	if err != nil {
		bd.logf("cannot get access token(%v): %v", appId, err.Error())
		return nil, err
	}
	bdReqURL, _ := url.Parse(`https://aip.baidubce.com/rest/2.0/ocr/v1/idcard`)
	bdReqQuery := bdReqURL.Query()
	bdReqQuery.Set("access_token", accessToken.AccessToken)
	bdReqURL.RawQuery = bdReqQuery.Encode()

	query := url.Values{
		"id_card_side": []string{string(side)},
		"image":        []string{base64.StdEncoding.EncodeToString(imageData)},
	}
	if detectRisk {
		query.Add("detect_risk", "true")
	}
	bdRespData, err := http.PostData(bdReqURL.String(),
		http.MIMEApplicationForm,
		[]byte(query.Encode()),
	)

	if err != nil {
		bd.logf("bd err: %v", err)
		return nil, err
	}
	bdResp := &IdCardRecognizeResponse{}
	err = json.Unmarshal(bdRespData, bdResp)
	if err != nil {
		bd.logf("bd err: %v", err)
		return nil, err
	}
	return bdResp, nil
}
