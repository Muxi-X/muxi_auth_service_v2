package oauth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const casSubjectPrefix = "cas:"

// AccessPrincipal 表示从 access token 中解析出来的认证主体。
// 本地登录与 CAS 登录都会落到这里，但它们的身份来源是分开的。
type AccessPrincipal struct {
	Subject     string
	LocalUserID uint64
	CASUsername string
}

// BuildCASSubject 会把 CAS 用户名编码成独立 subject。
// 使用显式前缀可以避免和原有本地数值 user_id 冲突。
func BuildCASSubject(casUsername string) string {
	return casSubjectPrefix + casUsername
}

func ResolveAccessToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"]
		if !ok {
			return "", errors.New("Token not include `cap` field.")
		}
		return sub.(string), nil
	}
	return "", errors.New("Unknown error.")
}

// ResolvePrincipalFromSubject 负责把 subject 字符串解析为统一主体信息。
// 这个函数拆出来是为了让测试和业务代码都能直接验证 subject 语义，
// 不必每次都先伪造一整个 JWT。
func ResolvePrincipalFromSubject(subject string) (AccessPrincipal, error) {
	if strings.HasPrefix(subject, casSubjectPrefix) {
		casUsername := strings.TrimPrefix(subject, casSubjectPrefix)
		if casUsername == "" {
			return AccessPrincipal{}, errors.New("cas subject is empty")
		}

		return AccessPrincipal{
			Subject:     subject,
			CASUsername: casUsername,
		}, nil
	}

	localUserID, err := strconv.ParseUint(subject, 10, 64)
	if err != nil {
		return AccessPrincipal{}, err
	}

	return AccessPrincipal{
		Subject:     subject,
		LocalUserID: localUserID,
	}, nil
}

// ResolvePrincipalFromToken 会把 token subject 解析为统一主体信息。
// 数字 subject 视为原有本地用户，cas: 前缀视为 CAS 用户。
func ResolvePrincipalFromToken(token string) (AccessPrincipal, error) {
	subject, err := ResolveAccessToken(token)
	if err != nil {
		return AccessPrincipal{}, err
	}
	return ResolvePrincipalFromSubject(subject)
}

// BuildCASUserInfo 会为 CAS 主体构造一个最小用户信息响应。
// CAS 身份与本地 users 表脱离，因此这里只返回独立的 CAS 用户资料视图。
func BuildCASUserInfo(casUsername string) *model.UserInfo {
	return &model.UserInfo{
		Username: casUsername,
		Info:     "cas authenticated user",
	}
}

func ParseRequest(c *gin.Context) (AccessPrincipal, error) {
	token := c.GetHeader("token")
	if token == "" {
		return AccessPrincipal{}, errors.New("token is required")
	}
	return ResolvePrincipalFromToken(token)
}
