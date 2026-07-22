# KamaChat 项目 QA 问答记录

> 记录日期：2026-07-22
> 来源：项目代码走读问答

---

## 1. 系统有使用 Kafka 吗？

**有，但当前未激活。** 项目内置了两套消息转发模式，通过 `config.toml` 中 `kafkaConfig.messageMode` 切换：

- **Channel 模式（当前默认）**：Go Channel 进程内存队列，单进程无外部依赖，适合本地开发。
- **Kafka 模式（可选）**：使用 segmentio/kafka-go，消息写入 Kafka topic，由 `KafkaServer` 消费后持久化和推送，适合分布式多实例。

当前本地开发不启动 Kafka broker，实际运行是 Channel 模式。相关代码见 `internal/service/kafka/kafka_service.go` 和 `internal/service/chat/kafka_server.go`。

---

## 2. 文件上传为什么不走 WebSocket 而是 HTTP？

文件上传走 HTTP，WebSocket 只传文件消息的元数据（URL、文件名、大小）。原因：

1. **WebSocket 不适合传大文件**：默认缓冲区只有 2KB，大文件会占满连接、阻塞实时消息。
2. **HTTP 上传基础设施成熟**：浏览器原生 FormData + multipart，支持进度条。
3. **职责分离**：HTTP 负责文件传输，WebSocket 负责实时信令，互不干扰。
4. **即时通讯通用设计**：微信、Slack 都是这样做的。

流程：HTTP 上传文件 → 后端存磁盘返回 URL → WebSocket 发文件元数据 → 接收方通过 URL 下载。

---

## 3. 音视频通话是如何实现的？

**WebRTC + WebSocket 信令**架构：

- **WebSocket**（信令通道）：转发 `start_call`、`receive_call`、`reject_call`、SDP（offer/answer）、ICE Candidate、`PEER_LEAVE` 等控制命令。后端纯透传，不解析信令内容。
- **WebRTC**（媒体通道）：浏览器之间 P2P 直传音视频流（SRTP 加密），不经过服务器。

信令流程：A 发起 → start_call → B 收到通知 → B 点接收 → A 发 SDP Offer → B 回 SDP Answer → 交换 ICE Candidate → P2P 连接建立 → 媒体直传。

当前局限：仅支持 1v1，无 STUN/TURN 服务器配置（`ICE_CFG` 为空），跨网络可能连不通。

---

## 4. Channel 和 Kafka 模式的对比与优缺点

**为什么默认用 Channel**：本地开发零外部依赖，不需要部署 Kafka。

| 维度 | Channel | Kafka |
|---|---|---|
| 本质 | Go 进程内存队列 | 分布式消息队列 |
| 依赖 | 无 | Kafka broker |
| 延迟 | 微秒级（内存） | 毫秒级（网络+磁盘） |
| 持久化 | ❌ 进程崩溃丢失 | ✅ 磁盘持久化 + 副本 |
| 水平扩展 | ❌ 单进程 | ✅ 多实例消费 |
| 复杂度 | 极低 | 高 |

项目设计上两套模式共享相同的消息处理逻辑（存 MySQL → 推在线用户 → 写 Redis），切换只需改一行配置。

---

## 5. 什么是 WebRTC + WebSocket 信令架构？

```
信令服务器 (WebSocket)
   /              \
 控制命令 /        \ 控制命令
 /                  \
A ──── P2P 媒体流 ──── B
     (WebRTC 直连)
```

- **WebSocket 传送的内容**：呼叫控制（发起/接收/拒绝/挂断）+ SDP 媒体协商（编码、分辨率、带宽）+ ICE Candidate（NAT 穿透地址）。
- **WebRTC 传送的内容**：实际的音视频二进制流，浏览器底层 SRTP 加密。

连接建立后，媒体流完全不经过 Go 后端。

---

## 6. P2P 是什么意思？

**Peer to Peer，点对点直连**——两台设备直接通信，不经过中间服务器。

- **Client-Server**：所有数据必须经过服务器中转。
- **P2P**：服务器只帮忙"牵线"（交换地址和协商信息），连接建立后数据直传。

