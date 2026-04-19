package middleware

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/oauth"

	"github.com/gin-gonic/gin"
)

func LoginRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		principal, err := oauth.ParseRequest(c)
		if err != nil {
			handler.SendUnauthorized(c, errno.ErrTokenInvalid, nil, err.Error())
			c.Abort()
			return
		}
		c.Set("principal", principal)
		if principal.LocalUserID != 0 {
			c.Set("userID", principal.LocalUserID)
		}
		if principal.CASUsername != "" {
			c.Set("casUsername", principal.CASUsername)
		}

		c.Next()
	}
}
