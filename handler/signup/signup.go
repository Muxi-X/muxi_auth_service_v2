package signup

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
	"strings"
)

type UserSignupRequestData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignupResponseData struct {
	ID uint64 `json:"id"`
}

func UserSignup(c *gin.Context) {
	var data UserSignupRequestData
	if err := c.BindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, "Bad Request.")
		return
	}
	if !strings.Contains(data.Email, "@") || !strings.HasSuffix(data.Email, ".com") {
		handler.SendBadRequest(c, errno.ErrUserSignupEmailInvalid, nil, "Email invalid.")
		return
	}

	sameEmailChannel, sameUsernameChannel := make(chan bool), make(chan bool)
	defer close(sameEmailChannel)
	defer close(sameUsernameChannel)

	go func(email string) {
		_, err := model.GetUserByEmail(email)
		if err != nil { // user not found
			sameEmailChannel <- true
		} else {
			sameEmailChannel <- false
		}
	}(data.Email)
	go func(username string) {
		_, err := model.GetUserByUsername(username)
		if err != nil { // user not found
			sameUsernameChannel <- true
		} else {
			sameUsernameChannel <- false
		}
	}(data.Username)

	if !<-sameEmailChannel || !<-sameUsernameChannel {
		handler.SendError(c, errno.ErrUserExisted, nil, "Email or Username duplicated.")
		close(sameUsernameChannel)
		close(sameEmailChannel)
		return
	}

	password, err := model.UserPasswordDecoder(data.Password)
	if err != nil {
		handler.SendError(c, errno.ErrPasswordBase64Decode, nil, err.Error())
		return
	}

	newUser := model.UserModel{
		BaseModel:    model.BaseModel{},
		Email:        data.Email,
		Birthday:     "",
		Hometown:     "",
		Group:        "",
		Timejoin:     "",
		Timeleft:     "",
		Username:     data.Username,
		PasswordHash: model.GeneratePasswordHash(password),
		RoleID:       3,
		Left:         false,
		ResetT:       "",
		Info:         "",
		AvatarURL:    "",
		PersonalBlog: "",
		Github:       "",
		Flickr:       "",
		Weibo:        "",
		Zhihu:        "",
	}
	err = newUser.Create()
	if err != nil {
		handler.SendError(c, errno.ErrUserCreate, nil, err.Error())
		return
	}
	handler.SendResponse(c, nil, UserSignupResponseData{newUser.Id})
	return
}
