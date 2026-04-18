package handler

import (
	"net/http"

	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/logx"
	"github.com/Muxi-X/muxi_auth_service_v2/util"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendResponse(c *gin.Context, err error, data interface{}) {
	code, message := errno.DecodeErr(err)

	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func SendUnauthorized(c *gin.Context, err error, data interface{}, cause string) {
	code, message := errno.DecodeErr(err)
	logx.Info(message, "X-Request-Id", util.GetReqID(c), "cause", cause)
	c.JSON(http.StatusUnauthorized, Response{
		Code:    code,
		Message: message + ": " + cause,
		Data:    nil,
	})
}

func SendBadRequest(c *gin.Context, err error, data interface{}, cause string) {
	code, message := errno.DecodeErr(err)
	logx.Info(message, "X-Request-Id", util.GetReqID(c), "cause", cause)
	c.JSON(http.StatusBadRequest, Response{
		Code:    code,
		Message: message + ": " + cause,
		Data:    data,
	})
}

func SendNotFound(c *gin.Context, err error, data interface{}, cause string) {
	code, message := errno.DecodeErr(err)
	logx.Info(message, "X-Request-Id", util.GetReqID(c), "cause", cause)
	c.JSON(http.StatusNotFound, Response{
		Code:    code,
		Message: message + ": " + cause,
		Data:    data,
	})
}

func SendError(c *gin.Context, err error, data interface{}, cause string) {
	code, message := errno.DecodeErr(err)
	logx.Info(message, "X-Request-Id", util.GetReqID(c), "cause", cause)
	c.JSON(http.StatusInternalServerError, Response{
		Code:    code,
		Message: message + ": " + cause,
		Data:    data,
	})
}
