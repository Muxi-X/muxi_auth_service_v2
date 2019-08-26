package router

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler/check"
	"github.com/Muxi-X/muxi_auth_service_v2/handler/email"
	"github.com/Muxi-X/muxi_auth_service_v2/handler/signin"
	"github.com/Muxi-X/muxi_auth_service_v2/handler/signup"
	"net/http"

	"github.com/Muxi-X/muxi_auth_service_v2/handler/sd"
	"github.com/gin-gonic/gin"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares.
	g.Use(gin.Recovery())
	g.Use(mw...)
	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// loginRequired := middleware.LoginRequiredMiddleware()

	authRouter := g.Group("/auth/api")
	{
		authRouter.POST("/signup", signup.UserSignup)
		authRouter.POST("/signin", signin.UserSignin)
		authRouter.GET("/check_name", check.CheckUsernameExisted)
		authRouter.GET("/check_email", check.CheckEmailExisted)
		authRouter.GET("/email", email.GetEmailByUsername)
	}

	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	return g
}
