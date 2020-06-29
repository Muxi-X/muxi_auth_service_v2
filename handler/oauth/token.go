package oauth

import (
	"errors"
	"time"

	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	. "github.com/Muxi-X/muxi_auth_service_v2/pkg/oauth"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"gopkg.in/oauth2.v4"
)

type AccessTokenResponse struct {
	AccessToken    string `json:"token"`
	AccessExpired  int64  `json:"expired"` // 过期时间（s）
	RefreshToken   string `json:"refresh_token"`
	RefreshExpired int64  `json:"refresh_expired"`
}

// 请求token
// Params:
//   grant_type: authorization_code
//   response_type: token
//   client_id:
//   redirect_uri:
// Forms:
//   client_secret:
//   code:
func Token(c *gin.Context) {

	grantType, ok := c.GetQuery("grant_type")
	if !ok {
		err := errors.New("grant_type is required")
		handler.SendBadRequest(c, err, nil, err.Error())
		return
	} else if grantType != "authorization_code" {
		err := errors.New("auth code grant")
		handler.SendBadRequest(c, err, nil, err.Error())
		return
	}

	code, ok := c.GetPostForm("code")
	if !ok {
		err := errors.New("code")
		handler.SendBadRequest(c, err, nil, err.Error())
		return
	}

	clientID, clientSecret, err := OauthServer.Server.ClientInfoHandler(c.Request)
	if err != nil {
		handler.SendBadRequest(c, err, nil, err.Error())
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
		handler.SendError(c, errno.ErrGenerateAccessToken, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, AccessTokenResponse{
		AccessToken:    tokenInfo.GetAccess(),
		AccessExpired:  int64(tokenInfo.GetAccessExpiresIn() / time.Second),
		RefreshToken:   tokenInfo.GetRefresh(),
		RefreshExpired: int64(tokenInfo.GetRefreshExpiresIn() / time.Second),
	})
}

/*
tokenInfo:

		"ClientID": "test",
        "UserID": "123",
        "RedirectURI": "",
        "Scope": "",
        "Code": "",
        "CodeCreateAt": "0001-01-01T00:00:00Z",
        "CodeExpiresIn": 0,
        "Access": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0IiwiZXhwIjoxNTkyNjUzMzI0LCJzdWIiOiIxMjMifQ.QAYq3BgTIgBHkDYqdYkn5RA3sNZJ_03AzkNbE_uYtJiuwlwEiEF1xnpUZZbpR9lrzvrE2YMKxPDT9wWyEyrmyQ",
        "AccessCreateAt": "2020-06-19T19:42:04.261935483+08:00",
        "AccessExpiresIn": 86400000000000,
        "Refresh": "EJ0QVMY_VIO0FQYKE4SJNG",
        "RefreshCreateAt": "2020-06-19T19:42:04.261935483+08:00",
        "RefreshExpiresIn": 172800000000000
*/
