package baidu

import (
	"encoding/base64"
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/json"
	"github.com/brickman-source/golang-utilities/log"
	"net/url"
)

type FaceSetUserAddRequest struct {
	Image           string `json:"image,omitempty"`
	ImageType       string `json:"image_type,omitempty"`
	GroupID         string `json:"group_id,omitempty"`
	UserID          string `json:"user_id,omitempty"`
	UserInfo        string `json:"user_info,omitempty"`
	QualityControl  string `json:"quality_control,omitempty"`
	LivenessControl string `json:"liveness_control,omitempty"`
}
type FaceSetUserAddResponse struct {
	ErrorCode int    `json:"error_code" xml:"error_code"`
	ErrorMsg  string `json:"error_msg" xml:"error_msg"`
	LogID     uint64 `json:"log_id"`
	FaceToken string `json:"face_token"`
	Location  struct {
		Left     int `json:"left"`
		Top      int `json:"top"`
		Width    int `json:"width"`
		Height   int `json:"height"`
		Rotation int `json:"rotation"`
	} `json:"location"`
}

func (bd *Baidu) FaceSetUserAdd(userGroupId, userId string, userInfo string, imageData []byte, appId, appSecret string) (*FaceSetUserAddResponse, error) {
	accessToken, err := bd.GetAccessTokenByClient(appId, appSecret)
	if err != nil {
		log.Errorf("cannot get access token(%v): %v", appId, err.Error())
		return nil, err
	}
	bdReqURL, _ := url.Parse(`https://aip.baidubce.com/rest/2.0/face/v3/faceset/user/add`)
	bdReqQuery := bdReqURL.Query()
	bdReqQuery.Set("access_token", accessToken.AccessToken)
	bdReqURL.RawQuery = bdReqQuery.Encode()

	req := &FaceSetUserAddRequest{
		Image:           base64.StdEncoding.EncodeToString(imageData),
		ImageType:       "BASE64",
		GroupID:         userGroupId,
		UserID:          userId,
		UserInfo:        userInfo,
		QualityControl:  "NORMAL",
		LivenessControl: "NORMAL",
	}

	bdRespData, err := http.PostData(bdReqURL.String(),
		http.MIMEApplicationJSONCharsetUTF8,
		json.ShouldMarshal(req),
	)

	if err != nil {
		log.Infof("bd err: %v", err)
		return nil, err
	}
	bdResp := &FaceSetUserAddResponse{}
	err = json.Unmarshal(bdRespData, bdResp)
	if err != nil {
		log.Infof("bd err: %v", err)
		return nil, err
	}
	return bdResp, nil
}
