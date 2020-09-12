package aliapi

import (
	"encoding/base64"
	"encoding/json"
	"github.com/brickman-source/golang-utilities/http"
	gohttp "net/http"
	"time"
)

const (
	OCR_IDCARD_API_URL = `https://dm-51.data.aliyun.com/rest/160601/ocr/ocr_idcard.json`
)

func (api *AliApiHelper) OCRIDCardFace(data []byte) (*IDCardFace, error) {
	respData, err := api.ocrIDCard(data, "face")
	if err != nil {
		return nil, err
	}
	ret := &IDCardFace{}
	err = json.Unmarshal(respData, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (api *AliApiHelper) OCRIDCardBack(data []byte) (*IDCardBack, error) {
	respData, err := api.ocrIDCard(data, "back")
	if err != nil {
		return nil, err
	}
	ret := &IDCardBack{}
	err = json.Unmarshal(respData, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (api *AliApiHelper) ocrIDCard(data []byte, side string) (ret []byte, err error) {
	bt := base64.StdEncoding.EncodeToString(data)
	confBody, err := json.Marshal(&struct {
		Side string `json:"side"`
	}{
		Side: side,
	})
	jsonBody, err := json.Marshal(&struct {
		Image     string `json:"image"`
		Configure string `json:"configure"`
	}{
		Image:     bt,
		Configure: string(confBody),
	})
	println(string(jsonBody))
	_, respData, err := http.PostDataWithHeaders(
		OCR_IDCARD_API_URL,
		gohttp.Header{
			"Authorization": []string{"APPCODE " + api.cfg.GetString("aliapi.ocr.idcard.appcode")},
		},
		"application/json",
		jsonBody,
		time.Minute,
	)

	//fmt.Printf("%v\n", respHeader)
	//fmt.Println(string(respData))

	ret = respData

	return
}

type IDCardFace struct {
	Address     string `json:"address"`
	Birth       string `json:"birth"`
	ConfigStr   string `json:"config_str"`
	Name        string `json:"name"`
	Nationality string `json:"nationality"`
	Num         string `json:"num"`
	RequestID   string `json:"request_id"`
	Sex         string `json:"sex"`
	Success     bool   `json:"success"`
	FaceRect    struct {
		Angle  float64 `json:"angle"`
		Center struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"center"`
		Size struct {
			Width  float64 `json:"width"`
			Height float64 `json:"height"`
		} `json:"size"`
	} `json:"face_rect"`
	FaceRectVertices []struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"face_rect_vertices"`
}

type IDCardBack struct {
	ConfigStr string `json:"config_str"`
	EndDate   string `json:"end_date"`
	Issue     string `json:"issue"`
	RequestID string `json:"request_id"`
	StartDate string `json:"start_date"`
	Success   bool   `json:"success"`
}
