package auth

import (
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	payload := TokenPayload{
		ID:     1,
		Role:   1,
		Expire: time.Hour,
	}

	tokenString, err := GenerateToken(payload)
	if err != nil {
		t.Error(err)
		return
	}
	resolve, err := ResolveToken(tokenString)
	if err != nil {
		t.Error(err)
		return
	}
	if resolve.ID != payload.ID || resolve.Role != payload.Role {
		t.Errorf("decoded token not equal to previous")
		return
	}
}
