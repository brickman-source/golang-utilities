package baidu

import (
	"encoding/base64"
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/json"
	"github.com/brickman-source/golang-utilities/log"
	"net/url"
	"strings"
)

type FaceMultiSearchRequest struct {
	Image           string `json:"image,omitempty"`
	ImageType       string `json:"image_type,omitempty"`
	GroupIDList     string `json:"group_id_list,omitempty"`
	MaxFaceNum      int    `json:"max_face_num,omitempty"`
	MatchThreshold  int    `json:"match_threshold,omitempty"`
	MaxUserNum      int    `json:"max_user_num,omitempty"`
	QualityControl  string `json:"quality_control,omitempty"`
	LivenessControl string `json:"liveness_control,omitempty"`
}

type FaceMultiSearchResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	LogID     int    `json:"log_id"`
	Timestamp int    `json:"timestamp"`
	Cached    int    `json:"cached"`
	Result    struct {
		FaceNum  int `json:"face_num"`
		FaceList []struct {
			FaceToken string `json:"face_token"`
			Location  struct {
				Left     float64 `json:"left"`
				Top      float64 `json:"top"`
				Width    int     `json:"width"`
				Height   int     `json:"height"`
				Rotation int     `json:"rotation"`
			} `json:"location"`
			UserList []struct {
				GroupID  string  `json:"group_id"`
				UserID   string  `json:"user_id"`
				UserInfo string  `json:"user_info"`
				Score    float64 `json:"score"`
			} `json:"user_list"`
		} `json:"face_list"`
	} `json:"result"`
}

func (bd *Baidu) FaceMultiSearch(
	userGroupIdList []string,
	imageData []byte,
	qualityControl FaceControlLevel,
	livenessControl FaceControlLevel,
	appId, appSecret string) (*FaceMultiSearchResponse, error) {
	accessToken, err := bd.GetAccessTokenByClient(appId, appSecret)
	if err != nil {
		log.Errorf("cannot get access token(%v): %v", appId, err.Error())
		return nil, err
	}
	bdReqURL, _ := url.Parse(`https://aip.baidubce.com/rest/2.0/face/v3/multi-search`)
	bdReqQuery := bdReqURL.Query()
	bdReqQuery.Set("access_token", accessToken.AccessToken)
	bdReqURL.RawQuery = bdReqQuery.Encode()

	req := &FaceMultiSearchRequest{
		Image:           base64.StdEncoding.EncodeToString(imageData),
		ImageType:       "BASE64",
		GroupIDList:     strings.Join(userGroupIdList, ","),
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
	bdResp := &FaceMultiSearchResponse{}
	err = json.Unmarshal(bdRespData, bdResp)
	if err != nil {
		log.Infof("bd err: %v", err)
		return nil, err
	}
	return bdResp, nil
}
