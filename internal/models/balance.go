package models

type BalanceEntry struct {
	PlaceID    int64   `json:"place_id"`
	CurrencyID int64   `json:"currency_id"`
	Sum        float64 `json:"sum"`
}
