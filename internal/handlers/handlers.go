package handlers

import (
	"github.com/drebedengi-rest/internal/soap"
	"github.com/go-chi/chi/v5"
)

// Handler holds dependencies shared by all HTTP handlers.
type Handler struct {
	SOAP     soap.Caller
	APIId    string
	Login    string
	Password string
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

	return r
}
