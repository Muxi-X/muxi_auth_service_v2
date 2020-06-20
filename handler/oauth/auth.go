package oauth

import (
	"strconv"

	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/handler/signin"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	. "github.com/Muxi-X/muxi_auth_service_v2/pkg/oauth"
	"github.com/Muxi-X/muxi_auth_service_v2/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// 授权&授权码
// Params:
//   response_type: code
//   client_id:
//   redirect_uri:
//   token_exp(可选):
// Json:
//   username:
//   password:
func Auth(c *gin.Context) {
	// 登录

	var data signin.UserSigninRequestData
	if err := c.BindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, err.Error())
		return
	}
	// 并发检查user是否存在
	user := service.CheckUserNotExisted(data.Username)

	if user == nil {
		handler.SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}
	// 校验密码
	if !user.CheckPassword(data.Password) {
		handler.SendResponse(c, errno.ErrUserPasswordIncorrect, nil)
		return
	}

	// userID := user.Id

	// 获取授权码

	req, err := OauthServer.Server.ValidationAuthorizeRequest(c.Request)
	if err != nil {

		return
	}

	req.UserID = strconv.Itoa(int(user.Id))

	// 可设置token过期时间，纳秒
	// if tokenExp, ok := c.GetQuery("token_exp"); ok {
	// 	exp, err := strconv.ParseInt(tokenExp, 10, 64)
	// 	if err == nil {
	// 		req.AccessTokenExp = time.Duration(exp)
	// 	}
	// }

	tokenInfo, err := OauthServer.Server.GetAuthorizeToken(c, req)
	if err != nil {
		log.Error("generate auth code error", err)
		handler.SendError(c, errno.ErrGenerateAuthCode, nil, err.Error())
		return
		// redirect
	}

	handler.SendResponse(c, nil, tokenInfo)
	// redirect
}
