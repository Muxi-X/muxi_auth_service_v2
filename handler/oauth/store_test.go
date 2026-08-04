package oauth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
)

func TestStoreRejectsUnguardedRegistration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/api/oauth/store", Store)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/api/oauth/store", strings.NewReader(`{"domain":"https://attacker.example.com"}`))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d: %s", w.Code, w.Body.String())
	}

	var response handler.Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}

	if response.Code != errno.ErrOAuthClientRegistrationDisabled.Code {
		t.Fatalf("expected error code %d, got %d", errno.ErrOAuthClientRegistrationDisabled.Code, response.Code)
	}
	if response.Data != nil {
		t.Fatalf("expected no client credentials in response, got %#v", response.Data)
	}
}
