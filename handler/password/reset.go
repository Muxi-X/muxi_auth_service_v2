package password

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
)

type PasswordResetRequest struct {
	Captcha     string `json:"captcha" binding:"required"`
	Email       string `json:"email" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// 重设密码
func PasswordReset(c *gin.Context) {
	data := PasswordResetRequest{}

	if err := c.BindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, err.Error())
		return
	}

	user, err := model.GetUserByEmail(data.Email)
	if err != nil {
		handler.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	if !user.VerifyCaptcha(data.Captcha) {
		handler.SendResponse(c, errno.ErrUserVerifyFail, nil)
		return
	}
	password, err := model.UserPasswordDecoder(data.NewPassword)

	if err != nil {
		handler.SendError(c, errno.ErrPasswordBase64Decode, nil, err.Error())
		return
	}
	user.PasswordHash = model.GeneratePasswordHash(password)

	if err := user.Update(); err != nil {
		handler.SendError(c, errno.ErrUserUpdate, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, nil)
	return
}
