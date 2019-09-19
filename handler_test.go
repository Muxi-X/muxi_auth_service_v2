package main

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/errno"
	"github.com/Muxi-X/muxi_auth_service_v2/handler/signin"
	_ "github.com/Muxi-X/muxi_auth_service_v2/handler/signup"
	"github.com/Muxi-X/muxi_auth_service_v2/router"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/config"
	"github.com/spf13/pflag"

	"os"
	"bytes"
	"encoding/base64"
	"testing"
	"net/http"
	"encoding/json"
	"net/http/httptest"
)

var (
	testRouter *gin.Engine
	token      string
	captcha    string
)

func TestMain(m *testing.M) {
	pflag.Parse()

	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	model.DB.Init()
	defer model.DB.Close()

	testRouter = gin.Default()
	router.Load(testRouter)

	os.Exit(m.Run())
}

/* func Test_B_SignUp(t *testing.T) {
    w := httptest.NewRecorder()
    signupMock := signup.UserSignupRequestData {
		Username: "testMockUser2",
		Email: "testUser2@mock.com",
		Password: base64.StdEncoding.EncodeToString([]byte("testMockPassword2")),
	}
    jsonify, _ := json.Marshal(signupMock)
    req, _ := http.NewRequest("POST", "/auth/api/signup", bytes.NewBuffer(jsonify))

    testRouter.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
} */

func sendTestRequest(method, path string, data interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	jsonify, _ := json.Marshal(data)
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonify))
	testRouter.ServeHTTP(w, req)
	return w
}

func getValueFromResponse(t *testing.T, bytes []byte, data interface{}, key string) interface{} {
	err := json.Unmarshal(bytes, &data)
	assert.Equal(t, err, nil)

	return data.(map[string]interface{})["data"].(map[string]interface{})[key]
}

func getCodeFromError(t *testing.T, bytes []byte) int {
	responseErr := errno.Err{}
	err := json.Unmarshal(bytes, &responseErr)
	assert.Equal(t, err, nil)

	return responseErr.Code
}

func Test_C_SignIn(t *testing.T) {
	signinMock := signin.UserSigninRequestData{
		Username: "testMockUser2",
		Password: base64.StdEncoding.EncodeToString([]byte("testMockPassword2")),
	}
	w := sendTestRequest("POST", "/auth/api/signin", signinMock)
	assert.Equal(t, 200, w.Code)

	var data handler.Response
	token = getValueFromResponse(t, w.Body.Bytes(), data, "token").(string)
}

func Test_D_A_CheckEmail(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/api/check_email", nil)
	query := req.URL.Query()
	query.Add("email", "testUser2@mock.com")
	req.URL.RawQuery = query.Encode()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
}

func Test_D_B_CheckEmail(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/api/check_email", nil)
	query := req.URL.Query()
	query.Add("email", "testUser2@mmm.com")
	req.URL.RawQuery = query.Encode()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_E_A_CheckUsername(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/api/check_name", nil)
	query := req.URL.Query()
	query.Add("username", "testMockUser2")
	req.URL.RawQuery = query.Encode()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var data struct{
		Code    int      `json:"code"`
		Message string   `json:"message"`
		Data    bool     `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &data)
	assert.Equal(t, err, nil)
	assert.Equal(t, true, data.Data)
}

func Test_E_B_CheckUsername(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/api/check_name", nil)
	query := req.URL.Query()
	query.Add("username", "testMmmUser2")
	req.URL.RawQuery = query.Encode()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_F_CheckToken(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/api/check_token", nil)
	query := req.URL.Query()
	query.Add("email", "testUser2@mock.com")
	query.Add("token", token)
	req.URL.RawQuery = query.Encode()

    testRouter.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_G_A_GetEmailByName(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/api/email", nil)
	query := req.URL.Query()
	query.Add("username", "testMockUser2")
	req.URL.RawQuery = query.Encode()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_G_B_GetEmailByName(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/api/email", nil)
	query := req.URL.Query()
	query.Add("username", "testMmmUser2")
	req.URL.RawQuery = query.Encode()

	testRouter.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var data struct{
		Code    int           `json:"code"`
		Message string        `json:"message"`
	}

	err := json.Unmarshal(w.Body.Bytes(), &data)
	assert.Equal(t, err, nil)
	assert.Equal(t, 20102, data.Code)
}