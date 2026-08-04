package oauth

import (
	"strings"

	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	. "github.com/Muxi-X/muxi_auth_service_v2/pkg/oauth"
	"github.com/Muxi-X/muxi_auth_service_v2/service"

	"github.com/gin-gonic/gin"
	"gopkg.in/oauth2.v4/models"
)

type StoreRequest struct {
	Domain string `json:"domain"`
}

type StoreResponse struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// 注册客户端
func Store(c *gin.Context) {
	handler.SendForbidden(c, errno.ErrOAuthClientRegistrationDisabled, nil, "client registration must be approved by an administrator")
}

// AdminStore 只应挂在管理员鉴权后的路由下，用于受控创建 OAuth 客户端。
func AdminStore(c *gin.Context) {
	var rq StoreRequest
	if err := c.BindJSON(&rq); err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, err.Error())
		return
	}

	domain, err := service.NormalizeOAuthClientDomain(rq.Domain)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "Domain is invalid: "+err.Error())
		return
	}
	domainLookupKey, err := service.OAuthClientDomainLookupKey(domain)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "Domain is invalid: "+err.Error())
		return
	}

	// 域名是否已存在
	clientInfo, err := OauthServer.ClientStore.GetByDomain(domainLookupKey)
	if err != nil {
		handler.SendError(c, errno.ErrOAuthClientCreate, nil, err.Error())
		return
	}
	if clientInfo != nil && strings.TrimSpace(clientInfo.GetID()) != "" {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "Domain has already been registered.")
		return
	}

	clientID, secret := service.GenerateClientIDAndSecret()

	if err := OauthServer.ClientStore.Create(&models.Client{
		ID:     clientID,
		Secret: secret,
		Domain: domain,
	}); err != nil {
		handler.SendError(c, errno.ErrOAuthClientCreate, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, StoreResponse{
		ClientId:     clientID,
		ClientSecret: secret,
	})
}
