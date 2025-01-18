package respond

type GetGroupMemberListRespond struct {
	UserId   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
