package request

type ChatMessageRequest struct {
	SessionId  string `json:"session_id"`
	Type       int8   `json:"type"`
	Content    string `json:"content"`
	Url        string `json:"url"`
	SendId     string `json:"send_id"`
	SendName   string `json:"send_name"`
	SendAvatar string `json:"send_avatar"`
	ReceiveId  string `json:"receive_id"`
	FileSize   string `json:"file_size"`
	FileType   string `json:"file_type"`
	FileName   string `json:"file_name"`
	AVdata     string `json:"av_data"`
}
