package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	cas "gopkg.in/cas.v2"
)

// fakeCASTicketValidator 用于在测试里替换真实 CAS 服务，
// 让测试只聚焦在我们自己的回调编排逻辑上。
type fakeCASTicketValidator struct {
	validateFunc func(serviceURL *url.URL, ticket string) (*cas.AuthenticationResponse, error)
}

func (f fakeCASTicketValidator) ValidateTicket(serviceURL *url.URL, ticket string) (*cas.AuthenticationResponse, error) {
	return f.validateFunc(serviceURL, ticket)
}

// fakeAuthorizeCodeGenerator 用于隔离真正的 OAuth 发码逻辑，
// 这样我们可以精确断言 CAS 回调阶段到底把什么信息喂给了发码层。
type fakeAuthorizeCodeGenerator struct {
	generateFunc func(ctx context.Context, request AuthorizeCodeRequest) (AuthorizeCodeResult, error)
}

func (f fakeAuthorizeCodeGenerator) GenerateAuthorizeCode(ctx context.Context, request AuthorizeCodeRequest) (AuthorizeCodeResult, error) {
	return f.generateFunc(ctx, request)
}

// TestCASOAuthFlowHandleCallbackSuccess 验证整条 CAS 回调成功路径：
// 1. 正确使用当前 callback URL 验票
// 2. 找到本地用户映射
// 3. 复用 OAuth 层生成 auth code
// 4. 拼装最终的回跳地址
func TestCASOAuthFlowHandleCallbackSuccess(t *testing.T) {
	flow := NewCASOAuthFlow(
		"",
		fakeCASTicketValidator{
			validateFunc: func(serviceURL *url.URL, ticket string) (*cas.AuthenticationResponse, error) {
				if serviceURL.Scheme != "http" || serviceURL.Host != "oauth.example.com" || serviceURL.Path != "/auth/api/oauth/cas/callback" {
					t.Fatalf("unexpected service url: %s", serviceURL.String())
				}

				if serviceURL.Query().Get("client_id") != "client-a" {
					t.Fatalf("expected client id client-a, got %s", serviceURL.Query().Get("client_id"))
				}

				if serviceURL.Query().Get("callback_url") != "https://client.example.com/cb" {
					t.Fatalf("expected callback url https://client.example.com/cb, got %s", serviceURL.Query().Get("callback_url"))
				}

				if serviceURL.Query().Get("token_exp") != "120" {
					t.Fatalf("expected token_exp 120, got %s", serviceURL.Query().Get("token_exp"))
				}

				if ticket != "ST-1" {
					t.Fatalf("expected ticket ST-1, got %s", ticket)
				}

				return &cas.AuthenticationResponse{User: "casuser"}, nil
			},
		},
		fakeAuthorizeCodeGenerator{
			generateFunc: func(ctx context.Context, request AuthorizeCodeRequest) (AuthorizeCodeResult, error) {
				if request.ClientID != "client-a" {
					t.Fatalf("expected client id client-a, got %s", request.ClientID)
				}
				if request.UserID != "cas:casuser" {
					t.Fatalf("expected user id cas:casuser, got %s", request.UserID)
				}
				if request.CallbackURL != "https://client.example.com/cb" {
					t.Fatalf("expected callback url https://client.example.com/cb, got %s", request.CallbackURL)
				}
				if request.AccessTokenExp != 120*time.Second {
					t.Fatalf("expected access token exp 120s, got %s", request.AccessTokenExp)
				}

				return AuthorizeCodeResult{
					Code:      "auth-code-1",
					ExpiresIn: 30 * time.Minute,
				}, nil
			},
		},
	)

	req := httptest.NewRequest(http.MethodGet, "http://oauth.example.com/auth/api/oauth/cas/callback?ticket=ST-1&client_id=client-a&callback_url=https%3A%2F%2Fclient.example.com%2Fcb&token_exp=120", nil)

	result, err := flow.HandleCallback(context.Background(), req)
	if err != nil {
		t.Fatalf("HandleCallback() returned error: %v", err)
	}

	expectedRedirectURL := "https://client.example.com/cb?code=auth-code-1"
	if result.RedirectURL != expectedRedirectURL {
		t.Fatalf("expected redirect url %s, got %s", expectedRedirectURL, result.RedirectURL)
	}

	if result.Code != "auth-code-1" {
		t.Fatalf("expected code auth-code-1, got %s", result.Code)
	}
}

