-- Active: 1717141082295@@192.168.1.27@3306@moremei-ai-saas-tenant-01
CREATE TABLE user_check_in_record(
    id BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    user_id BIGINT(20) UNSIGNED NOT NULL COMMENT '用户ID',
    checked_time VARCHAR(32) NOT NULL COMMENT '签到时间, yyyy-mm-dd',
    city VARCHAR(64) COMMENT '签到当前所在城市',
    latitude VARCHAR(32) COMMENT '纬度',
    longitude VARCHAR(32) COMMENT '经度',
    weather VARCHAR(32) COMMENT '天气',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY unique_user_time (user_id, checked_time)
) COMMENT '签到记录';
