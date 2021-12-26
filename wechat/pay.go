package wechat

import (
	"encoding/xml"
	"fmt"
	"github.com/brickman-source/golang-utilities/crypto"
	"strings"
	"time"

	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/brickman-source/golang-utilities/rand"
	xmlHelper "github.com/brickman-source/golang-utilities/xml"
)

var (
	WxUnifiedOrderUrl = "https://api.mch.weixin.qq.com/pay/unifiedorder"
)

type UnifiedOrderRequest struct {
	XMLName        xml.Name `xml:"xml"`
	AppId          string   `xml:"appid"`                 //公众账号ID	appid	是	String(32)	wxd678efh567hg6787	微信分配的公众账号ID（企业号corpid即为此appId）
	Attach         string   `xml:"attach,omitempty"`      //附加数据	attach	否	String(127)	深圳分店	附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
	Body           string   `xml:"body"`                  //商品描述	body	是	String(128)	Ipad mini  16G  白色	商品或支付单简要描述
	Detail         string   `xml:"detail,omitempty"`      //商品详情	detail	否	String(8192)	Ipad mini  16G  白色	商品名称明细列表
	DeviceInfo     string   `xml:"device_info,omitempty"` //设备号	device_info	否	String(32)	013467007045764	终端设备号(门店号或收银设备ID)，注意：PC网页或公众号内支付请传"WEB"
	FeeType        string   `xml:"fee_type,omitempty"`    //货币类型	fee_type	否	String(16)	CNY	符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	GoodsTag       string   `xml:"goods_tag,omitempty"`   //商品标记	goods_tag	否	String(32)	WXG	商品标记，代金券或立减优惠功能的参数，说明详见代金券或立减优惠
	LimitPay       string   `xml:"limit_pay,omitempty"`   //指定支付方式	limit_pay	否	String(32)	no_credit	no_credit--指定不能使用信用卡支付
	MchId          string   `xml:"mch_id"`                //商户号	mch_id	是	String(32)	1230000109	微信支付分配的商户号
	NonceStr       string   `xml:"nonce_str"`             //随机字符串	nonce_str	是	String(32)	5K8264ILTKCH16CQ2502SI8ZNMTM67VS	随机字符串，不长于32位。推荐随机数生成算法
	NotifyURL      string   `xml:"notify_url"`            //通知地址	notify_url	是	String(256)	http://www.weixin.qq.com/wxpay/pay.php	接收微信支付异步通知回调地址，通知url必须为直接可访问的url，不能携带参数。
	Openid         string   `xml:"openid,omitempty"`      //用户标识	openid	否	String(128)	oUpF8uMuAJO_M2pxb1Q9zNjWeS6o	trade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识。openid如何获取，可参考【获取openid】。企业号请使用【企业号OAuth2.0接口】获取企业号内成员userid，再调用【企业号userid转openid接口】进行转换
	OutTradeNo     string   `xml:"out_trade_no"`          //商户订单号	out_trade_no	是	String(32)	20150806125346	商户系统内部的订单号,32个字符内、可包含字母, 其他说明见商户订单号
	ProductId      string   `xml:"product_id,omitempty"`  //商品ID	product_id	否	String(32)	12235413214070356458058	trade_type=NATIVE，此参数必传。此id为二维码中包含的商品ID，商户自行定义。
	Sign           string   `xml:"sign"`                  //签名	sign	是	String(32)	C380BEC2BFD727A4B6845133519F3AD6	签名，详见签名生成算法
	SpbillCreateIp string   `xml:"spbill_create_ip"`      //终端IP	spbill_create_ip	是	String(16)	123.12.12.123	APP和网页支付提交用户端ip，Native支付填调用微信支付API的机器IP。
	TimeExpire     string   `xml:"time_expire,omitempty"` //交易结束时间	time_expire	否	String(14)	20091227091010
	//订单失效时间，格式为yyyyMMddHHmmss，如2009年12月27日9点10分10秒表示为20091227091010。其他详见时间规则
	//注意：最短失效时间间隔必须大于5分钟

	TimeStart string `xml:"time_start,omitempty"` //交易起始时间	time_start	否	String(14)	20091225091010
	// 订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则
	TotalFee  string `xml:"total_fee"`  //总金额	total_fee	是	Int	888	订单总金额，单位为分，详见支付金额
	TradeType string `xml:"trade_type"` //交易类型	trade_type	是	String(16)	JSAPI	取值如下：JSAPI，NATIVE，APP，详细说明见参数规定

	XML string `xml:"-"`
}

func (v *UnifiedOrderRequest) Xml() error {
	output, err := xml.Marshal(v)
	xmlResult := string(output)
	v.XML = xmlResult
	return err
}

