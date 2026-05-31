package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	pkgoauth "github.com/Muxi-X/muxi_auth_service_v2/pkg/oauth"
	"github.com/spf13/viper"
	cas "gopkg.in/cas.v2"
	oauth2 "gopkg.in/oauth2.v4"
	oauthserver "gopkg.in/oauth2.v4/server"
)

var (
	// ErrMissingCASTicket 表示 CAS callback 没有携带 ticket。
	ErrMissingCASTicket = errors.New("missing cas ticket")
	// ErrMissingOAuthClientID 表示 CAS callback 没有说明要给哪个 OAuth 客户端发码。
	ErrMissingOAuthClientID = errors.New("missing oauth client id")
	// ErrMissingCallbackURL 表示 callback 缺失最终要跳回的业务地址。
	ErrMissingCallbackURL = errors.New("missing callback url")
	// ErrInvalidCASTicket 表示 CAS 票据校验失败，常见于无效、过期或重复消费 ticket。
	ErrInvalidCASTicket = errors.New("invalid cas ticket")
)

// CASTicketValidator 抽象出 CAS 验票能力，便于测试替换真实网络交互。
type CASTicketValidator interface {
	ValidateTicket(serviceURL *url.URL, ticket string) (*cas.AuthenticationResponse, error)
}

// AuthorizeCodeGenerator 负责复用现有 OAuth 授权码发放逻辑。
type AuthorizeCodeGenerator interface {
	GenerateAuthorizeCode(ctx context.Context, request AuthorizeCodeRequest) (AuthorizeCodeResult, error)
}

// AuthorizeCodeRequest 描述生成 auth code 所需的最小输入。
type AuthorizeCodeRequest struct {
	ClientID       string
	UserID         string
	CallbackURL    string
	AccessTokenExp time.Duration
}

// AuthorizeCodeResult 表示授权码生成结果。
type AuthorizeCodeResult struct {
	Code      string
	ExpiresIn time.Duration
}

// CASCallbackResult 表示 callback 成功处理后的对外结果。
type CASCallbackResult struct {
	Code        string
	RedirectURL string
}

// CASOAuthFlow 负责把“CAS 登录成功”翻译成“OAuth auth code 已签发”。
// 这层只负责编排，不关心 Gin handler 如何组织响应。
type CASOAuthFlow struct {
	callBackBaseURL string //这个其实是当前服务的域名,为什么需要这个呢?因为没办法直接从请求中获取到域名,但是cas验证是需要校验整个service的完整性的,所以需要域名
	ticketValidator CASTicketValidator
	codeGenerator   AuthorizeCodeGenerator
}

// NewCASOAuthFlow 创建一条可测试、可替换依赖的 CAS -> OAuth 回调流程。
func NewCASOAuthFlow(
	callBackBaseURL string,
	ticketValidator CASTicketValidator,
	codeGenerator AuthorizeCodeGenerator,
) *CASOAuthFlow {
	return &CASOAuthFlow{
		callBackBaseURL: strings.TrimRight(callBackBaseURL, "/"),
		ticketValidator: ticketValidator,
		codeGenerator:   codeGenerator,
	}
}

// HandleCallback 会完成以下事情：
// 1. 读取 callback 参数
// 2. 还原出 CAS service URL 并验票
// 3. 把 CAS 用户映射为本地用户
// 4. 复用现有 OAuth server 生成授权码
// 5. 组装最终回跳地址
func (f *CASOAuthFlow) HandleCallback(ctx context.Context, request *http.Request) (*CASCallbackResult, error) {
	// 获取签发的 ticket
	ticket := strings.TrimSpace(request.URL.Query().Get("ticket"))
	if ticket == "" {
		return nil, ErrMissingCASTicket
	}
	// 获取用户传输的 client_id
	clientID := strings.TrimSpace(request.URL.Query().Get("client_id"))
	if clientID == "" {
		return nil, ErrMissingOAuthClientID
	}
	// 获取用户传输的重定向 url
	callbackURL := strings.TrimSpace(request.URL.Query().Get("callback_url"))
	if callbackURL == "" {
		return nil, ErrMissingCallbackURL
	}
	// 构建完整的登陆url,为什么要携带这个呢?
	serviceURL, err := f.buildServiceURL(request)
	if err != nil {
		return nil, err
	}

	// 验证当前 ticket 的合法性
	authenticationResponse, err := f.ticketValidator.ValidateTicket(serviceURL, ticket)
	if err != nil || authenticationResponse == nil || authenticationResponse.User == "" {
		return nil, fmt.Errorf("%w: %v", ErrInvalidCASTicket, err)
	}
	//从请求中去获取 code 的过期时间
	accessTokenExp, err := parseOptionalTokenExp(request.URL.Query().Get("token_exp"))
	if err != nil {
		return nil, err
	}

	// 这里直接把 CAS 用户名编码为独立 subject，
	// 不再尝试映射回本地 users.id，从而保证 CAS 身份和原有本地用户体系完全解耦。
	casSubject := pkgoauth.BuildCASSubject(authenticationResponse.User)

	// 签发code
	codeResult, err := f.codeGenerator.GenerateAuthorizeCode(ctx, AuthorizeCodeRequest{
		ClientID:       clientID,
		UserID:         casSubject,
		CallbackURL:    callbackURL,
		AccessTokenExp: accessTokenExp,
	})
	if err != nil {
		return nil, err
	}

	redirectURL, err := appendCodeToCallbackURL(callbackURL, codeResult.Code)
	if err != nil {
		return nil, err
	}

	return &CASCallbackResult{
		Code:        codeResult.Code,
		RedirectURL: redirectURL,
	}, nil
}

