# 修复进度

## 会话：2026-07-22 — Redis 缓存失效策略统一重构

### 完成项
- ✅ 创建 `pkg/constants/cache_keys.go` — 12 个 Key 构造函数 + 7 个批量前缀
- ✅ 新增 `ScanAndDelete` 方法到 Redis 服务层（使用 SCAN 替代 KEYS）
- ✅ 所有 8 个 service 文件中的硬编码 Key 已替换为常量函数
- ✅ 所有被注释的缓存写回（SetKeyEx）和缓存失效（DelKeys*）代码已取消注释
- ✅ 所有写操作统一为"先更新 MySQL → 再删除 Redis 缓存"模式
- ✅ 修复 2 个 Bug（空格 + 重复删除）
- ✅ `go build ./...` 通过
- ✅ `go vet ./...` 通过

### 涉及文件（10 个）
1. pkg/constants/cache_keys.go（新建）
2. internal/service/redis/redis_service.go
3. internal/service/sms/auth_code_service.go
4. internal/service/gorm/user_info_service.go
5. internal/service/gorm/user_contact_service.go
6. internal/service/gorm/group_info_service.go
7. internal/service/gorm/session_service.go
8. internal/service/gorm/message_service.go
9. internal/service/chat/server.go
10. internal/service/chat/kafka_server.go

### 缓存策略模式
- **读操作**：查 Redis → 未命中查 MySQL → 写回 Redis（Cache-Aside Read）
- **写操作**：更新 MySQL → 删除 Redis 对应 Key（Cache-Aside Write，Cache Invalidation）
- **批量失效**：使用 `ScanAndDelete(prefix:*)` 匹配前缀批量删除
