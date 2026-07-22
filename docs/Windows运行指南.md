# Windows + Docker 运行指南

## 环境

- Windows 10/11
- Docker Desktop（Linux 容器模式）
- Go 1.20+
- Node.js 与 npm

## 1. 启动 MySQL 和 Redis

在项目根目录执行：

```powershell
docker compose up -d
docker compose ps
```

项目使用独立端口，避免与本机已有服务冲突：

- MySQL：`127.0.0.1:13306`
- Redis：`127.0.0.1:16379`
- MySQL 数据库：`kama_chat`
- MySQL root 密码：`kama_chat`

配置位于 `configs/config.toml`。Docker 数据保存在命名卷中，普通的 `docker compose down` 不会删除数据。

如果 Docker Hub 拉取失败并提示连接 `127.0.0.1:7890`，请在 Docker Desktop 的代理设置中关闭失效代理，或先启动对应的本地代理软件。

## 2. 启动后端

在项目根目录执行：

```powershell
go mod download
go run .\cmd\kama_chat_server
```

后端地址为 `http://127.0.0.1:8000`，WebSocket 地址为 `ws://127.0.0.1:8000/wss`。首次启动时 GORM 会自动创建数据表，并自动创建 `logs`、`static/avatars` 和 `static/files` 目录。

如需使用其他配置文件，可设置：

```powershell
$env:KAMA_CHAT_CONFIG = "C:\path\to\config.toml"
go run .\cmd\kama_chat_server
```

本地配置自带开发用 JWT 密钥，可以直接运行。部署到其他机器或生产环境前，必须改用独立的强随机密钥：

```powershell
$env:KAMA_CHAT_JWT_SECRET = "请替换为至少32字符的强随机密钥"
go run .\cmd\kama_chat_server
```

登录、注册和验证码登录成功后，前端会保存 JWT，并自动为 HTTP 请求添加 Bearer Token；WebSocket 握手也会校验同一令牌。升级前已经存在的明文密码账号无需重建，首次使用正确密码登录后会自动迁移为 bcrypt 哈希。修改鉴权代码后，旧浏览器会话需要重新登录一次。

## 3. 启动前端

打开另一个 PowerShell：

```powershell
Set-Location .\web\chat-server
npm.cmd install
npm.cmd run serve
```

浏览器访问 `http://127.0.0.1:8080`。

前端已经从 Vuex 迁移到 Pinia。API 地址可通过 Vue CLI 环境变量覆盖：

```powershell
$env:VUE_APP_BACKEND_URL = "http://127.0.0.1:8000"
$env:VUE_APP_WS_URL = "ws://127.0.0.1:8000"
npm.cmd run serve
```

## 4. 本地验证码

默认配置开启本地验证码模式：

```toml
[authCodeConfig]
devMode = true
devCode = "123456"
```

点击发送验证码后，使用 `123456` 即可完成注册或短信登录，不会调用阿里云短信。

如需接入真实短信，将 `devMode` 改为 `false`，并配置有效的阿里云 `accessKeyID`、`accessKeySecret`、`signName` 和 `templateCode`。

## 5. 停止环境

在后端和前端终端按 `Ctrl+C`，然后执行：

```powershell
docker compose down
```

如确实需要连同本项目的 MySQL/Redis 数据一起删除：

```powershell
docker compose down -v
```
