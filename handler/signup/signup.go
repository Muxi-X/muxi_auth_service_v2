package signup

import (
	"strings"

	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/Muxi-X/muxi_auth_service_v2/service"
	"github.com/gin-gonic/gin"
)

type UserSignupRequestData struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserSignupResponseData struct {
	ID uint64 `json:"id"`
}

// 用户注册
func UserSignup(c *gin.Context) {
	// 声明接收JSON数据的变量
	var data UserSignupRequestData

	// 校验输入
	if err := c.BindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "Bad Request.")
		return
	}
	// 校验输入邮箱格式
	if !strings.Contains(data.Email, "@") || !strings.HasSuffix(data.Email, ".com") {
		handler.SendBadRequest(c, errno.ErrUserSignupEmailInvalid, nil, "Email invalid.")
		return
	}

	// 并发检查用户名和邮箱是否重复
	if service.CheckUserExisted(data.Username, data.Email) {
		handler.SendResponse(c, errno.ErrUserExisted, nil)
		return
	}

	password, err := model.UserPasswordDecoder(data.Password)
	if err != nil {
		handler.SendResponse(c, errno.ErrPasswordBase64Decode, nil)
		return
	}

	newUser := model.UserModel{
		Email:        data.Email,
		Username:     data.Username,
		PasswordHash: model.GeneratePasswordHash(password),
		RoleID:       3,
		Left:         false,
	}
	// 创建记录
	err = newUser.Create()
	if err != nil {
		handler.SendError(c, errno.ErrUserCreate, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, UserSignupResponseData{newUser.Id})
	return
}
