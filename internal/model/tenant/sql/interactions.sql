CREATE TABLE `content_interactions` (
    `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `user_id` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '用户ID',
    `cnt_type` VARCHAR(64) NOT NULL COMMENT '所属内容类型',
    `cnt_id` VARCHAR(64) NOT NULL COMMENT '所属内容ID',
    `cnt_extra` VARCHAR(64) NOT NULL COMMENT '所属内容额外信息,用于匹配',
    `type` TINYINT NOT NULL COMMENT '类型(1:支持,2:关心,3:鼓励,4:提醒,5:确认,6:短语)',
    `content` VARCHAR(255) NOT NULL COMMENT '内容',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `uk_user_id_cnt_type_cnt_id_extra_type_content` (`user_id`, `cnt_type`, `cnt_id`, `cnt_extra`, `type`, `content`)
) COMMENT '内容互动';
