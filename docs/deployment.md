# 部署文档

## 概述

本文档介绍如何在生产环境中部署美术作品收集系统。系统采用 Go + Gin + MySQL + Redis 架构，支持多种部署方式。

## 系统要求

### 硬件要求

**最低配置**:
- CPU: 2 核
- 内存: 4GB
- 磁盘: 20GB（不包括上传文件存储）

**推荐配置**:
- CPU: 4 核
- 内存: 8GB
- 磁盘: 100GB SSD（包括上传文件存储）

### 软件要求

- **操作系统**: Linux (Ubuntu 20.04+, CentOS 7+) 或其他支持 Docker 的系统
- **Go**: 1.21+ (如果从源码编译)
- **MySQL**: 8.0+
- **Redis**: 7.0+
- **Nginx**: 1.18+ (可选，用于反向代理)
- **Docker**: 20.10+ (可选，用于容器化部署)

## 部署方式

### 方式一：传统部署（推荐用于生产环境）

#### 1. 准备服务器

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装必要工具
sudo apt install -y git curl wget vim
```

#### 2. 安装 MySQL

```bash
# 安装 MySQL 8.0
sudo apt install -y mysql-server

# 启动 MySQL
sudo systemctl start mysql
sudo systemctl enable mysql

# 安全配置
sudo mysql_secure_installation

# 创建数据库和用户
sudo mysql -u root -p
```

在 MySQL 中执行：

```sql
CREATE DATABASE art_collection CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'artcollection'@'localhost' IDENTIFIED BY 'your_strong_password';
GRANT ALL PRIVILEGES ON art_collection.* TO 'artcollection'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

#### 3. 安装 Redis

```bash
# 安装 Redis
sudo apt install -y redis-server

# 配置 Redis（可选：设置密码）
sudo vim /etc/redis/redis.conf
# 找到 # requirepass foobared 并修改为：
# requirepass your_redis_password

# 重启 Redis
sudo systemctl restart redis
sudo systemctl enable redis

# 测试连接
redis-cli ping
```

#### 4. 安装 Go（如果从源码编译）

```bash
# 下载 Go 1.21
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz

# 解压
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

#### 5. 部署应用

```bash
# 创建应用目录
sudo mkdir -p /opt/art-collection
sudo chown $USER:$USER /opt/art-collection
cd /opt/art-collection

# 克隆代码（或上传编译好的二进制文件）
git clone <your-repo-url> .

# 安装依赖
go mod download

# 编译应用
go build -o bin/server cmd/server/main.go

# 创建必要目录
mkdir -p uploads logs config

# 复制配置文件
cp config/config.example.yaml config/config.yaml
```

#### 6. 配置应用

编辑 `config/config.yaml`：

```yaml
server:
  port: 8080
  mode: release  # 生产环境使用 release 模式

database:
  mysql:
    host: localhost
    port: 3306
    user: artcollection
    password: your_strong_password
    dbname: art_collection
    max_idle_conns: 10
    max_open_conns: 100
  redis:
    host: localhost
    port: 6379
    password: your_redis_password
    db: 0

jwt:
  secret: your-very-long-and-random-secret-key-at-least-32-bytes
  expire_hours: 24

upload:
  path: /opt/art-collection/uploads
  max_size: 10485760  # 10MB

email:
  smtp_host: smtp.example.com
  smtp_port: 587
  username: noreply@example.com
  password: your_email_password
  from: noreply@example.com

log:
  level: info
  file: /opt/art-collection/logs/app.log
```

**重要**: 
- 修改 `jwt.secret` 为强随机密钥（至少 32 字节）
- 修改数据库密码
- 配置邮件服务器信息

#### 7. 初始化数据库

```bash
# 运行初始化脚本
mysql -h localhost -u artcollection -p art_collection < scripts/init_db.sql

# 修改默认管理员密码（重要！）
# 使用工具生成新密码哈希
go run scripts/generate_password.go "YourNewAdminPassword"

# 更新数据库中的管理员密码
mysql -h localhost -u artcollection -p art_collection
```

在 MySQL 中执行：

```sql
UPDATE users SET password = '$2a$12$...' WHERE email = 'admin@example.com';
EXIT;
```

#### 8. 创建 Systemd 服务

创建服务文件 `/etc/systemd/system/art-collection.service`：

```ini
[Unit]
Description=Art Collection System
After=network.target mysql.service redis.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/art-collection
ExecStart=/opt/art-collection/bin/server
Restart=on-failure
RestartSec=5s

# 环境变量（可选）
Environment="GIN_MODE=release"

# 日志
StandardOutput=append:/opt/art-collection/logs/stdout.log
StandardError=append:/opt/art-collection/logs/stderr.log

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
# 重新加载 systemd
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start art-collection

