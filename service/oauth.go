package service

import (
	"net/url"

	"github.com/Muxi-X/muxi_auth_service_v2/util"
)

func GenerateClientIDAndSecret() (string, string) {
	clientID := util.GenerateUUID()
	secret := util.GenerateUUID()
	return clientID, secret
}

// CheckDomain ... 检验域名是否合理有效
func CheckDomain(s string) bool {
	_, err := url.Parse(s)
	if err != nil {
		return false
	}
	return true
}
