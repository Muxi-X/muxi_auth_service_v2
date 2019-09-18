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
	var err, newErr error
	if err = c.BindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, err.Error())
		return
	}
	var user *model.UserModel
	user, err = model.GetUserByUsername(data.Username)
	if err != nil {
		user, newErr = model.GetUserByEmail(data.Username)
		if newErr != nil {
			handler.SendError(c, errno.ErrUserNotFound, nil, err.Error())
			return
		}
	}
	if !user.CheckPassword(data.Password) {
		handler.SendError(c, errno.ErrUserPasswordIncorrect, nil, "Password not match.")
		return
	}
	token, err := auth.GenerateToken(auth.TokenPayload{
		ID:     user.Id,
		Expire: 604800,
	})
	if err != nil {
		handler.SendError(c, errno.ErrToken, nil, err.Error())
		return
	}
	handler.SendResponse(c, nil, UserSigninResponseData{
		UserID: user.Id,
		Token:  token,
	})
	return
}
