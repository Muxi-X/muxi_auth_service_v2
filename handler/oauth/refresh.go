package oauth

import (
	"time"

	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	. "github.com/Muxi-X/muxi_auth_service_v2/pkg/oauth"

	"github.com/gin-gonic/gin"
	"gopkg.in/oauth2.v4"
)

// 更新 access token
// Params:
//   grant_type: refresh_token
//   client_id:
// Forms:
//   client_secret:
//   refresh_token:
func Refresh(c *gin.Context) {
	grantType, ok := c.GetQuery("grant_type")
	if !ok || grantType != "refresh_token" {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "grant_type is required and must be refresh_token")
		return
	}

	refreshToken, ok := c.GetPostForm("refresh_token")
	if !ok {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "refresh_token is required")
		return
	}

	clientID, clientSecret, err := OauthServer.Server.ClientInfoHandler(c.Request)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, err.Error())
		return
	}

	tgr := &oauth2.TokenGenerateRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Request:      c.Request,
		Refresh:      refreshToken,
	}

	tokenInfo, err := OauthServer.Server.GetAccessToken(c, oauth2.GrantType(grantType), tgr)
	if err != nil {
		handler.SendError(c, errno.ErrRefreshToken, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, AccessTokenResponse{
		AccessToken:    tokenInfo.GetAccess(),
		AccessExpired:  int64(tokenInfo.GetAccessExpiresIn() / time.Second),
		RefreshToken:   tokenInfo.GetRefresh(),
		RefreshExpired: int64(tokenInfo.GetRefreshExpiresIn() / time.Second),
	})
}
