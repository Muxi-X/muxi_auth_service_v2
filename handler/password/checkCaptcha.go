package password

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
)

type CheckCaptchaRequest struct {
	Captcha string `json:"captcha" binding:"required"`
	Email   string `json:"email" binding:"required"`
}

// 检查验证码
func CheckCaptcha(c *gin.Context) {
	data := CheckCaptchaRequest{}

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
	} else {
		handler.SendResponse(c, nil, nil)
	}
	return
}
