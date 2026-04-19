package user

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/oauth"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	principal := c.MustGet("principal").(oauth.AccessPrincipal)

	// CAS 主体不再映射回本地 users 表，直接返回独立的 CAS 用户信息视图。
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
