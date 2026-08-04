package middleware

import (
	"testing"

	"github.com/Muxi-X/muxi_auth_service_v2/model"
)

func TestCanManageOAuthClientsAllowsAdministratorRole(t *testing.T) {
	if !canManageOAuthClients(&model.UserModel{RoleID: 2}) {
		t.Fatalf("expected administrator role to manage oauth clients")
	}
}

func TestCanManageOAuthClientsRejectsNilUser(t *testing.T) {
	if canManageOAuthClients(nil) {
		t.Fatalf("expected nil user to be rejected")
	}
}
