package oauth

import (
	"fmt"

	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	. "github.com/Muxi-X/muxi_auth_service_v2/pkg/oauth"
	"github.com/Muxi-X/muxi_auth_service_v2/util"

	"github.com/gin-gonic/gin"
	"gopkg.in/oauth2.v4/models"
)

type StoreRequest struct {
	ClientId string `json:"client_id"`
	Domain   string `json:"domain"`
}

type StoreResponse struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func Store(c *gin.Context) {
	var rq StoreRequest
	if err := c.BindJSON(&rq); err != nil {
		handler.SendBadRequest(c, err, nil, err.Error())
		return
	}
	fmt.Println(rq)

	if _, err := OauthServer.ClientStore.GetByID(c, rq.ClientId); err != nil {
		// 找到，已存在
		handler.SendBadRequest(c, err, nil, err.Error())
		return
	}

	secret := util.GenerateUUID()
	OauthServer.ClientStore.Set(rq.ClientId, &models.Client{
		ID:     rq.ClientId,
		Secret: secret,
		Domain: rq.Domain,
		UserID: "",
	})

	handler.SendResponse(c, nil, StoreResponse{
		ClientId:     rq.ClientId,
		ClientSecret: secret,
	})
}
