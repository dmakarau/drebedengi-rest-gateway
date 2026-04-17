package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func mustUnmarshalSlice(t *testing.T, body []byte) []map[string]any {
	t.Helper()
	var s []map[string]any
	if err := json.Unmarshal(body, &s); err != nil {
		t.Fatalf("response body is not a valid JSON array: %v\nbody: %s", err, body)
	}
	return s
}

func balanceHandlerXML(entries []struct{ placeID, currencyID, sum string }) []byte {
	var sb strings.Builder
	sb.WriteString(`<getBalanceResponse><getBalanceReturn>`)
	for _, e := range entries {
		sb.WriteString(`<item>`)
		sb.WriteString(`<item><key>place_id</key><value>` + e.placeID + `</value></item>`)
		sb.WriteString(`<item><key>currency_id</key><value>` + e.currencyID + `</value></item>`)
		sb.WriteString(`<item><key>sum</key><value>` + e.sum + `</value></item>`)
		sb.WriteString(`</item>`)
	}
	sb.WriteString(`</getBalanceReturn></getBalanceResponse>`)
	return []byte(sb.String())
}

func TestGetBalance_OK(t *testing.T) {
	h := newTestHandler(
		balanceHandlerXML([]struct{ placeID, currencyID, sum string }{
			{"12345", "643", "1000"},
		}),
		nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/balance", nil)
	w := httptest.NewRecorder()
	h.GetBalance(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	entries := mustUnmarshalSlice(t, w.Body.Bytes())
	if len(entries) != 1 {
		t.Fatalf("len = %d, want 1", len(entries))
	}
	if entries[0]["place_id"] != float64(12345) {
		t.Errorf("place_id = %v, want 12345", entries[0]["place_id"])
	}
	if entries[0]["currency_id"] != float64(643) {
		t.Errorf("currency_id = %v, want 643", entries[0]["currency_id"])
	}
	if entries[0]["sum"] != float64(10) {
		t.Errorf("sum = %v, want 10 (1000 hundredths / 100)", entries[0]["sum"])
	}
}

func TestGetBalance_Empty(t *testing.T) {
	h := newTestHandler(
		[]byte(`<getBalanceResponse><getBalanceReturn/></getBalanceResponse>`),
		nil,
	)

	req := httptest.NewRequest(http.MethodGet, "/balance", nil)
	w := httptest.NewRecorder()
	h.GetBalance(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	entries := mustUnmarshalSlice(t, w.Body.Bytes())
	if len(entries) != 0 {
		t.Errorf("len = %d, want 0", len(entries))
	}
}

func TestGetBalance_NetworkError(t *testing.T) {
	h := newTestHandler(nil, networkErr())

	req := httptest.NewRequest(http.MethodGet, "/balance", nil)
	w := httptest.NewRecorder()
	h.GetBalance(w, req)

	if w.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadGateway)
	}
	body := mustUnmarshalString(t, w.Body.Bytes())
	if body["error"] == "" {
		t.Error("expected non-empty error message in body")
	}
}
