package util

import (
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/constvar"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/teris-io/shortid"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func GenShortId() (string, error) {
	return shortid.Generate()
}

func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}
	if requestID, ok := v.(string); ok {
		return requestID
	}
	return ""
}

func SendTestRequest(method, path string, data interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	jsonify, _ := json.Marshal(data)
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonify))
	constvar.TestRouter.ServeHTTP(w, req)
	return w
}

func GetValueFromResponse(t *testing.T, bytes []byte, data interface{}, key string) interface{} {
	err := json.Unmarshal(bytes, &data)
	assert.Equal(t, err, nil)
	if key == "" {
		return data.(map[string]interface{})["data"]
	}

	return data.(map[string]interface{})["data"].(map[string]interface{})[key]
}

func GetCodeFromError(t *testing.T, bytes []byte) int {
	responseErr := errno.Err{}
	err := json.Unmarshal(bytes, &responseErr)
	assert.Equal(t, err, nil)

	return responseErr.Code
}

func CheckEmailAddress(email string) (b bool) {
	if m, _ := regexp.MatchString("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+", email); !m {
		return false
	}
	return true
}
