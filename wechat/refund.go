package wechat

import (
	"encoding/xml"
	"fmt"
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/brickman-source/golang-utilities/rand"
)

var (
	WxRefundURL = "https://api.mch.weixin.qq.com/secapi/pay/refund"
)

type RefundRequest struct {
	XMLName       xml.Name `xml:"xml"`
	AppId         string   `xml:"appid"`                     //must
	DeviceInfo    string   `xml:"device_info,omitempty"`     //maybe
	MchId         string   `xml:"mch_id"`                    //must
	NonceStr      string   `xml:"nonce_str"`                 //must
	OpUserId      string   `xml:"op_user_id"`                //maybe
	OutRefundNo   string   `xml:"out_refund_no"`             //must
	OutTradeNo    string   `xml:"out_trade_no"`              //must
	RefundFee     string   `xml:"refund_fee"`                //must
	RefundFeeType string   `xml:"refund_fee_type,omitempty"` //maybe
	Sign          string   `xml:"sign"`                      //must
	TotalFee      string   `xml:"total_fee"`                 //must
	TransactionId string   `xml:"transaction_id"`            //maybe
	RequestXML    string   `xml:"-"`                         //最终请求串
}

type RefundReponse struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`

	//以下字段在return_code为SUCCESS的时候有返回
	AppId             string `xml:"appid"`                          //must
	MchId             string `xml:"mch_id"`                         //must
	DeviceInfo        string `xml:"device_info,omitempty"`          //maybe
	NonceStr          string `xml:"nonce_str"`                      //must
	Sign              string `xml:"sign"`                           //must
	TransactionId     string `xml:"transaction_id"`                 //must
	OutTradeNo        string `xml:"out_trade_no"`                   //must
	OutRefundNo       string `xml:"out_refund_no"`                  //must
	RefundId          string `xml:"refund_id"`                      //must
	RefundFee         string `xml:"refund_fee"`                     //must
	RefundFeeType     string `xml:"refund_fee_type,omitempty"`      //must
	TotalFee          string `xml:"total_fee"`                      //must
	FeeType           string `xml:"fee_type,omitempty"`             //must
	CashFee           string `xml:"cash_fee"`                       //must
	CashFeeType       string `xml:"cash_fee_type,omitempty"`        //must
	CashRefundFee     string `xml:"cash_refund_fee,omitempty"`      //must
	CashRefundFeeType string `xml:"cash_refund_fee_type,omitempty"` //must

	// Coupon_refund_fee string `xml:"coupon_refund_fee,omitempty"` //must
	// Coupon_count string `xml:"coupon_count,omitempty"` //must
	// 	coupon_batch_id_$n  string `xml:"coupon_batch_id_$n,omitempty"`            //must
	// coupon_id_$n  string `xml:"coupon_id_$n,omitempty"`            //must
	// coupon_fee_$n  string `xml:"coupon_fee_$n,omitempty"`            //must

	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code,omitempty"`
	ErrCodeDes string `xml:"err_code_des,omitempty"`

	XML string `xml:"-"` //结果串
}

func (v *RefundRequest) Xml() error {
	output, err := xml.Marshal(v)
	xmlResult := string(output)
	v.RequestXML = xmlResult

	return err
}

func (wx *Wechat) Refund(orderId string, openId string, totalFee int) error {
	req := &RefundRequest{
		AppId: wx.config.GetString("wechat.offiaccount.appId"),
		MchId: wx.config.GetString("wechat.payment.mchId"),
	}
	req.NonceStr = rand.GetNonceString(32)
	req.OutTradeNo = orderId
	req.OutRefundNo = fmt.Sprintf("refund_%s", orderId)
	req.TotalFee = fmt.Sprintf("%d", totalFee)
	req.RefundFee = fmt.Sprintf("%d", totalFee)
	req.OpUserId = openId

	signMd5, _ := SignXmlMd5(*req, wx.config.GetString("wechat.payment.appKey"))
	req.Sign = signMd5

	reqData, err := xml.Marshal(req)
	if err != nil {
		return err
	}

	response := RefundReponse{}

	log.Infof( "refund req:%v", WxRefundURL)

	data, err := http.PostDataSsl(WxRefundURL, reqData, []byte(wx.config.GetString("wechat.payment.certStr")), []byte(wx.config.GetString("wechat.payment.keyStr")))

	log.Infof( "refund resp:%v", string(data))

	err = xml.Unmarshal(data, &response)

	if err != nil {
		log.Errorf( "refund err:%s", err.Error())
		response.XML = string(data)
	}

	return nil
}
