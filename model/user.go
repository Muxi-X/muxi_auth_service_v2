package model

import (
	"encoding/base64"
	"fmt"

	"github.com/Muxi-X/muxi_auth_service_v2/util/captcha"
	"github.com/ShiinaOrez/GoSecurity/security"
)

// User represents a registered user.
type UserModel struct {
	BaseModel
	Email        string `json:"email" column:"email"`
	Birthday     string `json:"birthday" column:"birthday"`
	Hometown     string `json:"hometown" column:"hometown"`
	Group        string `json:"group" column:"group"`
	Timejoin     string `json:"timejoin" column:"timejoin"`
	Timeleft     string `json:"timeleft" column:"timeleft"`
	Username     string `json:"username" column:"username"`
	PasswordHash string `json:"password_hash" column:"password_hash"`
	RoleID       uint64 `json:"role_id" column:"role_id"`
	Left         bool   `json:"left" column:"left"`
	ResetT       string `json:"reset_t" column:"reste_t"`
	Info         string `json:"info" column:"info"`
	AvatarURL    string `json:"avatar_url" column:"avatar_url"`
	PersonalBlog string `json:"personal_blog" column:"personal_blog"`
	Github       string `json:"github" column:"github"`
	Flickr       string `json:"flickr" column:"flickr"`
	Weibo        string `json:"weibo" column:"weibo"`
	Zhihu        string `json:"zhihu" column:"zhihu"`
}

func (c *UserModel) TableName() string {
	return "users"
}

// Create creates a new user account.
func (u *UserModel) Create() error {
	return DB.Self.Create(&u).Error
}

// DeleteUser deletes the user by the user identifier.
func DeleteUser(id uint64) error {
	user := UserModel{}
	user.BaseModel.Id = id
	return DB.Self.Delete(&user).Error
}

// Update updates an user account information.
func (u *UserModel) Update() error {
	return DB.Self.Save(u).Error
}

func (u *UserModel) IsAdmin() bool {
	return u.RoleID == 2
}

func (user *UserModel) CheckPassword(passwordBase64 string) bool {
	password, err := UserPasswordDecoder(passwordBase64)
	if err != nil {
		return false
	}
	return security.CheckPasswordHash(password, user.PasswordHash)
}

func UserPasswordDecoder(passwordBase64 string) (string, error) {
	passwordBytes, err := base64.StdEncoding.DecodeString(passwordBase64)
	if err != nil {
		return "", err
	}
	return string(passwordBytes), err
}

func GeneratePasswordHash(password string) string {
	return security.GeneratePasswordHash(password)
}

func GetUserByID(id uint64) (*UserModel, error) {
	user := &UserModel{}
	d := DB.Self.Where("id = ?", id).First(&user)
	return user, d.Error
}

func GetUserByEmail(email string) (*UserModel, error) {
	user := &UserModel{}
	d := DB.Self.Where("email = ?", email).First(&user)
	return user, d.Error
}

func GetUserByUsername(username string) (*UserModel, error) {
	user := &UserModel{}
	d := DB.Self.Where("username = ?", username).First(&user)
	return user, d.Error
}

func GetEmailByUsername(username string) (string, error) {
	user := &UserModel{}
	d := DB.Self.Select("email").Where("username = ?", username).First(&user)
	return user.Email, d.Error
}

func (user *UserModel) VerifyCaptcha(newCap string) bool {
	oldCap, err := captcha.ResolveCaptchaToken(user.ResetT)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	fmt.Println(oldCap, newCap)
	return oldCap == newCap
}
