# 美术作品收集系统 (Art Collection System)

基于活动的美术作品管理平台，使用 Go + Gin + Gorm + MySQL + Redis 构建。

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
│   ├── models/          # 数据模型
│   ├── repository/      # 数据访问层
│   ├── service/         # 业务逻辑层
│   ├── handler/         # HTTP 处理器
│   ├── middleware/      # 中间件
│   ├── utils/           # 工具函数
│   └── database/        # 数据库连接
├── config/              # 配置文件
├── uploads/             # 文件上传目录
└── logs/                # 日志目录
```

## 快速开始

### 前置要求

- Go 1.21+
- MySQL 8.0+
- Redis 7.0+

### 安装依赖

```bash
go mod download
```

### 配置

复制配置文件模板并修改：

```bash
cp config/config.yaml config/config.local.yaml
```

编辑 `config/config.local.yaml` 并设置数据库连接信息、JWT 密钥等。

### 运行

```bash
go run cmd/server/main.go
```

## 开发状态

项目当前处于初始化阶段，正在按照实现计划逐步开发功能。

## 文档

- [需求文档](.kiro/specs/art-collection-system/requirements.md)
- [设计文档](.kiro/specs/art-collection-system/design.md)
- [实现计划](.kiro/specs/art-collection-system/tasks.md)

## 许可证

待定
