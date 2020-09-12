package wechat

import (
	"encoding/xml"
)

type JsApiRequest struct {
	XMLName xml.Name `xml:"xml" json:"-"`

	AppId     string `xml:"appId" json:"appId"`         //公众号id	appId	是	String(16)	wx8888888888888888	商户注册具有支付权限的公众号成功后即可获得
	NonceStr  string `xml:"nonceStr" json:"nonceStr"`   //随机字符串	nonceStr	是	String(32)	5K8264ILTKCH16CQ2502SI8ZNMTM67VS	随机字符串，不长于32位。推荐随机数生成算法
	Package   string `xml:"package" json:"package"`     //订单详情扩展字符串	package	是	String(128)	prepay_id=123456789	统一下单接口返回的prepay_id参数值，提交格式如：prepay_id=***
	SignType  string `xml:"signType" json:"signType"`   //签名方式	signType	是	String(32)	MD5	签名算法，暂支持MD5
	TimeStamp string `xml:"timeStamp" json:"timeStamp"` //时间戳	timeStamp	是	String(32)	1414561699	当前的时间，其他详见时间戳规则

	PaySign string `xml:"-" json:"paySign"` //最终请求串
}

func (v *JsApiRequest) SignMd5(appKey string) bool {
	sign, _ := SignXmlMd5(*v, appKey) //wx.config.GetString("wechat.payment.appKey"))
	v.PaySign = sign
	return true
}
