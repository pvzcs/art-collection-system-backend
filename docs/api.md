# API 文档

## 概述

美术作品收集系统 RESTful API 文档。所有 API 端点使用 JSON 格式进行数据交换。

**Base URL**: `http://localhost:8080/api/v1`

## 认证

大部分 API 端点需要 JWT 认证。在请求头中包含 `Authorization` 字段：

```
Authorization: Bearer <your_jwt_token>
```

JWT 令牌通过登录接口获取，有效期为 24 小时。

## 统一响应格式

### 成功响应

```json
{
  "code": 0,
  "message": "success",
  "data": {
    // 响应数据
  }
}
```

### 错误响应

```json
{
  "code": 400,
  "message": "错误描述"
}
```

## 错误码说明

| HTTP 状态码 | 说明 |
|------------|------|
| 200 | 请求成功 |
| 400 | 请求参数错误或业务逻辑错误 |
| 401 | 未授权，需要登录或 token 无效 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

## API 端点

### 认证相关

#### 1. 发送验证码

发送注册验证码到指定邮箱。

**端点**: `POST /auth/send-code`

**请求头**: 无需认证

**请求体**:
```json
{
  "email": "user@example.com"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "验证码已发送"
}
```

**错误**:
- `400`: 邮箱格式错误
- `429`: 发送频率过快（每分钟最多 1 次）

**速率限制**: 每个邮箱每分钟最多 1 次

---

#### 2. 用户注册

使用邮箱验证码注册新用户。

**端点**: `POST /auth/register`

**请求头**: 无需认证

**请求体**:
```json
{
  "email": "user@example.com",
  "code": "123456",
  "password": "Password123",
  "nickname": "用户昵称"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "注册成功",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "nickname": "用户昵称",
    "role": "user",
    "created_at": "2025-10-21T10:00:00Z"
  }
}
```

**错误**:
- `400`: 参数验证失败、验证码错误或已过期、邮箱已存在
- `500`: 服务器错误

**字段说明**:
- `email`: 邮箱地址，必须唯一
- `code`: 6 位验证码
- `password`: 密码，最小长度 8 位
- `nickname`: 用户昵称

---

#### 3. 用户登录

使用邮箱和密码登录，获取 JWT 令牌。

**端点**: `POST /auth/login`

**请求头**: 无需认证

**请求体**:
```json
{
  "email": "user@example.com",
  "password": "Password123"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "nickname": "用户昵称",
      "role": "user"
    }
  }
}
```

**错误**:
- `400`: 邮箱或密码错误
- `429`: 登录尝试过于频繁（每个 IP 每分钟最多 5 次）

**速率限制**: 每个 IP 每分钟最多 5 次

---

#### 4. 用户登出

登出当前用户，将 JWT 令牌加入黑名单。

**端点**: `POST /auth/logout`

**请求头**: 需要认证

**请求体**: 无

**响应**:
```json
{
  "code": 0,
  "message": "登出成功"
}
```

**错误**:
- `401`: 未授权

---

### 用户相关

#### 5. 获取个人信息

获取当前登录用户的个人信息。

**端点**: `GET /user/profile`

**请求头**: 需要认证

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "nickname": "用户昵称",
    "role": "user",
    "created_at": "2025-10-21T10:00:00Z",
    "updated_at": "2025-10-21T10:00:00Z"
  }
}
```

**错误**:
- `401`: 未授权

---

#### 6. 更新个人信息

更新当前登录用户的昵称。

**端点**: `PUT /user/profile`

**请求头**: 需要认证

**请求体**:
```json
{
  "nickname": "新昵称"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "更新成功"
}
```

**错误**:
- `400`: 参数验证失败
- `401`: 未授权

---

#### 7. 修改密码

修改当前登录用户的密码。

**端点**: `PUT /user/password`

**请求头**: 需要认证

**请求体**:
```json
{
  "old_password": "OldPassword123",
  "new_password": "NewPassword123"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "密码修改成功"
}
```

**错误**:
- `400`: 旧密码错误或新密码格式不正确
- `401`: 未授权

---

#### 8. 获取用户作品列表

获取指定用户的所有作品列表（个人空间）。

**端点**: `GET /users/:id/artworks`

**请求头**: 需要认证

**路径参数**:
- `id`: 用户 ID

**查询参数**:
- `page`: 页码，默认 1
- `page_size`: 每页数量，默认 20

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "artworks": [
      {
        "id": 1,
        "activity_id": 1,
        "user_id": 1,
        "file_name": "artwork.jpg",
        "review_status": "approved",
        "created_at": "2025-10-21T10:00:00Z",
        "activity": {
          "id": 1,
          "name": "活动名称"
        }
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

**权限**:
- 用户只能访问自己的作品列表
- 管理员可以访问任何用户的作品列表

**错误**:
- `401`: 未授权
- `403`: 权限不足（访问他人的作品列表）
- `404`: 用户不存在

---

### 活动相关

#### 9. 获取活动列表

获取所有未删除的活动列表。

**端点**: `GET /activities`

**请求头**: 需要认证

**查询参数**:
- `page`: 页码，默认 1
- `page_size`: 每页数量，默认 20

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "activities": [
      {
        "id": 1,
        "name": "活动名称",
        "deadline": "2025-12-31T23:59:59Z",
        "description": "活动详情（Markdown 格式）",
        "max_uploads_per_user": 5,
        "created_at": "2025-10-21T10:00:00Z",
        "updated_at": "2025-10-21T10:00:00Z"
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

**错误**:
- `401`: 未授权

---

#### 10. 获取活动详情

获取指定活动的详细信息。

**端点**: `GET /activities/:id`

**请求头**: 需要认证

**路径参数**:
- `id`: 活动 ID

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "活动名称",
    "deadline": "2025-12-31T23:59:59Z",
    "description": "活动详情（Markdown 格式）",
    "max_uploads_per_user": 5,
    "created_at": "2025-10-21T10:00:00Z",
    "updated_at": "2025-10-21T10:00:00Z"
  }
}
```

