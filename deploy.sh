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

# 自动创建数据库
echo "Creating database 'kama_chat_server'..."
sudo mysql -u root -p <<EOF
CREATE DATABASE kama_chat_server;
EOF

# 安装 Redis
echo "Installing Redis..."
sudo apt install redis-server -y

# 配置 Redis
# sudo nano /etc/redis/redis.conf  # 修改 bind 127.0.0.1 改为 bind 0.0.0.0（如果需要外部访问）

# 启动并启用 Redis 服务
sudo systemctl restart redis
sudo systemctl enable redis

# 卸载旧版本 Node.js 和 npm
# echo "Uninstalling previous versions of Node.js and npm..."
# sudo apt remove --purge -y nodejs npm

# # 安装 Node.js 版本管理工具（nvm）
# echo "Installing Node Version Manager (nvm)..."
# rm -rf ~/.nvm
# export NVM_NODE_MIRROR=https://npmmirror.com/mirrors/node/
# export NVM_NPM_MIRROR=https://npmmirror.com/mirrors/npm/
# curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash



# # 设置 NVM_DIR 环境变量（避免重复写入 ~/.bashrc）
# # if ! grep -q "export NVM_DIR=~/.nvm" ~/.bashrc; then
# #     echo 'export NVM_DIR="$HOME/.nvm"' >> ~/.bashrc
# #     echo '[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # This loads nvm' >> ~/.bashrc
# #     echo '[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # This loads nvm bash_completion' >> ~/.bashrc
# # fi


# # 手动设置 NVM_DIR 并加载 nvm
# export NVM_DIR="$HOME/.nvm"
# [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"  # 加载 nvm
# [ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"  # 加载 bash_completion

# source ~/.bashrc

# # 检查 nvm 是否存在
# if ! command -v nvm &> /dev/null; then
#     echo "nvm could not be found. Please ensure it is installed and added to your PATH."
#     exit 1
# fi

# # 安装指定版本的 Node.js（例如 v16.x）
# echo "Installing Node.js v16.x..."
# nvm install 16
# nvm use 16


# 加载环境变量
# source ~/.bashrc

# 安装 Go
#  echo "Installing Go..."
#  wget https://mirrors.aliyun.com/golang/go1.20.linux-amd64.tar.gz
#  sudo tar -C /usr/local -xzf go1.20.linux-amd64.tar.gz

#  cp -r /usr/local/go/bin/* /usr/bin
 # 设置 Go 环境变量
#  echo "Configuring Go environment..."

# export PATH=$PATH:/usr/local/go/bin

 # 设置 Go 环境变量（避免重复写入 ~/.bashrc）
#  if ! grep -q "export GOPATH=$HOME/go" ~/.bashrc; then
#      echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
#  fi
#  source ~/.bashrc


# # 配置 Go 代理
# echo "Configuring Go proxy..."
#  go env -w GOPROXY=https://goproxy.cn,direct

# 安装 Vue.js 开发环境
echo "Installing Vue.js development environment..."
sudo apt install npm -y



# 方案1：使用 npm 安装 Yarn
# sudo npm install -g yarn

# 方案2：使用cnpm 安装 Yarn
sudo npm install -g cnpm --registry=https://registry.npmjs.org
sudo cnpm install -g yarn

# 安装 Vue CLI
sudo cnpm install -g @vue/cli

# 重新安装项目依赖
cd ~/project/KamaChat/web/chat-server

# yarn cache clean
# rm -rf node_modules

yarn install

