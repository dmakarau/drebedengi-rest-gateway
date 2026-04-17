package soap

import (
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

func (m *mockCaller) Call(method string, params []Param) ([]byte, error) {
	m.method = method
	m.params = params
	return m.response, m.err
}

func TestGetAccessStatus_Success(t *testing.T) {
	mock := &mockCaller{
		response: []byte(`<getAccessStatusResponse><getAccessStatusReturn>1</getAccessStatusReturn></getAccessStatusResponse>`),
	}

	status, err := GetAccessStatus(mock, "api", "login", "pass")
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
	mock := &mockCaller{err: fmt.Errorf("connection refused")}

	_, err := GetAccessStatus(mock, "api", "login", "pass")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetAccessStatus_BadXML(t *testing.T) {
	mock := &mockCaller{response: []byte(`not valid xml`)}

	_, err := GetAccessStatus(mock, "api", "login", "pass")
	if err == nil {
		t.Fatal("expected error for bad XML")
	}
}

func TestGetCurrentRevision_Success(t *testing.T) {
	mock := &mockCaller{
		response: []byte(`<getCurrentRevisionResponse><getCurrentRevisionReturn>99999</getCurrentRevisionReturn></getCurrentRevisionResponse>`),
	}

	rev, err := GetCurrentRevision(mock, "api", "login", "pass")
	if err != nil {
		t.Fatal(err)
	}
	if rev != 99999 {
		t.Errorf("revision = %d, want 99999", rev)
	}
}

func TestGetExpireDate_Success(t *testing.T) {
	mock := &mockCaller{
		response: []byte(`<getExpireDateResponse><getExpireDateReturn>2037-01-02</getExpireDateReturn></getExpireDateResponse>`),
	}

	date, err := GetExpireDate(mock, "api", "login", "pass")
	if err != nil {
		t.Fatal(err)
	}
	if date != "2037-01-02" {
		t.Errorf("date = %q, want %q", date, "2037-01-02")
	}
}

func TestGetSubscriptionStatus_Success(t *testing.T) {
	mock := &mockCaller{
		response: []byte(`<getSubscriptionStatusResponse><getSubscriptionStatusReturn>1</getSubscriptionStatusReturn></getSubscriptionStatusResponse>`),
	}

	status, err := GetSubscriptionStatus(mock, "api", "login", "pass")
	if err != nil {
		t.Fatal(err)
	}
	if status != "1" {
		t.Errorf("status = %q, want %q", status, "1")
	}
}

func TestGetRightAccess_Success(t *testing.T) {
	mock := &mockCaller{
		response: []byte(`<getRightAccessResponse><getRightAccessReturn>0</getRightAccessReturn></getRightAccessResponse>`),
	}

	access, err := GetRightAccess(mock, "api", "login", "pass")
	if err != nil {
		t.Fatal(err)
	}
	if access != "0" {
		t.Errorf("access = %q, want %q", access, "0")
	}
}

func TestGetUserIdByLogin_Success(t *testing.T) {
	mock := &mockCaller{
		response: []byte(`<getUserIdByLoginResponse><getUserIdByLoginReturn>1000000000539</getUserIdByLoginReturn></getUserIdByLoginResponse>`),
	}

	uid, err := GetUserIdByLogin(mock, "api", "login", "pass")
	if err != nil {
		t.Fatal(err)
	}
	if uid != "1000000000539" {
		t.Errorf("uid = %q, want %q", uid, "1000000000539")
	}
}

func TestAuthParams(t *testing.T) {
	params := authParams("myapi", "mylogin", "mypass")
	if len(params) != 3 {
		t.Fatalf("len = %d, want 3", len(params))
	}
	if params[0].Name != "apiId" || params[0].Value != "myapi" {
		t.Errorf("params[0] = %+v", params[0])
	}
	if params[1].Name != "login" || params[1].Value != "mylogin" {
		t.Errorf("params[1] = %+v", params[1])
	}
	if params[2].Name != "pass" || params[2].Value != "mypass" {
		t.Errorf("params[2] = %+v", params[2])
	}
}
