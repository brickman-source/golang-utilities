package alipay

type TradeAppPayExtUserInfo struct {
	Name          string `json:"name,omitempty" xml:"name"`
	Mobile        string `json:"mobile,omitempty" xml:"mobile"`
	CertType      string `json:"cert_type,omitempty" xml:"cert_type"`
	CertNo        string `json:"cert_no,omitempty" xml:"cert_no"`
	MinAge        string `json:"min_age,omitempty" xml:"min_age"`
	FixBuyer      string `json:"fix_buyer,omitempty" xml:"fix_buyer"`
	NeedCheckInfo string `json:"need_check_info,omitempty" xml:"need_check_info"`
}
type TradeAppPayAgreementSignParams struct {
	ExternalLogonID string `json:"external_logon_id,omitempty" xml:"external_logon_id"`
	AccessParams    struct {
		Channel string `json:"channel,omitempty" xml:"channel"`
	} `json:"access_params,omitempty" xml:"access_params"`
	SubMerchant struct {
		SubMerchantID                 string `json:"sub_merchant_id,omitempty" xml:"sub_merchant_id"`
		SubMerchantName               string `json:"sub_merchant_name,omitempty" xml:"sub_merchant_name"`
		SubMerchantServiceName        string `json:"sub_merchant_service_name,omitempty" xml:"sub_merchant_service_name"`
		SubMerchantServiceDescription string `json:"sub_merchant_service_description,omitempty" xml:"sub_merchant_service_description"`
	} `json:"sub_merchant,omitempty" xml:"sub_merchant"`
	PeriodRuleParams struct {
		PeriodType    string  `json:"period_type,omitempty" xml:"period_type"`
		Period        int     `json:"period,omitempty" xml:"period"`
		ExecuteTime   string  `json:"execute_time,omitempty" xml:"execute_time"`
		SingleAmount  float64 `json:"single_amount,omitempty" xml:"single_amount"`
		TotalAmount   int     `json:"total_amount,omitempty" xml:"total_amount"`
		TotalPayments int     `json:"total_payments,omitempty" xml:"total_payments"`
	} `json:"period_rule_params,omitempty" xml:"period_rule_params"`
	PersonalProductCode string `json:"personal_product_code,omitempty" xml:"personal_product_code"`
	SignScene           string `json:"sign_scene,omitempty" xml:"sign_scene"`
	ExternalAgreementNo string `json:"external_agreement_no,omitempty" xml:"external_agreement_no"`
}
type TradeAppPayExtendParams struct {
	CardType             string `json:"card_type,omitempty" xml:"card_type"`
	SysServiceProviderID string `json:"sys_service_provider_id,omitempty" xml:"sys_service_provider_id"`
	HbFqNum              string `json:"hb_fq_num,omitempty" xml:"hb_fq_num"`
	HbFqSellerPercent    string `json:"hb_fq_seller_percent,omitempty" xml:"hb_fq_seller_percent"`
	IndustryRefluxInfo   string `json:"industry_reflux_info,omitempty" xml:"industry_reflux_info"`
}
type TradeAppPay struct {
	Body               string `json:"body,omitempty" xml:"body"`
	BusinessParams     string `json:"business_params,omitempty" xml:"business_params"`
	DisablePayChannels string `json:"disable_pay_channels,omitempty" xml:"disable_pay_channels"`
	EnablePayChannels  string `json:"enable_pay_channels,omitempty" xml:"enable_pay_channels"`
	GoodsType          string `json:"goods_type,omitempty" xml:"goods_type"`
	MerchantOrderNo    string `json:"merchant_order_no,omitempty" xml:"merchant_order_no"`
	OutTradeNo         string `json:"out_trade_no,omitempty" xml:"out_trade_no"`
	PassbackParams     string `json:"passback_params,omitempty" xml:"passback_params"`
	ProductCode        string `json:"product_code,omitempty" xml:"product_code"`
	PromoParams        string `json:"promo_params,omitempty" xml:"promo_params"`
	SpecifiedChannel   string `json:"specified_channel,omitempty" xml:"specified_channel"`
	SellerId           string `json:"seller_id,omitempty" xml:"seller_id"`
	StoreID            string `json:"store_id,omitempty" xml:"store_id"`
	Subject            string `json:"subject,omitempty" xml:"subject"`
	TimeoutExpress     string `json:"timeout_express,omitempty" xml:"timeout_express"`
	TotalAmount        string `json:"total_amount,omitempty" xml:"total_amount"`
	GoodsDetail        []struct {
		Quantity       int    `json:"quantity,omitempty" xml:"quantity"`
		Price          int    `json:"price,omitempty" xml:"price"`
		GoodsCategory  string `json:"goods_category,omitempty" xml:"goods_category"`
		GoodsID        string `json:"goods_id,omitempty" xml:"goods_id"`
		AlipayGoodsID  string `json:"alipay_goods_id,omitempty" xml:"alipay_goods_id"`
		GoodsName      string `json:"goods_name,omitempty" xml:"goods_name"`
		CategoriesTree string `json:"categories_tree,omitempty" xml:"categories_tree"`
		Body           string `json:"body,omitempty" xml:"body"`
		ShowURL        string `json:"show_url,omitempty" xml:"show_url"`
	} `json:"goods_detail,omitempty" xml:"goods_detail"`
	ExtUserInfo         *TradeAppPayExtUserInfo         `json:"ext_user_info,omitempty" xml:"ext_user_info"`
	AgreementSignParams *TradeAppPayAgreementSignParams `json:"agreement_sign_params,omitempty" xml:"agreement_sign_params"`
	ExtendParams        *TradeAppPayExtendParams        `json:"extend_params,omitempty" xml:"extend_params"`
}
type RefundResponse struct {
	Sign                      string `json:"sign,omitempty" xml:"sign"`
	AlipayTradeRefundResponse struct {
		Code    string `json:"code,omitempty" xml:"code"`
		Msg     string `json:"msg,omitempty" xml:"msg"`
		SubCode string `json:"sub_code,omitempty" xml:"sub_code"`
		SubMsg  string `json:"sub_msg,omitempty" xml:"sub_msg"`

		BuyerLogonID string  `json:"buyer_logon_id,omitempty" xml:"buyer_logon_id"`
		BuyerUserID  string  `json:"buyer_user_id,omitempty" xml:"buyer_user_id"`
		FundChange   string  `json:"fund_change,omitempty" xml:"fund_change"`
		GmtRefundPay string  `json:"gmt_refund_pay,omitempty" xml:"gmt_refund_pay"`
		OpenID       string  `json:"open_id,omitempty" xml:"open_id"`
		OutTradeNo   string  `json:"out_trade_no,omitempty" xml:"out_trade_no"`
		RefundFee    float64 `json:"refund_fee,omitempty" xml:"refund_fee"`
		StoreName    string  `json:"store_name,omitempty" xml:"store_name"`
		TradeNo      string  `json:"trade_no,omitempty" xml:"trade_no"`

		RefundDetailItemList []struct {
			Amount      int     `json:"amount,omitempty" xml:"amount"`
			FundChannel string  `json:"fund_channel,omitempty" xml:"fund_channel"`
			RealAmount  float64 `json:"real_amount,omitempty" xml:"real_amount"`
		} `json:"refund_detail_item_list,omitempty" xml:"refund_detail_item_list"`
	} `json:"alipay_trade_refund_response,omitempty" xml:"alipay_trade_refund_response"`
}

