package check

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
)

// 检查邮箱是否存在
func CheckEmailExisted(c *gin.Context) {
	flag := false
	if email, ok := c.GetQuery("email"); !ok {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "Bad Request: Required email in query string.")
		return
	} else {
		_, err := model.GetUserByEmail(email)
		if err == nil {
			flag = true
		}
	}
	handler.SendResponse(c, nil, flag)
	return
}
