package model

import (
	"fmt"

	"github.com/Muxi-X/muxi_auth_service_v2/pkg/constvar"
)

type Role struct {
	BaseModel
	Name        string `json:"name" column:"name"`
	Default     bool   `json:"default" column:"default"`
	Permissions uint64 `json:"permissions" column:"permissions"`
}

func (role *Role) TableName() string {
	return "roles"
}

func (role *Role) Create() error {
	return DB.Self.Create(&role).Error
}

func (role *Role) Update() error {
	return DB.Self.Save(role).Error
}

func (role *Role) GetUsers(offset, limit int) ([]*UserModel, uint64, error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}
	var count uint64

	where := fmt.Sprintf("role_id = %d", role.Id)
	if err := DB.Self.Where(where).Count(&count).Error; err != nil {
		return nil, count, err
	}
	userSlice := make([]*UserModel, limit)[:0]
	if err := DB.Self.Where(where).Offset(offset).Limit(limit).Find(&userSlice).Error; err != nil {
		return nil, count, err
	}
	return userSlice, count, nil
}