#打包项目成dist，放到/var/www/html/，此时就可以通过云服务器的公网ip看到前端页面了
rm -rf /var/www/html/* 
rm -rf /root/project/KamaChat/web/chat-server/dist
yarn build
sudo cp -r /root/project/KamaChat/web/chat-server/dist/* /var/www/html # 改成自己的项目路径
sudo chmod -R 755 /var/www/html
sudo chown -R www-data:www-data /var/www/html

cd ~/project/KamaChat

# 安装 ssl 模块
echo "Installing ssl..."
sudo apt-get install openssl
sudo apt-get install libssl-dev

# # 创建根密钥，生成证书签名请求 (CSR)，创建根证书
# openssl genrsa -out /etc/ssl/private/root.key 2048
# openssl req -new -key /etc/ssl/private/root.key -out /etc/ssl/certs/root.csr
# openssl x509 -req -in /etc/ssl/certs/root.csr -out /etc/ssl/certs/root.crt -signkey /etc/ssl/private/root.key -CAcreateserial -days 3650

# # 生成服务器密钥，生成服务器证书签名请求 (CSR)，创建服务器证书扩展文件
# openssl genrsa -out /etc/ssl/private/server.key 2048
# openssl req -new -key /etc/ssl/private/server.key -out /etc/ssl/certs/server.csr
# sudo nano v3.ext
# # 内容如下
# # authorityKeyIdentifier=keyid,issuer
# # basicConstraints=CA:FALSE
# # keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
# # subjectAltName = @alt_names
# # [alt_names]
# # IP.1 = 123.56.164.220 # 你的云服务器地址

# # 使用根证书为服务器证书签名
# openssl x509 -req -in /etc/ssl/certs/server.csr -CA /etc/ssl/certs/root.crt -CAkey /etc/ssl/private/root.key -CAcreateserial -out /etc/ssl/certs/server.crt -days 500 -sha256 -extfile v3.ext


# 打开Apache2配置文件
sudo nano /etc/apache2/sites-enabled/000-default.conf
# 添加如下内容
# <VirtualHost *:443>
#     ServerAdmin webmaster@localhost
#     DocumentRoot /var/www/html

#     SSLEngine on
#     SSLProxyEngine on

#     # 替换为您的自签名证书路径
#     SSLCertificateFile /etc/ssl/certs/server.crt
#     SSLCertificateKeyFile /etc/ssl/private/server.key

#     # 如果有中间证书，添加以下行
#     # SSLCertificateChainFile /path/to/your_intermediate.crt

#     ErrorLog ${APACHE_LOG_DIR}/error.log
#     CustomLog ${APACHE_LOG_DIR}/access.log combined

#     # 以下配置可选，用于启用 HTTP 到 HTTPS 重定向，也可以把这段添加到80端口那儿
#     <IfModule mod_rewrite.c>
#         RewriteEngine On
#         RewriteCond %{HTTPS} off
#         RewriteRule ^/?(.*) https://%{SERVER_NAME}%{REQUEST_URI} [R=301,L]
#     </IfModule>
# </VirtualHost>

# 启用ssl模块，启用ssl站点，重启服务
sudo a2enmod ssl
sudo a2ensite 000-default.conf
sudo systemctl restart apache2


# 受信任的证书颁发机构生成证书
# echo "Installing certbot..."
# sudo apt install certbot python3-certbot-apache


# sudo certbot --apache -d kamachat.top

# # 自动续期
# sudo certbot renew --dry-run

# 配置turn服务器
# echo "Installing coturn..."
# sudo apt install coturn
# sudo nano /etc/coturn/coturn.conf
# # 配置以下参数
# listening-ip=0.0.0.0

# # 外部 IP 地址（替换为你的服务器公网 IP）
# external-ip=123.56.164.220

# # 监听端口
# listening-port=3478

# # 用户名和密码（替换为你的用户名和密码）
# user=root:dancernumber1

# # SSL 证书路径（如果需要加密通信）
# tls-certificate=/etc/ssl/certs/server.crt
# tls-private-key=/etc/ssl/private/server.key

# sudo systemctl start coturn
# sudo systemctl enable coturn


# 将后端打包部署
# cd ~/project/KamaChat/cmd/kama_chat_server
# go build -o kama_chat_backend main.go
# sudo cp kama_chat_backend /usr/local/bin/

# sudo nano /etc/systemd/system/kama_chat_backend.service
# 配置以下内容
# [Unit]
# Description=kama chat service
# After=network.target

# [Service]
# User=kama_chat  # 替换为你的用户名
# Group=kama_chat  # 替换为你的用户名
# WorkingDirectory=/root/project/KamaChat/cmd/kama_chat_server  # 替换为你的项目路径
# ExecStart=/usr/local/bin/kama_chat_backend  # 替换为你的可执行文件路径
# Restart=on-failure
# RestartSec=5

# [Install]
# WantedBy=multi-user.target

# sudo systemctl daemon-reload
# sudo systemctl start kama_chat_backend
# sudo systemctl enable kama_chat_backend

# 输出完成信息
echo "Deployment complete!"

