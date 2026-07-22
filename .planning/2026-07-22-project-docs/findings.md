# 项目调研发现

> 仅记录从本地代码和配置中核实的事实。

## 已知上下文
- 项目为 KamaChat 分布式聊天室，前后端分离。
- 当前工作区已做 Windows + Docker 运行适配，前端状态管理已从 Vuex 迁移到 Pinia。
- 本地已有 bufan 演示账号、20 位 Mock 好友、7 个群和 40 条主题消息。

## 基础架构与依赖
- 后端：Go 1.20、Gin 1.10、GORM/MySQL、go-redis、Gorilla WebSocket、Zap + Lumberjack；go.mod 仍包含 kafka-go、阿里云短信 SDK 等依赖。
- 前端：Vue 3、Vue CLI 5、Vue Router 4、Pinia 2、Element Plus、Axios；不是 Vite 项目。
- 本地基础设施：Docker Compose 只启动 MySQL 8.0 和 Redis 8.2；MySQL 映射 13306，Redis 映射 16379。
- 后端监听 `127.0.0.1:8000`，数据库名 `kama_chat`。
- 消息模式当前配置为 `channel`，因此本地完整演示不要求 Kafka；也可以切为 `kafka` 并配置 9092。
- 短信当前为开发模式，固定验证码 `123456`；真实短信需配置阿里云凭证、签名和模板。
- 静态头像和文件分别保存到 `./static/avatars`、`./static/files`。
- 前端路由包含登录、验证码登录、注册、个人资料、联系人、聊天、会话列表和管理员后台，并通过 Pinia 中的用户 UUID 做基础路由守卫。

## 启动与后端分层
- Windows 启动顺序：`docker compose up -d` → 根目录 `go run .\cmd\kama_chat_server` → `web/chat-server` 执行 `npm.cmd run serve`，访问 8080。
- 后端启动时 GORM `AutoMigrate` 自动创建 6 张核心表：用户、群、用户联系人、会话、联系人申请、消息。
- Gin 在初始化时创建静态目录，并将 `/static/avatars`、`/static/files` 映射为静态资源路由。
- 代码分层：`api/v1` 控制器；`internal/dto` 请求响应；`internal/service/gorm` 业务与数据库；`internal/model` 模型；`internal/service/chat` 实时通信；`internal/service/redis` 缓存；`pkg` 常量、枚举、日志与工具。
- 接口覆盖用户注册登录/资料/管理、好友申请与拉黑、群创建与成员管理、会话管理、消息历史、文件上传、在线联系人和 WebSocket。
- 退出时会关闭 Kafka（若启用）、聊天服务，并清空 Redis 键。
- 当前 CORS 允许任意来源，接口层未见统一鉴权中间件，这是后续安全优化重点。

## 数据模型与实时消息
- 6 张核心表：`user_info`、`group_info`、`user_contact`、`contact_apply`、`session`、`message`；多数业务删除采用 GORM 软删除。
- ID 通过首字母区分实体：用户 U、群 G、会话 S、消息 M、申请 A，后端多处通过首字符判断单聊或群聊。
- 联系关系有正常、拉黑/被拉黑、删除/被删除、禁言、退出群、被踢等状态；好友申请有申请中、通过、拒绝、拉黑状态。
- 消息类型定义为文本 0、语音 1、文件 2、通话 3；当前前端主要渲染文本和文件，语音未完成。
- Channel 模式的聊天服务器维护在线客户端 Map，并使用 Login、Logout、Transmit 三个缓冲 Channel；读写由 WebSocket 客户端协程处理。
- 单聊消息持久化 MySQL 后分别推送给在线接收者和发送者回显；离线接收者稍后通过历史接口读取。
- 群消息根据 `group_info.members` JSON 数组遍历在线成员并广播；消息同时持久化。
- Redis 缓存联系人、群详情、会话列表和消息历史，TTL 配置常量为分钟级；数据变化时需要主动删除相关键，否则会出现短时旧数据。
- 项目保留 Kafka 消息转发实现，作为 Channel 模式的可切换替代方案。
- 数据层存在 `group_info.members` JSON 与 `user_contact` 两份群成员信息，需要关注一致性和事务边界。

