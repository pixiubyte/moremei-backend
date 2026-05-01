-- Active: 1717141082295@@192.168.1.27@3306@moremei-ai-saas-tenant-01
CREATE TABLE auth_provider_user(
    `id` BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `provider` VARCHAR(32) NOT NULL COMMENT '第三方登录提供商',
    `app` VARCHAR(64) NOT NULL COMMENT '系统内app标识' DEFAULT "",
    `userid` BIGINT(20) UNSIGNED NOT NULL COMMENT '系统用户ID',
    `openid` VARCHAR(64) NOT NULL COMMENT '第三方用户ID',
    `name` VARCHAR(64) NOT NULL COMMENT '三方平台用户名',
    `unionid` VARCHAR(64) NOT NULL COMMENT '三方平台用户统一ID',
    `email` VARCHAR(128) NOT NULL COMMENT '三方平台用户邮箱',
    `avatar` VARCHAR(256) NOT NULL COMMENT '三方平台用户头像',
    `phone` VARCHAR(16) NOT NULL COMMENT '三方平台用户手机号',
    `extra` JSON NOT NULL COMMENT '三方平台用户扩展信息',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `provider_app_openid_key` (`provider`, `app`, `openid`),
    UNIQUE KEY `provider_app_user_open_unionid_key` (`provider`, `app`, `userid`, `openid`, `unionid`),
    KEY `provider_userid_key` (`provider`, `userid`),
    KEY `provider_app_unionid_key` (`provider`, `app`, `unionid`)
) COMMENT '第三方登录用户';