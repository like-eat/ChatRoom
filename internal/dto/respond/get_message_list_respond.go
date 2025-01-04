package respond

type GetMessageListRespond struct {
	SendId    string `json:"send_id"`
	SendName  string `json:"send_name"`
	ReceiveId string `json:"receive_id"`
	Type      int8   `json:"type"`
	Content   string `json:"content"`
	Url       string `json:"url"`
	FileType  string `json:"file_type"`
	FileName  string `json:"file_name"`
	FileSize  int    `json:"file_size"`
	CreatedAt string `json:"created_at"` // 先用CreatedAt排序，后面考虑改成SentAt
}