type QueryOrderResponse struct {
	AlipayTradeQueryResponse struct {
		TradeNo        string `json:"trade_no,omitempty" xml:"trade_no"`
		TradeStatus    string `json:"trade_status,omitempty" xml:"trade_status"`
		SubCode        string `json:"sub_code,omitempty" xml:"sub_code"`
		BuyerLogonID   string `json:"buyer_logon_id,omitempty" xml:"buyer_logon_id"`
		BuyerUserID    string `json:"buyer_user_id,omitempty" xml:"buyer_user_id"`
		Code           string `json:"code,omitempty" xml:"code"`
		OutTradeNo     string `json:"out_trade_no,omitempty" xml:"out_trade_no"`
		TotalAmount    string `json:"total_amount,omitempty" xml:"total_amount"`
		SubMsg         string `json:"sub_msg,omitempty" xml:"sub_msg"`
		PointAmount    string `json:"point_amount,omitempty" xml:"point_amount"`
		SendPayDate    string `json:"send_pay_date,omitempty" xml:"send_pay_date"`
		OpenID         string `json:"open_id,omitempty" xml:"open_id"`
		ReceiptAmount  string `json:"receipt_amount,omitempty" xml:"receipt_amount"`
		Msg            string `json:"msg,omitempty" xml:"msg"`
		BuyerPayAmount string `json:"buyer_pay_amount,omitempty" xml:"buyer_pay_amount"`
		InvoiceAmount  string `json:"invoice_amount,omitempty" xml:"invoice_amount"`
	} `json:"alipay_trade_query_response,omitempty" xml:"alipay_trade_query_response"`
	Sign string `json:"sign,omitempty" xml:"sign"`
}