**错误**:
- `401`: 未授权
- `404`: 活动不存在或已删除

---

#### 11. 创建活动（管理员）

创建新的活动。

**端点**: `POST /admin/activities`

**请求头**: 需要认证（管理员）

**请求体**:
```json
{
  "name": "活动名称",
  "deadline": "2025-12-31T23:59:59Z",
  "description": "活动详情（Markdown 格式）",
  "max_uploads_per_user": 5
}
```

**响应**:
```json
{
  "code": 0,
  "message": "创建成功",
  "data": {
    "id": 1,
    "name": "活动名称",
    "deadline": "2025-12-31T23:59:59Z",
    "description": "活动详情（Markdown 格式）",
    "max_uploads_per_user": 5,
    "created_at": "2025-10-21T10:00:00Z"
  }
}
```

**字段说明**:
- `name`: 活动名称，必填
- `deadline`: 截止日期，可选（null 表示无截止日期）
- `description`: 活动详情，支持 Markdown 格式
- `max_uploads_per_user`: 单用户最大上传数量，默认 5

**错误**:
- `400`: 参数验证失败
- `401`: 未授权
- `403`: 权限不足（非管理员）

---

#### 12. 更新活动（管理员）

更新指定活动的信息。

**端点**: `PUT /admin/activities/:id`

**请求头**: 需要认证（管理员）

**路径参数**:
- `id`: 活动 ID

**请求体**:
```json
{
  "name": "新活动名称",
  "deadline": "2025-12-31T23:59:59Z",
  "description": "新活动详情",
  "max_uploads_per_user": 10
}
```

**响应**:
```json
{
  "code": 0,
  "message": "更新成功"
}
```

**错误**:
- `400`: 参数验证失败
- `401`: 未授权
- `403`: 权限不足（非管理员）
- `404`: 活动不存在

---

#### 13. 删除活动（管理员）

软删除指定活动。

**端点**: `DELETE /admin/activities/:id`

**请求头**: 需要认证（管理员）

**路径参数**:
- `id`: 活动 ID

**响应**:
```json
{
  "code": 0,
  "message": "删除成功"
}
```

**错误**:
- `401`: 未授权
- `403`: 权限不足（非管理员）
- `404`: 活动不存在

**注意**: 这是软删除，活动数据不会从数据库中物理删除。

---

### 作品相关

#### 14. 上传作品

上传美术作品到指定活动。

**端点**: `POST /artworks`

**请求头**: 
- 需要认证
- `Content-Type: multipart/form-data`

**表单字段**:
- `activity_id`: 活动 ID（整数）
- `file`: 作品文件（图片）

**响应**:
```json
{
  "code": 0,
  "message": "上传成功",
  "data": {
    "id": 1,
    "activity_id": 1,
    "user_id": 1,
    "file_name": "artwork.jpg",
    "review_status": "pending",
    "created_at": "2025-10-21T10:00:00Z"
  }
}
```

**文件限制**:
- 最大文件大小: 10MB
- 允许的文件类型: 图片格式（JPEG, PNG, GIF, WebP）

