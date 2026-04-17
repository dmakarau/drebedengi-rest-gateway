package soap

import (
	"encoding/xml"
	"fmt"
)

// authParams returns the standard 3 auth params for every API call.
func authParams(apiId, login, pass string) []Param {
	return []Param{
		{Name: "apiId", Value: apiId},
		{Name: "login", Value: login},
		{Name: "pass", Value: pass},
	}
}

func GetAccessStatus(c Caller, apiId, login, pass string) (int, error) {
	body, err := c.Call("getAccessStatus", authParams(apiId, login, pass))
	if err != nil {
		return 0, err
	}
	var resp getAccessStatusResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("parsing getAccessStatus response: %w", err)
	}
	return resp.Return, nil
}

func GetCurrentRevision(c Caller, apiId, login, pass string) (int, error) {
	body, err := c.Call("getCurrentRevision", authParams(apiId, login, pass))
	if err != nil {
		return 0, err
	}
	var resp getCurrentRevisionResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("parsing getCurrentRevision response: %w", err)
	}
	return resp.Return, nil
}

func GetExpireDate(c Caller, apiId, login, pass string) (string, error) {
	body, err := c.Call("getExpireDate", authParams(apiId, login, pass))
	if err != nil {
		return "", err
	}
	var resp getExpireDateResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("parsing getExpireDate response: %w", err)
	}
	return resp.Return, nil
}

func GetSubscriptionStatus(c Caller, apiId, login, pass string) (string, error) {
	body, err := c.Call("getSubscriptionStatus", authParams(apiId, login, pass))
	if err != nil {
		return "", err
	}
	var resp getSubscriptionStatusResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("parsing getSubscriptionStatus response: %w", err)
	}
	return resp.Return, nil
}

func GetRightAccess(c Caller, apiId, login, pass string) (string, error) {
	body, err := c.Call("getRightAccess", authParams(apiId, login, pass))
	if err != nil {
		return "", err
	}
	var resp getRightAccessResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("parsing getRightAccess response: %w", err)
	}
	return resp.Return, nil
}

func GetUserIdByLogin(c Caller, apiId, login, pass string) (string, error) {
	body, err := c.Call("getUserIdByLogin", authParams(apiId, login, pass))
	if err != nil {
		return "", err
	}
	var resp getUserIdByLoginResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("parsing getUserIdByLogin response: %w", err)
	}
	return resp.Return, nil
}
