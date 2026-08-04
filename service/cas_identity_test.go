package service

import "testing"

func TestBuildCASLocalUsername(t *testing.T) {
	got := BuildCASLocalUsername(" Alice.Z ")
	if got != "cas_alice_z_50f08eb7" {
		t.Fatalf("unexpected local username: %s", got)
	}
}

func TestBuildCASLocalUsernameHandlesNonASCII(t *testing.T) {
	got := BuildCASLocalUsername("木犀")
	if got != "cas_user_05745470" {
		t.Fatalf("unexpected local username: %s", got)
	}
}
