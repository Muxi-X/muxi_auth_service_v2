package oauth

import (
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
	var rq StoreRequest
	if err := c.BindJSON(&rq); err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, err.Error())
		return
	}

	// 域名是否有效
	if ok := service.CheckDomain(rq.Domain); !ok {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "Domain is invalid.")
		return
	}

	// 域名是否已存在
	if _, err := OauthServer.ClientStore.GetByDomain(rq.Domain); err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, err.Error())
		return
	}

	clientID, secret := service.GenerateClientIDAndSecret()

	OauthServer.ClientStore.Create(&models.Client{
		ID:     clientID,
		Secret: secret,
		Domain: rq.Domain,
	})

	handler.SendResponse(c, nil, StoreResponse{
		ClientId:     clientID,
		ClientSecret: secret,
	})
}
