package respond

type LoadMyGroupRespond struct {
	GroupID   string `json:"group_id"`
	GroupName string `json:"group_name"`
	Avatar    string `json:"avatar"`
}
