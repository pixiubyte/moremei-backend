CREATE TABLE `mei_product_usage_plan_step`(
    `id` BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `plan_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '计划ID',
    `title` VARCHAR(64) NOT NULL COMMENT '标题',
    `sort` INT(11) NOT NULL COMMENT '排序',
    `plan_at` DATE NOT NULL COMMENT '计划时间',
    `use_at` DATE NOT NULL COMMENT '实际使用时间',
    `options` JSON NOT NULL COMMENT '选项',
    `opted` JSON COMMENT '已选择选项',
    `images` JSON COMMENT '反馈图片',
    `description` TEXT NOT NULL COMMENT '反馈描述',
    `manager_remark` TEXT NOT NULL COMMENT '健管师备注',
    `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '状态 0: 还未反馈 1: 用户已反馈 -1: 已过期',
    `extra` JSON COMMENT '扩展字段',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT '用户使用计划步骤';