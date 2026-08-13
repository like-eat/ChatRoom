package constants

// Redis Key 命名规范：{业务模块}:{主标识}[:{子标识}]
//
// 缓存失效策略（Cache-Aside）：
//   写操作：先更新 MySQL → 再主动删除/刷新 Redis
//   读操作：查 Redis → Miss 则查 MySQL → 写回 Redis
//
// 批量失效使用前缀通配符：{业务模块}:*

import "fmt"

// ════════════════ 用户信息 ════════════════

func CacheKeyUserInfo(uuid string) string {
	return fmt.Sprintf("user_info:%s", uuid)
}

// ════════════════ 联系人/群组列表 ════════════════

func CacheKeyContactUserList(ownerId string) string {
	return fmt.Sprintf("contact_user_list:%s", ownerId)
}

func CacheKeyMyGroupList(ownerId string) string {
	return fmt.Sprintf("contact_mygroup_list:%s", ownerId)
}

func CacheKeyMyJoinedGroupList(ownerId string) string {
	return fmt.Sprintf("my_joined_group_list:%s", ownerId)
}

// ════════════════ 群组信息 ════════════════

func CacheKeyGroupInfo(groupId string) string {
	return fmt.Sprintf("group_info:%s", groupId)
}

func CacheKeyGroupMemberList(groupId string) string {
	return fmt.Sprintf("group_memberlist:%s", groupId)
}

// ════════════════ 会话 ════════════════

func CacheKeySessionList(ownerId string) string {
	return fmt.Sprintf("session_list:%s", ownerId)
}

func CacheKeyGroupSessionList(ownerId string) string {
	return fmt.Sprintf("group_session_list:%s", ownerId)
}

func CacheKeySession(sendId, receiveId string) string {
	return fmt.Sprintf("session:%s:%s", sendId, receiveId)
}

// ════════════════ 消息记录 ════════════════

// 注意：key 必须把两个 user id 按固定顺序（字典序）拼接，
// 否则"发送者:接收者"和"我:对方"会形成两个不同 key，
// 导致缓存里只存了单方向消息，刷新历史时对端消息丢失。
func CacheKeyMessageList(userOneId, userTwoId string) string {
	if userOneId < userTwoId {
		return fmt.Sprintf("message_list:%s:%s", userOneId, userTwoId)
	}
	return fmt.Sprintf("message_list:%s:%s", userTwoId, userOneId)
}

func CacheKeyGroupMessageList(groupId string) string {
	return fmt.Sprintf("group_messagelist:%s", groupId)
}

// ════════════════ 批量前缀（用于 Delete 操作） ════════════════

const (
	PrefixContactUserList    = "contact_user_list:*"
	PrefixMyGroupList        = "contact_mygroup_list:*"
	PrefixMyJoinedGroupList  = "my_joined_group_list:*"
	PrefixGroupInfo          = "group_info:*"
	PrefixGroupMemberList    = "group_memberlist:*"
	PrefixSessionList        = "session_list:*"
	PrefixGroupSessionList   = "group_session_list:*"
)
