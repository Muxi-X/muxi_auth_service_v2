package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/util"
	"github.com/jinzhu/gorm"
	cas "gopkg.in/cas.v2"
)

const casIdentityProvider = "cas"

type defaultCASUserResolver struct{}

func (r *defaultCASUserResolver) ResolveCASUser(_ context.Context, authenticationResponse *cas.AuthenticationResponse) (uint64, error) {
	casUsername := strings.TrimSpace(authenticationResponse.User)
	if casUsername == "" {
		return 0, errors.New("cas username is empty")
	}

	identity, err := model.GetUserIdentity(casIdentityProvider, casUsername)
	if err == nil && identity.UserID != 0 {
		return identity.UserID, nil
	}
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return 0, err
	}

	email := strings.TrimSpace(authenticationResponse.Attributes.Get("email"))
	tx := model.DB.Self.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	if existingUser, err := findExistingLocalUserForCAS(tx, casUsername, email); err != nil {
		tx.Rollback()
		return 0, err
	} else if existingUser != nil {
		identity = &model.UserIdentity{
			UserID:          existingUser.Id,
			Provider:        casIdentityProvider,
			ProviderSubject: casUsername,
			Email:           email,
		}
		if err := tx.Create(identity).Error; err != nil {
			tx.Rollback()
			return 0, err
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return 0, err
		}
		return existingUser.Id, nil
	}

	user := &model.UserModel{
		Email:        email,
		Username:     BuildCASLocalUsername(casUsername),
		PasswordHash: model.GeneratePasswordHash(util.GenerateUUID()),
		RoleID:       3,
		Left:         false,
		Info:         "cas authenticated user",
	}
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	identity = &model.UserIdentity{
		UserID:          user.Id,
		Provider:        casIdentityProvider,
		ProviderSubject: casUsername,
		Email:           email,
	}
	if err := tx.Create(identity).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	return user.Id, nil
}

func findExistingLocalUserForCAS(tx *gorm.DB, casUsername, email string) (*model.UserModel, error) {
	if email != "" {
		user := &model.UserModel{}
		err := tx.Where("email = ? AND SUBSTR(username, 1, 4) <> ?", email, "cas_").Order("id ASC").First(user).Error
		if err == nil {
			return user, nil
		}
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
	}

	user := &model.UserModel{}
	err := tx.Where("username = ? AND SUBSTR(username, 1, 4) <> ?", casUsername, "cas_").Order("id ASC").First(user).Error
	if err == nil {
		return user, nil
	}
	if !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return nil, nil
}

func BuildCASLocalUsername(casUsername string) string {
	normalized := strings.ToLower(strings.TrimSpace(casUsername))
	var b strings.Builder
	for _, r := range normalized {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}

	base := strings.Trim(b.String(), "_")
	if base == "" {
		base = "user"
	}

	hash := sha1.Sum([]byte(casUsername))
	return fmt.Sprintf("cas_%s_%s", base, hex.EncodeToString(hash[:])[:8])
}
