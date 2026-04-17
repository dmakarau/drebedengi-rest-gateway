package models

type AccessStatus struct {
	Status int `json:"status"`
}

type CurrentRevision struct {
	Revision int `json:"revision"`
}

type ExpireDate struct {
	ExpireDate string `json:"expire_date"`
}

type SubscriptionStatus struct {
	Status string `json:"status"`
}

type RightAccess struct {
	Access string `json:"access"`
}

type UserID struct {
	UserID string `json:"user_id"`
}
