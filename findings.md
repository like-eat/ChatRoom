# 研究发现：Redis 缓存问题

## 问题 1：注释遗留的缓存代码
详见上方分析，涉及 5 个 service 文件，约 20+ 处被注释。

## 问题 2：Key 散落
11 种 Key 模式分布在 5 个文件中，以硬编码字符串拼接。

## 问题 3：模式不统一
部分写操作后有缓存失效，部分没有（被注释）。

## 附加 Bug
- group_info_service.go:315 — `"my_joined_group_list_ "` 末尾多余空格
- group_info_service.go:462-467 — `DelKeysWithPrefix("group_session_list")` 重复调用
- redis_service.go 全局 — `KEYS` 命令不适合生产，应优先使用 `SCAN`

## Redis Key 清单
| Key 模式 | 用途 | TTL |
|---|---|---|
| `auth_code:{telephone}` | 短信验证码 | 5min |
| `user_info:{uuid}` | 用户信息 | 1min |
| `contact_user_list:{ownerId}` | 联系人列表 | 1min |
| `contact_mygroup_list:{ownerId}` | 我创建的群 | 1min |
| `my_joined_group_list:{ownerId}` | 我加入的群 | 1min |
| `group_info:{groupId}` | 群信息 | 1min |
| `group_memberlist:{groupId}` | 群成员 | 1min |
| `session_list:{ownerId}` | 用户会话列表 | 1min |
| `group_session_list:{ownerId}` | 群聊会话列表 | 1min |
| `session:{sendId}:{receiveId}` | 单个会话 | 1min |
| `message_list:{userOneId}:{userTwoId}` | 私聊消息 | 1min |
| `group_messagelist:{groupId}` | 群聊消息 | 1min |
