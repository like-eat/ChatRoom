package contact_status_enum

const (
	NORMAL        = iota // 0 正常
	// 1 拉黑、2 被拉黑（已删除）
	BE_DELETE     = 3 // 3 被删除好友
	DELETE        = 4 // 4 删除好友
	SILENCE       = 5 // 5 被禁言
	QUIT_GROUP    = 6 // 6 退出群聊
	KICK_OUT_GROUP = 7 // 7 被踢出群聊
)
