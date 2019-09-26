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
	// 声明接收JSON数据的变量
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

	// 声明用于检查邮箱、用户名是否重复的通信信道；用于标识检查流程是否结束的信道
	sameEmailChannel, sameUsernameChannel, done := make(chan bool), make(chan bool), make(chan struct{})
	defer close(sameEmailChannel)
	defer close(sameUsernameChannel)
	// 自动关闭done信道，这样就可以以return来代替close()方法
	defer close(done)

	// 并发检查邮箱
	go func(email string) {
		_, err := model.GetUserByEmail(email)
		// 判断检查是否已经结束
		select {
		case <-done:
			return
		default:
			{
				if err != nil { // email not found
					sameEmailChannel <- true
				} else {
					sameEmailChannel <- false
				}
			}
		}
	}(data.Email)

	// 并发检查同户名
	go func(username string) {
		_, err := model.GetUserByUsername(username)
		// 判断检查是否已经结束
		select {
		case <-done:
			return
		default:
			{
				if err != nil { // user not found
					sameUsernameChannel <- true
				} else {
					sameUsernameChannel <- false
				}
			}
		}
	}(data.Username)

	// 用于标识用户名和邮箱重复检查的状态，false为没有重复
	userExisted := false

	for round := 0; !userExisted && round < 2; round++ {
		select {
		case emailResult := <-sameEmailChannel:
			{
				if !emailResult {
					userExisted = true
					break
				}
			}
		case usernameResult := <-sameUsernameChannel:
			{
				if !usernameResult {
					userExisted = true
					break
				}
			}
		}
	}

	if userExisted {
		handler.SendResponse(c, errno.ErrUserExisted, nil)
		// 关闭done信道和两个检查信道，检查未结束的goroutine不会向信道发送数据
		return
	}

	// 正常的逻辑
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
