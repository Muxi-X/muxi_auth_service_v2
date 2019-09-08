package check

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/auth"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
)

func CheckToken(c *gin.Context) {
	token := c.Param("token")
	email := c.Param("email")
	if token == "" || email == "" {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "Token and email are required.")
		return
	}
	tokenResolve, err := auth.ResolveToken(token)
	if err != nil {
		handler.SendError(c, errno.ErrTokenInvalid, nil, err.Error())
		return
	}

	user, err := model.GetUserByEmail(email)
	if err != nil {
		handler.SendNotFound(c, errno.ErrUserNotFound, nil, err.Error())
		return
	}

	if tokenResolve.ID != user.Id {
		handler.SendError(c, errno.ErrTokenInvalid, nil, "ID not match.")
		return
	}

	handler.SendResponse(c, nil, "OK")
	return
}