**错误**:
- `400`: 参数错误、活动不存在或已过期、超过上传数量限制、文件格式或大小不符合要求
- `401`: 未授权
- `429`: 上传频率过快（每个用户每分钟最多 10 次）

**速率限制**: 每个用户每分钟最多 10 次

---

#### 15. 获取作品信息

获取指定作品的详细信息。

**端点**: `GET /artworks/:id`

**请求头**: 需要认证

**路径参数**:
- `id`: 作品 ID

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "activity_id": 1,
    "user_id": 1,
    "file_name": "artwork.jpg",
    "review_status": "approved",
    "created_at": "2025-10-21T10:00:00Z",
    "updated_at": "2025-10-21T10:00:00Z",
    "activity": {
      "id": 1,
      "name": "活动名称"
    },
    "user": {
      "id": 1,
      "nickname": "用户昵称"
    }
  }
}
```

**权限**:
- 管理员可以访问所有作品
- 未审核作品（pending）仅管理员可见
- 已审核作品（approved）仅作者和管理员可见

**错误**:
- `401`: 未授权
- `403`: 权限不足
- `404`: 作品不存在

---

#### 16. 获取作品图片

通过代理接口获取作品图片文件。

**端点**: `GET /artworks/:id/image`

**请求头**: 需要认证

**路径参数**:
- `id`: 作品 ID

**响应**: 图片文件（二进制流）

**响应头**:
- `Content-Type`: 图片 MIME 类型（如 `image/jpeg`）
- `Content-Disposition`: `inline; filename="artwork.jpg"`

**权限**: 与"获取作品信息"接口相同

**错误**:
- `401`: 未授权
- `403`: 权限不足
- `404`: 作品或文件不存在

**注意**: 不能通过直接 URL 访问图片文件，必须通过此代理接口。

---

#### 17. 删除作品

删除自己上传的作品。

**端点**: `DELETE /artworks/:id`

**请求头**: 需要认证

**路径参数**:
- `id`: 作品 ID

**响应**:
```json
{
  "code": 0,
  "message": "删除成功"
}
```

**权限**:
- 用户只能删除自己上传的作品
- 管理员可以删除任何作品

**错误**:
- `401`: 未授权
- `403`: 权限不足（删除他人作品）
- `404`: 作品不存在

**注意**: 这是物理删除，作品记录和文件都会被删除。

---

### 管理员相关

#### 18. 获取审核队列

获取所有待审核作品列表。

**端点**: `GET /admin/review-queue`

**请求头**: 需要认证（管理员）

**查询参数**:
- `page`: 页码，默认 1
- `page_size`: 每页数量，默认 20

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "artworks": [
      {
        "id": 1,
        "activity_id": 1,
        "user_id": 1,
        "file_name": "artwork.jpg",
        "review_status": "pending",
        "created_at": "2025-10-21T10:00:00Z",
        "activity": {
          "id": 1,
          "name": "活动名称"
        },
        "user": {
          "id": 1,
          "nickname": "用户昵称",
          "email": "user@example.com"
        }
      }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
  }
}
```

**排序**: 按上传时间升序排列（最早上传的在前）

**错误**:
- `401`: 未授权
- `403`: 权限不足（非管理员）

---

#### 19. 审核作品

审核单个作品，更新审核状态。

**端点**: `PUT /admin/artworks/:id/review`

**请求头**: 需要认证（管理员）

**路径参数**:
- `id`: 作品 ID

**请求体**:
```json
{
  "approved": true
}
```

**响应**:
```json
{
  "code": 0,
  "message": "审核成功"
}
```

**字段说明**:
- `approved`: true 表示通过审核，false 表示保持未审核状态

**错误**:
- `400`: 参数错误
- `401`: 未授权
- `403`: 权限不足（非管理员）
- `404`: 作品不存在

---

#### 20. 批量审核作品

批量更新多个作品的审核状态。

**端点**: `PUT /admin/artworks/batch-review`

**请求头**: 需要认证（管理员）

**请求体**:
```json
{
  "artwork_ids": [1, 2, 3, 4, 5],
  "approved": true
}
```

**响应**:
```json
{
  "code": 0,
  "message": "批量审核成功",
  "data": {
    "updated_count": 5
  }
}
```

**字段说明**:
- `artwork_ids`: 作品 ID 数组
- `approved`: true 表示通过审核，false 表示保持未审核状态

**错误**:
- `400`: 参数错误
- `401`: 未授权
- `403`: 权限不足（非管理员）

---

#### 21. 获取用户列表

获取所有用户列表。

**端点**: `GET /admin/users`

**请求头**: 需要认证（管理员）

