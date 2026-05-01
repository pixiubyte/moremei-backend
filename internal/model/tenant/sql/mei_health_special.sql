CREATE TABLE `mei_health_special`(
    id BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `type` VARCHAR(64) NOT NULL COMMENT '类型',
    `title` VARCHAR(64) NOT NULL COMMENT '标题',
    `profile` TEXT NOT NULL COMMENT '专题简介',
    `theme` VARCHAR(1024) NOT NULL COMMENT '专题主题(颜色)',
    `background` VARCHAR(1024) NOT NULL COMMENT '专题背景图',
    `endorsers` JSON NOT NULL COMMENT '背书专家',
    `videos` JSON NOT NULL COMMENT '专题教育视频',
    `carers` JSON NOT NULL COMMENT '健康管理师',
    `archive_fields` JSON NOT NULL COMMENT '本专题显示的档案字段',
    `questionnaire` VARCHAR(64) NOT NULL COMMENT '健康问卷UUID',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态',
    `extra` JSON COMMENT '扩展字段',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY `type_key` (`type`)
) COMMENT '健康专题';

INSERT INTO `mei_health_special` (`type`, `title`, `profile`, `theme`, `background`, `endorsers`, `videos`, `carers`, `archive_fields`, `questionnaire`, `status`) VALUES 
('nxjk', '女性生殖健康', '帮助女性朋友关注女性健康', '#27C8CB', '#CDF1EE', '[1]', '[]', '[]', '{"birth_num": {"name": "生育次数"}, "fetaion_num": {"name": "怀孕次数"}}', '1740277532109836291', '1')
