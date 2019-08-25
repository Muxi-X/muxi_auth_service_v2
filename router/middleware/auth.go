package middleware

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/auth"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
	"time"
)

func LoginRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if res, err := auth.ParseRequest(c); err != nil {
			handler.SendUnauthorized(c, errno.ErrTokenInvalid, nil, "Token parse failed.")
			c.Abort()
			return
		} else if res.Expired < int64((time.Now().Nanosecond())) {
			handler.SendUnauthorized(c, errno.ErrTokenInvalid, nil, "Token expired.")
			c.Abort()
			return
		} else {
			c.Set("tokenID", res.ID)
		}
		c.Next()
	}
}
