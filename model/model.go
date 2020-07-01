package model

import (
	"sync"
)

type BaseModel struct {
	Id uint64 `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	// CreatedAt time.Time `gorm:"column:createdAt" json:"-"`
	// UpdatedAt time.Time `gorm:"column:updatedAt" json:"-"`
	// DeletedAt *time.Time `gorm:"column:deletedAt" sql:"index" json:"-"`
}

type UserList struct {
	Lock  *sync.Mutex
	IdMap map[uint64]*UserModel
}

// Token represents a JSON web token.
type Token struct {
	Token string `json:"token"`
}

type UserInfo struct {
	Email        string `json:"email"`
	Birthday     string `json:"birthday"`
	Hometown     string `json:"hometown"`
	Group        string `json:"group"`
	Timejoin     string `json:"timejoin"`
	Timeleft     string `json:"timeleft"`
	Username     string `json:"username"`
	RoleID       uint64 `json:"role_id"`
	Left         bool   `json:"left"`
	Info         string `json:"info"`
	AvatarURL    string `json:"avatar_url"`
	PersonalBlog string `json:"personal_blog"`
	Github       string `json:"github"`
	Flickr       string `json:"flickr"`
	Weibo        string `json:"weibo"`
	Zhihu        string `json:"zhihu"`
}
