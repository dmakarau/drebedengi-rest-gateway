package soap

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
)

type BalanceItem struct {
	PlaceID    int64
	CurrencyID int64
	Sum        int64 // hundredths
}

func GetBalance(ctx context.Context, c Caller, restDate string, isWithAccum, isWithDuty bool) ([]BalanceItem, error) {
	body, err := c.Call(ctx, "getBalance", []Param{
		{Name: "params", Value: map[string]any{
			"rest_date":     restDate,
			"is_with_accum": isWithAccum,
			"is_with_duty":  isWithDuty,
		}},
	})
	if err != nil {
		return nil, err
	}
	var resp getBalanceResponse
	if err := xml.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing getBalance response: %w", err)
	}
	items := make([]BalanceItem, 0, len(resp.Return))
	for _, entry := range resp.Return {
		var item BalanceItem
		for _, kv := range entry.Items {
			switch kv.Key {
			case "place_id":
				item.PlaceID, _ = strconv.ParseInt(kv.Value, 10, 64)
			case "currency_id":
				item.CurrencyID, _ = strconv.ParseInt(kv.Value, 10, 64)
			case "sum":
				item.Sum, _ = strconv.ParseInt(kv.Value, 10, 64)
			}
		}
		items = append(items, item)
	}
	return items, nil
}
