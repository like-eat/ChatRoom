package request

type ChatMessageRequest struct {
	SessionId string `json:"session_id"`
	Type      int8   `json:"type"`
	Content   string `json:"content"`
	SendId    string `json:"send_id"`
	ReceiveId string `json:"receive_id"`
}
