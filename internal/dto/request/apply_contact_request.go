package request

type ApplyContactRequest struct {
	OwnerId   string `json:"owner_id"`
	ContactId string `json:"contact_id"`
	Message   string `json:"message"`
}
