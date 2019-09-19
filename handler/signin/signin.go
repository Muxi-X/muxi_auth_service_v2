package signin

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/auth"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
)

type UserSigninRequestData struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserSigninResponseData struct {
	UserID uint64 `json:"user_id"`
	Token  string `json:"token"`
}

// 用户登录
func UserSignin(c *gin.Context) {
	var (
		data UserSigninRequestData
		err  error
		user *model.UserModel
	)

	if err = c.BindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, err.Error())
		return
	}

	user, err = model.GetUserByUsername(data.Username)
	if err != nil {
		user, err = model.GetUserByEmail(data.Username)
		if err != nil {
			handler.SendResponse(c, errno.ErrUserNotFound, nil)
			return
		}
	}

	if !user.CheckPassword(data.Password) {
		handler.SendResponse(c, errno.ErrUserPasswordIncorrect, nil)
		return
	}

	token, err := auth.GenerateToken(auth.TokenPayload{
		ID:     user.Id,
		Expire: 7 * 60 * 60 * 24, // 有效时间七天
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
