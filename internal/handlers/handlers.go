package handlers

import (
	"errors"
	"net/http"

	"github.com/drebedengi-rest/internal/respond"
	"github.com/drebedengi-rest/internal/soap"
	"github.com/go-chi/chi/v5"
)

// Handler holds dependencies shared by all HTTP handlers.
type Handler struct {
	SOAP soap.Caller
}

// NewRouter creates the /api/v1 route tree.
func NewRouter(h *Handler) chi.Router {
	r := chi.NewRouter()

	r.Route("/account", func(r chi.Router) {
		r.Get("/status", h.GetAccessStatus)
		r.Get("/revision", h.GetCurrentRevision)
		r.Get("/expire", h.GetExpireDate)
		r.Get("/subscription", h.GetSubscriptionStatus)
		r.Get("/access", h.GetRightAccess)
		r.Get("/userid", h.GetUserIdByLogin)
	})

	r.Get("/balance", h.GetBalance)

	return r
}

// soapErr writes the appropriate HTTP error response for a SOAP call failure.
// SOAP Client faults map to 400/401; all other errors map to 502.
func soapErr(w http.ResponseWriter, err error) {
	var f *soap.Fault
	if errors.As(err, &f) {
		respond.Error(w, f.HTTPStatus(), f.String)
	} else {
		respond.Error(w, http.StatusBadGateway, err.Error())
	}
}
