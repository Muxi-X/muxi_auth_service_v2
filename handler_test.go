package main

import (
	"github.com/Muxi-X/muxi_auth_service_v2/handler"
	"github.com/Muxi-X/muxi_auth_service_v2/handler/signin"
	"github.com/Muxi-X/muxi_auth_service_v2/handler/signup"
	"github.com/Muxi-X/muxi_auth_service_v2/pkg/constvar"
	"github.com/Muxi-X/muxi_auth_service_v2/util"
	"github.com/stretchr/testify/assert"

	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_B_SignUp(t *testing.T) {
	signupMock := signup.UserSignupRequestData{
		Username: "testMockUser2",
		Email:    "testUser2@mock.com",
		Password: base64.StdEncoding.EncodeToString([]byte("testMockPassword2")),
	}
	w := util.SendTestRequest("POST", "/auth/api/signup", signupMock)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, 0, util.GetCodeFromError(t, w.Body.Bytes()))
}

func Test_C_A_SignIn(t *testing.T) {
	signinMock := signin.UserSigninRequestData{
		Username: "testMockUser2",
		Password: base64.StdEncoding.EncodeToString([]byte("testMockPassword2")),
	}
	w := util.SendTestRequest("POST", "/auth/api/signin", signinMock)
	assert.Equal(t, 0, util.GetCodeFromError(t, w.Body.Bytes()))

	var data handler.Response
	constvar.Token = util.GetValueFromResponse(t, w.Body.Bytes(), data, "token").(string)
}

func Test_C_B_Signin(t *testing.T) {
	signinMock := signin.UserSigninRequestData{
		Username: "testMockUser2",
		Password: base64.StdEncoding.EncodeToString([]byte("testMoooooockPassword2")),
	}
	w := util.SendTestRequest("POST", "/auth/api/signin", signinMock)
	assert.Equal(t, 20301, util.GetCodeFromError(t, w.Body.Bytes()))
}

func Test_C_C_Signin(t *testing.T) {
	signinMock := signin.UserSigninRequestData{
		Username: "testMo00000ckUser2",
		Password: base64.StdEncoding.EncodeToString([]byte("testMoooooockPassword2")),
	}
	w := util.SendTestRequest("POST", "/auth/api/signin", signinMock)
	assert.Equal(t, 20102, util.GetCodeFromError(t, w.Body.Bytes()))
}

func Test_D_A_CheckEmail(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/api/check_email", nil)
	query := req.URL.Query()
	query.Add("email", "testUser2@mock.com")
	req.URL.RawQuery = query.Encode()

	constvar.TestRouter.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var data handler.Response
	result := util.GetValueFromResponse(t, w.Body.Bytes(), data, "").(bool)
	assert.Equal(t, true, result)
}

func Test_D_B_CheckEmail(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/api/check_email", nil)
	query := req.URL.Query()
	query.Add("email", "testUser2@mmm.com")
	req.URL.RawQuery = query.Encode()

	constvar.TestRouter.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var data handler.Response
	result := util.GetValueFromResponse(t, w.Body.Bytes(), data, "").(bool)
	assert.Equal(t, false, result)
}

func Test_E_A_CheckUsername(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/api/check_name", nil)
	query := req.URL.Query()
	query.Add("username", "testMockUser2")
	req.URL.RawQuery = query.Encode()
	constvar.TestRouter.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var data handler.Response
	result := util.GetValueFromResponse(t, w.Body.Bytes(), data, "").(bool)
	assert.Equal(t, true, result)
}

func Test_E_B_CheckUsername(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/api/check_name", nil)
	query := req.URL.Query()
	query.Add("username", "testMmmUser2")
	req.URL.RawQuery = query.Encode()

	constvar.TestRouter.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var data handler.Response
	result := util.GetValueFromResponse(t, w.Body.Bytes(), data, "").(bool)
	assert.Equal(t, false, result)
}

func Test_F_CheckToken(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/api/check_token", nil)
	query := req.URL.Query()
	query.Add("email", "testUser2@mock.com")
	query.Add("token", constvar.Token)
	req.URL.RawQuery = query.Encode()

	constvar.TestRouter.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_G_A_GetEmailByName(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/api/email", nil)
	query := req.URL.Query()
	query.Add("username", "testMockUser2")
	req.URL.RawQuery = query.Encode()

	constvar.TestRouter.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func Test_G_B_GetEmailByName(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/api/email", nil)
	query := req.URL.Query()
	query.Add("username", "testMmmUser2")
	req.URL.RawQuery = query.Encode()
	constvar.TestRouter.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	assert.Equal(t, 20102, util.GetCodeFromError(t, w.Body.Bytes()))
}
