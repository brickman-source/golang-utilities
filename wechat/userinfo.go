package wechat

import (
	"errors"
	"github.com/brickman-source/golang-utilities/http"
	"net/url"
)

type UserInfoResult struct {
	Errcode    int      `json:"errcode,omitempty" xml:"errcode,omitempty" bson:"-"`
	Errmsg     string   `json:"errmsg,omitempty" xml:"errmsg,omitempty" bson:"-"`
	City       string   `json:"city" bson:"city" xml:"city"`
	Country    string   `json:"country" bson:"country" xml:"country"`
	HeadImgURL string   `json:"headimgurl" bson:"headImgURL" xml:"headimgurl"`
	Language   string   `json:"language" bson:"language" xml:"language"`
	Nickname   string   `json:"nickname" bson:"nickname" xml:"nickname"`
	Openid     string   `json:"openid" bson:"openId" xml:"openid"`
	Privilege  []string `json:"privilege" bson:"privilege" xml:"privilege"`
	Province   string   `json:"province" bson:"province" xml:"province"`
	Sex        int      `json:"sex" bson:"sex" xml:"sex"`
}

func (wx *Wechat) GetUserInfo(openId, accessTtoken string) (*UserInfoResult, error) {
	requestURL, _ := url.Parse("https://api.weixin.qq.com/sns/userinfo")
	parameters := requestURL.Query()

	parameters.Set("access_token", accessTtoken)
	parameters.Set("openid", openId)
	parameters.Set("lang", "zh_CN")

	requestURL.RawQuery = parameters.Encode()

	ret := &UserInfoResult{}

	err := http.GetJSON(requestURL.String(), ret)

	if err != nil {
		return nil, err
	}
	if ret.Errcode != 0 {
		return nil, errors.New(ret.Errmsg)
	}
	return ret, nil
}
