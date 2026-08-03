package service

import "testing"

func TestNormalizeOAuthClientDomainAcceptsHTTPSOrigin(t *testing.T) {
	got, err := NormalizeOAuthClientDomain(" HTTPS://Pass.MuxiXYZ.com ")
	if err != nil {
		t.Fatalf("NormalizeOAuthClientDomain() returned error: %v", err)
	}
	if got != "https://pass.muxixyz.com" {
		t.Fatalf("expected normalized origin, got %s", got)
	}
}

func TestOAuthClientDomainLookupKey(t *testing.T) {
	got, err := OAuthClientDomainLookupKey("https://Pass.MuxiXYZ.com:8443")
	if err != nil {
		t.Fatalf("OAuthClientDomainLookupKey() returned error: %v", err)
	}
	if got != "pass.muxixyz.com:8443" {
		t.Fatalf("expected lookup key pass.muxixyz.com:8443, got %s", got)
	}
}

func TestNormalizeOAuthClientDomainRejectsUnsafeValues(t *testing.T) {
	tests := []string{
		"http://pass.muxixyz.com",
		"https://localhost:8080",
		"https://127.0.0.1",
		"https://*.muxixyz.com",
		"https://pass.muxixyz.com/callback",
		"https://pass.muxixyz.com?debug=true",
	}

	for _, tt := range tests {
		if _, err := NormalizeOAuthClientDomain(tt); err == nil {
			t.Fatalf("expected %q to be rejected", tt)
		}
	}
}

func TestValidateOAuthCallbackURLForDomain(t *testing.T) {
	if err := ValidateOAuthCallbackURLForDomain("https://pass.muxixyz.com", "https://pass.muxixyz.com/auth/callback?state=ok"); err != nil {
		t.Fatalf("expected callback url to be allowed: %v", err)
	}

	tests := []string{
		"http://pass.muxixyz.com/auth/callback",
		"https://evil.muxixyz.com/auth/callback",
		"https://pass.muxixyz.com.evil.example/auth/callback",
		"https://pass.muxixyz.com/auth/callback#token",
	}

	for _, tt := range tests {
		if err := ValidateOAuthCallbackURLForDomain("https://pass.muxixyz.com", tt); err == nil {
			t.Fatalf("expected callback %q to be rejected", tt)
		}
	}
}
