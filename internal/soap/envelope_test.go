package soap

import (
	"net/http"
	"strings"
	"testing"
)

func TestBuildEnvelope_SimpleMethod(t *testing.T) {
	envelope, err := BuildEnvelope("getAccessStatus", []Param{
		{Name: "apiId", Value: "demo_api"},
		{Name: "login", Value: "demo@example.com"},
		{Name: "pass", Value: "demo"},
	})
	if err != nil {
		t.Fatal(err)
	}

	s := string(envelope)

	// Check XML declaration
	if !strings.HasPrefix(s, `<?xml version="1.0" encoding="UTF-8"?>`) {
		t.Error("missing XML declaration")
	}

	// Check all required namespaces
	for _, ns := range []string{
		`xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/"`,
		`xmlns:ns1="urn:ddengi"`,
		`xmlns:xsd="http://www.w3.org/2001/XMLSchema"`,
		`xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"`,
		`xmlns:ns2="http://xml.apache.org/xml-soap"`,
		`xmlns:SOAP-ENC="http://schemas.xmlsoap.org/soap/encoding/"`,
		`SOAP-ENV:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"`,
	} {
		if !strings.Contains(s, ns) {
			t.Errorf("missing namespace: %s", ns)
		}
	}

	// Check method wrapper
	if !strings.Contains(s, `<ns1:getAccessStatus>`) {
		t.Error("missing method open tag")
	}
	if !strings.Contains(s, `</ns1:getAccessStatus>`) {
		t.Error("missing method close tag")
	}

	// Check params
	if !strings.Contains(s, `<apiId xsi:type="xsd:string">demo_api</apiId>`) {
		t.Error("missing apiId param")
	}
	if !strings.Contains(s, `<login xsi:type="xsd:string">demo@example.com</login>`) {
		t.Error("missing login param")
	}
	if !strings.Contains(s, `<pass xsi:type="xsd:string">demo</pass>`) {
		t.Error("missing pass param")
	}
}

func TestBuildEnvelope_WithMapParam(t *testing.T) {
	envelope, err := BuildEnvelope("getRecordList", []Param{
		{Name: "apiId", Value: "demo_api"},
		{Name: "login", Value: "demo@example.com"},
		{Name: "pass", Value: "demo"},
		{Name: "params", Value: map[string]any{
			"r_period": 8,
		}},
		{Name: "idList", Value: nil},
	})
	if err != nil {
		t.Fatal(err)
	}

	s := string(envelope)

	if !strings.Contains(s, `<ns1:getRecordList>`) {
		t.Error("missing method tag")
	}
	if !strings.Contains(s, `<params xsi:type="ns2:Map">`) {
		t.Error("missing Map param")
	}
	if !strings.Contains(s, `<key xsi:type="xsd:string">r_period</key><value xsi:type="xsd:int">8</value>`) {
		t.Error("missing r_period in Map")
	}
	if !strings.Contains(s, `<idList xsi:nil="true"/>`) {
		t.Error("missing nil idList")
	}
}

func TestBuildEnvelope_NoParams(t *testing.T) {
	envelope, err := BuildEnvelope("someMethod", nil)
	if err != nil {
		t.Fatal(err)
	}
	s := string(envelope)
	if !strings.Contains(s, `<ns1:someMethod></ns1:someMethod>`) {
		t.Errorf("expected empty method body, got: %s", s)
	}
}

func TestParseResponse_Success(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Body>
    <getAccessStatusResponse>
      <getAccessStatusReturn xsi:type="xsd:integer">1</getAccessStatusReturn>
    </getAccessStatusResponse>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	body, fault, err := ParseResponse([]byte(xml))
	if err != nil {
		t.Fatal(err)
	}
	if fault != nil {
		t.Fatalf("unexpected fault: %v", fault)
	}
	if !strings.Contains(string(body), "getAccessStatusReturn") {
		t.Errorf("body should contain response element, got: %s", string(body))
	}
}

func TestParseResponse_Fault(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
  <SOAP-ENV:Body>
    <SOAP-ENV:Fault>
      <faultcode>SOAP-ENV:Client</faultcode>
      <faultstring>Access denied</faultstring>
    </SOAP-ENV:Fault>
  </SOAP-ENV:Body>
</SOAP-ENV:Envelope>`

	_, fault, err := ParseResponse([]byte(xml))
	if err != nil {
		t.Fatal(err)
	}
	if fault == nil {
		t.Fatal("expected fault, got nil")
	}
	if fault.Code != "SOAP-ENV:Client" {
		t.Errorf("fault code = %q, want %q", fault.Code, "SOAP-ENV:Client")
	}
	if fault.String != "Access denied" {
		t.Errorf("fault string = %q, want %q", fault.String, "Access denied")
	}
}

func TestParseResponse_MalformedXML(t *testing.T) {
	_, _, err := ParseResponse([]byte("not xml at all"))
	if err == nil {
		t.Error("expected error for malformed XML")
	}
}

func TestFault_Error(t *testing.T) {
	f := &Fault{Code: "Server", String: "internal error"}
	got := f.Error()
	want := "SOAP Fault [Server]: internal error"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFault_HTTPStatus(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		msg      string
		wantCode int
	}{
		{"client auth error", "SOAP-ENV:Client", "Access denied", http.StatusUnauthorized},
		{"client denied error", "Client", "permission denied for user", http.StatusUnauthorized},
		{"client bad input", "SOAP-ENV:Client", "invalid parameter value", http.StatusBadRequest},
		{"server error", "SOAP-ENV:Server", "internal server error", http.StatusBadGateway},
		{"unknown error", "Unknown", "something went wrong", http.StatusBadGateway},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fault{Code: tt.code, String: tt.msg}
			got := f.HTTPStatus()
			if got != tt.wantCode {
				t.Errorf("HTTPStatus() = %d, want %d", got, tt.wantCode)
			}
		})
	}
}
