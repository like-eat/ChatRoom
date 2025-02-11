# KamaChat

【代码随想录知识星球】项目分享-基于go+vue实现的聊天室+仿微信项目

# 项目概述
1. 简介：KamaChat 是一个前后端分离的即时通讯项目，具备后台管理、单聊群聊、联系人管理、多种消息（文本 / 文件 / 视频）处理、离线消息处理以及音视频通话等功能，旨在打造类似微信的聊天体验。
2. 技术栈：
   + 前端：Vue3、Vue Router、Vuex、WebSocket、Element - UI 等。
   + 后端：Go、Gin、GORM、GoRedis、WebSocket、Kafka、WebRTC、Zap 日志库等。


# 功能特性
1. 即时通讯功能
   + 单聊与群聊：支持一对一私密聊天和群组聊天，消息实时推送。
   + 联系人管理：可添加、删除、拉黑联系人，处理好友申请等。
   + 消息类型：支持文本、文件、音视频等多种类型消息的发送与接收。
   + 离线消息处理：确保用户离线时消息不丢失，上线后可正常接收。
2. 音视频通话：基于 WebRTC 实现 1 对 1 音视频通话，包括发起、拒绝、接收、挂断通话等功能。
3. 后台管理：具备后台管理界面，靓号用户可进行人员管控等维护操作。
4. 安全与验证：登录注册采用 SMS 短信验证方式，并支持 SSL 加密，保障用户信息安全。
5. 后台mysql数据库：使用 GORM 进行数据库操作，确保数据持久化存储。
6. 日志记录：使用 Zap 日志库记录系统运行日志，便于问题排查与性能监控。
7. 消息队列：使用 Kafka 处理消息队列，确保消息的高效传输与处理。
8. redis缓存：使用 GoRedis 进行缓存操作，提高系统性能。
9. WebSocket：使用 WebSocket 实现实时消息推送，保证消息的实时性。

# 安装与运行
## 后端

### 环境要求

Go 环境（作者环境为1.20.11 windows/amd64），使用goland来运行后端项目

### 克隆项目

git clone git@github.com:youngyangyang04/KamaChat.git

### 安装依赖

进入后端项目目录，执行 go mod tidy

### 配置文件

```toml
[mainConfig]
appName = "your app name"
host = "127.0.0.1"
port = 8000

[mysqlConfig]
host = "127.0.0.1"
port = 3306
user = "root"
password = "123456"
databaseName = "your database name"

[redisConfig]
host = "127.0.0.1"
port = 6379
password = "123456"
db = 0

[authCodeConfig]
accessKeyID = "your accessKeyID in alibaba cloud"
accessKeySecret = "your accessKeySecret in alibaba cloud"
signName = "阿里云短信测试"
templateCode = "SMS_154950909"

[logConfig]
logPath = "your log path"

[kafkaConfig]
messageMode = "channel"# 消息模式 channel or kafka
hostPort = "127.0.0.1:9092" # "127.0.0.1:9092,127.0.0.1:9093,127.0.0.1:9094" 多个kafka服务器
loginTopic = "login"
chatTopic = "chat_message"
logoutTopic = "logout"
partition = 0 # kafka partition
timeout = 1 # 单位秒

[staticSrcConfig]
staticAvatarPath = "./static/avatars"
staticFilePath = "./static/files"
```

进入/configs/config.toml文件，修改自己的配置，包括后端服务器配置，mysql配置，redis配置，短信验证配置（这里使用到阿里云的短信验证功能，具体可以看知识星球的操作流程），输入日志的配置，kafka配置，静态文件路径配置。

### 运行项目

通过goland来运行项目，你可以看到如下输出：

![image-20250211114212775](C:\Users\HP\AppData\Roaming\Typora\typora-user-images\image-20250211114212775.png)



## 前端

### 环境准备

在开始安装和运行项目之前，你需要确保本地环境已经安装了以下工具：

