package auth

import "testing"

func TestAuth(t *testing.T) {
	password := "passwordTest"
	encrypted, err := Encrypt(password)
	if err != nil {
		t.Error(err)
		return
	}
	if err = Compare(encrypted, password); err != nil {
		t.Error(err)
		return
	}
}
