package email

import (
    "github.com/Muxi-X/muxi_auth_service_v2/handler"
    "github.com/Muxi-X/muxi_auth_service_v2/model"
    "github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
    "github.com/gin-gonic/gin"
)

// 根据用户名获取邮箱
func GetEmailByUsername(c *gin.Context) {
    if username, ok := c.GetQuery("username"); !ok {
        handler.SendBadRequest(c, errno.ErrBadRequest, nil, "Bad Request: Required username in query string.")
        return
    } else {
        email, err := model.GetEmailByUsername(username)
        if err != nil {
            handler.SendResponse(c, errno.ErrUserNotFound, nil)
            return
        } else {
            handler.SendResponse(c, nil, struct {
                Email string `json:"email"`
            }{email})
        }
    }
    return
}
