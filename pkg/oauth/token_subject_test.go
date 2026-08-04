package oauth

import "testing"

// TestBuildCASSubject 确保历史 CAS subject 仍保持 cas: 前缀，
// 避免兼容旧 token 时和本地数值 user_id 混淆。
func TestBuildCASSubject(t *testing.T) {
	if subject := BuildCASSubject("alice"); subject != "cas:alice" {
		t.Fatalf("expected cas:alice, got %s", subject)
	}
}

// TestResolvePrincipalFromTokenLocalAndCAS 覆盖主体解析的两条主路径：
// 1. 原有本地数值 user id
// 2. 历史 cas:username 外部身份
func TestResolvePrincipalFromTokenLocalAndCAS(t *testing.T) {
	localPrincipal, err := ResolvePrincipalFromSubject("42")
	if err != nil {
		t.Fatalf("ResolvePrincipalFromSubject(local) returned error: %v", err)
	}
	if localPrincipal.LocalUserID != 42 {
		t.Fatalf("expected local user id 42, got %d", localPrincipal.LocalUserID)
	}
	if localPrincipal.CASUsername != "" {
		t.Fatalf("expected empty cas username, got %s", localPrincipal.CASUsername)
	}

	casPrincipal, err := ResolvePrincipalFromSubject("cas:alice")
	if err != nil {
		t.Fatalf("ResolvePrincipalFromSubject(cas) returned error: %v", err)
	}
	if casPrincipal.CASUsername != "alice" {
		t.Fatalf("expected cas username alice, got %s", casPrincipal.CASUsername)
	}
	if casPrincipal.LocalUserID != 0 {
		t.Fatalf("expected empty local user id, got %d", casPrincipal.LocalUserID)
	}
}
