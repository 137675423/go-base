package unit

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/137675423/public/pay"
	"io/ioutil"
	"net/url"
	"sort"
	"strings"
)

const (
	AliPubKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAi7n+eWwwdU93XLzCiM9Z6n1avDg1koG8rz/kwwnhdKJzw4hOyIpDI/DfkhGwmBTDxDXIBm/z70jp9ptmxxYcz8XSMCC/g0+gVPZ1hVwX/1RB4oWzXsBBaYXfOsLLTri32r0LvK3zlGtaHM2BnfIKip/ttMRm17uemPCAWGjVm83j0eaSCq9qiei3qu7u+Q9s2xl53sSJ2FWYyHZ9ZGjxq2cpqdZXqHauNgyrnVmrwGinoz8nE5sE2OWCmgeEUN4ADkwcUHdWOwfKk+nc6SpvYQfxnxsyF7xA1jCZqFP9rR6GESuRxO7V4TB6/CZlt8riIItkyS7y03XxewJeEvQAVwIDAQAB"
)

var alog map[string]string = map[string]string{
	"MD2-RSA":       "MD2WithRSA",
	"MD5-RSA":       "MD5WithRSA",
	"SHA1-RSA":      "SHA1WithRSA",
	"SHA256-RSA":    "SHA256WithRSA",
	"SHA384-RSA":    "SHA384WithRSA",
	"SHA512-RSA":    "SHA512WithRSA",
	"SHA256-RSAPSS": "SHA256WithRSAPSS",
	"SHA384-RSAPSS": "SHA384WithRSAPSS",
	"SHA512-RSAPSS": "SHA512WithRSAPSS",
}

func md5V2(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func GetCertRootSn(certPath string) (string, error) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		return "", err
	}
	strs := strings.Split(string(certData), "-----END CERTIFICATE-----")

	var cert bytes.Buffer
	for i := 0; i < len(strs); i++ {
		if strs[i] == "" {
			continue
		}
		if blo, _ := pem.Decode([]byte(strs[i] + "-----END CERTIFICATE-----")); blo != nil {
			c, err := x509.ParseCertificate(blo.Bytes)
			if err != nil {
				continue
			}
			if _, ok := alog[c.SignatureAlgorithm.String()]; !ok {
				continue
			}
			fmt.Println(c.SerialNumber, c.DNSNames, c.Version)
			si := c.Issuer.String() + c.SerialNumber.String()
			if cert.String() == "" {
				cert.WriteString(md5V2(si))
			} else {
				cert.WriteString("_")
				cert.WriteString(md5V2(si))
			}
		}

	}
	return cert.String(), nil
}

func BuildAliPriKey() string {
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

func RsaVerySignWithSha1Base64(originalData, signData string) error {
	fmt.Println(1)
	fmt.Println(originalData)
	fmt.Println(signData)
	fmt.Println(2)
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}

	// 2、读取私钥文件，解析出私钥对象
	block, _ := pem.Decode([]byte(BuildAliPriKey()))
	if block == nil {
		return nil
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	fmt.Println(pub)
	hash := crypto.SHA256
	hashInstance := hash.New()
	hashInstance.Write([]byte(originalData))
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, hashInstance.Sum(nil), sign)
}

func p2() string {
	return `gmt_create=2020-09-08+15%3A54%3A07&charset=utf-8&gmt_payment=2020-09-08+15%3A54%3A11&notify_time=2020-09-08+15%3A54%3A11&subject=EChain-Buy%3A533&sign=FCmJH02KOjxbZvLpSyOTg8kENmWscA12K02723WcBp0DIYr5K8O5LLpxLnw7pEpaQZfgRFNEGdY%2BV1hB%2FjpePGcT4334sztF%2F9Dn940lViEEZanZab6y7wWjuY1tMQMdNqamyg4aZg2%2FEGYkX9%2FJe6XZUcS7AjYGcu8pMeaib9OBmIQ3scS7qt9DM%2FoXVmDNjdcY2gw5fBARK9yKGIatWtjbZfoVvgAtYIyalu0cbeMXf4KedQ8dshTorKXCAQewbWbjZe06tK%2FjQWWFLhVrdrnskpkvM57PmghiQ6bIMlnS5I2dyfAcu6%2F7HTuzu3dmI0N6KMw0vgEjQlLZdIHSKg%3D%3D&buyer_id=2088022904078000&passback_params=touchuantest&invoice_amount=0.01&version=1.0&notify_id=2020090800222155411078001450703526&fund_bill_list=%5B%7B%22amount%22%3A%220.01%22%2C%22fundChannel%22%3A%22PCREDIT%22%7D%5D&notify_type=trade_status_sync&out_trade_no=533&total_amount=0.01&trade_status=TRADE_SUCCESS&trade_no=2020090822001478001409530357&auth_app_id=2021001190635424&receipt_amount=0.01&point_amount=0.00&buyer_pay_amount=0.01&app_id=2021001190635424&sign_type=RSA2&seller_id=2088931366545619`
}

func AliAppBack() {
	s := p2()

	var j []string
	p := strings.Split(s, "&")
	m := make(map[string]string)
	for _, v := range p {
		pa := strings.Split(v, "=")
		if len(pa) == 2 {
			val, _ := url.QueryUnescape(pa[1])
			m[pa[0]] = val
			if pa[0] != "sign" && pa[0] != "sign_type" && val != "" {
				j = append(j, pa[0]+"="+val)
			}
		}
	}
	sort.Strings(j)
	for k, v := range m {
		fmt.Println(k, ":", v)
	}
	err := RsaVerySignWithSha1Base64(strings.Join(j, "&"), m["sign"])
	fmt.Println(err)

}

func WX() {
	rsp := `<appid><![CDATA[wx7ab73e42bb32e23e]]></appid>
<bank_type><![CDATA[OTHERS]]></bank_type>
<cash_fee><![CDATA[1]]></cash_fee>
<fee_type><![CDATA[CNY]]></fee_type>
<is_subscribe><![CDATA[N]]></is_subscribe>
<mch_id><![CDATA[1602411149]]></mch_id>
<nonce_str><![CDATA[WechatCharset]]></nonce_str>
<openid><![CDATA[o6k790sgyygAdFBMFXy3j6nrosG0]]></openid>
<out_trade_no><![CDATA[585]]></out_trade_no>
<result_code><![CDATA[SUCCESS]]></result_code>
<return_code><![CDATA[SUCCESS]]></return_code>
<sign><![CDATA[39FA887715989628E79F755A4EDA6C1F]]></sign>
<time_end><![CDATA[20200908182740]]></time_end>
<total_fee>1</total_fee>
<trade_type><![CDATA[NATIVE]]></trade_type>
<transaction_id><![CDATA[4200000763202009086444101791]]></transaction_id>`

	m := pay.WeChatXmlParse(rsp)
	var key []string
	var val = make(map[string]interface{})
	for k, v := range m {
		if k != "sign" {
			key = append(key, k)
			val[k] = v
		}
	}
	fmt.Println(key)
	sort.Strings(key)
	fmt.Println(key)
	sign := pay.BuildWeChatSign(key, val)
	fmt.Println(sign == m["sign"])
}
