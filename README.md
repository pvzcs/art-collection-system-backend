# 美术作品投稿系统 (Art Collection System)

基于活动的美术作品管理平台，允许管理员创建活动并管理用户提交的美术作品。系统采用严格的审核机制和访问控制，确保只有授权用户（管理员和作品作者）可以查看已审核的作品。

## 功能特性

### 用户功能

- **用户注册与认证**: 邮箱验证码注册，JWT 令牌认证
- **个人信息管理**: 更新昵称、修改密码
- **作品上传**: 上传美术作品到指定活动，支持上传数量限制
- **个人空间**: 查看自己的所有作品（包括已审核和未审核）
- **作品管理**: 删除自己上传的作品

### 管理员功能

- **活动管理**: 创建、更新、删除活动，设置截止日期和上传限制
- **作品审核**: 审核用户提交的作品，支持单个和批量审核
- **审核队列**: 查看所有待审核作品列表
- **用户管理**: 查看用户列表、更新用户角色、查看用户统计信息
- **全局访问**: 访问所有作品和用户数据

### 安全特性

- **密码加密**: 使用 Bcrypt (cost 12) 加密存储密码
- **JWT 认证**: 24 小时有效期，支持登出黑名单
- **权限控制**: 细粒度的访问控制，未审核作品仅管理员可见
- **文件代理访问**: 通过权限验证的代理接口访问图片，禁止直接 URL 访问
- **速率限制**: 验证码发送、登录尝试、文件上传的频率限制
- **文件验证**: 限制文件大小和类型，验证文件内容

## 技术栈

- **Web 框架**: Gin v1.11+
- **ORM**: Gorm v1.31+
- **数据库**: MySQL 8.0+
- **缓存**: Redis 7.0+
- **密码加密**: Bcrypt
- **JWT**: golang-jwt/jwt v5
- **配置管理**: Viper
- **日志**: Zap

## 项目结构

```
art-collection-system/
├── cmd/
│   └── server/          # 应用入口
├── internal/
│   ├── config/          # 配置管理
│   ├── models/          # 数据模型 (User, Activity, Artwork)
│   ├── repository/      # 数据访问层
│   ├── service/         # 业务逻辑层
│   ├── handler/         # HTTP 处理器
│   ├── middleware/      # 中间件 (认证、权限、CORS、日志、速率限制)
│   ├── utils/           # 工具函数 (JWT, 密码, 验证, 响应)
│   └── database/        # 数据库连接 (MySQL, Redis)
├── config/              # 配置文件
├── uploads/             # 文件上传目录 (私有)
├── logs/                # 日志目录
├── scripts/             # 脚本和工具
│   ├── init_db.sql      # 数据库初始化脚本
│   ├── docker-compose.yml # Docker 开发环境
│   └── QUICKSTART.md    # 快速开始指南
└── docs/                # 文档
    ├── api.md           # API 文档
    └── deployment.md    # 部署文档
```

## 快速开始

### 前置要求

- Go 1.21+
- MySQL 8.0+
- Redis 7.0+

### 使用 Docker 快速启动（推荐）

最简单的方式是使用 Docker Compose 启动开发环境：

```bash
# 启动 MySQL 和 Redis
cd scripts
docker-compose up -d

# 等待数据库启动完成
sleep 10

# 初始化数据库（创建默认管理员账户）
mysql -h 127.0.0.1 -P 3306 -u root -proot123456 art_collection < init_db.sql

# 返回项目根目录
cd ..

# 安装依赖
go mod download

# 运行应用
go run cmd/server/main.go
```

详细的快速开始指南请参考 [scripts/QUICKSTART.md](scripts/QUICKSTART.md)。

### 手动安装

#### 1. 安装依赖

```bash
go mod download
```

#### 2. 配置数据库

创建 MySQL 数据库：

```sql
CREATE DATABASE art_collection CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 3. 配置应用

复制配置文件模板并修改：

```bash
cp config/config.example.yaml config/config.yaml
```

编辑 `config/config.yaml` 并设置：

- 数据库连接信息（MySQL 和 Redis）
- JWT 密钥（至少 32 字节）
- 文件上传路径和大小限制
- 日志配置

#### 4. 初始化数据库

运行初始化脚本创建默认管理员账户：

```bash
mysql -h localhost -u root -p art_collection < scripts/init_db.sql
```

默认管理员账户：

- 邮箱: `admin@example.com`
- 密码: `Admin123456`

#### 5. 运行应用

```bash
go run cmd/server/main.go
```

应用将在 `http://localhost:8080` 启动。

