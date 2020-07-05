package user

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)

	user, err := model.GetUserInfoByID(userID)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, user)
}
