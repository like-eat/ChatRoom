# 项目文档整理计划

## 目标
在项目根目录生成两份中文 Markdown：一份作为后续 AI/开发者快速接手上下文，另一份作为开发复盘与面试问答材料。

## 阶段
- [complete] 1. 检查项目结构、配置、依赖与启动方式
- [complete] 2. 梳理后端架构、数据模型、通信链路与功能
- [complete] 3. 梳理前端架构、状态管理、路由、页面与功能完成度
- [complete] 4. 汇总 Windows/Docker/Pinia 等本地改造和已知问题
- [complete] 5. 编写两份根目录文档
- [complete] 6. 交叉校验文档准确性与可读性
- [complete] 7. 设计密码迁移、JWT 和 WebSocket 鉴权方案
- [complete] 8. 实现后端密码哈希与兼容迁移
- [complete] 9. 实现 JWT 签发、HTTP 中间件和管理员保护
- [complete] 10. 实现 WebSocket JWT 身份校验
- [complete] 11. 更新 Pinia、Axios 和 WebSocket 前端认证链路
- [complete] 12. 运行后端测试、前端构建和本地集成验证

## 交付文件
- `项目快速了解.md`
- `项目开发与面试指南.md`

## 错误记录
| 错误 | 原因 | 处理 |
|---|---|---|
| 首次检查规划文件的 PowerShell 命令解析失败 | `$name:` 被解析成带冒号的变量引用 | 改用格式化字符串后成功 |
| rg 的 Windows 通配路径报错 | `package*.json` 作为参数在 Windows 被当作非法路径 | 改为显式指定 package.json 和 package-lock.json |
| 首次运行 `go test` 无法写入默认 GOCACHE | 当前沙箱无权写默认 Go 构建缓存 | 将 GOCACHE 临时指向系统临时目录后通过 |
| 首次运行 `go mod tidy` 无法写入 Go 模块缓存 | 依赖整理需要访问沙箱外模块缓存 | 经授权在沙箱外执行后成功 |
