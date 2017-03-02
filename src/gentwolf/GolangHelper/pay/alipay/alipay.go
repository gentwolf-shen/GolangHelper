package alipay

import (
	"net/url"
	"sort"
	"strings"

	"gentwolf/GolangHelper/crypto"
	"gentwolf/GolangHelper/httpHelper"
)

var (
	Gateway     = "https://mapi.alipay.com/gateway.do"
	Partner     = ""
	AppKey      = ""
	SellerEmail = ""
)

type AlipayParam struct {
	OutTradeNo string
	Subject    string
	TotalFee   string
	ReturnUrl  string
	NotifyUrl  string
}

func GetPayFormData(param *AlipayParam) (string, map[string]string) {
	formData := make(map[string]string, 12)

	formData["out_trade_no"] = param.OutTradeNo
	formData["subject"] = param.Subject
	formData["payment_type"] = "1"
	formData["total_fee"] = param.TotalFee

	formData["service"] = "create_direct_pay_by_user"
	formData["partner"] = Partner
	formData["seller_email"] = SellerEmail
	formData["_input_charset"] = "utf-8"
	formData["notify_url"] = param.NotifyUrl
	formData["return_url"] = param.ReturnUrl

	formData["sign"] = sign(formData)
	formData["sign_type"] = "MD5"

	return Gateway, formData
}

func VerifyNotifyId(notifyId string) bool {
	params := url.Values{}

	params.Add("service", "notify_verify")
	params.Add("partner", Partner)
	params.Add("notify_id", notifyId)

	b, err := httpHelper.Get(Gateway, params, nil)
	return err == nil && string(b) == "true"
}

func Verify(params map[string]string) bool {
	originSign := params["sign"]
	if originSign == "" {
		return false
	}

	delete(params, "sign")
	delete(params, "sign_type")

	return originSign == sign(params)
}

func sign(params map[string]string) string {
	paramsLength := len(params)

	keys := make([]string, paramsLength)
	i := 0
	for k, _ := range params {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	tmp := make([]string, paramsLength)
	for i, k := range keys {
		tmp[i] = k + "=" + params[k]
	}
	str := strings.Join(tmp, "&") + AppKey

	return crypto.Md5(str)
}
