# 数据库初始化脚本

本目录包含用于初始化美术作品收集系统数据库的脚本和工具。

## 文件说明

- `init_db.sql` - 数据库初始化 SQL 脚本
- `generate_password.go` - 密码哈希生成工具
- `verify_password.go` - 密码哈希验证工具
- `docker-compose.yml` - Docker Compose 配置文件
- `QUICKSTART.md` - 快速启动指南
- `README.md` - 本文档

## 使用方法

### 方式 1: 使用 Docker Compose（推荐用于开发环境）

如果您使用 Docker，可以快速启动 MySQL 和 Redis：

```bash
# 启动服务（会自动执行初始化脚本）
docker-compose -f scripts/docker-compose.yml up -d

# 查看日志
docker-compose -f scripts/docker-compose.yml logs -f

# 停止服务
docker-compose -f scripts/docker-compose.yml down

# 停止并删除数据卷（⚠️ 会删除所有数据）
docker-compose -f scripts/docker-compose.yml down -v
```

Docker Compose 配置：
- MySQL 端口: 3306
- Redis 端口: 6379
- MySQL root 密码: rootpassword
- MySQL 用户: artuser / artpassword
- 数据库名: art_collection

### 方式 2: 使用本地 MySQL

使用 MySQL 客户端执行初始化脚本：

```bash
mysql -u root -p < scripts/init_db.sql
```

或者登录 MySQL 后执行：

```sql
source scripts/init_db.sql;
```

### 2. 默认管理员账户

初始化脚本会创建一个默认管理员账户：

- **邮箱**: `admin@example.com`
- **密码**: `Admin123456`
- **角色**: `admin`

⚠️ **重要安全提示**: 请在生产环境中立即修改默认密码！

### 3. 示例活动数据

初始化脚本会创建两个示例活动：

1. **2024 春季美术作品征集** - 截止日期：30 天后
2. **夏日创意绘画大赛** - 截止日期：60 天后

这些示例数据可以帮助您快速测试系统功能。

## 生成自定义密码哈希

如果您需要创建其他管理员账户或修改默认密码，可以使用密码哈希生成工具：

```bash
go run scripts/generate_password.go YourPassword123
```

输出示例：

```
密码哈希生成成功！
原始密码: YourPassword123
密码哈希: $2a$12$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

您可以在 SQL 脚本中使用此哈希值。
```

### 验证密码哈希

如果您想验证密码哈希是否正确，可以使用验证工具：

```bash
go run scripts/verify_password.go Admin123456 '$2a$12$p9iQ0R6DnMVB4U50Fk/45ejxqc.dl3XielMdYytWxu/f/R7BD/y1C'
```

输出示例：

```
✅ 密码验证成功！
密码: Admin123456
哈希: $2a$12$p9iQ0R6DnMVB4U50Fk/45ejxqc.dl3XielMdYytWxu/f/R7BD/y1C
```

### 创建新管理员账户

然后您可以使用生成的哈希值在数据库中创建或更新用户：

```sql
INSERT INTO users (email, password, nickname, role, created_at, updated_at)
VALUES (
  'newadmin@example.com',
  '$2a$12$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx',
  'New Administrator',
  'admin',
  NOW(),
  NOW()
);
```

## 数据库表结构

### users 表

| 字段       | 类型                 | 说明               |
| ---------- | -------------------- | ------------------ |
| id         | bigint unsigned      | 主键，自增         |
| email      | varchar(255)         | 邮箱地址，唯一索引 |
| password   | varchar(255)         | 密码哈希（bcrypt） |
| nickname   | varchar(100)         | 用户昵称           |
| role       | enum('user','admin') | 用户角色           |
| created_at | datetime(3)          | 创建时间           |
| updated_at | datetime(3)          | 更新时间           |

### activities 表

| 字段                 | 类型            | 说明                      |
| -------------------- | --------------- | ------------------------- |
| id                   | bigint unsigned | 主键，自增                |
| name                 | varchar(255)    | 活动名称                  |
| deadline             | datetime(3)     | 截止日期（可为空）        |
| description          | text            | 活动详情（Markdown 格式） |
| max_uploads_per_user | int             | 单用户最大上传数量        |
| is_deleted           | tinyint(1)      | 软删除标记                |
| created_at           | datetime(3)     | 创建时间                  |
| updated_at           | datetime(3)     | 更新时间                  |

### artworks 表

| 字段          | 类型                       | 说明          |
| ------------- | -------------------------- | ------------- |
| id            | bigint unsigned            | 主键，自增    |
| activity_id   | bigint unsigned            | 活动 ID，外键 |
| user_id       | bigint unsigned            | 用户 ID，外键 |
| file_path     | varchar(500)               | 文件存储路径  |
| file_name     | varchar(255)               | 文件名        |
| review_status | enum('pending','approved') | 审核状态      |
| created_at    | datetime(3)                | 创建时间      |
| updated_at    | datetime(3)                | 更新时间      |

## 索引说明

- `users.email` - 唯一索引，用于登录查询
- `activities.is_deleted` - 普通索引，用于过滤已删除活动
- `artworks.activity_id` - 普通索引，用于按活动查询
- `artworks.user_id` - 普通索引，用于按用户查询
- `artworks.(user_id, activity_id)` - 复合索引，用于检查上传限制
- `artworks.review_status` - 普通索引，用于审核队列查询
- `artworks.created_at` - 普通索引，用于按时间排序

## 外键约束

- `artworks.activity_id` → `activities.id`
- `artworks.user_id` → `users.id`

## 注意事项

1. **字符集**: 所有表使用 `utf8mb4` 字符集和 `utf8mb4_unicode_ci` 排序规则，支持完整的 Unicode 字符（包括 emoji）

2. **时间精度**: 所有时间字段使用 `datetime(3)` 类型，支持毫秒级精度

3. **密码安全**: 密码使用 bcrypt 算法加密，cost 参数为 12

4. **软删除**: 活动表使用 `is_deleted` 字段实现软删除，不会物理删除数据

5. **枚举类型**: `role` 和 `review_status` 字段使用 MySQL 的 ENUM 类型，确保数据一致性

## 重置数据库

如果需要完全重置数据库（⚠️ 会删除所有数据）：

```sql
DROP DATABASE IF EXISTS art_collection;
```

然后重新执行初始化脚本。

## 备份建议

在生产环境中，建议定期备份数据库：

```bash
# 备份数据库
mysqldump -u root -p art_collection > backup_$(date +%Y%m%d_%H%M%S).sql

# 恢复数据库
mysql -u root -p art_collection < backup_20240101_120000.sql
```

## 故障排查

### 问题：外键约束错误

如果遇到外键约束错误，请确保：

1. MySQL 版本支持 InnoDB 引擎
2. 表的创建顺序正确（先创建父表，再创建子表）
3. 外键字段的类型和父表主键类型完全一致

### 问题：字符集问题

如果遇到中文乱码，请检查：

1. MySQL 服务器字符集配置
2. 客户端连接字符集
3. 数据库和表的字符集设置

可以在 MySQL 配置文件中添加：

```ini
[mysqld]
character-set-server=utf8mb4
collation-server=utf8mb4_unicode_ci

[client]
default-character-set=utf8mb4
```
