CREATE TABLE `in_store_messages` (
    `id` BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    `uuid` VARCHAR(64) NOT NULL COMMENT 'uuid',
    `conversation_id` VARCHAR(64) NOT NULL COMMENT '会话id',
    `parent_id` BIGINT(20) unsigned NOT NULL DEFAULT 0 COMMENT '父消息ID',
    `author_type` VARCHAR(12) NOT NULL COMMENT '所属人类型',
    `author_id` VARCHAR(128) NOT NULL DEFAULT "" COMMENT '所属人ID',
    `type` int(11) NOT NULL DEFAULT '0' COMMENT '类型：0-未知，1-文本，2-图片，3-语音，5-视频，7-小程序，8-链接，9-文件，15-引用',
    `content` TEXT NOT NULL COMMENT '会话内容ID',
    `q_flag` tinyint(4) NOT NULL DEFAULT '2' COMMENT '问题标记：0-非问题，1-问题，2-未知',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `uuid_key` (`uuid`)
) COMMENT '线下场景消息表';

CREATE TABLE `in_store_conversations` (
    `id` BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    `uuid` VARCHAR(64) NOT NULL COMMENT 'uuid',
    `sn` VARCHAR(64) NOT NULL COMMENT '设备SN',
    `ai_robot_id` BIGINT(20) unsigned NOT NULL COMMENT 'AI机器人ID',
    `cid` VARCHAR(64) NOT NULL COMMENT 'AI会话ID',
    `title` VARCHAR(255) NOT NULL COMMENT '会话标题',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `uuid_key` (`uuid`)
) COMMENT '线下场景会话表';

CREATE TABLE `in_store_sn_whitelist` (
    `id` BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT '主键',
    `sn` VARCHAR(64) NOT NULL COMMENT '设备SN',
    `description` VARCHAR(1024) NOT NULL COMMENT '描述',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `sn_key` (`sn`)
) COMMENT '线下场景设备SN白名单表';