func TestCASOAuthFlowHandleCallbackKeepsServiceParamsWhenTicketComesFirst(t *testing.T) {
	flow := NewCASOAuthFlow(
		"",
		fakeCASTicketValidator{
			validateFunc: func(serviceURL *url.URL, ticket string) (*cas.AuthenticationResponse, error) {
				if ticket != "ST-1" {
					t.Fatalf("expected ticket ST-1, got %s", ticket)
				}
				if got := serviceURL.RawQuery; got != "callback_url=https%3A%2F%2Fclient.example.com%2Fcb&client_id=client-a&token_exp=120" {
					t.Fatalf("unexpected service raw query: %s", got)
				}
				if got := serviceURL.Query().Get("ticket"); got != "" {
					t.Fatalf("service url should not include ticket, got %s", got)
				}
				return &cas.AuthenticationResponse{User: "casuser"}, nil
			},
		},
		fakeAuthorizeCodeGenerator{
			generateFunc: func(ctx context.Context, request AuthorizeCodeRequest) (AuthorizeCodeResult, error) {
				return AuthorizeCodeResult{
					Code:      "auth-code-1",
					ExpiresIn: 30 * time.Minute,
				}, nil
			},
		},
	)

	req := httptest.NewRequest(http.MethodGet, "http://oauth.example.com/auth/api/oauth/cas/callback?ticket=ST-1&callback_url=https%3A%2F%2Fclient.example.com%2Fcb&client_id=client-a&token_exp=120", nil)

	result, err := flow.HandleCallback(context.Background(), req)
	if err != nil {
		t.Fatalf("HandleCallback() returned error: %v", err)
	}
	if result.RedirectURL != "https://client.example.com/cb?code=auth-code-1" {
		t.Fatalf("unexpected redirect url: %s", result.RedirectURL)
	}
}

// TestCASOAuthFlowHandleCallbackRejectsMissingTicket 确保 callback 参数不完整时会被明确拒绝，
// 避免服务继续执行到更深层逻辑后才报出难懂的错误。
func TestCASOAuthFlowHandleCallbackRejectsMissingTicket(t *testing.T) {
	flow := NewCASOAuthFlow("", fakeCASTicketValidator{}, fakeAuthorizeCodeGenerator{})
	req := httptest.NewRequest(http.MethodGet, "http://oauth.example.com/auth/api/oauth/cas/callback?client_id=client-a&callback_url=https%3A%2F%2Fclient.example.com%2Fcb", nil)

	_, err := flow.HandleCallback(context.Background(), req)
	if !errors.Is(err, ErrMissingCASTicket) {
		t.Fatalf("expected ErrMissingCASTicket, got %v", err)
	}
}

// TestCASOAuthFlowHandleCallbackRejectsInvalidTicket 确保 ticket 校验失败时，
// 不会继续进入用户映射和发码逻辑。
func TestCASOAuthFlowHandleCallbackRejectsInvalidTicket(t *testing.T) {
	flow := NewCASOAuthFlow(
		"",
		fakeCASTicketValidator{
			validateFunc: func(serviceURL *url.URL, ticket string) (*cas.AuthenticationResponse, error) {
				return nil, errors.New("invalid ticket")
			},
		},
		fakeAuthorizeCodeGenerator{},
	)

	req := httptest.NewRequest(http.MethodGet, "http://oauth.example.com/auth/api/oauth/cas/callback?ticket=ST-1&client_id=client-a&callback_url=https%3A%2F%2Fclient.example.com%2Fcb", nil)

	_, err := flow.HandleCallback(context.Background(), req)
	if !errors.Is(err, ErrInvalidCASTicket) {
		t.Fatalf("expected ErrInvalidCASTicket, got %v", err)
	}
}