# 设置开机自启
sudo systemctl enable art-collection

# 查看状态
sudo systemctl status art-collection

# 查看日志
sudo journalctl -u art-collection -f
```

#### 9. 配置 Nginx 反向代理

安装 Nginx：

```bash
sudo apt install -y nginx
```

创建 Nginx 配置文件 `/etc/nginx/sites-available/art-collection`：

```nginx
# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

# HTTPS 配置
server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL 证书配置
    ssl_certificate /etc/ssl/certs/your-domain.crt;
    ssl_certificate_key /etc/ssl/private/your-domain.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # 日志
    access_log /var/log/nginx/art-collection-access.log;
    error_log /var/log/nginx/art-collection-error.log;

    # 客户端最大上传大小
    client_max_body_size 10M;

    # 反向代理到 Gin 应用
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 禁止直接访问上传目录
    location /uploads {
        deny all;
        return 403;
    }

    # Gzip 压缩
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
}
```

启用配置：

```bash
# 创建软链接
sudo ln -s /etc/nginx/sites-available/art-collection /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重启 Nginx
sudo systemctl restart nginx
sudo systemctl enable nginx
```

#### 10. 配置 SSL 证书（使用 Let's Encrypt）

```bash
# 安装 Certbot
sudo apt install -y certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

#### 11. 配置防火墙

```bash
# 允许 HTTP 和 HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# 允许 SSH（如果还没有）
sudo ufw allow 22/tcp

# 启用防火墙
sudo ufw enable

# 查看状态
sudo ufw status
```

---

### 方式二：Docker 部署

#### 1. 创建 Dockerfile

在项目根目录创建 `Dockerfile`：

```dockerfile
# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/server .
COPY --from=builder /app/config ./config

# 创建必要目录
RUN mkdir -p uploads logs

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./server"]
```

#### 2. 创建 docker-compose.yml

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: art-collection-mysql
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: root123456
      MYSQL_DATABASE: art_collection
      MYSQL_USER: artcollection
      MYSQL_PASSWORD: artcollection123
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./scripts/init_db.sql:/docker-entrypoint-initdb.d/init.sql
    command: --default-authentication-plugin=mysql_native_password --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

  redis:
    image: redis:7-alpine
    container_name: art-collection-redis
    restart: always
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes

  app:
    build: .
    container_name: art-collection-app
    restart: always
    ports:
      - "8080:8080"
    depends_on:
      - mysql
      - redis
    environment:
      - GIN_MODE=release
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=artcollection
      - DB_PASSWORD=artcollection123
      - DB_NAME=art_collection
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    volumes:
      - ./uploads:/app/uploads
      - ./logs:/app/logs
      - ./config/config.yaml:/app/config/config.yaml

  nginx:
    image: nginx:alpine
    container_name: art-collection-nginx
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - app

volumes:
  mysql_data:
  redis_data:
```

#### 3. 部署

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f app

# 停止服务
docker-compose down

# 重启服务
docker-compose restart app
```

---

### 方式三：Docker Swarm 集群部署

适用于高可用和负载均衡场景。

#### 1. 初始化 Swarm

```bash
# 在主节点初始化
docker swarm init --advertise-addr <MANAGER-IP>

# 在工作节点加入集群
docker swarm join --token <TOKEN> <MANAGER-IP>:2377
```

#### 2. 创建 docker-compose-swarm.yml

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root123456
      MYSQL_DATABASE: art_collection
      MYSQL_USER: artcollection
      MYSQL_PASSWORD: artcollection123
    volumes:
      - mysql_data:/var/lib/mysql
    deploy:
      replicas: 1
      placement:
        constraints:
          - node.role == manager
    networks:
      - backend

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    deploy:
      replicas: 1
    networks:
      - backend

  app:
    image: your-registry/art-collection:latest
    environment:
      - GIN_MODE=release
      - DB_HOST=mysql
      - REDIS_HOST=redis
    volumes:
      - uploads:/app/uploads
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure
    networks:
      - backend
      - frontend

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
    deploy:
      replicas: 2
      placement:
        constraints:
          - node.role == worker
    networks:
      - frontend

volumes:
  mysql_data:
  redis_data:
  uploads:

networks:
  frontend:
  backend:
```

#### 3. 部署到 Swarm

```bash
# 部署 stack
docker stack deploy -c docker-compose-swarm.yml art-collection

# 查看服务
docker stack services art-collection

# 查看日志
docker service logs -f art-collection_app

# 扩容
docker service scale art-collection_app=5

