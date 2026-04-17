package soap

import (
	"context"
	"fmt"
	"testing"
)

// mockCaller is a test double for soap.Caller.
type mockCaller struct {
	response []byte
	err      error
	// captured call info
	method string
	params []Param
}

func (m *mockCaller) Call(_ context.Context, method string, params []Param) ([]byte, error) {
	m.method = method
	m.params = params
	return m.response, m.err
}

// successXML returns a minimal valid SOAP body for the given method.
func successXML(method string, value string) []byte {
	return []byte(`<` + method + `Response><` + method + `Return>` + value + `</` + method + `Return></` + method + `Response>`)
}

var badXMLErr = fmt.Errorf("connection refused")

// --- GetAccessStatus ---

func TestGetAccessStatus_Success(t *testing.T) {
	mock := &mockCaller{response: successXML("getAccessStatus", "1")}

	status, err := GetAccessStatus(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if status != 1 {
		t.Errorf("status = %d, want 1", status)
	}
	if mock.method != "getAccessStatus" {
		t.Errorf("method = %q, want %q", mock.method, "getAccessStatus")
	}
}

func TestGetAccessStatus_CallError(t *testing.T) {
	mock := &mockCaller{err: badXMLErr}
	_, err := GetAccessStatus(context.Background(), mock)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetAccessStatus_BadXML(t *testing.T) {
	mock := &mockCaller{response: []byte(`not valid xml`)}
	_, err := GetAccessStatus(context.Background(), mock)
	if err == nil {
		t.Fatal("expected error for bad XML")
	}
}

// --- GetCurrentRevision ---

func TestGetCurrentRevision_Success(t *testing.T) {
	mock := &mockCaller{response: successXML("getCurrentRevision", "99999")}

	rev, err := GetCurrentRevision(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if rev != 99999 {
		t.Errorf("revision = %d, want 99999", rev)
	}
	if mock.method != "getCurrentRevision" {
		t.Errorf("method = %q, want %q", mock.method, "getCurrentRevision")
	}
}

func TestGetCurrentRevision_CallError(t *testing.T) {
	mock := &mockCaller{err: badXMLErr}
	_, err := GetCurrentRevision(context.Background(), mock)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetCurrentRevision_BadXML(t *testing.T) {
	mock := &mockCaller{response: []byte(`not valid xml`)}
	_, err := GetCurrentRevision(context.Background(), mock)
	if err == nil {
		t.Fatal("expected error for bad XML")
	}
}

// --- GetExpireDate ---

func TestGetExpireDate_Success(t *testing.T) {
	mock := &mockCaller{response: successXML("getExpireDate", "2037-01-02")}

	date, err := GetExpireDate(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if date != "2037-01-02" {
		t.Errorf("date = %q, want %q", date, "2037-01-02")
	}
	if mock.method != "getExpireDate" {
		t.Errorf("method = %q, want %q", mock.method, "getExpireDate")
	}
}

func TestGetExpireDate_CallError(t *testing.T) {
	mock := &mockCaller{err: badXMLErr}
	_, err := GetExpireDate(context.Background(), mock)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetExpireDate_BadXML(t *testing.T) {
	mock := &mockCaller{response: []byte(`not valid xml`)}
	_, err := GetExpireDate(context.Background(), mock)
	if err == nil {
		t.Fatal("expected error for bad XML")
	}
}

// --- GetSubscriptionStatus ---

func TestGetSubscriptionStatus_Success(t *testing.T) {
	mock := &mockCaller{response: successXML("getSubscriptionStatus", "1")}

	status, err := GetSubscriptionStatus(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if status != "1" {
		t.Errorf("status = %q, want %q", status, "1")
	}
	if mock.method != "getSubscriptionStatus" {
		t.Errorf("method = %q, want %q", mock.method, "getSubscriptionStatus")
	}
}

func TestGetSubscriptionStatus_CallError(t *testing.T) {
	mock := &mockCaller{err: badXMLErr}
	_, err := GetSubscriptionStatus(context.Background(), mock)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetSubscriptionStatus_BadXML(t *testing.T) {
	mock := &mockCaller{response: []byte(`not valid xml`)}
	_, err := GetSubscriptionStatus(context.Background(), mock)
	if err == nil {
		t.Fatal("expected error for bad XML")
	}
}

// --- GetRightAccess ---

func TestGetRightAccess_Success(t *testing.T) {
	mock := &mockCaller{response: successXML("getRightAccess", "0")}

	access, err := GetRightAccess(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if access != "0" {
		t.Errorf("access = %q, want %q", access, "0")
	}
	if mock.method != "getRightAccess" {
		t.Errorf("method = %q, want %q", mock.method, "getRightAccess")
	}
}

func TestGetRightAccess_CallError(t *testing.T) {
	mock := &mockCaller{err: badXMLErr}
	_, err := GetRightAccess(context.Background(), mock)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetRightAccess_BadXML(t *testing.T) {
	mock := &mockCaller{response: []byte(`not valid xml`)}
	_, err := GetRightAccess(context.Background(), mock)
	if err == nil {
		t.Fatal("expected error for bad XML")
	}
}

// --- GetUserIdByLogin ---

func TestGetUserIdByLogin_Success(t *testing.T) {
	mock := &mockCaller{response: successXML("getUserIdByLogin", "1000000000539")}

	uid, err := GetUserIdByLogin(context.Background(), mock)
	if err != nil {
		t.Fatal(err)
	}
	if uid != "1000000000539" {
		t.Errorf("uid = %q, want %q", uid, "1000000000539")
	}
	if mock.method != "getUserIdByLogin" {
		t.Errorf("method = %q, want %q", mock.method, "getUserIdByLogin")
	}
}

func TestGetUserIdByLogin_CallError(t *testing.T) {
	mock := &mockCaller{err: badXMLErr}
	_, err := GetUserIdByLogin(context.Background(), mock)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetUserIdByLogin_BadXML(t *testing.T) {
	mock := &mockCaller{response: []byte(`not valid xml`)}
	_, err := GetUserIdByLogin(context.Background(), mock)
	if err == nil {
		t.Fatal("expected error for bad XML")
	}
}