### 验证安装

测试健康检查端点：

```bash
curl http://localhost:8080/health
```

测试登录：

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"Admin123456"}'
```

## API 文档

完整的 API 文档请参考 [docs/api.md](docs/api.md)。

### 主要端点

#### 认证

- `POST /api/v1/auth/send-code` - 发送验证码
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户登出

#### 用户

- `GET /api/v1/user/profile` - 获取个人信息
- `PUT /api/v1/user/profile` - 更新个人信息
- `PUT /api/v1/user/password` - 修改密码
- `GET /api/v1/users/:id/artworks` - 获取用户作品列表

#### 活动

- `GET /api/v1/activities` - 获取活动列表
- `GET /api/v1/activities/:id` - 获取活动详情
- `POST /api/v1/admin/activities` - 创建活动（管理员）
- `PUT /api/v1/admin/activities/:id` - 更新活动（管理员）
- `DELETE /api/v1/admin/activities/:id` - 删除活动（管理员）

#### 作品

- `POST /api/v1/artworks` - 上传作品
- `GET /api/v1/artworks/:id` - 获取作品信息
- `GET /api/v1/artworks/:id/image` - 获取作品图片
- `DELETE /api/v1/artworks/:id` - 删除作品

#### 管理员

- `GET /api/v1/admin/review-queue` - 获取审核队列
- `PUT /api/v1/admin/artworks/:id/review` - 审核作品
- `PUT /api/v1/admin/artworks/batch-review` - 批量审核作品
- `GET /api/v1/admin/users` - 获取用户列表
- `PUT /api/v1/admin/users/:id/role` - 更新用户角色
- `GET /api/v1/admin/users/:id/statistics` - 获取用户统计

## 开发

### 生成密码哈希

使用提供的工具生成 Bcrypt 密码哈希：

```bash
go run scripts/generate_password.go "YourPassword123"
```

### 验证密码

验证密码是否匹配哈希值：

```bash
go run scripts/verify_password.go "YourPassword123" "$2a$12$..."
```

### 运行测试

```bash
go test ./...
```

### 代码格式化

```bash
go fmt ./...
```

## 部署

生产环境部署指南请参考 [docs/deployment.md](docs/deployment.md)。

### 使用 Nginx 反向代理

系统提供了 Nginx 配置示例 `nginx.conf.example`，包括：

- 反向代理到 Gin 应用
- 禁止直接访问 `/uploads` 目录
- HTTPS 配置
- Gzip 压缩

### 环境变量

可以通过环境变量覆盖配置文件中的设置：

```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=password
export DB_NAME=art_collection
export REDIS_HOST=localhost
export REDIS_PORT=6379
export JWT_SECRET=your-secret-key
```

## 安全建议

1. **修改默认管理员密码**: 首次登录后立即修改默认管理员密码
2. **使用强 JWT 密钥**: 生成至少 32 字节的随机密钥
3. **启用 HTTPS**: 生产环境必须使用 HTTPS
4. **配置防火墙**: 限制数据库和 Redis 的访问
5. **定期备份**: 定期备份 MySQL 数据库和上传文件
6. **监控日志**: 监控应用日志和访问日志，及时发现异常
7. **更新依赖**: 定期更新 Go 依赖包，修复安全漏洞

## 故障排查

### 数据库连接失败

检查 MySQL 是否运行：

```bash
mysql -h localhost -u root -p
```

检查配置文件中的数据库连接信息是否正确。

### Redis 连接失败

检查 Redis 是否运行：

```bash
redis-cli ping
```

### 文件上传失败

检查 `uploads` 目录是否存在且有写权限：

```bash
mkdir -p uploads
chmod 755 uploads
```

### JWT 验证失败

确保 JWT 密钥在配置文件中正确设置，且与生成 token 时使用的密钥一致。

## 文档

- [API 文档](docs/api.md) - 完整的 API 接口文档
- [部署文档](docs/deployment.md) - 生产环境部署指南
- [快速开始](scripts/QUICKSTART.md) - 数据库快速开始指南

## 贡献

欢迎提交 Issue 和 Pull Request。

## 许可证

MIT License
