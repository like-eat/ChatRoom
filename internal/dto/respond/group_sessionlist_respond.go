package respond

type GroupSessionListRespond struct {
	SessionId string `json:"session_id"`
	GroupName string `json:"group_name"`
	GroupId   string `json:"group_id"`
	Avatar    string `json:"avatar"`
}
