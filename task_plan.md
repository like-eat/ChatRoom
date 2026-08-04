# Redis 缓存失效策略统一重构

## 目标
修复 Redis 缓存三大问题：取消注释遗留的缓存代码、建立集中的缓存 Key 管理机制、统一写操作为"先更新 MySQL → 再删除/刷新 Redis"模式。

## 阶段
- [complete] 1. 创建集中式 Cache Key 管理文件（`pkg/constants/cache_keys.go`）
- [complete] 2. 增强 Redis 服务层（添加 `ScanAndDelete` 方法）
- [complete] 3. 修复 user_info_service.go（取消注释 + 统一模式）
- [complete] 4. 修复 user_contact_service.go（重构 Key + 统一模式）
- [complete] 5. 修复 group_info_service.go（取消注释 + 修复 Bug + 统一模式）
- [complete] 6. 修复 session_service.go（取消注释 + 统一模式）
- [complete] 7. 修复 message_service.go（取消注释 + 统一模式）
- [complete] 8. 修复 chat/server.go 和 chat/kafka_server.go（遗漏文件）
- [complete] 9. 修复 sms/auth_code_service.go（遗漏文件）
- [complete] 10. 编译验证 + 测试

## 修改文件清单
| 文件 | 操作 | 说明 |
|------|------|------|
| `pkg/constants/cache_keys.go` | 新建 | 12 个 Key 函数 + 7 个批量前缀常量 |
| `internal/service/redis/redis_service.go` | 修改 | 新增 `ScanAndDelete` 方法 |
| `internal/service/sms/auth_code_service.go` | 修改 | 替换硬编码 Key |
| `internal/service/gorm/user_info_service.go` | 修改 | 取消注释 6 处 + 统一 Key |
| `internal/service/gorm/user_contact_service.go` | 修改 | 替换 5 处 Key |
| `internal/service/gorm/group_info_service.go` | 修改 | 取消注释 10+ 处 + 修复空格 Bug + 修复重复删除 |
| `internal/service/gorm/session_service.go` | 修改 | 取消注释 1 处 + 替换 5 处 Key |
| `internal/service/gorm/message_service.go` | 修改 | 取消注释 2 处 + 添加 time 导入 |
| `internal/service/chat/server.go` | 修改 | 替换 8 处 Key |
| `internal/service/chat/kafka_server.go` | 修改 | 替换 8 处 Key |

## Bug 修复
1. `group_info_service.go:315` — `"my_joined_group_list_ "` 末尾多余空格已修正
2. `group_info_service.go:462-467` — 重复的 `DelKeysWithPrefix("group_session_list")` 已删除
