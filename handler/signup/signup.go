package signup

import (
	"strings"

	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
)

type UserSignupRequestData struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserSignupResponseData struct {
	ID uint64 `json:"id"`
}

// 用户注册
func UserSignup(c *gin.Context) {
	var data UserSignupRequestData

	// 校验输入
	if err := c.BindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "Bad Request.")
		return
	}
	// 校验输入邮箱格式
	if !strings.Contains(data.Email, "@") || !strings.HasSuffix(data.Email, ".com") {
		handler.SendBadRequest(c, errno.ErrUserSignupEmailInvalid, nil, "Email invalid.")
		return
	}

	sameEmailChannel, sameUsernameChannel, done := make(chan bool), make(chan bool), make(chan struct{})
	defer close(sameEmailChannel)
	defer close(sameUsernameChannel)

	go func(email string) {
		_, err := model.GetUserByEmail(email)
		select {
		case <-done:
			return
		default: {
			if err != nil { // email not found
				sameEmailChannel <- true
			} else {
				sameEmailChannel <- false
			}
		}
		}
	}(data.Email)

	go func(username string) {
		_, err := model.GetUserByUsername(username)
		select {
		case <-done:
			return
		default: {
			if err != nil { // user not found
				sameUsernameChannel <- true
			} else {
				sameUsernameChannel <- false
			}
		}
		}
	}(data.Username)

	userExisted := false

	for round:=0; !userExisted && round < 2; round++ {
		select {
		case emailResult := <- sameEmailChannel: {
			if !emailResult {
				userExisted = true
				break
			}
		}
		case usernameResult := <- sameUsernameChannel: {
			if !usernameResult {
				userExisted = true
				break
			}
		}
		}
	}

	if userExisted {
		close(done)
		handler.SendResponse(c, errno.ErrUserExisted, nil)
		return
	} else {
		defer close(done)
	}

	password, err := model.UserPasswordDecoder(data.Password)
	if err != nil {
		handler.SendResponse(c, errno.ErrPasswordBase64Decode, nil)
		return
	}

	newUser := model.UserModel{
		Email:        data.Email,
		Username:     data.Username,
		PasswordHash: model.GeneratePasswordHash(password),
		RoleID:       3,
		Left:         false,
	}
	err = newUser.Create()
	if err != nil {
		handler.SendError(c, errno.ErrUserCreate, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, UserSignupResponseData{newUser.Id})
	return
}
