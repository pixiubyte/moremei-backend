CREATE TABLE `tracking_points_data` (
    `id` BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `user_uuid` VARCHAR(255) NOT NULL COMMENT '用户 UUID',
    `page_path` VARCHAR(1024) NOT NULL COMMENT '页面路径',
    `page_params` TEXT NOT NULL COMMENT '页面参数',
    `page_duration` INT NOT NULL COMMENT '页面停留时间',
    `event_name` VARCHAR(255) NOT NULL COMMENT '事件名称',
    `event_params` TEXT NOT NULL COMMENT '事件参数',
    `event_time` DATETIME NOT NULL COMMENT '事件时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX `idx_user_uuid` (`user_uuid`)
) COMMENT='埋点数据';