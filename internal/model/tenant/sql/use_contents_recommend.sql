CREATE TABLE user_contents_template_modules(
    `id` BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `type` VARCHAR(64) NOT NULL COMMENT '内容分类',
    `class` VARCHAR(64) NOT NULL COMMENT '细分类型（模块）',
    `content` JSON NOT NULL COMMENT '内容',
    `is_hide` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否隐藏',
    `is_public` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否公开',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `userid_type_class_key` (`type`, `class`)
) COMMENT '用户内容模板表';

CREATE TABLE user_contents_recommend(
    `id` BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `userid` BIGINT(20) UNSIGNED NOT NULL COMMENT '用户ID',
    `type` VARCHAR(64) NOT NULL COMMENT '内容分类',
    `class` VARCHAR(64) NOT NULL COMMENT '细分类型（模块）',
    `content` JSON NOT NULL COMMENT '内容',
    `is_hide` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否隐藏',
    `is_public` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否公开',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `userid_type_class_key` (`user_id`, `type`, `class`)
) COMMENT '用户内容推荐表';

CREATE TABLE user_contents_recommend_history(
    `id` BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `userid` BIGINT(20) UNSIGNED NOT NULL COMMENT '用户ID',
    `type` VARCHAR(64) NOT NULL COMMENT '内容分类',
    `class` VARCHAR(64) NOT NULL COMMENT '细分类型（模块）',
    `content` JSON NOT NULL COMMENT '内容',
    `is_hide` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否隐藏',
    `is_public` TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否公开',
    `history_type` VARCHAR(64) NOT NULL COMMENT '历史类型',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `userid_type_class_key` (`user_id`, `type`, `class`)
) COMMENT '用户内容推荐历史表';