package soap

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

func balanceXML(entries []struct{ placeID, currencyID, sum string }) []byte {
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

func TestGetBalance_Success(t *testing.T) {
	mock := &mockCaller{
		response: balanceXML([]struct{ placeID, currencyID, sum string }{
			{"12345", "643", "500000"},
		}),
	}

	items, err := GetBalance(context.Background(), mock, "", false, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("len = %d, want 1", len(items))
	}
	if items[0].PlaceID != 12345 {
		t.Errorf("PlaceID = %d, want 12345", items[0].PlaceID)
	}
	if items[0].CurrencyID != 643 {
		t.Errorf("CurrencyID = %d, want 643", items[0].CurrencyID)
	}
	if items[0].Sum != 500000 {
		t.Errorf("Sum = %d, want 500000", items[0].Sum)
	}
	if mock.method != "getBalance" {
		t.Errorf("method = %q, want %q", mock.method, "getBalance")
	}
}

func TestGetBalance_Empty(t *testing.T) {
	mock := &mockCaller{
		response: []byte(`<getBalanceResponse><getBalanceReturn/></getBalanceResponse>`),
	}

	items, err := GetBalance(context.Background(), mock, "", false, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 0 {
		t.Errorf("len = %d, want 0", len(items))
	}
}

func TestGetBalance_CallError(t *testing.T) {
	mock := &mockCaller{err: fmt.Errorf("connection refused")}
	_, err := GetBalance(context.Background(), mock, "", false, false)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetBalance_BadXML(t *testing.T) {
	mock := &mockCaller{response: []byte(`not valid xml`)}
	_, err := GetBalance(context.Background(), mock, "", false, false)
	if err == nil {
		t.Fatal("expected error for bad XML")
	}
}
