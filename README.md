# ChatRoom

仿微信即时通讯系统 — 前后端分离，Go + Vue 3 + WebSocket

## 技术栈

| 层级 | 技术 |
|---|---|
| 后端 | Go 1.20 + Gin + GORM + Gorilla WebSocket |
| 数据库 | MySQL 8 + Redis |
| 消息队列 | Go Channel（默认）/ Kafka（可选切换） |
| 前端 | Vue 3 + Pinia + Element Plus + WebRTC |
| 部署 | Docker Compose |

## 功能

- 密码注册/登录、短信验证码登录、JWT 鉴权
- 单聊、群聊、文件消息、音视频通话（WebRTC 信令）
- 好友申请/审批/拉黑/删除、群创建/成员管理/解散
- 后台管理：用户和群启用/禁用/删除、管理员设置
- Channel 与 Kafka 双模式消息转发

## 本地运行

### 环境

- Docker Desktop（已启动）
- Go 1.20+、Node.js + npm

### 第一步：启动基础设施

```powershell
docker compose up -d
```

### 第二步：启动后端

```powershell
go run .\cmd\kama_chat_server
# 后端地址：http://127.0.0.1:8000
```

### 第三步：启动前端

```powershell
cd web/chat-server
npm.cmd install
npm.cmd run serve
# 前端地址：http://127.0.0.1:8080
```

本地验证码：`123456`

### 停止

```powershell
docker compose down
```

## 目录结构

```text
api/v1/             Gin 控制器
cmd/                程序入口
configs/            配置文件
internal/
  ├── config/       配置加载
  ├── dao/          GORM 初始化
  ├── dto/          请求/响应 DTO
  ├── model/        数据模型
  ├── service/      业务逻辑 + 聊天服务 + Redis + Kafka + 短信
  ├── middleware/    JWT 鉴权中间件
  └── https_server/  路由注册
pkg/                工具包（日志、枚举、常量）
static/             静态资源
test/               测试
web/chat-server/    Vue 前端
docs/               文档
```

## 配置

`configs/config.toml`：

| 配置项 | 说明 |
|---|---|
| `mysqlConfig` | MySQL 连接（默认 127.0.0.1:13306） |
| `redisConfig` | Redis 连接（默认 127.0.0.1:16379） |
| `authCodeConfig.devMode` | `true` 时验证码固定 `123456` |
| `kafkaConfig.messageMode` | `channel`（默认）/ `kafka` |
| `jwtConfig.expireHours` | Token 有效期 |

可通过 `KAMA_CHAT_CONFIG` 环境变量指定其他配置文件。

## 相关文档

- [项目快速了解](项目快速了解.md)
- [Windows 运行指南](docs/Windows运行指南.md)
- [业务逻辑说明](docs/业务逻辑.md)
- [QA 问答记录](docs/QA问答记录.md)
- [开发与面试指南](项目开发与面试指南.md)
