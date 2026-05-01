-- Active: 1717141082295@@192.168.1.27@3306@moremei-ai-saas-tenant-01
CREATE TABLE auth_provider_conf(
    `id` BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `type` VARCHAR(32) NOT NULL DEFAULT "" COMMENT '第三方提供商类型',
    `provider` VARCHAR(32) NOT NULL COMMENT '第三方登录提供商',
    `app` VARCHAR(64) NOT NULL COMMENT '系统内app标识' DEFAULT "",
    `config` JSON NOT NULL COMMENT '三方平台应用配置信息',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `provider_type_app_key` (`provider`, `type`, `app`)
) COMMENT '第三方登录平台配置';