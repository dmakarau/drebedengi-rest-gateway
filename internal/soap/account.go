package soap

import (
	"context"
	"encoding/xml"
	"fmt"
)

func GetAccessStatus(ctx context.Context, c Caller) (int, error) {
	body, err := c.Call(ctx, "getAccessStatus", nil)
	if err != nil {
		return 0, err
	}
	var resp getAccessStatusResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("parsing getAccessStatus response: %w", err)
	}
	return resp.Return, nil
}

func GetCurrentRevision(ctx context.Context, c Caller) (int, error) {
	body, err := c.Call(ctx, "getCurrentRevision", nil)
	if err != nil {
		return 0, err
	}
	var resp getCurrentRevisionResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return 0, fmt.Errorf("parsing getCurrentRevision response: %w", err)
	}
	return resp.Return, nil
}

func GetExpireDate(ctx context.Context, c Caller) (string, error) {
	body, err := c.Call(ctx, "getExpireDate", nil)
	if err != nil {
		return "", err
	}
	var resp getExpireDateResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("parsing getExpireDate response: %w", err)
	}
	return resp.Return, nil
}

func GetSubscriptionStatus(ctx context.Context, c Caller) (string, error) {
	body, err := c.Call(ctx, "getSubscriptionStatus", nil)
	if err != nil {
		return "", err
	}
	var resp getSubscriptionStatusResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("parsing getSubscriptionStatus response: %w", err)
	}
	return resp.Return, nil
}

func GetRightAccess(ctx context.Context, c Caller) (string, error) {
	body, err := c.Call(ctx, "getRightAccess", nil)
	if err != nil {
		return "", err
	}
	var resp getRightAccessResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("parsing getRightAccess response: %w", err)
	}
	return resp.Return, nil
}

func GetUserIdByLogin(ctx context.Context, c Caller) (string, error) {
	body, err := c.Call(ctx, "getUserIdByLogin", nil)
	if err != nil {
		return "", err
	}
	var resp getUserIdByLoginResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("parsing getUserIdByLogin response: %w", err)
	}
	return resp.Return, nil
}