P2P 的难点是 NAT 穿透：多数设备在路由器后没有公网 IP，需要 STUN（获取公网地址）/TURN（实在打不通就中继）来辅助。本项目 `ICE_CFG` 为空，未配置 STUN/TURN，跨网络大概率失败。

---

## 7. 项目如何处理高并发？Go 语言如何处理高并发？

### Go 的并发模型：GMP

Go 使用 **goroutine**（用户态轻量线程，初始 ~2KB vs OS 线程 ~1MB），通过 GMP 调度器把海量 goroutine 分配到少量 OS 线程上执行。

- **G（Goroutine）**：任务
- **M（Machine）**：操作系统线程
- **P（Processor）**：调度器，连接 G 和 M

当一个 goroutine 阻塞（等数据库、读 channel），调度器立即切走换下一个——CPU 不空转。

### 项目中的并发设计

- **每个 HTTP 请求一个 goroutine**（Gin 框架自动）。
- **每个 WebSocket 连接两个 goroutine**：`Read()` 持续读消息、`Write()` 持续写消息。
- **单线程事件循环**：Login/Logout/Transmit 三个 Channel 汇入一个 `select` 循环，串行化对在线用户 Map 的操作。
- **缓冲 Channel 做背压**：`make(chan []byte, CHANNEL_SIZE)`，满时才拒绝新消息。
- **sync.Mutex 保护共享数据**：在线用户 Map 的读写加锁。

### 与线程池的关键区别

```
线程池 200 线程 + 10 万请求：
  只有 200 个请求能"发出数据库查询"
  其余 99,800 个在线程池队列排队，连数据库都没见到

Goroutine 10 万：
  全部 10 万个同时发出数据库查询
  谁的先返回就先处理谁
```

区别不是 CPU 利用率，而是**能同时处于"等待中"的请求数量不受线程数限制**。

---

## 8. Redis 和 MySQL 的数据一致性如何？

**不是强一致性，是最终一致性（且实现不完整）。**

### 当前模式：Cache-Aside

- **写**：先写 MySQL → 再写 Redis（无事务保证）。
- **读**：先查 Redis → 未命中查 MySQL → 回写 Redis（很多回写被注释掉了）。

### 主要问题

1. **大量缓存更新代码被注释**：更新用户信息后，旧缓存未失效，等 1 分钟 TTL 自然过期才能看到新数据。
2. **无事务包裹**：MySQL 写入成功 + Redis 写入失败时，缺少回滚或重试。
3. **缓存失效分散**：手动在各处删除，容易漏删。
4. **群成员双写**：存在 `group_info.members` JSON 和 `user_contact` 表两处，无事务保证一致性。

### 改进方向

- 取消注释恢复缓存更新代码。
- 写操作：MySQL 成功后至少主动删除 Redis 旧缓存（不等 TTL 过期）。
- 多步业务操作加数据库事务。
- 关系表加联合唯一索引防止重复。

---

## 9. Go 语言本身是幂等的吗？

**不是。** Go 提供的是并发安全（Mutex、Channel 保证多 goroutine 不打架），不是幂等。

- **并发安全**：多个 goroutine 同时读写同一数据不出脏数据 → Go 能帮忙。
- **幂等**：同一条消息处理 1 次和 10 次结果一样 → 纯业务逻辑，任何语言都要自己实现。

本项目缺少幂等：同一条消息重发 3 次，MySQL 存 3 行，对方收 3 条。需要给每条消息带客户端唯一 ID，处理前查重。

---

## 10. Java/Python 线程池和 Go Goroutine 的区别

**线程池不是"等完再处理下一个"**——OS 也会切换线程，CPU 不空转。

**真正的区别**：同时能"飞行中"的请求数量。

| | 线程池（200 线程） | Goroutine（10 万） |
|---|---|---|
| 能同时查数据库的请求 | 200 个 | 全部 10 万个 |
| 其余请求 | 排在线程池队列等空位 | 已全部发出，等返回 |
| 瓶颈 | 线程数限制并发度 | 数据库/Redis 连接数 |

Java 近年也加入了 Virtual Threads（JDK 21），原理和 goroutine 类似。

10 万个 goroutine，同一瞬间只有 CPU 核数个在真正执行，其余全在等待。Go runtime 在任何一个 goroutine 要"等"的时候立刻切走换下一个，把所有等待时间利用起来。
