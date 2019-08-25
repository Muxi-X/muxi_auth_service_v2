package auth

import (
	"encoding/json"
	"time"
)

/*
  001 -User       011 -User or Company       111 -All Role
  010 -Company    110 -Company or Reviewer
  100 -Reviewer   101 -Reviewer or User
  用户生成时，不会生成 011 110 101 111 这些角色，这些角色只是用在LoginRequried的limit中来限制。
*/

type Payload struct {
	Data TokenResolve `json:"data"`
}

type TokenPayload struct {
	ID     uint64        `json:"id"`
	Role   int           `json:"type"`
	Expire time.Duration `json:"expire"`
}

type TokenResolve struct {
	ID      uint64 `json:"id"`
	Role    int    `json:"role"`
	Expired int64  `json:"expired"`
}

func GenerateToken(payload TokenPayload) (string, error) {
	expired := int64(time.Now().Unix()) + int64(payload.Expire.Seconds())
	//	fmt.Println(expired)
	newToken := TokenResolve{
		payload.ID,
		payload.Role,
		expired,
	}
	js, err := json.Marshal(map[string]TokenResolve{"data": newToken})
	if err != nil {
		return "", err
	}
	aesString, err := AesEncrypt(string(js))
	if err != nil {
		return "", err
	}
	return aesString, nil
}

func ResolveToken(aesString string) (TokenResolve, error) {
	res := &Payload{}
	js, err := AesDecrypt(aesString)
	if err != nil {
		return TokenResolve{}, err
	}
	err = json.Unmarshal([]byte(js), res)
	if err != nil {
		return TokenResolve{}, err
	}
	return (*res).Data, nil
}
