# 数据库快速启动指南

## 选择启动方式

### 🐳 方式 A: 使用 Docker（最简单，推荐）

适合：开发环境、快速测试

```bash
# 一键启动 MySQL + Redis
docker-compose -f scripts/docker-compose.yml up -d

# 等待几秒让服务启动完成
sleep 5

# 验证服务状态
docker-compose -f scripts/docker-compose.yml ps
```

✅ 完成！数据库已自动初始化，跳转到"默认账户信息"部分。

连接信息：
- MySQL Host: localhost:3306
- MySQL User: artuser
- MySQL Password: artpassword
- Redis Host: localhost:6379

### 💻 方式 B: 使用本地 MySQL

适合：生产环境、已有 MySQL 服务

#### 步骤 1: 确保 MySQL 正在运行

```bash
# 检查 MySQL 服务状态
mysql --version

# 如果未安装，请先安装 MySQL 8.0+
```

#### 步骤 2: 执行初始化脚本

```bash
# 方法 1: 直接执行（推荐）
mysql -u root -p < scripts/init_db.sql

# 方法 2: 登录后执行
mysql -u root -p
mysql> source scripts/init_db.sql;
mysql> exit;
```

#### 步骤 3: 验证初始化

```bash
mysql -u root -p -e "USE art_collection; SELECT email, nickname, role FROM users WHERE role='admin';"
```

预期输出：

```
+---------------------+------------------------+-------+
| email               | nickname               | role  |
+---------------------+------------------------+-------+
| admin@example.com   | System Administrator   | admin |
+---------------------+------------------------+-------+
```

## 默认账户信息

初始化完成后，您可以使用以下账户登录：

```
邮箱: admin@example.com
密码: Admin123456
角色: 管理员
```

⚠️ **安全警告**: 请在首次登录后立即修改默认密码！

## 修改管理员密码

### 方法 1: 通过 API 修改（推荐）

1. 启动应用服务器
2. 使用默认账户登录
3. 调用修改密码 API

### 方法 2: 通过数据库修改

```bash
# 1. 生成新密码的哈希值
go run scripts/generate_password.go YourNewPassword123

# 2. 更新数据库
mysql -u root -p art_collection -e "UPDATE users SET password='生成的哈希值' WHERE email='admin@example.com';"
```

## 创建额外的管理员账户

```bash
# 1. 生成密码哈希
go run scripts/generate_password.go SecurePassword456

# 2. 插入新管理员
mysql -u root -p art_collection << EOF
INSERT INTO users (email, password, nickname, role, created_at, updated_at)
VALUES (
  'admin2@example.com',
  '生成的哈希值',
  'Second Administrator',
  'admin',
  NOW(),
  NOW()
);
EOF
```

## 查看示例活动

```bash
mysql -u root -p art_collection -e "SELECT id, name, deadline, max_uploads_per_user FROM activities WHERE is_deleted=0;"
```

## 常见问题

### Q: 数据库已存在，如何重新初始化？

```bash
# ⚠️ 警告：这将删除所有数据！
mysql -u root -p -e "DROP DATABASE IF EXISTS art_collection;"
mysql -u root -p < scripts/init_db.sql
```

### Q: 如何只插入默认管理员，不创建示例活动？

编辑 `scripts/init_db.sql`，删除或注释掉示例活动的 INSERT 语句。

### Q: 忘记管理员密码怎么办？

```bash
# 重置为默认密码 Admin123456
mysql -u root -p art_collection -e "UPDATE users SET password='$2a$12$p9iQ0R6DnMVB4U50Fk/45ejxqc.dl3XielMdYytWxu/f/R7BD/y1C' WHERE email='admin@example.com';"
```

### Q: 如何检查数据库连接？

```bash
mysql -u root -p -e "SELECT VERSION(); SHOW DATABASES LIKE 'art_collection';"
```

## 下一步

数据库初始化完成后：

1. 配置 `config/config.yaml` 中的数据库连接信息
2. 启动应用服务器：`go run cmd/server/main.go`
3. 使用默认管理员账户登录
4. 修改默认密码
5. 开始使用系统！
