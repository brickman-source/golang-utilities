package wechat

import (
	"encoding/xml"
	"fmt"
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/brickman-source/golang-utilities/rand"
	"net/url"
	"time"
)

type Req struct {
	ActName      string `json:"act_name" xml:"act_name"`             //act_name        活动名称
	ClientIP     string `json:"client_ip" xml:"client_ip"`           //client_ip       IP地址
	ConsumeMchID string `json:"consume_mch_id" xml:"consume_mch_id"` //consume_mch_id  资金授权商户号
	MchBillno    string `json:"mch_billno" xml:"mch_billno"`         //mch_billno      商户订单号
	MchID        string `json:"mch_id" xml:"mch_id"`                 //mch_id          商户号
	NonceStr     string `json:"nonce_str" xml:"nonce_str"`           //nonce_str       随机字符串
	ReOpenid     string `json:"re_openid" xml:"re_openid"`           //re_openid       用户openid
	Remark       string `json:"remark" xml:"remark"`                 //remark          备注
	RiskInfo     string `json:"risk_info" xml:"risk_info"`           //risk_info       活动信息
	SceneID      string `json:"scene_id" xml:"scene_id"`             //scene_id        场景id
	SendName     string `json:"send_name" xml:"send_name"`           //send_name       商户名称
	Sign         string `json:"sign" xml:"sign"`                     //sign            签名
	TotalAmount  string `json:"total_amount" xml:"total_amount"`     //total_amount    付款金额
	TotalNum     string `json:"total_num" xml:"total_num"`           //total_num       红包发放总人数
	Wishing      string `json:"wishing" xml:"wishing"`               //wishing         红包祝福语
	Wxappid      string `json:"wxappid" xml:"wxappid"`               //wxappid         公众账号appid
}
type Resp struct {
	ErrCode     string `json:"err_code" xml:"err_code"`         //err_code 错误代码
	ErrCodeDes  string `json:"err_code_des" xml:"err_code_des"` //err_code_des 错误代码描述
	MchBillno   string `json:"mch_billno" xml:"mch_billno"`     //mch_billno 商户订单号
	MchID       string `json:"mch_id" xml:"mch_id"`             //mch_id 商户号
	ReOpenid    string `json:"re_openid" xml:"re_openid"`       //re_openid 用户openid
	ResultCode  string `json:"result_code" xml:"result_code"`   //result_code 业务结果
	ReturnCode  string `json:"return_code" xml:"return_code"`   //return_code 返回状态码
	ReturnMsg   string `json:"return_msg" xml:"return_msg"`     //return_msg 返回信息
	SendListid  string `json:"send_listid" xml:"send_listid"`   //send_listid 微信单号
	Sign        string `json:"sign" xml:"sign"`                 //sign 签名
	TotalAmount string `json:"total_amount" xml:"total_amount"` //total_amount 付款金额
	Wxappid     string `json:"wxappid" xml:"wxappid"`           //wxappid 公众账号appid
}

func (wx *Wechat) SendRedPack(extId, actName, remark, wishing, openId string, amount int) (*Resp, error) {
	requestURL := `https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack`

	req := &Req{
		ActName:      actName, //"充值送微信红包",
		ClientIP:     wx.config.GetString("wechat.payment.clientIp"),
		ConsumeMchID: "",
		MchBillno:    extId,
		MchID:        wx.config.GetString("wechat.payment.mchId"),
		NonceStr:     rand.GetNonceString(16),
		ReOpenid:     openId,
		Remark:       remark, //"充值红包",
		RiskInfo:     url.QueryEscape(fmt.Sprintf("posttime=%d&clientversion=1", time.Now().Unix())),
		SceneID:      "PRODUCT_1", //PRODUCT_1:商品促销 1至499元
		SendName:     wx.config.GetString("wechat.payment.senderName"),
		TotalAmount:  fmt.Sprintf("%d", amount),
		TotalNum:     "1",
		Wishing:      wishing, //"感谢您参加充值活动，祝您开心！",
		Wxappid:      wx.config.GetString("wechat.appId"),
	}

	signMd5, _ := SignXmlMd5(*req, wx.config.GetString("wechat.payment.appKey"))
	req.Sign = signMd5

	output, err := xml.Marshal(req)

	if err != nil {
		return nil, err
	}

	resp := &Resp{}

	data, err := http.PostDataSsl(requestURL, output, []byte(wx.config.GetString("wechat.payment.certStr")), []byte(wx.config.GetString("wechat.payment.keyStr")))

	if err != nil {
		log.Errorf(nil, "[SendRedPack]PostDataSsl err:%s", err.Error())
		return nil, err
	}
	log.Debugf(nil, string(data))
	err = xml.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
