package oauth

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	. "github.com/Muxi-X/muxi_auth_service_v2/pkg/oauth"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"gopkg.in/oauth2.v4"
	e "gopkg.in/oauth2.v4/errors"
)

type AccessTokenResponse struct {
	AccessToken    string `json:"access_token"`
	AccessExpired  int64  `json:"access_expired"` // 过期时间（s）
	RefreshToken   string `json:"refresh_token"`
	RefreshExpired int64  `json:"refresh_expired"`
}

// 请求token
// Params:
//   grant_type: authorization_code
//   response_type: token
//   client_id:
// Forms:
//   client_secret:
//   code:
func Token(c *gin.Context) {
	grantType, ok := c.GetQuery("grant_type")
	if !ok || grantType != "authorization_code" {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "grant_type is required and must be authorization_code")
		return
	}

	code, ok := c.GetPostForm("code")
	if !ok {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "code is required")
		return
	}

	clientID, clientSecret, err := OauthServer.Server.ClientInfoHandler(c.Request)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "client_id and client_secret are required")
		return
	}

	tgr := &oauth2.TokenGenerateRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Request:      c.Request,
		Code:         code,
		// UserID:    userID, 在生成token时根据code获取，userID会被加入token的claim中
	}

	tokenInfo, err := OauthServer.Server.GetAccessToken(c, oauth2.GrantType(grantType), tgr)
	if err != nil {
		log.Error("GetAccessToken error", err)
		errCase := err.Error()
		if err == e.ErrInvalidGrant {
			errCase = "The code is invalid or has expired"
		}
		handler.SendError(c, errno.ErrGenerateAccessToken, nil, errCase)
		return
	}

	handler.SendResponse(c, nil, AccessTokenResponse{
		AccessToken:    tokenInfo.GetAccess(),
		AccessExpired:  int64(tokenInfo.GetAccessExpiresIn().Seconds()),
		RefreshToken:   tokenInfo.GetRefresh(),
		RefreshExpired: int64(tokenInfo.GetRefreshExpiresIn().Seconds()),
	})
}
