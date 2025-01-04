package request

type BlackApplyRequest struct {
	OwnerId   string `json:"owner_id"`
	ContactId string `json:"contact_id"`
}
