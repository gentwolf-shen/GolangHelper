package ginHelper

import (
	"net/url"

	"gentwolf/GolangHelper/convert"
	"gentwolf/GolangHelper/dict"
	"gentwolf/GolangHelper/grace"
	"github.com/gin-gonic/gin"
)

func AllowCrossDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, X-Requested-With, Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,OPTIONS")
			c.Header("Access-Control-Max-Age", "3600")
			c.AbortWithStatus(200)
		}
	}
}

func AllowCrossDomainV2(domains []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := c.Request.Header.Get("Origin")
		bl := false
		for _, domain := range domains {
			if domain == host {
				bl = true
				break
			}
		}

		if bl {
			c.Header("Access-Control-Allow-Origin", host)
			c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, X-Requested-With, Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
			if c.Request.Method == "OPTIONS" {
				c.Header("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,OPTIONS")
				c.Header("Access-Control-Max-Age", "3600")
				c.AbortWithStatus(200)
			}
		}
	}
}

func ShowNoContent(c *gin.Context) {
	c.Status(204)
}

func ShowNetError(c *gin.Context) {
	ShowError(c, 5000000)
}

func ShowNoAuth(c *gin.Context) {
	ShowError(c, 4010000)
}

func ShowParamError(c *gin.Context) {
	ShowError(c, 4000001)
}

func ShowNotFound(c *gin.Context) {
	ShowError(c, 4040000)
}

func ShowError(c *gin.Context, errorCode int) {
	if errorCode == 0 {
		ShowSuccess(c, nil)
	} else {
		msg := ErrorMessage{}
		msg.Code = errorCode
		msg.Message = dict.Get(convert.ToStr(errorCode))

		ShowMsg(c, errorCode/10000, msg)
	}
}

func ShowErrorMsg(c *gin.Context, errorCode int, errMsg interface{}) {
	msg := ErrorMessage{}
	msg.Code = errorCode
	msg.Message = errMsg

	ShowMsg(c, errorCode/10000, msg)
}

func ShowSuccess(c *gin.Context, msg interface{}) {
	if msg == nil {
		msg := ErrorMessage{}
		msg.Code = 0
		msg.Message = "success"
		ShowMsg(c, 200, msg)
	} else {
		ShowMsg(c, 200, msg)
	}
}

func ShowMsg(c *gin.Context, httpCode int, msg interface{}) {
	c.JSON(httpCode, msg)
}

func RestartGin(c *gin.Context) {
	if c.Query("code") == dict.Get("restartWebCode") {
		if err := grace.SrvControl("restart"); err != nil {
			c.String(200, "restart %s", err.Error())
		} else {
			c.String(200, "restart succeed")
		}
	}
}

func GetReferer(c *gin.Context) string {
	referer := c.Request.Header.Get("Referer")
	r, _ := url.Parse(referer)
	return r.Scheme + "://" + r.Host
}

func GetHost(c *gin.Context, isAddHttp bool) string {
	host := c.Request.Host
	if isAddHttp {
		host = "https://" + host
	}
	return host
}

type ErrorMessage struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}
