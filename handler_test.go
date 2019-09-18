package main

import (
	_ "fmt"
	"github.com/Muxi-X/muxi_auth_service_v2/handler/signin"
	_ "github.com/Muxi-X/muxi_auth_service_v2/handler/signup"
	"github.com/Muxi-X/muxi_auth_service_v2/router"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/Muxi-X/muxi_auth_service_v2/model"
	"github.com/Muxi-X/muxi_auth_service_v2/config"
	"github.com/spf13/pflag"

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
func closeDB() {
	model.DB.Close()
}

func Test_A_Build(t *testing.T) {
	pflag.Parse()

	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	model.DB.Init()

	testRouter = gin.Default()
	router.Load(testRouter)
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

func Test_C_SignIn(t *testing.T) {
	w := httptest.NewRecorder()
	singinMock := signin.UserSigninRequestData{
		Username: "testMockUser2",
		Password: base64.StdEncoding.EncodeToString([]byte("testMockPassword2")),
	}
	jsonify, _ := json.Marshal(singinMock)
	req, _ := http.NewRequest("POST", "/auth/api/signin", bytes.NewBuffer(jsonify))

	testRouter.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var data struct{
		Code    int                           `json:"code"`
		Message string                        `json:"message"`
		Data    signin.UserSigninResponseData `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &data)
	assert.Equal(t, err, nil)

	token = data.Data.Token
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

	assert.Equal(t, 500, w.Code)
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

	assert.Equal(t, 404, w.Code)
}

func Test_Z_CloseDB(t *testing.T) {
	closeDB()
}