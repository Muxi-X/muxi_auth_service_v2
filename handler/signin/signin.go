package signin

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/auth"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
)

type UserSigninRequestData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserSigninResponseData struct {
	UserID uint64 `json:"user_id"`
	Token  string `json:"token"`
}

func UserSignin(c *gin.Context) {
	var data UserSigninRequestData
	if err := c.BindJSON(&data); err != nil {
		handler.SendBadRequest(c, nil, errno.ErrBadRequest, err.Error())
		return
	}

	user, err := model.GetUserByUsername(data.Username)
	if err != nil {
		handler.SendError(c, nil, errno.ErrUserNotFound, err.Error())
		return
	}
	if !user.CheckPassword(data.Password) {
		handler.SendError(c, nil, errno.ErrUserPasswordIncorrect, "Password not match.")
		return
	}
	token, err := auth.GenerateToken(auth.TokenPayload{
		ID:     user.Id,
		Expire: 604800,
	})
	if err != nil {
		handler.SendError(c, nil, errno.ErrToken, err.Error())
		return
	}
	handler.SendResponse(c, nil, UserSigninResponseData{
		UserID: user.Id,
		Token:  token,
	})
	return
}
