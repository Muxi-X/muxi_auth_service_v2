package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestOAuthStoreRequiresAdminToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := Load(gin.New())

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/api/oauth/store", strings.NewReader(`{"domain":"https://client.example.com"}`))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", w.Code, w.Body.String())
	}
}
