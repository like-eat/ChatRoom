#!/bin/bash

# 更新系统软件包
sudo apt update && sudo apt upgrade -y

# 安装 MySQL
echo "Installing MySQL..."
sudo apt install mysql-server -y

# 配置 MySQL 安全
sudo mysql_secure_installation

# 启动并启用 MySQL 服务
sudo systemctl start mysql
sudo systemctl enable mysql

# 安装 Redis
echo "Installing Redis..."
sudo apt install redis-server -y

# 配置 Redis
sudo nano /etc/redis/redis.conf  # 修改 bind 127.0.0.1 改为 bind 0.0.0.0（如果需要外部访问）

# 启动并启用 Redis 服务
sudo systemctl restart redis
sudo systemctl enable redis

# 安装 Go
echo "Installing Go..."
sudo apt install golang-go -y

# 设置 Go 环境变量
echo "Configuring Go environment..."
mkdir -p ~/go
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 Vue.js 开发环境
echo "Installing Vue.js development environment..."
sudo apt install nodejs npm -y

# 安装 Yarn
sudo npm install -g yarn

# 安装 Vue CLI
sudo npm install -g @vue/cli

# 输出完成信息
echo "Deployment complete!"