# 删除 stack
docker stack rm art-collection
```

---

## 环境变量配置

应用支持通过环境变量覆盖配置文件中的设置：

| 环境变量 | 说明 | 默认值 |
|---------|------|--------|
| `SERVER_PORT` | 服务器端口 | 8080 |
| `GIN_MODE` | Gin 运行模式 | debug |
| `DB_HOST` | MySQL 主机 | localhost |
| `DB_PORT` | MySQL 端口 | 3306 |
| `DB_USER` | MySQL 用户名 | root |
| `DB_PASSWORD` | MySQL 密码 | - |
| `DB_NAME` | 数据库名 | art_collection |
| `REDIS_HOST` | Redis 主机 | localhost |
| `REDIS_PORT` | Redis 端口 | 6379 |
| `REDIS_PASSWORD` | Redis 密码 | - |
| `JWT_SECRET` | JWT 密钥 | - |
| `UPLOAD_PATH` | 上传目录 | ./uploads |
| `UPLOAD_MAX_SIZE` | 最大文件大小（字节） | 10485760 |
| `LOG_LEVEL` | 日志级别 | info |
| `LOG_FILE` | 日志文件路径 | ./logs/app.log |

示例：

```bash
export DB_HOST=192.168.1.100
export DB_PASSWORD=strong_password
export JWT_SECRET=your-very-long-secret-key
./server
```

---

## 数据备份

### MySQL 备份

#### 手动备份

```bash
# 备份数据库
mysqldump -h localhost -u artcollection -p art_collection > backup_$(date +%Y%m%d_%H%M%S).sql

# 恢复数据库
mysql -h localhost -u artcollection -p art_collection < backup_20251021_100000.sql
```

#### 自动备份脚本

创建 `/opt/scripts/backup-mysql.sh`：

```bash
#!/bin/bash

BACKUP_DIR="/opt/backups/mysql"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/art_collection_$TIMESTAMP.sql"

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份数据库
mysqldump -h localhost -u artcollection -p'your_password' art_collection > $BACKUP_FILE

# 压缩备份
gzip $BACKUP_FILE

# 删除 7 天前的备份
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete

echo "Backup completed: $BACKUP_FILE.gz"
```

添加到 crontab：

```bash
# 每天凌晨 2 点备份
0 2 * * * /opt/scripts/backup-mysql.sh >> /var/log/mysql-backup.log 2>&1
```

### 文件备份

```bash
# 备份上传文件
tar -czf uploads_backup_$(date +%Y%m%d).tar.gz /opt/art-collection/uploads

# 使用 rsync 同步到远程服务器
rsync -avz /opt/art-collection/uploads/ user@backup-server:/backups/uploads/
```

---

## 监控和日志

### 应用日志

日志文件位置：
- 应用日志: `/opt/art-collection/logs/app.log`
- Systemd 日志: `journalctl -u art-collection`
- Nginx 访问日志: `/var/log/nginx/art-collection-access.log`
- Nginx 错误日志: `/var/log/nginx/art-collection-error.log`

查看实时日志：

```bash
# 应用日志
tail -f /opt/art-collection/logs/app.log

# Systemd 日志
sudo journalctl -u art-collection -f

