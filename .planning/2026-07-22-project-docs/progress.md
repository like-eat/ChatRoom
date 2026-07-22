# 工作进度

## 2026-07-22
- 已读取 planning-with-files 技能说明。
- 已确认根目录不存在旧的规划文件。
- 已创建本次文档整理的独立规划目录。
- 当前进行：项目结构与配置调研。
- 已完成 README、Docker Compose、后端配置、Go 依赖、前端依赖和路由入口检查。
- 已完成 Windows 启动指南、后端入口、Gin 路由、GORM 初始化和项目目录分层检查。
- 已完成核心模型、实时 Channel/Kafka 服务、联系人/群/会话业务和 Redis 缓存策略检查。
- 已核实密码、验证码、WebSocket 客户端协程、消息状态和安全边界。
- 已完成 Pinia、路由、页面功能、接口使用和未完成功能检查。
- 已通过 git diff 核实 Windows、Docker、HTTP、配置路径、短信开发模式和 Pinia 迁移等现有改造。
- 已生成根目录 `项目快速了解.md` 和 `项目开发与面试指南.md`。
- 当前进行：文档结构、事实和 Markdown 格式校验。
- 已核对两份文档的标题结构、关键事实、功能完成度和本地运行参数。
- `git diff --check` 未发现 Markdown 空白格式问题。
- 两份交付文档已完成。

## 2026-07-22 安全升级
- 用户要求解决明文密码、缺少 JWT、WebSocket 信任 client_id 三项问题。
- 设计目标：bcrypt 新密码；旧明文账号登录后自动升级；HTTP Bearer JWT；WebSocket 使用 JWT 并由服务端从令牌获取用户 UUID。
- 当前进行：核对认证 DTO、控制器、路由与前端登录链路。
- 已完成 bcrypt、旧密码渐进迁移、JWT 签发解析、HTTP/管理员中间件、WebSocket 令牌握手与发送者身份覆盖。
- 已完成 Pinia token、Axios Bearer、三个登录入口、刷新重连和退出流程改造。
- 当前进行：格式化、测试和本地集成验证。
- 已新增密码哈希、JWT 签发解析与算法限制单元测试。
- 已新增 HTTP/WebSocket 集成测试，覆盖无令牌 401、有效令牌访问、非管理员 403、WebSocket 子协议鉴权和旧明文密码自动升级。
- 已通过 `go test ./...`，所有 Go 包和鉴权集成测试通过。
- 已通过 `npm.cmd run build`；仅有原项目包体积和 Browserslist 数据较旧警告。
- 安全升级与验证全部完成。
