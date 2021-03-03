package baidu

import (
	"encoding/base64"
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/json"
	"github.com/brickman-source/golang-utilities/log"
	"net/url"
)

type PersonVerifyRequest struct {
	Image           string `json:"image"`
	ImageType       string `json:"image_type"`
	IDCardNumber    string `json:"id_card_number"`
	Name            string `json:"name"`
	QualityControl  string `json:"quality_control"`
	LivenessControl string `json:"liveness_control"`
}

type PersonVerifyResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	LogID     int    `json:"log_id"`
	Timestamp int    `json:"timestamp"`
	Cached    int    `json:"cached"`
	Result    struct {
		Score float64 `json:"score"`
	} `json:"result"`
}
func (bd *Baidu) PersonVerify(
	imageData []byte,
	fullName string,
	idCardNo string,
	qualityControl FaceControlLevel,
	livenessControl FaceControlLevel,
	appId, appSecret string) (*PersonVerifyResponse, error) {
	accessToken, err := bd.GetAccessTokenByClient(appId, appSecret)
	if err != nil {
		log.Errorf("cannot get access token(%v): %v", appId, err.Error())
		return nil, err
	}
	bdReqURL, _ := url.Parse(`https://aip.baidubce.com/rest/2.0/face/v3/person/verify`)
	bdReqQuery := bdReqURL.Query()
	bdReqQuery.Set("access_token", accessToken.AccessToken)
	bdReqURL.RawQuery = bdReqQuery.Encode()

	req := &PersonVerifyRequest{
		Image:           base64.StdEncoding.EncodeToString(imageData),
		ImageType:       "BASE64",
		IDCardNumber:    idCardNo,
		Name:            fullName,
		QualityControl:  string(qualityControl),
		LivenessControl: string(livenessControl),
	}

	bdRespData, err := http.PostData(bdReqURL.String(),
		http.MIMEApplicationJSONCharsetUTF8,
		json.ShouldMarshal(req),
	)

	if err != nil {
		log.Infof("bd err: %v", err)
		return nil, err
	}
	bdResp := &PersonVerifyResponse{}
	err = json.Unmarshal(bdRespData, bdResp)
	if err != nil {
		log.Infof("bd err: %v", err)
		return nil, err
	}
	return bdResp, nil
}
