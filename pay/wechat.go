package pay

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	//微信商户号
	WechatMid = ""
	//微信appid
	WechatAppId = ""
	//微信appkey
	WechatKey = ""
	//微信签名加密类型
	WechatSignType = "MD5"
	//微信统一下单接口
	WechatPayUrl = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	//微信App支付结果通知地址
	WechatNotifyUrl = "http://47.106.179.241:7336/callBack/wechatPayBack"
)

type WechatClient struct {
	Url, AppId, AppPrivateKey, Format, Charset, WechatPayPublicKey, SignType string
}

type WechatSignGet struct {
	//商户网站唯一订单号
	OutTradeNo int64 `json:"outTradeNo"`
	//支付渠道 1支付宝2微信
	PayChannel int64 `json:"payChannel"`
}

func NewWechatSign(p AppSignParam) *WechatSign {
	as := new(WechatSign)
	as.SignKey = []string{"appid", "mch_id", "nonce_str", "sign", "sign_type", "body", "out_trade_no", "total_fee", "spbill_create_ip", "notify_url", "trade_type"}
	as.PayKey = []string{"appid", "partnerid", "prepayid", "package", "noncestr", "timestamp"}
	sort.Strings(as.SignKey)
	sort.Strings(as.PayKey)
	as.SignData = map[string]interface{}{
		"appid":            WechatAppId,
		"mch_id":           WechatMid,
		"nonce_str":        "WechatCharset",
		"sign_type":        WechatSignType,
		"body":             "EChain-Buy",
		"spbill_create_ip": p.ClientIp,
		"notify_url":       WechatNotifyUrl,
		"out_trade_no":     fmt.Sprintf("%d", p.OutTradeNo),
		"total_fee":        p.TotalAmount * 100,
	}
	return as
}

//微信签名参数
type WechatSign struct {
	//统一下单参与签名的参数
	SignKey []string
	//统一下单签名数据
	SignData map[string]interface{}
	//统一下单签名返回数据
	SignRsp map[string]string
	//App发起支付与签名的参数
	PayKey []string
	//App发起支付签名数据
	PayData map[string]interface{}
}

//构建微信签名
func BuildWeChatSign(key []string, val map[string]interface{}) string {
	var signArr []string
	for _, v := range key {
		vak, ok := val[v]
		if ok && vak != "" {
			signArr = append(signArr, fmt.Sprintf("%s=%v", v, vak))
		}
	}
	sta := strings.Join(signArr, "&") + fmt.Sprintf("&key=%s", WechatKey)
	return strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(sta))))
}

func (as *WechatSign) ToXml() (s string) {
	var arr []string
	arr = append(arr, `<xml>`)
	for k, v := range as.SignData {
		arr = append(arr, fmt.Sprintf(`<%s>%v</%s>`, k, v, k))
	}
	arr = append(arr, `</xml>`)

	return strings.Join(arr, "\n")
}

func (as *WechatSign) ParseXml(data []byte) {
	as.SignRsp = WeChatXmlParse(string(data))
	as.PayData = map[string]interface{}{
		"appid":     WechatAppId,
		"partnerid": WechatMid,
		"prepayid":  as.SignRsp["prepay_id"],
		"package":   "Sign=WXPay",
		"noncestr":  as.SignRsp["sign"],
		"timestamp": fmt.Sprintf("%v", time.Now().Unix()),
	}
}

func WeChatXmlParse(s string) map[string]string {
	rsp := strings.Trim(s, `</xml>`)
	p := strings.Split(rsp, "\n")
	var m = make(map[string]string)
	for _, v := range p {
		line := strings.Trim(v, `</>`)
		one := strings.Split(line, "><")
		if len(one) > 2 {
			m[one[0]] = strings.Trim(strings.ReplaceAll(one[1], "CDATA", ""), `![[]]`)
		} else {
			line := strings.Split(v, "/")
			one := strings.Split(line[0], ">")
			m[strings.Trim(one[0], "<>")] = strings.Trim(one[1], "<>")
		}
	}
	return m
}
