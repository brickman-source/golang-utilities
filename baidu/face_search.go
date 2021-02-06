package baidu

import (
	"encoding/base64"
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/json"
	"github.com/brickman-source/golang-utilities/log"
	"net/url"
	"strings"
)

type FaceSearchRequest struct {
	Image           string `json:"image"`
	ImageType       string `json:"image_type"`
	GroupIDList     string `json:"group_id_list"`
	QualityControl  string `json:"quality_control"`
	LivenessControl string `json:"liveness_control"`
}

type FaceSearchResponse struct {
	ErrorCode      int    `json:"error_code" xml:"error_code"`
	ErrorMsg       string `json:"error_msg" xml:"error_msg"`
	FaceToken string `json:"face_token"`
	UserList  []struct {
		GroupID  string  `json:"group_id"`
		UserID   string  `json:"user_id"`
		UserInfo string  `json:"user_info"`
		Score    float64 `json:"score"`
	} `json:"user_list"`
}

func (bd *Baidu) FaceSearch(userGroupIdList []string, imageData []byte, appId, appSecret string) (*FaceSearchResponse, error) {
	accessToken, err := bd.GetAccessTokenByClient(appId, appSecret)
	if err != nil {
		log.Errorf("cannot get access token(%v): %v", appId, err.Error())
		return nil, err
	}
	bdReqURL, _ := url.Parse(`https://aip.baidubce.com/rest/2.0/face/v3/search`)
	bdReqQuery := bdReqURL.Query()
	bdReqQuery.Set("access_token", accessToken.AccessToken)
	bdReqURL.RawQuery = bdReqQuery.Encode()

	req := &FaceSearchRequest{
		Image:       base64.StdEncoding.EncodeToString(imageData),
		ImageType:   "BASE64",
		GroupIDList: strings.Join(userGroupIdList, ","),
	}

	bdRespData, err := http.PostData(bdReqURL.String(),
		http.MIMEApplicationJSONCharsetUTF8,
		json.ShouldMarshal(req),
	)

	if err != nil {
		log.Infof("bd err: %v", err)
		return nil, err
	}
	bdResp := &FaceSearchResponse{}
	err = json.Unmarshal(bdRespData, bdResp)
	if err != nil {
		log.Infof("bd err: %v", err)
		return nil, err
	}
	return bdResp, nil
}
