package soap

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

const (
	nsSOAPEnv = "http://schemas.xmlsoap.org/soap/envelope/"
	nsDDengi  = "urn:ddengi"
	nsXSD     = "http://www.w3.org/2001/XMLSchema"
	nsXSI     = "http://www.w3.org/2001/XMLSchema-instance"
	nsMap     = "http://xml.apache.org/xml-soap"
	nsSOAPEnc = "http://schemas.xmlsoap.org/soap/encoding/"
)

// Fault represents a SOAP Fault element parsed from a response.
type Fault struct {
	Code   string `xml:"faultcode"`
	String string `xml:"faultstring"`
}

func (f *Fault) Error() string {
	return fmt.Sprintf("SOAP Fault [%s]: %s", f.Code, f.String)
}

// HTTPStatus maps the SOAP fault code to an appropriate HTTP status code.
func (f *Fault) HTTPStatus() int {
	code := strings.ToLower(f.Code)
	msg := strings.ToLower(f.String)
	if strings.Contains(code, "client") {
		if strings.Contains(msg, "access") || strings.Contains(msg, "denied") || strings.Contains(msg, "auth") {
			return http.StatusUnauthorized
		}
		return http.StatusBadRequest
	}
	return http.StatusBadGateway
}

// BuildEnvelope constructs the full SOAP XML request for a given method and parameters.
func BuildEnvelope(method string, params []Param) ([]byte, error) {
	var b strings.Builder

	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	b.WriteString(`<SOAP-ENV:Envelope`)
	b.WriteString(` xmlns:SOAP-ENV="` + nsSOAPEnv + `"`)
	b.WriteString(` xmlns:ns1="` + nsDDengi + `"`)
	b.WriteString(` xmlns:xsd="` + nsXSD + `"`)
	b.WriteString(` xmlns:xsi="` + nsXSI + `"`)
	b.WriteString(` xmlns:ns2="` + nsMap + `"`)
	b.WriteString(` xmlns:SOAP-ENC="` + nsSOAPEnc + `"`)
	b.WriteString(` SOAP-ENV:encodingStyle="` + nsSOAPEnc + `"`)
	b.WriteString(`>`)

	b.WriteString(`<SOAP-ENV:Body>`)
	b.WriteString(`<ns1:` + method + `>`)

	for _, p := range params {
		b.WriteString(encodeParam(p))
	}

	b.WriteString(`</ns1:` + method + `>`)
	b.WriteString(`</SOAP-ENV:Body>`)
	b.WriteString(`</SOAP-ENV:Envelope>`)

	return []byte(b.String()), nil
}

// envelopeResponse is used to parse the SOAP response envelope.
type envelopeResponse struct {
	XMLName xml.Name     `xml:"Envelope"`
	Body    bodyResponse `xml:"Body"`
}

type bodyResponse struct {
	Fault   *Fault `xml:"Fault,omitempty"`
	Content []byte `xml:",innerxml"`
}

// ParseResponse parses a SOAP response, returning the inner Body XML or a Fault.
func ParseResponse(data []byte) ([]byte, *Fault, error) {
	var env envelopeResponse
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, nil, fmt.Errorf("failed to parse SOAP response: %w", err)
	}
	if env.Body.Fault != nil && env.Body.Fault.Code != "" {
		return nil, env.Body.Fault, nil
	}
	return env.Body.Content, nil, nil
}