type UnifiedOrderResponse struct {
	XMLName    xml.Name        `xml:"xml"`
	ReturnCode xmlHelper.CData `xml:"return_code"`
	ReturnMsg  xmlHelper.CData `xml:"return_msg"`

	AppId      xmlHelper.CData `xml:"appid"`
	MchId      xmlHelper.CData `xml:"mch_id"`
	DeviceInfo xmlHelper.CData `xml:"device_info,omitempty"`
	NonceStr   xmlHelper.CData `xml:"nonce_str"`
	Sign       xmlHelper.CData `xml:"sign"`
	ResultCode xmlHelper.CData `xml:"result_code"`
	ErrCode    xmlHelper.CData `xml:"err_code,omitempty"`
	ErrCodeDes xmlHelper.CData `xml:"err_code_des,omitempty"`

	TradeType xmlHelper.CData `xml:"trade_type,omitempty"`
	PrepayId  xmlHelper.CData `xml:"prepay_id,omitempty"`
	CodeURL   xmlHelper.CData `xml:"code_url,omitempty"`

	ResponseXML string `xml:"-"`
}

type WxTradeType string

const (
	WxTradeTypeJSAPI  WxTradeType = "JSAPI"
	WxTradeTypeNATIVE WxTradeType = "NATIVE"
	WxTradeTypeAPP    WxTradeType = "APP"
)

func (wx *Wechat) UnifiedOrder(
	appId, mchId, apiKey string,
	productId, productName, orderId string,
	expire int64, openId string,
	totalFee int, tradeType WxTradeType,
	clientIp string, notifyURL string) (prepayId string, err error) {
	req := &UnifiedOrderRequest{
		AppId: appId, // wx.config.GetString("wechat.payment.appId"),
		MchId: mchId, // wx.config.GetString("wechat.payment.mchId"),
	}

	req.Body = productName
	req.Detail = ""
	req.GoodsTag = ""
	req.LimitPay = ""
	req.Attach = ""

	req.Openid = openId //payModel.Openid

	req.OutTradeNo = orderId
	req.ProductId = productId

	req.TimeStart = time.Now().Format("20060102150405")
	if expire == 0 {
		req.TimeExpire = time.Now().Add(time.Hour).Format("20060102150405")
	} else {
		req.TimeExpire = time.Unix(expire, 0).Format("20060102150405")
	}

	req.TotalFee = fmt.Sprintf("%d", totalFee)
	req.NonceStr = rand.GetNonceString(32)
	req.SpbillCreateIp = clientIp

	req.NotifyURL = notifyURL
	req.TradeType = string(tradeType) //"JSAPI"

	signMd5, _ := SignXmlMd5(*req, apiKey) // wx.config.GetString("wechat.payment.apiKey"))
	req.Sign = signMd5

	reqData, err := xml.Marshal(req)
	if err != nil {
		return "", err
	}

	response := &UnifiedOrderResponse{}

	log.Infof("Req Xml: %s", string(reqData))

	respData, err := http.PostData(WxUnifiedOrderUrl, "application/x-www-form-urlencoded", reqData)

	if err != nil {
		return "", err
	}

	log.Infof("Resp Xml: %s", string(respData))

	err = xml.Unmarshal(respData, response)
	if err != nil {
		return "", err
	}

	return response.PrepayId.Value, nil
}

// 针对小程序的统一下单
func (wx *Wechat) MinipUnifiedOrder(appId, mchId, apiKey,
	productId, productName, orderId string,
	expire int64, openId string,
	totalFee int, tradeType WxTradeType, clientIp, notifyURL string) (prepayId string, nonce string, timestamp int64, signStr string, err error) {
	prepayId, err = wx.UnifiedOrder(appId, mchId, apiKey, productId, productName, orderId, expire, openId, totalFee, tradeType, clientIp, notifyURL)
	if err != nil {
		log.Errorf("UnifiedOrder err: %v", err)
		return prepayId, "", 0, "", err
	}
	nonce = rand.GetNonceString(32)
	timestamp = time.Now().Unix()
	// paySign = MD5(appId=wxd678efh567hg6787&nonceStr=5K8264ILTKCH16CQ2502SI8ZNMTM67VS&package=prepay_id=wx2017033010242291fcfe0db70013231072&signType=MD5&timeStamp=1490840662&key=qazwsxedcrfvtgbyhnujmikolp111111) = 22D9B4E54AB1950F51E0649E8810ACD6
	toSign := fmt.Sprintf(
		`appId=%s&nonceStr=%s&package=prepay_id=%s&signType=MD5&timeStamp=%d&key=%s`,
		appId,
		nonce,
		prepayId,
		timestamp,
		apiKey,
	)
	signStr = strings.ToUpper(crypto.String(toSign).GetMd5String())
	return
}
