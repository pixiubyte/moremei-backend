CREATE TABLE `mei_product_usage_plan_template`(
    `id` BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT 'Primary Key',
    `product_id` VARCHAR(64) NOT NULL COMMENT '产品ID',
    `type` VARCHAR(64) NOT NULL COMMENT '使用计划产品类型',
    `notice` TEXT NOT NULL COMMENT '用前须知',
    `instruction` TEXT NOT NULL COMMENT '使用说明',
    `use_steps` TEXT NOT NULL COMMENT '使用步骤说明',
    `pre_instruction` TEXT NOT NULL COMMENT '使用前教育内容',
    `steps` JSON NOT NULL COMMENT '使用步骤',
    `videos` JSON NOT NULL COMMENT '教育视频',
    `expire_time` INT(11) NOT NULL DEFAULT -1 COMMENT '使用期限',
    `extra` JSON COMMENT '扩展字段',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间'
) COMMENT '产品使用计划模板';
