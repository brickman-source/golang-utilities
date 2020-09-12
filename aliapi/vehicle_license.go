package aliapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/brickman-source/golang-utilities/http"
	gohttp "net/http"
	"time"
)

const (
	OCR_VEHICLE_LICENSE_API_URL = `https://dm-53.data.aliyun.com/rest/160601/ocr/ocr_vehicle.json`
)

func (api *AliApiHelper) OCRVehicleLicenseFace(data []byte) (*VehicleLicenseFace, error) {
	respData, err := api.ocrVehicleLicense(data, "face")
	if err != nil {
		return nil, err
	}
	ret := &VehicleLicenseFace{}
	err = json.Unmarshal(respData, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
func (api *AliApiHelper) OCRVehicleLicenseBack(data []byte) (*VehicleLicenseBack, error) {
	respData, err := api.ocrVehicleLicense(data, "back")
	if err != nil {
		return nil, err
	}
	ret := &VehicleLicenseBack{}
	err = json.Unmarshal(respData, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (api *AliApiHelper) ocrVehicleLicense(data []byte, side string) (ret []byte, err error) {
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
	respHeader, respData, err := http.PostDataWithHeaders(
		OCR_VEHICLE_LICENSE_API_URL,
		gohttp.Header{
			"Authorization": []string{"APPCODE " + api.cfg.GetString("aliapi.ocr.idcard.appcode")},
		},
		"application/json",
		jsonBody,
		time.Minute,
	)

	fmt.Printf("%v\n", respHeader)
	fmt.Println(string(respData))

	ret = respData

	return
}

type VehicleLicenseFace struct {
	Addr         string `json:"addr" xml:"addr"`
	ConfigStr    string `json:"config_str" xml:"config_str"`
	EngineNum    string `json:"engine_num" xml:"engine_num"`
	IssueDate    string `json:"issue_date" xml:"issue_date"`
	Model        string `json:"model" xml:"model"`
	Owner        string `json:"owner" xml:"owner"`
	PlateNum     string `json:"plate_num" xml:"plate_num"`
	RegisterDate string `json:"register_date" xml:"register_date"`
	RequestID    string `json:"request_id" xml:"request_id"`
	Success      bool   `json:"success" xml:"success"`
	UseCharacter string `json:"use_character" xml:"use_character"`
	VehicleType  string `json:"vehicle_type" xml:"vehicle_type"`
	Vin          string `json:"vin" xml:"vin"`
}
type VehicleLicenseBack struct {
	AppprovedPassengerCapacity string `json:"appproved_passenger_capacity" xml:"appproved_passenger_capacity"`
	ApprovedLoad               string `json:"approved_load" xml:"approved_load"`
	ConfigStr                  string `json:"config_str" xml:"config_str"`
	FileNo                     string `json:"file_no" xml:"file_no"`
	GrossMass                  string `json:"gross_mass" xml:"gross_mass"`
	InspectionRecord           string `json:"inspection_record" xml:"inspection_record"`
	OverallDimension           string `json:"overall_dimension" xml:"overall_dimension"`
	PlateNum                   string `json:"plate_num" xml:"plate_num"`
	RequestID                  string `json:"request_id" xml:"request_id"`
	Success                    bool   `json:"success" xml:"success"`
	TractionMass               string `json:"traction_mass" xml:"traction_mass"`
	UnladenMass                string `json:"unladen_mass" xml:"unladen_mass"`
}
