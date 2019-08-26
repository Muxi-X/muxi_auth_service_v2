package password

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
)

func PasswordReset(c *gin.Context) {
	data := struct {
		Captcha     string `json:"captcha"`
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, err.Error())
		return
	}
	user, err := model.GetUserByEmail(data.Email)
	if err != nil {
		handler.SendNotFound(c, errno.ErrUserNotFound, nil, err.Error())
		return
	}
	if !user.VerifyCaptcha(data.Captcha) {
		handler.SendError(c, errno.ErrUserVerifyFail, nil, "Failed.")
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
	handler.SendResponse(c, nil, "OK")
	return
}
