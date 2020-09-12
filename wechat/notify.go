package wechat

import (
	"encoding/xml"

	xmlHelper "github.com/brickman-source/golang-utilities/xml"
)

type NotifyData struct {
	XMLName xml.Name `xml:"xml"`
	// Coupon_batch_id_$n   string `xml:"coupon_batch_id_$n,omitempty"`
	// Coupon_count   string `xml:"coupon_count,omitempty"`
	// Coupon_fee   string `xml:"coupon_fee,omitempty"`
	// Coupon_fee_$n   string `xml:"coupon_fee_$n,omitempty"`
	// Coupon_id_$n   string `xml:"coupon_id_$n,omitempty"`
	AppId         xmlHelper.CData `xml:"appid"` //must
	Attach        xmlHelper.CData `xml:"attach,omitempty"`
	BankType      xmlHelper.CData `xml:"bank_type"`
	CashFee       xmlHelper.CData `xml:"cash_fee"`
	CashFeeType   xmlHelper.CData `xml:"cash_fee_type,omitempty"`
	DeviceInfo    xmlHelper.CData `xml:"device_info,omitempty"` //maybe
	ErrCode       xmlHelper.CData `xml:"err_code,omitempty"`
	ErrCodeDes    xmlHelper.CData `xml:"err_code_des,omitempty"`
	FeeType       xmlHelper.CData `xml:"fee_type,omitempty"`
	IsSubscribe   xmlHelper.CData `xml:"is_subscribe"`
	MchId         xmlHelper.CData `xml:"mch_id"`    //must
	NonceStr      xmlHelper.CData `xml:"nonce_str"` //must
	OpenId        xmlHelper.CData `xml:"openid"`
	OutTradeNo    xmlHelper.CData `xml:"out_trade_no"`
	ResultCode    xmlHelper.CData `xml:"result_code"`
	ReturnCode    xmlHelper.CData `xml:"return_code"`
	ReturnMsg     xmlHelper.CData `xml:"return_msg"`
	Sign          xmlHelper.CData `xml:"sign"`     //must
	TimeEnd       xmlHelper.CData `xml:"time_end"` //支付完成时间
	TotalFee      xmlHelper.CData `xml:"total_fee"`
	TradeType     xmlHelper.CData `xml:"trade_type"`
	TransactionId xmlHelper.CData `xml:"transaction_id"`

	XML string `xml:"-"` //结果串
}

type NotifyResponse struct {
	XMLName    xml.Name     `xml:"xml"`
	ReturnCode xmlHelper.CData `xml:"return_code"`
	ReturnMsg  xmlHelper.CData `xml:"return_msg"`
	XML        string       `xml:"-"` //最终请求串
}

func (v *NotifyResponse) Xml() error {
	data, err := xml.Marshal(v)
	v.XML = string(data)

	return err

}
