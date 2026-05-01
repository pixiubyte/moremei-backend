CREATE TABLE `mei_product_usage_plan`(
    `id` BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `user_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '用户ID',
    `template_id` BIGINT(20) UNSIGNED NOT NULL COMMENT '模板ID',
    `product_name` VARCHAR(64) NOT NULL COMMENT '产品名称',
    `start_at` DATE NOT NULL COMMENT '预计开始时间',
    `instruction` TEXT NOT NULL COMMENT '使用说明',
    `use_steps` TEXT NOT NULL COMMENT '使用步骤说明',
    `steps` JSON NOT NULL COMMENT '使用步骤',
    `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '状态 0: 未激活 1: 激活 2: 已完成 -1: 已过期（未使用）',
    `expired_at` DATE DEFAULT NULL COMMENT '过期时间',
    `mark` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '备注',
    `extra` JSON COMMENT '扩展字段',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT '用户使用计划';