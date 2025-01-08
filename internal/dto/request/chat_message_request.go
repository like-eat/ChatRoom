package request

type ChatMessageRequest struct {
	SessionId string `json:"session_id"`
	Type      int8   `json:"type"`
	Content   string `json:"content"`
	Url       string `json:"url"`
	SendId    string `json:"send_id"`
	SendName  string `json:"send_name"`
	ReceiveId string `json:"receive_id"`
	FileSize  int    `json:"file_size"`
	FileType  string `json:"file_type"`
	FileName  string `json:"file_name"`
}
