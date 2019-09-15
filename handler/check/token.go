package check

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/auth"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
	"fmt"
)

func CheckToken(c *gin.Context) {
	var token, email string
	var ok bool
	if token, ok = c.GetQuery("token"); !ok {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "Token is required.")
		return
	}
	if email, ok = c.GetQuery("email"); !ok {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "Email is required.")
		return
	}

	tokenResolve, err := auth.ResolveToken(token)
	if err != nil {
		fmt.Println("[email]:", email)
		fmt.Println("[token]:", token)
		fmt.Println("[error]:", err.Error())
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
