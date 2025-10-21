-- 美术作品收集系统 - 数据库初始化脚本
-- 此脚本用于初始化数据库表结构和默认数据

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS art_collection CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE art_collection;

-- 创建用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `nickname` varchar(100) NOT NULL,
  `role` enum('user','admin') NOT NULL DEFAULT 'user',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建活动表
CREATE TABLE IF NOT EXISTS `activities` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `deadline` datetime(3) DEFAULT NULL,
  `description` text,
  `max_uploads_per_user` int NOT NULL DEFAULT '5',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_activities_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建作品表
CREATE TABLE IF NOT EXISTS `artworks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `activity_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `file_path` varchar(500) NOT NULL,
  `file_name` varchar(255) NOT NULL,
  `review_status` enum('pending','approved') NOT NULL DEFAULT 'pending',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_activity_id` (`activity_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_user_activity` (`user_id`,`activity_id`),
  KEY `idx_review_status` (`review_status`),
  KEY `idx_artworks_created_at` (`created_at`),
  CONSTRAINT `fk_activities_artworks` FOREIGN KEY (`activity_id`) REFERENCES `activities` (`id`),
  CONSTRAINT `fk_users_artworks` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认管理员账户
-- 邮箱: admin@example.com
-- 密码: Admin123456
-- 注意: 请在生产环境中修改默认密码！
-- 密码哈希使用 bcrypt cost=12 生成
INSERT INTO `users` (`email`, `password`, `nickname`, `role`, `created_at`, `updated_at`)
VALUES (
  'admin@example.com',
  '$2a$12$p9iQ0R6DnMVB4U50Fk/45ejxqc.dl3XielMdYytWxu/f/R7BD/y1C',
  'System Administrator',
  'admin',
  NOW(),
  NOW()
) ON DUPLICATE KEY UPDATE `email` = `email`;

-- 插入示例活动数据（可选）
INSERT INTO `activities` (`name`, `deadline`, `description`, `max_uploads_per_user`, `is_deleted`, `created_at`, `updated_at`)
VALUES (
  '2024春季美术作品征集',
  DATE_ADD(NOW(), INTERVAL 30 DAY),
  '# 2024春季美术作品征集活动\n\n欢迎参加我们的春季美术作品征集活动！\n\n## 活动主题\n春天的色彩\n\n## 作品要求\n- 原创作品\n- 图片格式：JPG、PNG\n- 文件大小：不超过10MB\n\n## 截止日期\n请在活动截止日期前提交您的作品。\n\n期待您的精彩作品！',
  5,
  0,
  NOW(),
  NOW()
) ON DUPLICATE KEY UPDATE `name` = `name`;

INSERT INTO `activities` (`name`, `deadline`, `description`, `max_uploads_per_user`, `is_deleted`, `created_at`, `updated_at`)
VALUES (
  '夏日创意绘画大赛',
  DATE_ADD(NOW(), INTERVAL 60 DAY),
  '# 夏日创意绘画大赛\n\n展现您的创意，用画笔描绘夏日的美好！\n\n## 活动说明\n本次大赛面向所有热爱绘画的朋友，不限年龄、不限风格。\n\n## 奖项设置\n- 一等奖：1名\n- 二等奖：3名\n- 三等奖：5名\n- 优秀奖：若干\n\n## 评选标准\n- 创意性：40%\n- 技巧性：30%\n- 主题契合度：30%\n\n祝您取得好成绩！',
  3,
  0,
  NOW(),
  NOW()
) ON DUPLICATE KEY UPDATE `name` = `name`;

-- 完成初始化
SELECT '数据库初始化完成！' AS message;
SELECT '默认管理员账户：' AS info;
SELECT 'Email: admin@example.com' AS email;
SELECT 'Password: Admin123456' AS password;
SELECT '⚠️  请在生产环境中立即修改默认密码！' AS warning;