## 认证与通信风险
- 密码当前以明文写入并直接比较，AES 文件中的函数未被业务调用；密码列仅 `char(18)`。应改为 bcrypt/Argon2 哈希并扩大字段。
- 系统没有 JWT/Session 认证；前端路由守卫只检查 Pinia 状态，后端管理接口也未见服务端角色中间件。
- WebSocket 仅依赖查询参数 `client_id` 标识用户，升级器允许所有 Origin，存在身份冒用风险。
- WebSocket Read/Write 各自独立协程，发送成功后将消息状态从未发送更新为已发送。
- 当服务端 Transmit Channel 满时先进入客户端 SendTo 缓冲；两级缓冲均满才向客户端返回繁忙提示。
- 开发验证码保存在 Redis，一分钟有效；注册和验证码登录读取该键。开发模式固定 123456，生产模式调用阿里云短信。
- 现有退出流程和断线清理、心跳、重连、ACK/重试、消息幂等仍有较大优化空间。

## 前端架构与功能完成度
- 全局状态已迁移为 Pinia：保存后端 URL、WebSocket URL、用户信息和 Socket；用户信息写入 `sessionStorage`，刷新后可恢复，Socket 在 App 挂载时重建。
- 未发现业务代码继续调用 Vuex API；`package-lock.json` 中的 Vuex CLI 插件是脚手架间接依赖，不代表项目仍使用 Vuex。
- 登录、注册和验证码登录成功后创建 WebSocket；App 刷新时也会尝试恢复连接。
- 联系人模块支持好友/群列表、搜索申请、审批/拒绝/拉黑申请、创建群、加群、取消拉黑和进入会话。
- 会话模块区分用户会话和群会话，并支持进入和删除会话。
- 聊天页支持文本消息、基础文件上传/下载、好友删除/拉黑、群资料、群成员、退群/解散、群主审批以及 WebRTC 音视频信令。
- 管理后台支持用户启用/禁用/删除、管理员设置，群启用/禁用/删除。
- 表情包按钮未实现，且当前错误调用下载方法；聊天记录与全文复制只有图标无事件；语音消息未实现。
- 文件上传基础链路已存在，但使用原始文件名、本地磁盘存储，缺少防重名、完整错误反馈、进度和分布式对象存储。
- `App.vue` 的被禁用用户退出分支引用未定义的 `data/router/ElMessage`，且状态检查存在异步时序问题，是明确待修复项。

## 当前工作区已完成的本地化改造
- 新增 Docker Compose，使用命名卷运行 MySQL/Redis，并避开常用端口冲突。
- 将原本硬编码 Linux 配置路径改为默认相对路径，并支持 `KAMA_CHAT_CONFIG` 环境变量；日志与静态目录自动创建。
- MySQL 从 Unix Socket 改为 TCP DSN，适配 Windows 主机连接 Docker 数据库。
- 本地开发由强依赖 Linux HTTPS 证书改为 HTTP 8000；前端从 HTTPS 443 改为 HTTP 8080。
- 增加短信开发模式，固定验证码可在无阿里云账号时完成注册和登录。
- Vuex 已迁移为 Pinia，所有视图和组件改用 `useAppStore`；API/WS 地址支持 Vue CLI 环境变量覆盖。
- 用户状态写入 sessionStorage，支持刷新恢复；刷新后重建 WebSocket。
- 聊天页外部背景已改为米白色，群标题改为单行自适应并支持省略号。
- 增加 Windows 运行指南、Mock 数据脚本、本地头像，以及 bufan 演示数据（20 好友、7 群、40 消息、管理员）。
- 测试数量较少，只有配置、日志和 GORM 创建用户等基础 Go 测试；前端未发现单元测试。

