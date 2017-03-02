package paypal

import (
	"gentwolf/GolangHelper/httpHelper"
	"net/url"
)

var (
	Host     = "https://www.paypal.com/cgi-bin/webscr"
	Business = ""
)

type PaypalParam struct {
	Invoice   string
	Fee       string
	Subject   string
	NotifyUrl string
	ReturnUrl string
}

func GetPayFormData(params *PaypalParam) (string, map[string]string) {
	formData := make(map[string]string, 10)
	formData["cmd"] = "_xclick"
	formData["business"] = Business
	formData["no_shipping"] = "1"
	formData["invoice"] = params.Invoice
	formData["charset"] = "utf-8"
	formData["currency_code"] = "USD"
	formData["amount"] = params.Fee
	formData["item_name"] = params.Subject
	formData["notify_url"] = params.NotifyUrl
	formData["return"] = params.ReturnUrl

	return Host, formData
}

func Verify(params url.Values) bool {
	params.Add("cmd", "_notify-validate")
	b, err := httpHelper.Post(Host, params, nil)
	return err == nil && string(b) == "VERIFIED"
}
