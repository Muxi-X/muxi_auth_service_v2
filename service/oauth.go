package service

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/Muxi-X/muxi_auth_service_v2/util"
)

func GenerateClientIDAndSecret() (string, string) {
	clientID := util.GenerateUUID()
	secret := util.GenerateUUID()
	return clientID, secret
}

// CheckDomain ... 检验域名是否合理有效
func CheckDomain(s string) bool {
	_, err := NormalizeOAuthClientDomain(s)
	return err == nil
}

// NormalizeOAuthClientDomain 校验并规范化 OAuth client 可登记的业务源地址。
// 生产登记值必须是 HTTPS origin，例如 https://forum.muxistudio.xyz。
func NormalizeOAuthClientDomain(raw string) (string, error) {
	u, err := parseOAuthURL(raw)
	if err != nil {
		return "", err
	}
	if u.Path != "" && u.Path != "/" {
		return "", fmt.Errorf("domain must be an origin without path")
	}
	if u.RawQuery != "" || u.Fragment != "" {
		return "", fmt.Errorf("domain must not contain query or fragment")
	}

	return oauthURLOrigin(u), nil
}

// OAuthClientDomainLookupKey 返回 oauth2-mysql-store 实际用于按域名查找 client 的 host key。
func OAuthClientDomainLookupKey(raw string) (string, error) {
	u, err := parseOAuthURL(raw)
	if err != nil {
		return "", err
	}
	return strings.ToLower(u.Host), nil
}

// ValidateOAuthCallbackURLForDomain 校验 callback_url 是否落在已登记的业务源地址内。
func ValidateOAuthCallbackURLForDomain(registeredDomain, callbackURL string) error {
	domainOrigin, err := NormalizeOAuthClientDomain(registeredDomain)
	if err != nil {
		return fmt.Errorf("registered domain is invalid: %w", err)
	}

	callback, err := parseOAuthURL(callbackURL)
	if err != nil {
		return err
	}
	if callback.Fragment != "" {
		return fmt.Errorf("callback_url must not contain fragment")
	}

	if oauthURLOrigin(callback) != domainOrigin {
		return fmt.Errorf("callback_url origin does not match registered domain")
	}

	return nil
}

func parseOAuthURL(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("url is required")
	}
	if strings.Contains(raw, "*") {
		return nil, fmt.Errorf("wildcard is not allowed")
	}
	if strings.ContainsAny(raw, " \t\r\n") {
		return nil, fmt.Errorf("url must not contain whitespace")
	}

	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("url is invalid: %w", err)
	}
	if strings.ToLower(u.Scheme) != "https" {
		return nil, fmt.Errorf("url must use https")
	}
	if u.User != nil {
		return nil, fmt.Errorf("url must not contain user info")
	}
	if u.Host == "" {
		return nil, fmt.Errorf("url host is required")
	}

	hostname := strings.ToLower(strings.TrimSpace(u.Hostname()))
	port := u.Port()
	if hostname == "" {
		return nil, fmt.Errorf("url host is required")
	}
	if strings.Contains(u.Host, ":") && port == "" {
		return nil, fmt.Errorf("url port is invalid")
	}
	if port != "" {
		portNumber, err := strconv.Atoi(port)
		if err != nil || portNumber < 1 || portNumber > 65535 {
			return nil, fmt.Errorf("url port is invalid")
		}
	}
	if hostname == "localhost" || strings.HasSuffix(hostname, ".localhost") {
		return nil, fmt.Errorf("localhost is not allowed")
	}
	if !strings.Contains(hostname, ".") {
		return nil, fmt.Errorf("host must be a fully qualified domain name")
	}
	if net.ParseIP(hostname) != nil {
		return nil, fmt.Errorf("ip address is not allowed")
	}

	u.Scheme = "https"
	u.Host = hostname
	if port != "" {
		u.Host = net.JoinHostPort(hostname, port)
	}
	return u, nil
}

func oauthURLOrigin(u *url.URL) string {
	return u.Scheme + "://" + strings.ToLower(u.Host)
}
