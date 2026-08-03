package user

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/oauth"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	principal := c.MustGet("principal").(oauth.AccessPrincipal)

	// 兼容旧的 cas:<username> access token；新签发的 CAS token 已统一为本地 user id。
	if principal.CASUsername != "" {
		handler.SendResponse(c, nil, oauth.BuildCASUserInfo(principal.CASUsername))
		return
	}

	user, err := model.GetUserInfoByID(principal.LocalUserID)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, user)
}
