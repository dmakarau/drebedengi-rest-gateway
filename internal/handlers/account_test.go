package handlers

import (
	"context"
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

func (m *mockCaller) Call(_ context.Context, _ string, _ []soap.Param) ([]byte, error) {
	return m.response, m.err
}

func newTestHandler(response []byte, err error) *Handler {
	return &Handler{
		SOAP: &mockCaller{response: response, err: err},
	}
}

// mustUnmarshal fails the test immediately if JSON decoding fails.
func mustUnmarshalString(t *testing.T, body []byte) map[string]string {
	t.Helper()
	var m map[string]string
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("response body is not valid JSON: %v\nbody: %s", err, body)
	}
	return m
}

func mustUnmarshalInt(t *testing.T, body []byte) map[string]int {
	t.Helper()
	var m map[string]int
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("response body is not valid JSON: %v\nbody: %s", err, body)
	}
	return m
}

// soapFaultResponse returns a minimal XML body for a given SOAP fault.
func networkErr() error { return fmt.Errorf("connection timeout") }

// --- GetAccessStatus ---

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
	body := mustUnmarshalInt(t, w.Body.Bytes())
	if body["status"] != 1 {
		t.Errorf("body status = %d, want 1", body["status"])
	}
}

func TestGetAccessStatus_SOAPFault_Auth(t *testing.T) {
	h := newTestHandler(nil, &soap.Fault{Code: "SOAP-ENV:Client", String: "Access denied"})

	req := httptest.NewRequest(http.MethodGet, "/account/status", nil)
	w := httptest.NewRecorder()
	h.GetAccessStatus(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
	body := mustUnmarshalString(t, w.Body.Bytes())
	if body["error"] == "" {
		t.Error("expected non-empty error message in body")
	}
}

func TestGetAccessStatus_SOAPFault_BadInput(t *testing.T) {
	h := newTestHandler(nil, &soap.Fault{Code: "SOAP-ENV:Client", String: "invalid parameter"})

	req := httptest.NewRequest(http.MethodGet, "/account/status", nil)
	w := httptest.NewRecorder()
	h.GetAccessStatus(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
	body := mustUnmarshalString(t, w.Body.Bytes())
	if body["error"] == "" {
		t.Error("expected non-empty error message in body")
	}
}

func TestGetAccessStatus_NetworkError(t *testing.T) {
	h := newTestHandler(nil, networkErr())

	req := httptest.NewRequest(http.MethodGet, "/account/status", nil)
	w := httptest.NewRecorder()
	h.GetAccessStatus(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadGateway)
	}
	body := mustUnmarshalString(t, w.Body.Bytes())
	if body["error"] == "" {
		t.Error("expected non-empty error message in body")
	}
}

// --- GetCurrentRevision ---

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
	body := mustUnmarshalInt(t, w.Body.Bytes())
	if body["revision"] != 42 {
		t.Errorf("revision = %d, want 42", body["revision"])
	}
}

func TestGetCurrentRevision_Error(t *testing.T) {
	h := newTestHandler(nil, networkErr())

	req := httptest.NewRequest(http.MethodGet, "/account/revision", nil)
	w := httptest.NewRecorder()
	h.GetCurrentRevision(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadGateway)
	}
	body := mustUnmarshalString(t, w.Body.Bytes())
	if body["error"] == "" {
		t.Error("expected non-empty error message in body")
	}
}

// --- GetExpireDate ---

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
	body := mustUnmarshalString(t, w.Body.Bytes())
	if body["expire_date"] != "2037-01-02" {
		t.Errorf("expire_date = %q, want %q", body["expire_date"], "2037-01-02")
	}
}

func TestGetExpireDate_Error(t *testing.T) {
	h := newTestHandler(nil, networkErr())

	req := httptest.NewRequest(http.MethodGet, "/account/expire", nil)
	w := httptest.NewRecorder()
	h.GetExpireDate(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadGateway)
	}
	body := mustUnmarshalString(t, w.Body.Bytes())
	if body["error"] == "" {
		t.Error("expected non-empty error message in body")
	}
}

// --- GetSubscriptionStatus ---

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
	body := mustUnmarshalString(t, w.Body.Bytes())
	if body["status"] != "1" {
		t.Errorf("status = %q, want %q", body["status"], "1")
	}
}

func TestGetSubscriptionStatus_Error(t *testing.T) {
	h := newTestHandler(nil, networkErr())

	req := httptest.NewRequest(http.MethodGet, "/account/subscription", nil)
	w := httptest.NewRecorder()
	h.GetSubscriptionStatus(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadGateway)
	}
	body := mustUnmarshalString(t, w.Body.Bytes())
	if body["error"] == "" {
		t.Error("expected non-empty error message in body")
	}
}

// --- GetRightAccess ---

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
	body := mustUnmarshalString(t, w.Body.Bytes())
	if body["access"] != "0" {
		t.Errorf("access = %q, want %q", body["access"], "0")
	}
}

func TestGetRightAccess_Error(t *testing.T) {
	h := newTestHandler(nil, networkErr())

	req := httptest.NewRequest(http.MethodGet, "/account/access", nil)
	w := httptest.NewRecorder()
	h.GetRightAccess(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadGateway)
	}
	body := mustUnmarshalString(t, w.Body.Bytes())
	if body["error"] == "" {
		t.Error("expected non-empty error message in body")
	}
}

// --- GetUserIdByLogin ---

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
	body := mustUnmarshalString(t, w.Body.Bytes())
	if body["user_id"] != "12345" {
		t.Errorf("user_id = %q, want %q", body["user_id"], "12345")
	}
}

func TestGetUserIdByLogin_SOAPFault_Server(t *testing.T) {
	h := newTestHandler(nil, &soap.Fault{Code: "SOAP-ENV:Server", String: "internal error"})

	req := httptest.NewRequest(http.MethodGet, "/account/userid", nil)
	w := httptest.NewRecorder()
	h.GetUserIdByLogin(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadGateway)
	}
	body := mustUnmarshalString(t, w.Body.Bytes())
	if body["error"] == "" {
		t.Error("expected non-empty error message in body")
	}
}
