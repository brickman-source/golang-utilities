package baidu

import (
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/json"
	"net/url"
)

//
//type IdCardSide string
//
//const (
//	IdCardSide_Front IdCardSide = "front"
//	IdCardSide_Back  IdCardSide = "back"
//)
//
//type IdCardRecognizeRequest struct {
//	Image string `json:"image"`
//	Side  string `json:"side"`
//}
//
//type IdCardRecognizeResponse struct {
//	LogID         int64  `json:"log_id"`
//	Direction     int    `json:"direction"`
//	ImageStatus   string `json:"image_status"`
//	Photo         string `json:"photo"`
//	PhotoLocation struct {
//		Width  int `json:"width"`
//		Top    int `json:"top"`
//		Left   int `json:"left"`
//		Height int `json:"height"`
//	} `json:"photo_location"`
//	WordsResult struct {
//		Address struct {
//			Location struct {
//				Left   int `json:"left"`
//				Top    int `json:"top"`
//				Width  int `json:"width"`
//				Height int `json:"height"`
//			} `json:"location"`
//			Words string `json:"words"`
//		} `json:"住址"`
//		IdCard struct {
//			Location struct {
//				Left   int `json:"left"`
//				Top    int `json:"top"`
//				Width  int `json:"width"`
//				Height int `json:"height"`
//			} `json:"location"`
//			Words string `json:"words"`
//		} `json:"公民身份号码"`
//		Birth struct {
//			Location struct {
//				Left   int `json:"left"`
//				Top    int `json:"top"`
//				Width  int `json:"width"`
//				Height int `json:"height"`
//			} `json:"location"`
//			Words string `json:"words"`
//		} `json:"出生"`
//		Name struct {
//			Location struct {
//				Left   int `json:"left"`
//				Top    int `json:"top"`
//				Width  int `json:"width"`
//				Height int `json:"height"`
//			} `json:"location"`
//			Words string `json:"words"`
//		} `json:"姓名"`
//		Sex struct {
//			Location struct {
//				Left   int `json:"left"`
//				Top    int `json:"top"`
//				Width  int `json:"width"`
//				Height int `json:"height"`
//			} `json:"location"`
//			Words string `json:"words"`
//		} `json:"性别"`
//		Nationality struct {
//			Location struct {
//				Left   int `json:"left"`
//				Top    int `json:"top"`
//				Width  int `json:"width"`
//				Height int `json:"height"`
//			} `json:"location"`
//			Words string `json:"words"`
//		} `json:"民族"`
//	} `json:"words_result"`
//	WordsResultNum int `json:"words_result_num"`
//}

type VoiceRecognizeResponse struct {
	CorpusNo string   `json:"corpus_no"`
	ErrMsg   string   `json:"err_msg"`
	ErrNo    int      `json:"err_no"`
	Result   []string `json:"result"`
	Sn       string   `json:"sn"`
}

func (bd *Baidu) VoiceRecognizePro(
	voiceData []byte,
	appId, appSecret string) (*VoiceRecognizeResponse, error) {
	accessToken, err := bd.GetAccessTokenOpenApiByClient(appId, appSecret)
	if err != nil {
		bd.logf("cannot get access token(%v): %v", appId, err.Error())
		return nil, err
	}

	bdReqURL, _ := url.Parse(`https://vop.baidu.com/pro_api`)
	bdReqQuery := bdReqURL.Query()
	bdReqQuery.Set("token", accessToken.AccessToken)
	bdReqQuery.Set("dev_pid", "80001")
	bdReqQuery.Set("cuid", "fd342f603890941ea5416a9508c75f8cd437b54d")

	bdReqURL.RawQuery = bdReqQuery.Encode()

	bdRespData, err := http.PostData(bdReqURL.String(),
		`audio/pcm;rate=16000`,
		voiceData,
	)

	if err != nil {
		bd.logf("bd err: %v", err)
		return nil, err
	}
	bdResp := &VoiceRecognizeResponse{}
	err = json.Unmarshal(bdRespData, bdResp)
	if err != nil {
		bd.logf("bd err: %v", err)
		return nil, err
	}
	return bdResp, nil
}
