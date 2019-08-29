package password

import (
	"fmt"
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/constvar"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/Muxi-X/muxi_auth_service_v2/util/captcha"
	"github.com/Muxi-X/muxi_auth_service_v2/util/smtpMail"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strings"
)

func GetCaptcha(c *gin.Context) {
	data := struct {
		Email string `json:"email"`
	}{}
	if err := c.BindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBadRequest, nil, err.Error())
		return
	} else {
		user, err := model.GetUserByEmail(data.Email)
		if err != nil {
			handler.SendNotFound(c, errno.ErrUserNotFound, nil, err.Error())
			return
		} else {
			email := user.Email
			captchaCode := captcha.GetCaptcha(6)
			mailContent := strings.Replace(strings.Replace(constvar.EmailTemp, "YourEmailAddress", email, 10), "TheCaptcha", captchaCode, 1)

			mailSendErrorChan, userUpdateErrorChan := make(chan error), make(chan error)
			defer close(mailSendErrorChan)
			defer close(userUpdateErrorChan)

			go func() {
				fmt.Println("Start to send email.", viper.GetString("authcode"))
				mailSendErrorChan <- smtpMail.SendMail("muxistudio@qq.com", viper.GetString("authcode"), []string{email}, smtpMail.Content{
					NickName:    "Muxi Studio: Auth Service",
					User:        "muxistudio@qq.com",
					Subject:     "Auth Code For Password Reseting: 密码重置",
					Body:        mailContent,
					ContentType: "Content-Type: text/html; charset=UTF-8",
				})
				fmt.Println("Sent")
			}()

			captchaToken, newErr := captcha.GenerateCaptchaToken(captchaCode)
			if newErr != nil {
				handler.SendError(c, errno.ErrGenerateCaptchaToken, nil, err.Error())
				return
			} else {
				go func() {
					user.ResetT = captchaToken
					userUpdateErrorChan <- user.Update()
				}()
			}

			if err := <-mailSendErrorChan; err != nil {
				handler.SendError(c, errno.ErrMailSend, nil, err.Error())
				return
			}
			if err := <-userUpdateErrorChan; err != nil {
				handler.SendError(c, errno.ErrUserUpdate, nil, err.Error())
				return
			}
		}
	}
	handler.SendResponse(c, nil, "Sent OK.")
	return
}