# Nginx 日志
tail -f /var/log/nginx/art-collection-access.log
```

### 日志轮转

创建 `/etc/logrotate.d/art-collection`：

```
/opt/art-collection/logs/*.log {
    daily
    rotate 30
    compress
    delaycompress
    notifempty
    create 0640 www-data www-data
    sharedscripts
    postrotate
        systemctl reload art-collection > /dev/null 2>&1 || true
    endscript
}
```

### 性能监控

使用 Prometheus + Grafana 监控（可选）：

1. 在应用中集成 Prometheus metrics
2. 配置 Prometheus 抓取 metrics
3. 在 Grafana 中创建仪表板

监控指标：
- API 响应时间
- 请求成功率
- 数据库连接池状态
- Redis 连接状态
- 系统资源使用（CPU、内存、磁盘）

---

## 性能优化

### 数据库优化

```sql
-- 添加索引
CREATE INDEX idx_artworks_activity_user ON artworks(activity_id, user_id);
CREATE INDEX idx_artworks_review_status ON artworks(review_status);
CREATE INDEX idx_artworks_created_at ON artworks(created_at);

-- 优化查询
ANALYZE TABLE users;
ANALYZE TABLE activities;
ANALYZE TABLE artworks;
```

### Redis 优化

编辑 `/etc/redis/redis.conf`：

```
# 最大内存
maxmemory 2gb

# 内存淘汰策略
maxmemory-policy allkeys-lru

# 持久化
save 900 1
save 300 10
save 60 10000
```

### Nginx 优化

```nginx
# 工作进程数
worker_processes auto;

# 连接数
events {
    worker_connections 4096;
    use epoll;
}

# 缓存
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=my_cache:10m max_size=1g inactive=60m;

# 启用缓存
location /api/v1/activities {
    proxy_cache my_cache;
    proxy_cache_valid 200 5m;
    proxy_pass http://127.0.0.1:8080;
}
```

---

## 安全加固

### 1. 系统安全

```bash
# 禁用 root SSH 登录
sudo vim /etc/ssh/sshd_config
# PermitRootLogin no

# 配置防火墙
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# 安装 fail2ban
sudo apt install -y fail2ban
sudo systemctl enable fail2ban
```

### 2. 数据库安全

```sql
-- 限制用户权限
REVOKE ALL PRIVILEGES ON *.* FROM 'artcollection'@'localhost';
GRANT SELECT, INSERT, UPDATE, DELETE ON art_collection.* TO 'artcollection'@'localhost';
FLUSH PRIVILEGES;

-- 禁用远程 root 登录
DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');
FLUSH PRIVILEGES;
```

### 3. 应用安全

- 使用强 JWT 密钥（至少 32 字节随机字符）
- 定期更新依赖包：`go get -u ./...`
- 启用 HTTPS
- 配置 CORS 白名单
- 实施速率限制
- 定期审计日志

---

## 故障排查

### 应用无法启动

```bash
# 检查日志
sudo journalctl -u art-collection -n 50

# 检查配置文件
cat /opt/art-collection/config/config.yaml

# 检查端口占用
sudo netstat -tlnp | grep 8080

# 检查文件权限
ls -la /opt/art-collection
```

### 数据库连接失败

```bash
# 测试 MySQL 连接
mysql -h localhost -u artcollection -p

# 检查 MySQL 状态
sudo systemctl status mysql

# 查看 MySQL 日志
sudo tail -f /var/log/mysql/error.log
```

### Redis 连接失败

```bash
# 测试 Redis 连接
redis-cli ping

# 检查 Redis 状态
sudo systemctl status redis

# 查看 Redis 日志
sudo tail -f /var/log/redis/redis-server.log
```

### 文件上传失败

```bash
# 检查上传目录权限
ls -la /opt/art-collection/uploads

# 修复权限
sudo chown -R www-data:www-data /opt/art-collection/uploads
sudo chmod -R 755 /opt/art-collection/uploads

# 检查磁盘空间
df -h
```

---

## 升级和回滚

### 升级应用

```bash
# 备份当前版本
cp /opt/art-collection/bin/server /opt/art-collection/bin/server.backup

# 停止服务
sudo systemctl stop art-collection

# 更新代码
cd /opt/art-collection
git pull origin main

# 编译新版本
go build -o bin/server cmd/server/main.go

# 运行数据库迁移（如果有）
# ./bin/server migrate

# 启动服务
sudo systemctl start art-collection

# 检查状态
sudo systemctl status art-collection
```

### 回滚

```bash
# 停止服务
sudo systemctl stop art-collection

# 恢复旧版本
cp /opt/art-collection/bin/server.backup /opt/art-collection/bin/server

# 恢复数据库（如果需要）
mysql -h localhost -u artcollection -p art_collection < backup.sql

# 启动服务
sudo systemctl start art-collection
```

---

## 高可用架构

### 架构图

```
                    ┌─────────────┐
                    │   用户      │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │  负载均衡   │
                    │  (Nginx)    │
                    └──────┬──────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   ┌────▼────┐       ┌────▼────┐       ┌────▼────┐
   │ App 1   │       │ App 2   │       │ App 3   │
   └────┬────┘       └────┬────┘       └────┬────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   ┌────▼────┐       ┌────▼────┐       ┌────▼────┐
   │ MySQL   │       │ Redis   │       │  NFS    │
   │ Master  │       │ Cluster │       │ Storage │
   └────┬────┘       └─────────┘       └─────────┘
        │
   ┌────▼────┐
   │ MySQL   │
   │ Slave   │
   └─────────┘
```

### 配置要点

1. **负载均衡**: 使用 Nginx 或 HAProxy
2. **应用集群**: 多个应用实例，无状态设计
3. **数据库主从**: MySQL 主从复制，读写分离
4. **Redis 集群**: Redis Sentinel 或 Cluster
5. **共享存储**: NFS 或对象存储（如 S3）存储上传文件

---

## 总结

本文档涵盖了美术作品收集系统的多种部署方式和最佳实践。根据实际需求选择合适的部署方式：

- **小型项目**: 传统部署（单服务器）
- **中型项目**: Docker Compose 部署
- **大型项目**: Docker Swarm 或 Kubernetes 集群部署

部署后务必：
1. 修改默认管理员密码
2. 配置强 JWT 密钥
3. 启用 HTTPS
4. 配置定期备份
5. 设置监控和告警
6. 定期更新依赖和系统补丁

如有问题，请参考故障排查章节或查看应用日志。