**Node.js**：Vue 项目依赖 Node.js 运行环境，你可以从 [Node.js 官方网站](https://nodejs.org/) 下载并安装适合你操作系统的版本。建议安装 LTS（长期支持）版本，安装完成后，在命令行中输入 `node -v` 和 `npm -v` 来验证安装是否成功。如果npm运行速度慢，可以使用cnpm。

**Yarn**：Yarn 是一个快速、可靠、安全的依赖管理工具，你可以使用以下命令来全局安装 Yarn：

```bash
npm install -g yarn
```

安装完成后，输入 `yarn -v` 验证 Yarn 是否安装成功。 

**Vue CLI**：Vue CLI 是一个基于 Vue.js 进行快速项目搭建的工具，使用以下命令全局安装 Vue CLI：

```bash
yarn global add @vue/cli
```

安装完成后，输入 `vue --version` 验证 Vue CLI 是否安装成功。

### 安装依赖

后端git clone后，找到其中的web项目目录。通过终端进入到该目录中。

使用 Yarn 安装项目所需的依赖。在命令行中执行以下命令：

```bash
yarn install
```

`yarn install` 命令会读取项目根目录下的 `package.json` 文件，根据其中定义的依赖信息，从 npm 注册表或其他指定的源下载并安装所有依赖项到项目的 `node_modules` 目录中。这个过程可能需要一些时间，具体取决于你的网络速度和依赖项的数量。

### 运行项目

安装完依赖并完成配置后，你可以使用以下命令启动开发服务器：

```bash
yarn serve
```

`yarn serve` 命令会调用 Vue CLI 的开发服务器，它会编译项目代码，并在本地启动一个开发服务器。启动成功后，你会在命令行中看到类似以下的输出：

![image-20250211114114697](C:\Users\HP\AppData\Roaming\Typora\typora-user-images\image-20250211114114697.png)

### 停止项目

如果你想停止开发服务器，可以在命令行中按下 `Ctrl + C` 组合键，然后选择 `y` 确认停止服务器。

### 新建项目

```bash
vue create kama-chat-server --registry=https://registry.npmmirror.com
```
vue-cli创建项目，选择默认配置，项目名称为kama-chat-server，使用npm镜像源。


```bash
Vue CLI v5.0.8
? Please pick a preset: Manually select features
? Check the features needed for your project: Babel, Router, Vuex, Linter 
? Choose a version of Vue.js that you want to start the project with 3.x
? Use history mode for router? (Requires proper server setup for index 
fallback in production) Yes
? Pick a linter / formatter config: Basic
? Pick additional lint features: Lint on save
? Where do you prefer placing config for Babel, ESLint, etc.? In dedicated 
config files
? Save this as a preset for future projects? No
```
以上为创建项目时的配置，跟着选就行了。

```bash
yarn add element-plus -g
yarn add axios -g
yarn add @element-plus/icons-vue -g
```
配置element-plus，axios，element-plus的图标等。

```bash
yarn serve
```
启动项目。



# 项目结构

## 后端

```
kama-chat-server/
├── api/
│   └── v1/
│       └── chatroom_controller.go
│       └── controller.go
│       └── group_info_controller.go
│       └── message_controller.go
│       └── session_controller.go
│       └── user_contact_controller.go
│       └── user_info_controller.go
│       └── ws_controller.go
├── cmd/
│   └── kama-chat-server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── dao/
│   │   └── gorm.go
│   ├── dto/
│   │   ├── request/
│   │   │   └── ......
│   │   └── respond/
│   │   │   └── ......
│   ├── https_server/
│   │   └── https_server.go
│   ├── model/
│   │   ├── contact_apply.go
│   │   ├── group_info.go
│   │   ├── message.go
│   │   ├── session.go
│   │   ├── user_contact.go
│   │   └── user_info.go
│   └── service/
│       ├── chat/
│       │   ├── client.go
│       │   ├── kafka_server.go
│       │   └── server.go
│       ├── gorm/
│       │   ├── chatroom_service.go
│       │   ├── group_info_service.go
│       │   ├── message_service.go
│       │   ├── session_service.go
│       │   ├── user_contact_service.go
│       │   └── user_info_service.go
│       ├── kafka/
│       │   └── kafka_service.go
│       ├── redis/
│       │   └── redis_service.go
│       └── sms/
│           ├── local/
│           │   └── user_info_service_local.go
│           └── auth_code_service.go
├── logs/
│   └── test.log
├── pkg/
│   ├── constants/
│   │   └── constants.go
│   ├── enum/
│   │   ├── contact/
│   │   ├── contact_apply/
│   │   ├── group_info/
│   │   ├── message/
│   │   ├── session/
│   │   └── user_info/
│   ├── ssl/
│   │   ├── xxx.pem
│   │   ├── xxx-key.pem
│   │   └── tls_handler.go
│   ├── util/
│   │   └── random/
│   │       └── random_int.go
│   └── zlog/
│       └── logger.go
├── configs/
│   ├── config.toml
│   └── config_local.toml
├── static/
│   ├── avatars/
│   │   └── ......
│   └── files/
│   │   └── ......
├── web/
│   └── (前端项目结构)
├── .gitignore
├── go.mod
├── go.sum
└── README.md

```

## 前端

```
web/chat-server/
├── src/
│   ├── assets/
│   │   ├── cert/
│   │   │   ├── xxx.pem
│   │   │   ├── xxx-key.pem
│   │   │   └── mkcert.exe
│   │   ├── css/
│   │   │   └── chat.css
│   │   ├── img/
│   │   │   └── chat_server_background.jpg
│   │   ├── js/
│   │   │   ├── random.js
│   │   │   └── valid.js
│   ├── components/
│   │   ├── ContactListModal.vue
│   │   ├── DeleteGroupModal.vue
│   │   ├── DeleteUserModal.vue
│   │   ├── DisableGroupModal.vue
│   │   ├── DisableUserModal.vue
│   │   ├── Modal.vue
│   │   ├── NavigationModal.vue
│   │   ├── SetAdminModal.vue
│   │   ├── SmallModal.vue
│   │   └── VideoModal.vue
│   ├── router/
│   │   └── index.js
│   ├── store/
│   │   └── index.js
│   ├── views/
│   │   ├── access/
│   │   │   ├── Login.vue
│   │   │   ├── Register.vue
│   │   │   └── SmsLogin.vue
│   │   ├── chat/
│   │   │   ├── contact/
│   │   │   │   ├── ContactChat.vue
│   │   │   │   └── ContactList.vue
│   │   │   ├── session/
│   │   │   │   └── SessionList.vue
│   │   │   └── user/
│   │   │       └── OwnInfo.vue
│   │   ├── manager/
│   │       └── Manager.vue
│   ├── App.vue
│   └── main.js
├── .gitignore
├── package.json
├── README.md
└── vue.config.js
```

# docs

在/docs/业务逻辑.md中，介绍了具体的业务设计，对业务有问题的同学可以查看了解一下。

# todoList

-  多对多群聊

-  nginx分布式部署

  
