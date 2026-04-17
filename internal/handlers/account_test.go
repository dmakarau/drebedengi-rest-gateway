package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drebedengi-rest/internal/soap"
)

// mockCaller is a test double for soap.Caller.
type mockCaller struct {
	response []byte
	err      error
}

func (m *mockCaller) Call(method string, params []soap.Param) ([]byte, error) {
	return m.response, m.err
}

func newTestHandler(response []byte, err error) *Handler {
	return &Handler{
		SOAP:     &mockCaller{response: response, err: err},
		APIId:    "test_api",
		Login:    "test@example.com",
		Password: "test",
	}
}

func TestGetAccessStatus_OK(t *testing.T) {
	h := newTestHandler(
		[]byte(`<getAccessStatusResponse><getAccessStatusReturn>1</getAccessStatusReturn></getAccessStatusResponse>`),
		nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/account/status", nil)
	w := httptest.NewRecorder()
	h.GetAccessStatus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var body map[string]int
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["status"] != 1 {
		t.Errorf("body status = %d, want 1", body["status"])
	}
}

func TestGetAccessStatus_SOAPError(t *testing.T) {
	h := newTestHandler(nil, fmt.Errorf("SOAP Fault [Client]: access denied"))

	req := httptest.NewRequest(http.MethodGet, "/account/status", nil)
	w := httptest.NewRecorder()
	h.GetAccessStatus(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadGateway)
	}

	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] == "" {
		t.Error("expected error message in body")
	}
}

func TestGetCurrentRevision_OK(t *testing.T) {
	h := newTestHandler(
		[]byte(`<getCurrentRevisionResponse><getCurrentRevisionReturn>42</getCurrentRevisionReturn></getCurrentRevisionResponse>`),
		nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/account/revision", nil)
	w := httptest.NewRecorder()
	h.GetCurrentRevision(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var body map[string]int
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["revision"] != 42 {
		t.Errorf("revision = %d, want 42", body["revision"])
	}
}

func TestGetExpireDate_OK(t *testing.T) {
	h := newTestHandler(
		[]byte(`<getExpireDateResponse><getExpireDateReturn>2037-01-02</getExpireDateReturn></getExpireDateResponse>`),
		nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/account/expire", nil)
	w := httptest.NewRecorder()
	h.GetExpireDate(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["expire_date"] != "2037-01-02" {
		t.Errorf("expire_date = %q, want %q", body["expire_date"], "2037-01-02")
	}
}

func TestGetSubscriptionStatus_OK(t *testing.T) {
	h := newTestHandler(
		[]byte(`<getSubscriptionStatusResponse><getSubscriptionStatusReturn>1</getSubscriptionStatusReturn></getSubscriptionStatusResponse>`),
		nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/account/subscription", nil)
	w := httptest.NewRecorder()
	h.GetSubscriptionStatus(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetRightAccess_OK(t *testing.T) {
	h := newTestHandler(
		[]byte(`<getRightAccessResponse><getRightAccessReturn>0</getRightAccessReturn></getRightAccessResponse>`),
		nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/account/access", nil)
	w := httptest.NewRecorder()
	h.GetRightAccess(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestGetUserIdByLogin_OK(t *testing.T) {
	h := newTestHandler(
		[]byte(`<getUserIdByLoginResponse><getUserIdByLoginReturn>12345</getUserIdByLoginReturn></getUserIdByLoginResponse>`),
		nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/account/userid", nil)
	w := httptest.NewRecorder()
	h.GetUserIdByLogin(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["user_id"] != "12345" {
		t.Errorf("user_id = %q, want %q", body["user_id"], "12345")
	}
}

func TestGetUserIdByLogin_SOAPError(t *testing.T) {
	h := newTestHandler(nil, fmt.Errorf("connection timeout"))

	req := httptest.NewRequest(http.MethodGet, "/account/userid", nil)
	w := httptest.NewRecorder()
	h.GetUserIdByLogin(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadGateway)
	}
}
