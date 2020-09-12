package alidayu

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"github.com/brickman-source/golang-utilities/config"
	"github.com/brickman-source/golang-utilities/errors"
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/brickman-source/golang-utilities/rand"
	"net/url"
	"strings"
	"time"
)

const (
	HMACSHA1 = "HMAC-SHA1"
)

// SmsReq ...
type SmsReq struct {
	//AccessKeyId	是否必填：是 说明：AccessKeyId
	AccessKeyId string `form:"AccessKeyId"`
	//Action 是否必填：是	说明：API的命名，固定值，如发送短信API的值为：SendSms
	Action string `form:"Action"`
	//Format	是否必填：否	说明：没传默认为JSON，可选填值：XML
	Format string `form:"Format"`
	//OutId 是否必填：否	说明：外部流水扩展字段
	OutId string `form:"OutId"`
	//PhoneNumbers 是否必填：是	说明：短信接收号码,支持以逗号分隔的形式进行批量调用，批量上限为1000个手机号码,批量调用相对于单条调用及时性稍有延迟,验证码类型的短信推荐使用单条调用的方式
	PhoneNumbers string `form:"PhoneNumbers"`
	//RegionId 是否必填：是	说明：API支持的RegionID，如短信API的值为：cn-hangzhou
	RegionId string `form:"RegionId"`
	//SignName 是否必填：是	说明：短信签名(如 云通信)
	SignName string `form:"SignName"`
	//Signature	是否必填：是	说明：最终生成的签名结果值
	Signature string `form:"Signature"`
	//SignatureMethod	是否必填：是	说明：建议固定值：HMAC-SHA1
	SignatureMethod string `form:"SignatureMethod"`
	//SignatureNonce	是否必填：是	说明：用于请求的防重放攻击，每次请求唯一，JAVA语言建议用：java.util.UUID.randomUUID()生成即可
	SignatureNonce string `form:"SignatureNonce"`
	//SignatureVersion	是否必填：是	说明：建议固定值：1.0
	SignatureVersion string `form:"SignatureVersion"`
	//TemplateCode 是否必填：是	说明：短信模板ID
	TemplateCode string `form:"TemplateCode"`
	//TemplateParam 是否必填：否	说明：短信模板变量替换JSON串,友情提示:如果JSON中需要带换行符,请参照标准的JSON协议。（如{“code”:”1234”,”product”:”ytx”} ）
	TemplateParam string `form:"TemplateParam"`
	//Timestamp	是否必填：是	说明：格式为：yyyy-MM-dd’T’HH:mm:ss’Z’；时区为：GMT
	Timestamp string `form:"Timestamp"`
	//Version 是否必填：是	说明：API的版本，固定值，如短信API的值为：2017-05-25
	Version string `form:"Version"`
}

type SmsHelper struct {
	cfg *config.Config
}

func NewSmsHelper(cfg *config.Config) *SmsHelper {
	return &SmsHelper{cfg: cfg}
}

// DefaultSmsReq ...
func (ss *SmsHelper) getSmsReq() *SmsReq {
	timeStamp := time.Now().In(time.UTC).Format("2006-01-02T15:04:05Z")

	return &SmsReq{
		AccessKeyId:      ss.cfg.GetString("alidayu.accessKeyID"),
		Action:           "SendSms",
		Format:           "JSON",
		OutId:            "",
		PhoneNumbers:     "",
		RegionId:         "cn-hangzhou",
		SignName:         "",
		Signature:        "",
		SignatureMethod:  "HMAC-SHA1",
		SignatureNonce:   rand.GetNonceString(32),
		SignatureVersion: "1.0",
		TemplateCode:     "",
		TemplateParam:    "",
		Timestamp:        timeStamp,
		Version:          "2017-05-25",
	}
}

func (ss *SmsHelper) req2Query(req *SmsReq) url.Values {
	query := url.Values{}
	query.Set("AccessKeyId", req.AccessKeyId)
	query.Set("Action", req.Action)
	query.Set("Format", req.Format)
	if req.OutId != "" {
		query.Set("OutId", req.OutId)
	}
	query.Set("PhoneNumbers", req.PhoneNumbers)
	query.Set("RegionId", req.RegionId)
	query.Set("SignName", req.SignName)
	query.Set("Signature", req.Signature)
	query.Set("SignatureMethod", req.SignatureMethod)
	query.Set("SignatureNonce", req.SignatureNonce)
	query.Set("SignatureVersion", req.SignatureVersion)
	query.Set("TemplateCode", req.TemplateCode)
	query.Set("TemplateParam", req.TemplateParam)
	query.Set("Timestamp", req.Timestamp)
	query.Set("Version", req.Version)
	return query
}

func (ss *SmsHelper) sign(req *SmsReq) string {
	query := ss.req2Query(req)
	query.Del("Signature")

	orig := query.Encode()
	orig = strings.Replace(orig, "+", "%20", -1)
	orig = strings.Replace(orig, "*", "%2A", -1)
	orig = strings.Replace(orig, "%7E", "~", -1)

	orig = "GET&" + url.QueryEscape("/") + "&" + url.QueryEscape(orig)

	req.Signature = encrypt([]byte(orig), []byte(ss.cfg.GetString("alidayu.accessSecret")+"&"), req.SignatureMethod)
	return req.Signature
}

func encrypt(s, secret []byte, method string) (h string) {
	switch method {
	case HMACSHA1:
		mac := hmac.New(sha1.New, secret)
		mac.Write(s)
		sum := mac.Sum(nil)
		h = base64.StdEncoding.EncodeToString(sum)
	default:
		panic("unsupported sign method!")
	}
	return h
}

type (
	SmsResp struct {
		RequestId string `json:"RequestId"`
		Code      string `json:"Code"`
		Message   string `json:"Message"`
		BizId     string `json:"BizId"`
	}
)

func (ss *SmsHelper) SendSms(phoneNumber string, signName, templateCode, templateParam, outId string) (*SmsResp, error) {
	smsReq := ss.getSmsReq()
	smsResp := new(SmsResp)

	smsReq.SignName = signName
	smsReq.PhoneNumbers = phoneNumber
	smsReq.TemplateCode = templateCode
	smsReq.TemplateParam = templateParam
	smsReq.OutId = outId
	ss.sign(smsReq)

	query := ss.req2Query(smsReq)

	sendSmsUrl, _ := url.Parse(ss.cfg.GetString("alidayu.sendSmsUrl"))
	sendSmsUrl.RawQuery = query.Encode()

	respData, err := http.GetData(sendSmsUrl.String())

	err = json.Unmarshal(respData, smsResp)
	if err != nil {
		return nil, err
	}

	if smsResp.Code != "OK" {
		log.Debugf(nil, "[SmsHelper.SendSms]send sms resp:%s", smsResp.Code)
		return smsResp, errors.WithCode(nil, smsResp.Code, smsResp.Message)
	}
	log.Debugf(nil, "[SmsHelper.SendSms]send sms resp 2:%s", smsResp.Code)

	return smsResp, nil
}
