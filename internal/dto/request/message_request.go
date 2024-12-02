package request

type MessageRequest struct {
	Type      int    `json:"type"`
	Content   string `json:"content"`
	Url       string `json:"url"`
	SendId    string `json:"send_id"`
	ReceiveId string `json:"receive_id"`
	FileType  string `json:"file_type"`
	FileName  string `json:"file_name"`
	FileSize  int    `json:"file_size"`
}
