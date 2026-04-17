package handlers

import (
	"net/http"

	"github.com/drebedengi-rest/internal/models"
	"github.com/drebedengi-rest/internal/respond"
	"github.com/drebedengi-rest/internal/soap"
)

func (h *Handler) GetAccessStatus(w http.ResponseWriter, r *http.Request) {
	status, err := soap.GetAccessStatus(h.SOAP, h.APIId, h.Login, h.Password)
	if err != nil {
		respond.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, models.AccessStatus{Status: status})
}

func (h *Handler) GetCurrentRevision(w http.ResponseWriter, r *http.Request) {
	rev, err := soap.GetCurrentRevision(h.SOAP, h.APIId, h.Login, h.Password)
	if err != nil {
		respond.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, models.CurrentRevision{Revision: rev})
}

func (h *Handler) GetExpireDate(w http.ResponseWriter, r *http.Request) {
	date, err := soap.GetExpireDate(h.SOAP, h.APIId, h.Login, h.Password)
	if err != nil {
		respond.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, models.ExpireDate{ExpireDate: date})
}

func (h *Handler) GetSubscriptionStatus(w http.ResponseWriter, r *http.Request) {
	status, err := soap.GetSubscriptionStatus(h.SOAP, h.APIId, h.Login, h.Password)
	if err != nil {
		respond.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, models.SubscriptionStatus{Status: status})
}

func (h *Handler) GetRightAccess(w http.ResponseWriter, r *http.Request) {
	access, err := soap.GetRightAccess(h.SOAP, h.APIId, h.Login, h.Password)
	if err != nil {
		respond.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, models.RightAccess{Access: access})
}

func (h *Handler) GetUserIdByLogin(w http.ResponseWriter, r *http.Request) {
	uid, err := soap.GetUserIdByLogin(h.SOAP, h.APIId, h.Login, h.Password)
	if err != nil {
		respond.Error(w, http.StatusBadGateway, err.Error())
		return
	}
	respond.JSON(w, http.StatusOK, models.UserID{UserID: uid})
}
