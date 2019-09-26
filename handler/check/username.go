package check

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
)

// 检查用户名是否存在
func CheckUsernameExisted(c *gin.Context) {
	flag := false
	if username, ok := c.GetQuery("username"); !ok {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "Bad Request: Required username in query string.")
		return
	} else {
		_, err := model.GetUserByUsername(username)
		if err == nil {
			flag = true
		}
	}

	handler.SendResponse(c, nil, flag)
	return
}
