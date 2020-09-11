package pay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
	"time"
)

const (
	AliAppId    = ""
	AliCharset  = "utf-8"
	AliVersion  = "1.0"
	AliSignType = "RSA2"
	AliFormat   = "json"
)

const (
	//阿里支付网关
	AliGateWay = "https://openapi.alipaydev.com/gateway.do"
	//APP公钥
	AppPubKey = ""
	//APP私钥内容串
	AppPriKeySub = ""
	//阿里云公钥
	AliPubKey = ""
)

//APP私钥
var AppPriKey = BuildAliPriKey()

func BuildAliPriKey() string {
	var sl []string
	sl = append(sl, "-----BEGIN RSA PRIVATE KEY-----")
	i := 0
	for {
		if (i+1)*64 < len(AppPriKeySub) {
			sl = append(sl, AppPriKeySub[i*64:(i+1)*64])
		} else {
			sl = append(sl, AppPriKeySub[i*64:])
			sl = append(sl, "-----END RSA PRIVATE KEY-----")
			return strings.Join(sl, "\n")
		}
		i++
	}
}

type AliClient struct {
	Url, AppId, AppPrivateKey, Format, Charset, AliPayPublicKey, SignType string
}

func NewAliClient() *AliClient {
	a := new(AliClient)
	a.Url = "https://openapi.alipay.com/gateway.do"
	a.AppId = AliAppId
	return a
}

var AliCertClient = NewAliCertClient()

type aliCertClient struct {
	ServerUrl        string `json:"server_url"`
	AppPrivateKey    string `json:"app_private_key"`
	AppId            string `json:"app_id"`
	Format           string `json:"format"`
	Charset          string `json:"charset"`
	SignType         string `json:"sign_type"`
	AppCertSn        string `json:"app_cert_sn"`
	AlipayRootCertSn string `json:"alipay_root_cert_sn"`
}

func NewAliCertClient() *aliCertClient {
	ac := new(aliCertClient)
	//设置网关地址
	ac.ServerUrl = "https://openapi.alipay.com/gateway.do"
	//设置应用Id
	ac.AppId = AliAppId
	//设置请求格式，固定值json
	ac.Format = AliFormat
	//设置字符集
	ac.Charset = AliCharset
	//设置签名类型
	ac.SignType = AliSignType
	return ac
}

//app签名接口入参
type AppSignGet struct {
	//商户网站唯一订单号
	OutTradeNo int64 `json:"outTradeNo"`
	//支付渠道 1支付宝2微信
	PayChannel int64 `json:"payChannel"`
	//pc支付成功跳转地址,app无需传此参数
	ReturnUrl string `json:"returnUrl"`
}

//app签名参数
type AppSignParam struct {
	//用户ID
	UserId string
	//客户端ip
	ClientIp string
	//商户网站唯一订单号
	OutTradeNo int64 `json:"outTradeNo"`
	//支付渠道 1支付宝2微信
	PayChannel int64 `json:"payChannel"`
	//支付总额 元
	TotalAmount float64 `json:"payChannel"`
	//支付方式 APP PC
	PayMode string `json:"payChannel"`
}

func NewAliSign() *AliSign {
	as := new(AliSign)
	as.SignKey = []string{"app_id", "biz_content", "format", "charset", "sign_type", "method", "timestamp", "version", "notify_url", "return_url"}
	as.SignData = map[string]string{
		"app_id":    AliAppId,
		"format":    AliFormat,
		"charset":   AliCharset,
		"sign_type": AliSignType,
		"method":    "alipay.trade.app.pay",
		"version":   AliVersion,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}
	return as
}
func NewAliBizContent() *AliBizContent {
	abc := new(AliBizContent)
	abc.ProductCode = "QUICK_MSECURITY_PAY"
	return abc
}

//阿里云签名返回结构
type AliSign struct {
	//参与签名的参数
	SignKey []string
	//参与签名的数据
	SignData map[string]string
}

//等待签名字符串
func (as *AliSign) WaitStr() string {
	return ""
}

// RSASign 私钥签名
func RSASign(data []byte) (string, error) {
	// 1、选择hash算法，对需要签名的数据进行hash运算
	myhash := crypto.SHA256
	hashInstance := myhash.New()
	hashInstance.Write(data)
	hashed := hashInstance.Sum(nil)
	// 2、读取私钥文件，解析出私钥对象
	block, perr := pem.Decode([]byte(AppPriKey))
	if block == nil {
		fmt.Println(perr)
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		return "", err
	}
	// 3、RSA数字签名（参数是随机数、私钥对象、哈希类型、签名文件的哈希串，生成bash64编码）
	bytes, err := rsa.SignPKCS1v15(rand.Reader, privateKey, myhash, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (as *AliSign) SignBuild(sour string) (s string) {

	s, e := RSASign([]byte(sour))
	if e != nil {
		fmt.Println(e)
	}
	return s
}

type AliBizContent struct {
	//商品的标题/交易标题/订单标题/订单关键字等。
	Subject string `json:"subject"`
	//商户网站唯一订单号
	OutTradeNo string `json:"out_trade_no"`
	//该笔订单允许的最晚付款时间，逾期将关闭交易。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。 注：若为空，则默认为15d。
	TimeoutExpress string `json:"timeout_express"`
	//订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000]
	TotalAmount string `json:"total_amount"`
	//销售产品码，商家和支付宝签约的产品码，APP支付功能中该值固定为： QUICK_MSECURITY_PAY
	ProductCode string `json:"product_code"`
	//商品主类型：0—虚拟类商品；1—实物类商品 注：虚拟类商品不支持使用花呗渠道
	GoodsType string `json:"goods_type"`
	//公用回传参数，如果请求时传递了该参数，则返回给商户时会回传该参数。支付宝会在异步通知时将该参数原样返回。本参数必须进行 UrlEncode 之后才可以发送给支付宝
	PassbackParams string `json:"passback_params"`
}

//导出非空字段
func (abc *AliBizContent) Explore() map[string]string {
	m := make(map[string]string)
	m["subject"] = abc.Subject
	m["out_trade_no"] = abc.OutTradeNo
	m["total_amount"] = abc.TotalAmount
	m["product_code"] = abc.ProductCode
	m["passback_params"] = "touchuantest"
	m["timeout_express"] = "1h"
	return m

}

//生成阿里公钥
func BuildAliPubKey() string {
	var sl []string
	sl = append(sl, "-----BEGIN PUBLIC KEY-----")
	i := 0
	for {
		if (i+1)*64 < len(AliPubKey) {
			sl = append(sl, AliPubKey[i*64:(i+1)*64])
		} else {
			sl = append(sl, AliPubKey[i*64:])
			sl = append(sl, "-----END PUBLIC KEY-----")
			return strings.Join(sl, "\n")
		}
		i++
	}
}

//阿里验签
func RsaVerySignWithSha1Base64(originalData, signData string) error {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}

	// 2、读取公钥对象
	block, _ := pem.Decode([]byte(BuildAliPubKey()))
	if block == nil {
		return nil
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	hash := crypto.SHA256
	hashInstance := hash.New()
	hashInstance.Write([]byte(originalData))
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, hashInstance.Sum(nil), sign)
}
