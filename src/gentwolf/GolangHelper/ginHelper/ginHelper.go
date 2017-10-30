package ginHelper

import (
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
			c.AbortWithStatus(200)
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

func ShowError(c *gin.Context, errorCode int) {
	msg := ErrorMessage{}
	msg.Code = errorCode
	msg.Message = dict.Get(convert.ToStr(errorCode))

	ShowMsg(c, errorCode/10000, msg)
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

type ErrorMessage struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}
