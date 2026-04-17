package handlers

import (
	"net/http"

	"github.com/drebedengi-rest/internal/models"
	"github.com/drebedengi-rest/internal/respond"
	"github.com/drebedengi-rest/internal/soap"
)

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	restDate := q.Get("rest_date")
	isWithAccum := q.Get("with_accum") == "true"
	isWithDuty := q.Get("with_duty") == "true"

	items, err := soap.GetBalance(r.Context(), h.SOAP, restDate, isWithAccum, isWithDuty)
	if err != nil {
		soapErr(w, err)
		return
	}

	entries := make([]models.BalanceEntry, len(items))
	for i, item := range items {
		entries[i] = models.BalanceEntry{
			PlaceID:    item.PlaceID,
			CurrencyID: item.CurrencyID,
			Sum:        float64(item.Sum) / 100,
		}
	}
	respond.JSON(w, http.StatusOK, entries)
}