**查询参数**:
- `page`: 页码，默认 1
- `page_size`: 每页数量，默认 20

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "users": [
      {
        "id": 1,
        "email": "user@example.com",
        "nickname": "用户昵称",
        "role": "user",
        "created_at": "2025-10-21T10:00:00Z",
        "updated_at": "2025-10-21T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

**错误**:
- `401`: 未授权
- `403`: 权限不足（非管理员）

---

#### 22. 更新用户角色

更新指定用户的角色。

**端点**: `PUT /admin/users/:id/role`

**请求头**: 需要认证（管理员）

**路径参数**:
- `id`: 用户 ID

**请求体**:
```json
{
  "role": "admin"
}
```

**响应**:
```json
{
  "code": 0,
  "message": "角色更新成功"
}
```

**字段说明**:
- `role`: 用户角色，可选值: `user`, `admin`

**错误**:
- `400`: 参数错误（无效的角色值）
- `401`: 未授权
- `403`: 权限不足（非管理员）
- `404`: 用户不存在

---

#### 23. 获取用户统计信息

获取指定用户的统计信息。

**端点**: `GET /admin/users/:id/statistics`

**请求头**: 需要认证（管理员）

**路径参数**:
- `id`: 用户 ID

**响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "user_id": 1,
    "total_artworks": 15,
    "approved_artworks": 10,
    "pending_artworks": 5,
    "activities_participated": 3
  }
}
```

**错误**:
- `401`: 未授权
- `403`: 权限不足（非管理员）
- `404`: 用户不存在

---

## 使用示例

### 完整的用户注册和登录流程

```bash
# 1. 发送验证码
curl -X POST http://localhost:8080/api/v1/auth/send-code \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'

# 2. 注册用户
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email":"user@example.com",
    "code":"123456",
    "password":"Password123",
    "nickname":"测试用户"
  }'

# 3. 登录获取 token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email":"user@example.com",
    "password":"Password123"
  }'

# 响应示例
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "nickname": "测试用户",
      "role": "user"
    }
  }
}
```

### 上传作品

```bash
# 使用 token 上传作品
curl -X POST http://localhost:8080/api/v1/artworks \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -F "activity_id=1" \
  -F "file=@/path/to/artwork.jpg"
```

### 管理员审核作品

```bash
# 获取审核队列
curl -X GET http://localhost:8080/api/v1/admin/review-queue \
  -H "Authorization: Bearer <admin_token>"

# 审核单个作品
curl -X PUT http://localhost:8080/api/v1/admin/artworks/1/review \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{"approved":true}'

# 批量审核作品
curl -X PUT http://localhost:8080/api/v1/admin/artworks/batch-review \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "artwork_ids":[1,2,3,4,5],
    "approved":true
  }'
```

## 速率限制

为了防止滥用，以下接口有速率限制：

| 接口 | 限制 |
|------|------|
| 发送验证码 | 每个邮箱每分钟 1 次 |
| 登录 | 每个 IP 每分钟 5 次 |
| 上传作品 | 每个用户每分钟 10 次 |

超过速率限制将返回 `429 Too Many Requests` 错误。

## 安全建议

1. **使用 HTTPS**: 生产环境必须使用 HTTPS 传输，保护 JWT 令牌和敏感数据
2. **保护 JWT 令牌**: 不要在 URL 中传递 token，不要在客户端存储明文 token
3. **定期刷新令牌**: JWT 有效期为 24 小时，过期后需要重新登录
4. **登出时清理令牌**: 用户登出后应清理客户端存储的 token
5. **验证文件类型**: 上传文件时验证文件类型和内容，防止恶意文件上传
6. **输入验证**: 所有用户输入都应进行验证和清理

## 附录

### 审核状态枚举

| 值 | 说明 |
|----|------|
| `pending` | 未审核 |
| `approved` | 已审核通过 |

### 用户角色枚举

| 值 | 说明 |
|----|------|
| `user` | 普通用户 |
| `admin` | 管理员 |

### 常见问题

**Q: JWT 令牌过期后如何处理？**

A: 令牌过期后会返回 401 错误，需要重新登录获取新的令牌。

**Q: 如何访问作品图片？**

A: 必须通过 `/artworks/:id/image` 接口访问，不能直接访问文件 URL。

**Q: 未审核的作品作者能看到吗？**

A: 不能。根据需求，未审核作品只有管理员可见。

**Q: 删除活动后作品会怎样？**

A: 活动是软删除，作品数据不会丢失，但用户无法再向该活动上传新作品。

**Q: 如何修改默认管理员密码？**

A: 使用管理员账户登录后，调用修改密码接口。
