package soap

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client holds credentials and HTTP transport for calling the drebedengi SOAP API.
// It implements the Caller interface.
type Client struct {
	APIId    string
	Login    string
	Password string
	URL      string
	http     *http.Client
}

func NewClient(apiId, login, password, url string) *Client {
	return &Client{
		APIId:    apiId,
		Login:    login,
		Password: password,
		URL:      url,
		http:     &http.Client{Timeout: 30 * time.Second},
	}
}

// Call sends a SOAP request for the given method and returns the inner Body XML.
func (c *Client) Call(method string, params []Param) ([]byte, error) {
	envelope, err := BuildEnvelope(method, params)
	if err != nil {
		return nil, fmt.Errorf("building envelope: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.URL, bytes.NewReader(envelope))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "text/xml; charset=UTF-8")
	req.Header.Set("SOAPAction", "urn:SoapAction")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	inner, fault, err := ParseResponse(body)
	if err != nil {
		return nil, err
	}
	if fault != nil {
		return nil, fault
	}

	return inner, nil
}