## 调研错误补充
- Windows 下对 `web/chat-server/package*.json` 使用 rg glob 触发路径语法错误；改为列出具体文件后完成 Vuex 残留检查。

## 2026-07-22 安全升级依据
- RFC 7519 定义 JWT 为紧凑、URL-safe 的 Claims 表示，并定义 `sub`、`exp`、`iat` 等注册 Claim。
- RFC 8725 要求验证时固定允许的算法集合，不能信任攻击者修改的 `alg`，并强调 HS256 密钥必须具有足够熵。
- `golang.org/x/crypto/bcrypt` 已存在于当前 go.mod 的依赖图，可直接使用 `GenerateFromPassword` 和 `CompareHashAndPassword`，不需要新增密码学依赖。
- 迁移方案：bcrypt 新密码；旧明文密码成功登录一次后原地升级；数据库密码列从 char(18) 扩到 varchar(100)。
- JWT 计划使用 HS256、固定算法检查、`iss`/`aud`/`sub`/`iat`/`exp` Claims，并允许 `KAMA_CHAT_JWT_SECRET` 覆盖本地配置密钥。
- 浏览器 WebSocket API 不能自定义 Authorization Header，因此通过 `Sec-WebSocket-Protocol` 携带 JWT；后端验证令牌后从 `sub` 获取用户 UUID，不再读取 `client_id`，也避免令牌出现在 URL 和常规访问日志中。
- 登录和注册响应 DTO 当前没有 token 字段，需要同时扩展 `LoginRespond`、`RegisterRespond`；短信登录复用 LoginRespond。
- 三个入口页面都在成功后立即创建 WebSocket，App 刷新恢复也会创建连接；这些位置都必须从 `client_id` 改为 token。
- Pinia 需要新增 token 状态、sessionStorage 持久化、setToken 和清理逻辑；Axios 可用全局请求拦截器统一添加 Bearer Header。
- HTTP 路由应分为公开组（登录、注册、发验证码、短信登录）和 JWT 保护组；用户/群管理接口再叠加管理员中间件。
- `NavigationModal.logout` 当前先清空用户再调用 wsLogout，但请求体使用组件副本，可工作；升级后仍应先保存 ownerId，完成请求后统一清除认证状态。
- `golang-jwt/jwt/v5` 当前文档展示 v5.3.1，但该版本 go.mod 要求 Go 1.21；项目声明 Go 1.20，因此选择兼容 Go 1.18+ 的 v5.2.2，并使用 `WithValidMethods`、issuer、audience 和 expiration 校验选项。

## 安全升级实现记录
- 已新增 JWT 配置，生产可通过 `KAMA_CHAT_JWT_SECRET` 覆盖；密钥短于 32 字符时拒绝签发和解析。
- JWT 使用 HS256，包含 `sub`、`iss`、`aud`、`iat`、`nbf`、`exp`、`token_use` 和管理员标记；解析固定 HS256 并校验签发者、受众和过期时间。
- HTTP 中间件不只验证 JWT，还查询数据库确认用户存在且未禁用；管理员中间件使用数据库中的实时角色，不依赖可能过期的 JWT 管理员 Claim。
- HTTP 路由已拆为公开、登录保护和管理员保护三组。
- WebSocket 握手通过 `Sec-WebSocket-Protocol` 提交 JWT，服务端从 JWT `sub` 查询有效用户并创建 Client；协商出的应用子协议固定为 `kama-chat`。
- WebSocket Client 保存后端查询到的昵称和头像；每条浏览器消息进入队列前强制覆盖 send_id/send_name/send_avatar，防止伪造发送者。
- Pinia 已新增 token 持久化、统一 Socket 建连和清理；Axios 全局拦截器添加 Bearer Header，401 时清理会话并回登录页。
- 登录、注册、验证码登录和刷新恢复均改为使用 JWT WebSocket；退出请求在清理 token 前发送。
