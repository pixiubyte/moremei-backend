-- Active: 1717141082295@@192.168.1.27@3306@moremei-ai-saas-tenant-01
CREATE TABLE app_user_profile(
    `id` BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `app` VARCHAR(32) NOT NULL COMMENT '应用',
    `uuid` VARCHAR(64) NOT NULL COMMENT '系统用户UUID',
    `key` VARCHAR(64) NOT NULL COMMENT '键',
    `value` JSON NOT NULL COMMENT '值，用户信息',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `app_uuid_key` (`app`, `uuid`, `key`)
) COMMENT '应用特定用户属性表，如应用设置等应用内用户信息';