// buildServiceURL 会构造出与 CAS 当初签发 ticket 时一致的 service URL。
// 注意这里必须剔除 ticket 本身，否则 serviceValidate 时会因为 service 不一致而失败。
func (f *CASOAuthFlow) buildServiceURL(request *http.Request) (*url.URL, error) {
	// 1. 获取最原始、未经 Go 库重新拼装的完整 URL 字符串
	// request.RequestURI 包含了最原始的路径和参数
	rawURI := request.RequestURI

	// 2. 从原始 query 中移除 ticket，同时保留其他参数的原始顺序和编码。
	// CAS 回调会把 ticket 放到 query 的任意位置，不能简单截断字符串。
	parsedClean, err := url.ParseRequestURI(rawURI)
	if err != nil {
		return nil, err
	}
	parsedClean.RawQuery = removeQueryParam(parsedClean.RawQuery, "ticket")

	// 3. 构建返回的 URL 对象
	// 依然保留你之前的 Scheme 和 Host 处理逻辑，确保验证请求发往正确的地址
	var scheme string
	if request.TLS != nil || request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	} else {
		scheme = "http"
	}

	host := request.Host
	if f.callBackBaseURL != "" {
		u, _ := url.Parse(f.callBackBaseURL)
		scheme = u.Scheme
		host = u.Host
	}

	return &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     parsedClean.Path,
		RawQuery: parsedClean.RawQuery, // 这里拿到的就是完全没被重排过的原始 Query 字符串
	}, nil
}

func removeQueryParam(rawQuery, key string) string {
	if rawQuery == "" {
		return ""
	}

	parts := strings.Split(rawQuery, "&")
	filtered := parts[:0]
	for _, part := range parts {
		if part == "" {
			continue
		}

		paramKey := part
		if idx := strings.Index(part, "="); idx != -1 {
			paramKey = part[:idx]
		}
		if paramKey == key {
			continue
		}
		filtered = append(filtered, part)
	}
	return strings.Join(filtered, "&")
}

// parseOptionalTokenExp 用来兼容现有 OAuth 接口的 token_exp 参数语义。
// 如果没有传值，就返回 0，表示沿用服务当前默认的 access token 过期时间。
func parseOptionalTokenExp(raw string) (time.Duration, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}

	expireSeconds, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid token_exp: %w", err)
	}

	if expireSeconds <= 0 {
		return 0, nil
	}

	return time.Duration(expireSeconds) * time.Second, nil
}

// appendCodeToCallbackURL 会把 code 安全地拼接回 callback_url。
// 这里显式使用 URL API，是为了兼容 callback_url 本身已有查询参数的情况。
func appendCodeToCallbackURL(callbackURL, code string) (string, error) {
	parsedURL, err := url.Parse(callbackURL)
	if err != nil {
		return "", err
	}

	query := parsedURL.Query()
	query.Set("code", code)
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

// defaultCASTicketValidator 是生产环境使用的 CAS 验票实现。
type defaultCASTicketValidator struct {
	validator *cas.ServiceTicketValidator
}

func (d *defaultCASTicketValidator) ValidateTicket(serviceURL *url.URL, ticket string) (*cas.AuthenticationResponse, error) {
	return d.validator.ValidateTicket(serviceURL, ticket)
}

// defaultAuthorizeCodeGenerator 复用现有 OAuth server 的发码逻辑，
// 避免新 CAS 流生成出来的 code 脱离原有 token 存储与过期策略。
type defaultAuthorizeCodeGenerator struct{}

func (g *defaultAuthorizeCodeGenerator) GenerateAuthorizeCode(ctx context.Context, request AuthorizeCodeRequest) (AuthorizeCodeResult, error) {
	authorizeRequest := &oauthserver.AuthorizeRequest{
		ResponseType:   oauth2.Code,
		ClientID:       request.ClientID,
		RedirectURI:    request.CallbackURL,
		UserID:         request.UserID,
		AccessTokenExp: request.AccessTokenExp,
	}

	tokenInfo, err := pkgoauth.OauthServer.Server.GetAuthorizeToken(ctx, authorizeRequest)
	if err != nil {
		return AuthorizeCodeResult{}, err
	}

	return AuthorizeCodeResult{
		Code:      tokenInfo.GetCode(),
		ExpiresIn: tokenInfo.GetCodeExpiresIn(),
	}, nil
}

// NewDefaultCASOAuthFlow 会基于配置创建生产环境使用的默认流程实现。
func NewDefaultCASOAuthFlow() (*CASOAuthFlow, error) {
	casServerURL, err := url.Parse(strings.TrimSpace(viper.GetString("cas.server_url")))
	if err != nil {
		return nil, err
	}

	if casServerURL.String() == "" {
		return nil, errors.New("cas.server_url is required")
	}

	validator := cas.NewServiceTicketValidator(http.DefaultClient, casServerURL)

	return NewCASOAuthFlow(
		viper.GetString("cas.callback_base_url"),
		&defaultCASTicketValidator{validator: validator},
		&defaultAuthorizeCodeGenerator{},
	), nil
}
