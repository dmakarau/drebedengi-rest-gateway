package handlers

import (
	"net/http"

	"github.com/drebedengi-rest/internal/models"
	"github.com/drebedengi-rest/internal/respond"
	"github.com/drebedengi-rest/internal/soap"
)

func (h *Handler) GetAccessStatus(w http.ResponseWriter, r *http.Request) {
	status, err := soap.GetAccessStatus(r.Context(), h.SOAP)
	if err != nil {
		soapErr(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, models.AccessStatus{Status: status})
}

func (h *Handler) GetCurrentRevision(w http.ResponseWriter, r *http.Request) {
	rev, err := soap.GetCurrentRevision(r.Context(), h.SOAP)
	if err != nil {
		soapErr(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, models.CurrentRevision{Revision: rev})
}

func (h *Handler) GetExpireDate(w http.ResponseWriter, r *http.Request) {
	date, err := soap.GetExpireDate(r.Context(), h.SOAP)
	if err != nil {
		soapErr(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, models.ExpireDate{ExpireDate: date})
}

func (h *Handler) GetSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	status, err := soap.GetSubscriptionStatus(r.Context(), h.SOAP)
	if err != nil {
		soapErr(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, models.SubscriptionStatus{Status: status})
}

func (h *Handler) GetRightAccess(w http.ResponseWriter, r *http.Request) {
	access, err := soap.GetRightAccess(r.Context(), h.SOAP)
	if err != nil {
		soapErr(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, models.RightAccess{Access: access})
}

func (h *Handler) GetUserIdByLogin(w http.ResponseWriter, r *http.Request) {
	uid, err := soap.GetUserIdByLogin(r.Context(), h.SOAP)
	if err != nil {
		soapErr(w, err)
		return
	}
	respond.JSON(w, http.StatusOK, models.UserID{UserID: uid})
}
