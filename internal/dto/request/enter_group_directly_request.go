package request

type EnterGroupDirectlyRequest struct {
	OwnerId   string `json:"owner_id"`
	ContactId string `json:"contact_id"`
}
